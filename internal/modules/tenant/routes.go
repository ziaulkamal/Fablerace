package tenant

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

// Routes lingkup USER (butuh JWT, tanpa konteks tenant): list & create —
// dipasang server di grup auth-only (RELASI: internal/platform/server).
func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
	}
}

// CurrentRoutes lingkup TENANT (butuh X-Tenant-ID + keanggotaan): dipasang
// server di dalam grup tenant pada path /tenants/current.
func (h *Handler) CurrentRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", h.GetCurrent)
		r.With(mw.RequirePermission("settings.manage")).Patch("/", h.UpdateCurrent)
	}
}
