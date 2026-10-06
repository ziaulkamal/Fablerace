package middleware

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// License menegakkan siklus langganan (docs/02-model-bisnis.md §1) di tiap
// request: trialing/active/past_due → penuh; grace/suspended → read-only;
// cancelled/tanpa lisensi → ditolak. Dipasang SETELAH middleware Tenant.
//
// RELASI: internal/modules/license — pemilik tabel licenses & siklus statusnya;
// middleware ini hanya konsumen baca.
// NEXT: cache status lisensi di Redis TTL 60s (docs/08-keamanan.md §5) agar
// tidak menambah satu query per request pada beban tinggi.
func License(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := TenantID(r.Context())
			if !ok {
				response.Error(w, apperror.Forbidden("tenant belum ter-resolve"))
				return
			}

			var status string
			err := pool.QueryRow(r.Context(),
				`SELECT status FROM licenses WHERE tenant_id = $1`, tenantID).Scan(&status)
			if errors.Is(err, pgx.ErrNoRows) {
				response.Error(w, apperror.New(http.StatusPaymentRequired,
					"license_missing", "tenant belum memiliki lisensi aktif"))
				return
			}
			if err != nil {
				response.Error(w, apperror.Internal(err))
				return
			}

			switch status {
			case "trialing", "active", "past_due":
				next.ServeHTTP(w, r)
			case "grace", "suspended":
				if r.Method == http.MethodGet || r.Method == http.MethodHead {
					next.ServeHTTP(w, r)
					return
				}
				response.Error(w, apperror.New(http.StatusPaymentRequired,
					"license_readonly", "lisensi menunggak — akses sementara read-only"))
			default: // cancelled
				response.Error(w, apperror.New(http.StatusPaymentRequired,
					"license_inactive", "lisensi tidak aktif"))
			}
		})
	}
}
