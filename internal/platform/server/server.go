// Package server merakit router chi: middleware global, endpoint sistem,
// dan mounting routes seluruh modul (peta endpoint di docs/lampiran-a-api.md).
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/ziaulkamal/fablerace/internal/modules/auth"
	"github.com/ziaulkamal/fablerace/internal/modules/content"
	"github.com/ziaulkamal/fablerace/internal/modules/keyword"
	"github.com/ziaulkamal/fablerace/internal/modules/license"
	"github.com/ziaulkamal/fablerace/internal/modules/media"
	"github.com/ziaulkamal/fablerace/internal/modules/taxonomy"
	"github.com/ziaulkamal/fablerace/internal/modules/tenant"
	"github.com/ziaulkamal/fablerace/internal/modules/user"
	"github.com/ziaulkamal/fablerace/internal/platform/config"
	"github.com/ziaulkamal/fablerace/internal/platform/middleware"
	"github.com/ziaulkamal/fablerace/internal/platform/storage"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// ModuleHandlers: seluruh handler modul, di-inject dari cmd/api/main.go —
// server tidak membangun dependensi sendiri (komposisi hanya di main).
type ModuleHandlers struct {
	Auth     *auth.Handler
	Tenant   *tenant.Handler
	User     *user.Handler
	Content  *content.Handler
	Taxonomy *taxonomy.Handler
	Media    *media.Handler
	Keyword  *keyword.Handler
	License  *license.Handler
}

func New(cfg *config.Config, pool *pgxpool.Pool, rdb *redis.Client, store *storage.Client, h ModuleHandlers) *http.Server {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(requestLogger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	// Endpoint sistem (docs/09-devops.md §4): liveness & readiness.
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		response.OK(w, map[string]string{"status": "ok"})
	})
	r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
		defer cancel()
		checks := map[string]string{"postgres": "ok", "redis": "ok", "minio": "ok"}
		healthy := true
		if err := pool.Ping(ctx); err != nil {
			checks["postgres"], healthy = err.Error(), false
		}
		if err := rdb.Ping(ctx).Err(); err != nil {
			checks["redis"], healthy = err.Error(), false
		}
		if err := store.Ping(ctx); err != nil {
			checks["minio"], healthy = err.Error(), false
		}
		if !healthy {
			response.JSON(w, http.StatusServiceUnavailable, checks)
			return
		}
		response.OK(w, checks)
	})

	r.Route("/api/v1", func(api chi.Router) {
		// Grup PUBLIK pra-login: hanya auth, dibungkus rate limit ketat
		// (RELASI: internal/modules/auth/routes.go).
		api.Group(func(pub chi.Router) {
			pub.Use(middleware.RateLimit(30))
			pub.Route("/auth", h.Auth.Routes())
		})

		// Grup USER: butuh JWT, belum butuh konteks tenant.
		api.Group(func(priv chi.Router) {
			priv.Use(middleware.Auth(cfg.JWTSecret))
			priv.Route("/tenants", h.Tenant.Routes())

			// Grup TENANT: + X-Tenant-ID terverifikasi + lisensi aktif.
			// Urutan middleware = kontrak docs/05-rbac.md §3.
			priv.Group(func(t chi.Router) {
				t.Use(middleware.Tenant(pool))
				t.Use(middleware.License(pool))

				t.Route("/tenants/current", h.Tenant.CurrentRoutes())
				t.Route("/users", h.User.Routes())
				t.Route("/contents", func(c chi.Router) {
					h.Content.Routes()(c)
					// Endpoint lintas-modul di bawah /contents — pemiliknya
					// modul taxonomy & keyword (lihat RELASI di handler masing-masing).
					c.With(middleware.RequirePermission("content.edit_own")).
						Put("/{id}/terms", h.Taxonomy.SetContentTerms)
					c.With(middleware.RequirePermission("content.edit_own")).
						Put("/{id}/keywords", h.Keyword.AttachToContent)
				})
				t.Route("/taxonomies", h.Taxonomy.Routes())
				t.Route("/media", h.Media.Routes())
				t.Route("/keywords", h.Keyword.Routes())
				t.Route("/license", h.License.Routes())
			})
		})

		// NEXT: delivery publik ber-API-key — /public/contents, /sitemap.xml,
		// /feed.xml (docs/lampiran-a-api.md) lewat middleware APIKey
		// (NEXT di internal/platform/middleware/auth.go).
		// NEXT: POST /billing/webhook publik (NEXT di internal/modules/license/handler.go).
	})

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

// requestLogger: structured log per request dengan request_id; tenant_id ikut
// tercatat bila sudah ter-resolve (docs/09-devops.md §4).
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", chimw.GetReqID(r.Context()),
		}
		if tenantID, ok := middleware.TenantID(r.Context()); ok {
			attrs = append(attrs, "tenant_id", tenantID.String())
		}
		slog.Info("http", attrs...)
	})
}
