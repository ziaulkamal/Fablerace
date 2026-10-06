package keyword

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("content.view")).Get("/", h.List)
		r.With(mw.RequirePermission("keyword.manage")).Post("/", h.Create)
		r.With(mw.RequirePermission("content.view")).Get("/{id}/coverage", h.Coverage)
	}
}
