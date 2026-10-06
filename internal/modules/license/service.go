package license

import (
	"context"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	Usage(ctx context.Context) (*UsageResponse, error)
	Activate(ctx context.Context, in ActivateRequest) (*License, error)
}

// PaymentProvider mengabstraksi gateway (Paddle global, Midtrans ID —
// docs/02-model-bisnis.md §2) agar penyedia bisa ditukar tanpa menyentuh service.
// NEXT: implementasi pertama = provider "manual" untuk beta (aktivasi oleh
// platform admin), lalu Paddle (v1.0, docs/10-roadmap.md).
type PaymentProvider interface {
	VerifyWebhook(payload []byte, signature string) (eventType string, providerRef string, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NEXT (docs/10-roadmap.md minggu 7 & v1.0):
//  - Usage: gabungkan lisensi + metering (seats dari tenant_users, storage dari
//    media.SumSizeByTenant — RELASI: internal/modules/media & internal/modules/user,
//    via interface service yang di-inject dari cmd/api/main.go).
//  - Activate: validasi key (format fbr_live_*, docs/08-keamanan.md §4),
//    pindahkan status trialing → active, catat billing_event.
//  - Webhook gateway: dunning past_due → grace → suspended sesuai state machine
//    docs/02 §1; job cron billing:metering agregasi harian.
//  - Pelanggaran kuota → status over_limit + tenggang 7 hari, BUKAN blokir
//    seketika (keputusan produk docs/02 §4).

func (s *service) Usage(ctx context.Context) (*UsageResponse, error) {
	return nil, apperror.NotImplemented("license.usage")
}

func (s *service) Activate(ctx context.Context, in ActivateRequest) (*License, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return nil, apperror.NotImplemented("license.activate")
}
