package media

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("content.view")).Get("/", h.List)
		r.With(mw.RequirePermission("media.upload")).Post("/presign", h.Presign)
		r.With(mw.RequirePermission("media.upload")).Post("/{id}/confirm", h.Confirm)
		r.With(mw.RequirePermission("media.manage")).Delete("/{id}", h.Delete)
	}
}
