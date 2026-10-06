// Command api adalah entrypoint Fablerace: satu-satunya tempat komposisi
// dependensi (config → infra → repository → service → handler → server).
// Modul tidak pernah saling membangun — semua interface di-inject dari sini
// (docs/lampiran-b-konvensi-kode.md §5).
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ziaulkamal/fablerace/internal/modules/auth"
	"github.com/ziaulkamal/fablerace/internal/modules/content"
	"github.com/ziaulkamal/fablerace/internal/modules/keyword"
	"github.com/ziaulkamal/fablerace/internal/modules/license"
	"github.com/ziaulkamal/fablerace/internal/modules/media"
	"github.com/ziaulkamal/fablerace/internal/modules/taxonomy"
	"github.com/ziaulkamal/fablerace/internal/modules/tenant"
	"github.com/ziaulkamal/fablerace/internal/modules/user"
	"github.com/ziaulkamal/fablerace/internal/platform/cache"
	"github.com/ziaulkamal/fablerace/internal/platform/config"
	"github.com/ziaulkamal/fablerace/internal/platform/database"
	"github.com/ziaulkamal/fablerace/internal/platform/server"
	"github.com/ziaulkamal/fablerace/internal/platform/storage"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// --- Infrastruktur ---
	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		return err
	}
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	rdb, err := cache.New(ctx, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		return err
	}
	defer rdb.Close()

	store, err := storage.New(ctx, storage.Options{
		Endpoint:  cfg.MinioEndpoint,
		AccessKey: cfg.MinioAccessKey,
		SecretKey: cfg.MinioSecretKey,
		Bucket:    cfg.MinioBucket,
		UseSSL:    cfg.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	// --- Komposisi modul: repository → service → handler ---
	handlers := server.ModuleHandlers{
		Auth:     auth.NewHandler(auth.NewService(auth.NewRepository(pool), cfg.JWTSecret)),
		Tenant:   tenant.NewHandler(tenant.NewService(tenant.NewRepository(pool))),
		User:     user.NewHandler(user.NewService(user.NewRepository(pool))),
		Content:  content.NewHandler(content.NewService(content.NewRepository(pool))),
		Taxonomy: taxonomy.NewHandler(taxonomy.NewService(taxonomy.NewRepository(pool))),
		Media:    media.NewHandler(media.NewService(media.NewRepository(pool), store)),
		Keyword:  keyword.NewHandler(keyword.NewService(keyword.NewRepository(pool))),
		License:  license.NewHandler(license.NewService(license.NewRepository(pool))),
	}
	// NEXT: saat lintas-modul dibutuhkan (mis. auth.Register memanggil
	// tenant.Service & license.Service), interface-nya ditambahkan ke konstruktor
	// service terkait dan dirangkai DI SINI — modul tetap tidak saling import
	// implementasi (docs/lampiran-b-konvensi-kode.md §3).

	srv := server.New(cfg, pool, rdb, store, handlers)

	// NEXT: proses worker jobs (flag --worker) berbagi binary ini —
	// scheduled publish, webhook deliver, sitemap, metering
	// (docs/03-arsitektur.md §5, NEXT di internal/platform/cache).

	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
