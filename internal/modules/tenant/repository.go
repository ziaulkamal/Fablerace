package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// Repository: tabel tenants & tenant_users — tanpa RLS karena menjadi pintu
// resolusi tenant itu sendiri (lihat migrations/0001_init.up.sql); isolasi
// ditegakkan filter eksplisit user_id/tenant_id di tiap query.
type Repository interface {
	Create(ctx context.Context, t *Tenant, ownerUserID uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Tenant, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) Create(ctx context.Context, t *Tenant, ownerUserID uuid.UUID) error {
	// NEXT: satu transaksi — insert tenants, tenant_users (role owner dari seed),
	// content_types default (post & page), taxonomies default (category & tag).
	return apperror.NotImplemented("tenant.create")
}

func (r *pgxRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]Tenant, error) {
	return nil, apperror.NotImplemented("tenant.list_by_user")
}

func (r *pgxRepository) FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return nil, apperror.NotImplemented("tenant.find_by_id")
}
