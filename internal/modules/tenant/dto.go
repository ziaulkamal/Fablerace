package tenant

import (
	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type CreateRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"` // kosong → dari nama (RELASI: internal/shared/slugutil)
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" || len(r.Name) > 120 {
		return apperror.BadRequest("name_invalid", "name wajib diisi, maksimal 120 karakter")
	}
	return nil
}

type UpdateSettingsRequest struct {
	Name         *string `json:"name"`
	CustomDomain *string `json:"custom_domain"`
}
