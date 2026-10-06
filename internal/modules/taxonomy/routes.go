package taxonomy

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("content.view")).Get("/", h.List)
		r.With(mw.RequirePermission("taxonomy.manage")).Post("/", h.Create)
		r.With(mw.RequirePermission("content.view")).Get("/{id}/terms", h.ListTerms)
		r.With(mw.RequirePermission("taxonomy.manage")).Post("/{id}/terms", h.CreateTerm)
	}
}
