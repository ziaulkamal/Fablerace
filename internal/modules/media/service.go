package media

import (
	"context"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/platform/storage"
	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
)

type Service interface {
	Presign(ctx context.Context, in PresignRequest) (*PresignResponse, error)
	Confirm(ctx context.Context, mediaID uuid.UUID) (*Media, error)
	List(ctx context.Context, p pagination.Params) ([]Media, int, error)
	Delete(ctx context.Context, mediaID uuid.UUID) error
}

type service struct {
	repo  Repository
	store *storage.Client
}

func NewService(repo Repository, store *storage.Client) Service {
	return &service{repo: repo, store: store}
}

// NEXT (docs/10-roadmap.md minggu 5) — alur presign→confirm (docs/07 §2):
//  - Presign: cek kuota storage tier via repo.SumSizeByTenant
//    (RELASI: internal/modules/license — pembaca batas plan), lalu
//    store.ObjectKey(tenantID, "media", yyyy, mm, mediaID, fileName) +
//    store.PresignPut 15 menit (RELASI: internal/platform/storage — pabrik
//    object key satu-satunya; tenant_id dari context, bukan payload).
//  - Confirm: store.Stat untuk verifikasi keberadaan & ukuran + deteksi
//    magic-byte, baru insert baris media; antrikan job media:process
//    (varian WebP 640/1280 — docs/07 §4).
//  - Delete: hapus baris + objek + varian; media yang masih direferensikan
//    konten ditolak (NEXT: tabel referensi media↔content).

func (s *service) Presign(ctx context.Context, in PresignRequest) (*PresignResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return nil, apperror.NotImplemented("media.presign")
}

func (s *service) Confirm(ctx context.Context, mediaID uuid.UUID) (*Media, error) {
	return nil, apperror.NotImplemented("media.confirm")
}

func (s *service) List(ctx context.Context, p pagination.Params) ([]Media, int, error) {
	return nil, 0, apperror.NotImplemented("media.list")
}

func (s *service) Delete(ctx context.Context, mediaID uuid.UUID) error {
	return apperror.NotImplemented("media.delete")
}
