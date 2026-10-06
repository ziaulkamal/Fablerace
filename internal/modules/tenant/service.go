package tenant

import (
	"context"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	Create(ctx context.Context, in CreateRequest) (*Tenant, error)
	ListMine(ctx context.Context) ([]Tenant, error)
	GetCurrent(ctx context.Context) (*Tenant, error)
	UpdateCurrent(ctx context.Context, in UpdateSettingsRequest) (*Tenant, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NEXT (docs/10-roadmap.md minggu 1–2):
//  - Create: slug unik global (shared/slugutil), provision default content_types
//    & taxonomies, lisensi trialing 14 hari.
//    RELASI: internal/modules/license — pembuatan lisensi trial dipanggil via
//    interface service dari sini (di-inject di cmd/api/main.go).
//    RELASI: internal/modules/auth.Register — pemanggil pertama alur ini.
//  - UpdateCurrent custom_domain: validasi DNS + provisioning TLS Caddy
//    on-demand (docs/09-devops.md §1).

func (s *service) Create(ctx context.Context, in CreateRequest) (*Tenant, error) {
	return nil, apperror.NotImplemented("tenant.create")
}

func (s *service) ListMine(ctx context.Context) ([]Tenant, error) {
	return nil, apperror.NotImplemented("tenant.list_mine")
}

func (s *service) GetCurrent(ctx context.Context) (*Tenant, error) {
	return nil, apperror.NotImplemented("tenant.get_current")
}

func (s *service) UpdateCurrent(ctx context.Context, in UpdateSettingsRequest) (*Tenant, error) {
	return nil, apperror.NotImplemented("tenant.update_current")
}
