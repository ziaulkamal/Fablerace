package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// Auth memverifikasi Bearer JWT (HS256) dan meletakkan user_id di context.
// Klaim sengaja minimal (sub/exp/iat) — role & permission TIDAK ada di token,
// selalu dibaca segar oleh middleware tenant agar pencabutan akses berlaku
// seketika (docs/08-keamanan.md §2).
//
// RELASI: internal/modules/auth — service login/refresh adalah penerbit token
// yang diverifikasi di sini; secret & TTL dari platform/config.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if raw == "" || raw == r.Header.Get("Authorization") {
				response.Error(w, apperror.Unauthorized("token tidak ditemukan"))
				return
			}

			tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !tok.Valid {
				response.Error(w, apperror.Unauthorized("token tidak valid atau kedaluwarsa"))
				return
			}
			sub, err := tok.Claims.GetSubject()
			if err != nil {
				response.Error(w, apperror.Unauthorized("token tanpa subjek"))
				return
			}
			userID, err := uuid.Parse(sub)
			if err != nil {
				response.Error(w, apperror.Unauthorized("subjek token tidak valid"))
				return
			}

			next.ServeHTTP(w, r.WithContext(withUserID(r.Context(), userID)))
		})
	}
}

// NEXT: jalur autentikasi kedua via X-API-Key (scope subset, untuk delivery
// publik — docs/05-rbac.md §4, docs/08-keamanan.md §4) sebagai middleware
// APIKey(pool) yang mengisi tenant_id+permissions langsung dari tabel api_keys.
