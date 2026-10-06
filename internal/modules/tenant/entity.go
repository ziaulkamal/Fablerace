// Package tenant: siklus hidup tenant (situs/workspace) & keanggotaannya.
// 1 tenant = 1 situs = 1 lisensi (docs/02-model-bisnis.md §1).
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID           uuid.UUID
	Name         string
	Slug         string
	CustomDomain *string
	Status       string // active | suspended | cancelled
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Membership = baris tenant_users; dasar seat metering (docs/02 §4) dan sumber
// kebenaran middleware tenant (RELASI: internal/platform/middleware/tenant.go).
type Membership struct {
	TenantID uuid.UUID
	UserID   uuid.UUID
	RoleID   uuid.UUID
}
