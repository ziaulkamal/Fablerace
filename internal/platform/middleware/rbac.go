package middleware

import (
	"net/http"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// RequirePermission adalah guard RBAC per-route: dipasang di routes.go tiap
// modul (matriks izin di docs/05-rbac.md §2; seed di migrations/0001_init.up.sql).
// Aturan gabungan permission+ownership (edit_own vs edit_any) BUKAN di sini —
// itu urusan service modul karena butuh data baris
// (RELASI: internal/modules/content/service.go — authorizeMutation).
//
// NEXT: evaluasi ulang Casbin bila muncul kebutuhan policy per-objek yang rumit,
// mis. izin per-folder media (keputusan di docs/05-rbac.md §3).
func RequirePermission(slug string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasPermission(r.Context(), slug) {
				response.Error(w, apperror.Forbidden("izin "+slug+" dibutuhkan"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
