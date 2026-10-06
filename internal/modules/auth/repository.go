package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// Repository: tabel users (kredensial) & refresh_tokens — keduanya tanpa RLS
// karena diakses sebelum konteks tenant ada (lihat migrations/0001_init.up.sql).
type Repository interface {
	FindCredentialsByEmail(ctx context.Context, email string) (userID uuid.UUID, passwordHash string, err error)
	StoreRefreshToken(ctx context.Context, t *RefreshToken) error
	RevokeByHash(ctx context.Context, tokenHash string) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) FindCredentialsByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	return uuid.Nil, "", apperror.NotImplemented("auth.find_credentials")
}

func (r *pgxRepository) StoreRefreshToken(ctx context.Context, t *RefreshToken) error {
	return apperror.NotImplemented("auth.store_refresh_token")
}

func (r *pgxRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	return apperror.NotImplemented("auth.revoke_refresh_token")
}
