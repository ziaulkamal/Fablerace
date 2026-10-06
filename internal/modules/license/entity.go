// Package license: lisensi berlangganan per tenant ($25/bln — docs/02-model-bisnis.md).
// Siklus status ditegakkan middleware License di tiap request.
// Status modul: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package license

import (
	"time"

	"github.com/google/uuid"
)

// Status mengikuti state machine docs/02-model-bisnis.md §1.
type Status string

const (
	StatusTrialing  Status = "trialing"
	StatusActive    Status = "active"
	StatusPastDue   Status = "past_due"
	StatusGrace     Status = "grace"
	StatusSuspended Status = "suspended"
	StatusCancelled Status = "cancelled"
)

type Plan string

const (
	PlanStarter Plan = "starter" // $9
	PlanPro     Plan = "pro"     // $25 — anchor
	PlanAgency  Plan = "agency"  // $20/lisensi, min 5
)

type License struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	LicenseKey         string
	Plan               Plan
	Seats              int
	Status             Status
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	Provider           string // paddle | midtrans | manual (beta)
	ProviderRef        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// RELASI: internal/platform/middleware/license.go — konsumen baca status ini
// di setiap request (trialing/active/past_due penuh; grace/suspended read-only).
