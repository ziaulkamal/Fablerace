package taxonomy

import (
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type CreateTaxonomyRequest struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	IsHierarchical bool   `json:"is_hierarchical"`
}

func (r *CreateTaxonomyRequest) Validate() error {
	if r.Name == "" || len(r.Name) > 120 {
		return apperror.BadRequest("name_invalid", "name wajib diisi, maksimal 120 karakter")
	}
	return nil
}

type CreateTermRequest struct {
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	ParentID *uuid.UUID `json:"parent_id"`
}

func (r *CreateTermRequest) Validate() error {
	if r.Name == "" || len(r.Name) > 160 {
		return apperror.BadRequest("name_invalid", "name wajib diisi, maksimal 160 karakter")
	}
	return nil
}

// SetContentTermsRequest mengganti seluruh relasi term sebuah konten (idempoten).
type SetContentTermsRequest struct {
	TermIDs []uuid.UUID `json:"term_ids"`
}
