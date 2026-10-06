package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// Tenant me-resolve X-Tenant-ID, MEMVERIFIKASI keanggotaan user di tenant itu,
// lalu memuat permission set role-nya ke context. Header mentah tidak pernah
// dipercaya tanpa verifikasi ini (docs/08-keamanan.md §3).
//
// Tabel tenant_users/roles/permissions sengaja TANPA RLS karena query ini
// berjalan sebelum app.tenant_id ter-set (lihat migrations/0001_init.up.sql).
//
// RELASI: internal/platform/database.WithTenant — tenant_id yang diinject di
// sini adalah satu-satunya nilai sah untuk transaksi ber-RLS di repository.
// RELASI: internal/modules/tenant — pengelolaan keanggotaan yang dibaca di sini.
func Tenant(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserID(r.Context())
			if !ok {
				response.Error(w, apperror.Unauthorized("butuh autentikasi"))
				return
			}
			tenantID, err := uuid.Parse(r.Header.Get("X-Tenant-ID"))
			if err != nil {
				response.Error(w, apperror.BadRequest("tenant_header_invalid",
					"header X-Tenant-ID wajib berisi UUID tenant"))
				return
			}

			// NEXT: resolusi via subdomain {slug}.fablerace.app / custom domain
			// (docs/03-arsitektur.md §4) sebagai alternatif header.
			// NEXT: cache permission set di Redis TTL 60s (docs/05-rbac.md §3);
			// satu query join di bawah cukup untuk fase awal.
			rows, err := pool.Query(r.Context(), `
				SELECT p.slug
				FROM tenant_users tu
				JOIN role_permissions rp ON rp.role_id = tu.role_id
				JOIN permissions p       ON p.id = rp.permission_id
				WHERE tu.tenant_id = $1 AND tu.user_id = $2`, tenantID, userID)
			if err != nil {
				response.Error(w, apperror.Internal(err))
				return
			}
			defer rows.Close()

			perms := make(map[string]struct{})
			for rows.Next() {
				var slug string
				if err := rows.Scan(&slug); err != nil {
					response.Error(w, apperror.Internal(err))
					return
				}
				perms[slug] = struct{}{}
			}
			if len(perms) == 0 {
				// Bukan anggota (atau role tanpa izin apa pun) — jangan bocorkan
				// eksistensi tenant: 403 generik.
				response.Error(w, apperror.Forbidden("tidak memiliki akses ke tenant ini"))
				return
			}

			ctx := withTenantID(r.Context(), tenantID)
			ctx = withPermissions(ctx, perms)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
