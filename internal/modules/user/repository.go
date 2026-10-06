package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Repository interface {
	ListMembers(ctx context.Context, tenantID uuid.UUID) ([]Member, error)
	UpsertMembership(ctx context.Context, tenantID, userID, roleID uuid.UUID) error
	RemoveMembership(ctx context.Context, tenantID, userID uuid.UUID) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) ListMembers(ctx context.Context, tenantID uuid.UUID) ([]Member, error) {
	return nil, apperror.NotImplemented("user.list_members")
}

func (r *pgxRepository) UpsertMembership(ctx context.Context, tenantID, userID, roleID uuid.UUID) error {
	return apperror.NotImplemented("user.upsert_membership")
}

func (r *pgxRepository) RemoveMembership(ctx context.Context, tenantID, userID uuid.UUID) error {
	return apperror.NotImplemented("user.remove_membership")
}
