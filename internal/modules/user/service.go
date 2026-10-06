package user

import (
	"context"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	ListMembers(ctx context.Context) ([]Member, error)
	Invite(ctx context.Context, in InviteRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NEXT (docs/10-roadmap.md minggu 1–2):
//  - Invite: user sudah ada → langsung upsert membership; belum ada → kirim
//    email undangan bertoken (butuh integrasi email transaksional).
//    Enforcement kuota seat per tier SEBELUM upsert (docs/02-model-bisnis.md §4;
//    RELASI: internal/modules/license — pembaca kuota plan).
//  - Jangan izinkan menurunkan role owner terakhir (lock-out guard).
// NEXT: role kustom per tenant (duplikasi role sistem — docs/05-rbac.md §1, v1.5).
// NEXT: endpoint ekspor & hapus akun (GDPR dasar — docs/08-keamanan.md §7).

func (s *service) ListMembers(ctx context.Context) ([]Member, error) {
	return nil, apperror.NotImplemented("user.list_members")
}

func (s *service) Invite(ctx context.Context, in InviteRequest) error {
	if err := in.Validate(); err != nil {
		return err
	}
	return apperror.NotImplemented("user.invite")
}
