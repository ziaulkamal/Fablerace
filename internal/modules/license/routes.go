package license

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("billing.view")).Get("/", h.Usage)
		r.With(mw.RequirePermission("billing.manage")).Post("/activate", h.Activate)
	}
}
