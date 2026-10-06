// Package taxonomy: taxonomy & term berjenjang ala WordPress (category/tag/
// kustom), plus relasi konten↔term (docs/04-desain-data.md §2.3).
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package taxonomy

import (
	"time"

	"github.com/google/uuid"
)

type Taxonomy struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Name           string
	Slug           string
	IsHierarchical bool // category = true, tag = false
	CreatedAt      time.Time
}

// Term berjenjang lewat ParentID — hanya valid bila taxonomy induknya
// IsHierarchical (ditegakkan service, bukan skema).
type Term struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	TaxonomyID uuid.UUID
	ParentID   *uuid.UUID
	Name       string
	Slug       string
	CreatedAt  time.Time
}

// RELASI: internal/modules/content — term_relationships menautkan term ke
// contents; ON DELETE CASCADE membersihkan relasi dari kedua arah (migrasi).
