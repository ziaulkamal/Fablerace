package license

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// Repository: tabel licenses & billing_events — tanpa RLS karena dibaca
// middleware sebelum konteks tenant penuh terbentuk (lihat migrasi).
type Repository interface {
	FindByTenant(ctx context.Context, tenantID uuid.UUID) (*License, error)
	FindByKey(ctx context.Context, licenseKey string) (*License, error)
	Update(ctx context.Context, l *License) error
	RecordBillingEvent(ctx context.Context, licenseID uuid.UUID, eventType string, payload []byte) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) FindByTenant(ctx context.Context, tenantID uuid.UUID) (*License, error) {
	return nil, apperror.NotImplemented("license.find_by_tenant")
}

func (r *pgxRepository) FindByKey(ctx context.Context, licenseKey string) (*License, error) {
	return nil, apperror.NotImplemented("license.find_by_key")
}

func (r *pgxRepository) Update(ctx context.Context, l *License) error {
	return apperror.NotImplemented("license.update")
}

func (r *pgxRepository) RecordBillingEvent(ctx context.Context, licenseID uuid.UUID, eventType string, payload []byte) error {
	return apperror.NotImplemented("license.record_billing_event")
}
