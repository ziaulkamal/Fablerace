package keyword

import (
	"context"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	Create(ctx context.Context, in CreateRequest) (*Keyword, error)
	List(ctx context.Context) ([]Keyword, error)
	AttachToContent(ctx context.Context, contentID uuid.UUID, in AttachRequest) error
	Coverage(ctx context.Context, keywordID uuid.UUID) ([]CoverageItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NEXT (docs/10-roadmap.md minggu 6) — urutan nilai-jual:
//  1. AttachToContent + Coverage (anti-kanibalisasi, docs/06 §1); verifikasi
//     kepemilikan konten via interface content.Service SEBELUM menulis relasi
//     (RELASI: internal/modules/content — pola yang sama dengan taxonomy).
//  2. Skor on-save 0–100: keyword di judul/slug/meta/heading + kepadatan body.
//  3. Automasi: auto meta description (potong cerdas 155 char), job
//     seo:rebuild_sitemap, feed RSS — dipicu event publish
//     (RELASI: internal/modules/content.Service.Publish; worker di
//     internal/platform/cache NEXT asynq).
// NEXT: AI assist add-on $10/bln — saran keyword/meta via LLM, selalu lewat
//  job queue + cache per revisi, kuota bulanan per tenant (docs/06 §5).

func (s *service) Create(ctx context.Context, in CreateRequest) (*Keyword, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return nil, apperror.NotImplemented("keyword.create")
}

func (s *service) List(ctx context.Context) ([]Keyword, error) {
	return nil, apperror.NotImplemented("keyword.list")
}

func (s *service) AttachToContent(ctx context.Context, contentID uuid.UUID, in AttachRequest) error {
	return apperror.NotImplemented("keyword.attach")
}

func (s *service) Coverage(ctx context.Context, keywordID uuid.UUID) ([]CoverageItem, error) {
	return nil, apperror.NotImplemented("keyword.coverage")
}
