package media

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
)

type Repository interface {
	Create(ctx context.Context, m *Media) error
	FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Media, error)
	List(ctx context.Context, tenantID uuid.UUID, p pagination.Params) ([]Media, int, error)
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	// SumSizeByTenant menyokong metering kuota storage (docs/02-model-bisnis.md §4).
	SumSizeByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) Create(ctx context.Context, m *Media) error {
	return apperror.NotImplemented("media.create")
}

func (r *pgxRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Media, error) {
	return nil, apperror.NotImplemented("media.find_by_id")
}

func (r *pgxRepository) List(ctx context.Context, tenantID uuid.UUID, p pagination.Params) ([]Media, int, error) {
	return nil, 0, apperror.NotImplemented("media.list")
}

func (r *pgxRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return apperror.NotImplemented("media.delete")
}

func (r *pgxRepository) SumSizeByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return 0, apperror.NotImplemented("media.sum_size")
}
