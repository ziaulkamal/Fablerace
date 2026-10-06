package auth

import (
	"github.com/go-chi/chi/v5"
)

// Routes auth bersifat publik (pra-login) — dilindungi rate limit ketat di
// server, bukan RBAC (RELASI: internal/platform/server — grup publik /auth
// dibungkus middleware.RateLimit).
func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
	}
}
