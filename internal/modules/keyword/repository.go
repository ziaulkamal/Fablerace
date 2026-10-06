package keyword

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Repository interface {
	Create(ctx context.Context, k *Keyword) error
	List(ctx context.Context, tenantID uuid.UUID) ([]Keyword, error)
	ReplaceContentKeywords(ctx context.Context, tenantID, contentID uuid.UUID, items []ContentKeyword) error
	Coverage(ctx context.Context, tenantID, keywordID uuid.UUID) ([]CoverageItem, error)
	UpsertSEOMeta(ctx context.Context, tenantID uuid.UUID, m *SEOMeta) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) Create(ctx context.Context, k *Keyword) error {
	return apperror.NotImplemented("keyword.create")
}

func (r *pgxRepository) List(ctx context.Context, tenantID uuid.UUID) ([]Keyword, error) {
	return nil, apperror.NotImplemented("keyword.list")
}

func (r *pgxRepository) ReplaceContentKeywords(ctx context.Context, tenantID, contentID uuid.UUID, items []ContentKeyword) error {
	return apperror.NotImplemented("keyword.replace_content_keywords")
}

func (r *pgxRepository) Coverage(ctx context.Context, tenantID, keywordID uuid.UUID) ([]CoverageItem, error) {
	return nil, apperror.NotImplemented("keyword.coverage")
}

func (r *pgxRepository) UpsertSEOMeta(ctx context.Context, tenantID uuid.UUID, m *SEOMeta) error {
	return apperror.NotImplemented("keyword.upsert_seo_meta")
}
