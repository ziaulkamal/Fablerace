package content

import (
	"github.com/go-chi/chi/v5"

	mw "github.com/ziaulkamal/fablerace/internal/platform/middleware"
)

// Routes memetakan endpoint modul content (docs/lampiran-a-api.md) dengan guard
// RBAC per route. Guard di sini memakai izin TERENDAH yang masuk akal
// (content.edit_own); nuansa own-vs-any ditegakkan service.authorizeMutation
// (RELASI: internal/platform/middleware/rbac.go — pembagian tanggung jawab).
func (h *Handler) Routes() func(chi.Router) {
	return func(r chi.Router) {
		r.With(mw.RequirePermission("content.view")).Get("/", h.List)
		r.With(mw.RequirePermission("content.create")).Post("/", h.Create)
		r.With(mw.RequirePermission("content.view")).Get("/{id}", h.Get)
		r.With(mw.RequirePermission("content.view")).Get("/{id}/children", h.Children)
		r.With(mw.RequirePermission("content.edit_own")).Patch("/{id}", h.Update)
		r.With(mw.RequirePermission("content.edit_own")).Put("/{id}/meta", h.SetMeta)
		r.With(mw.RequirePermission("content.publish")).Post("/{id}/publish", h.Publish)
		r.With(mw.RequirePermission("content.delete")).Delete("/{id}", h.Delete)
	}
}
