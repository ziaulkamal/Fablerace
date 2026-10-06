package user

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("user.manage")).Get("/", h.ListMembers)
		r.With(mw.RequirePermission("user.invite")).Post("/invites", h.Invite)
	}
}
