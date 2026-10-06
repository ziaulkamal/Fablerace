// Package database memegang pool PostgreSQL, runner migrasi embed, dan helper
// transaksi ber-RLS — satu-satunya tempat `app.tenant_id` di-set
// (docs/04-desain-data.md §4, docs/08-keamanan.md §3).
package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: buka pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}
	return pool, nil
}

// WithTenant menjalankan fn dalam satu transaksi dengan `app.tenant_id` ter-set
// (padanan SET LOCAL), sehingga policy RLS membatasi semua statement di dalamnya
// ke baris milik tenant tersebut. Repository modul WAJIB memakai helper ini
// untuk tabel ber-tenant; query di luar helper ini tidak lolos RLS.
//
// RELASI: internal/platform/middleware/tenant.go — sumber tenant_id yang sah
// (sudah diverifikasi keanggotaan) sebelum sampai ke repository.
func WithTenant(ctx context.Context, pool *pgxpool.Pool, tenantID uuid.UUID, fn func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("database: begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op setelah Commit sukses

	if _, err := tx.Exec(ctx,
		"SELECT set_config('app.tenant_id', $1, true)", tenantID.String()); err != nil {
		return fmt.Errorf("database: set tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("database: commit: %w", err)
	}
	return nil
}
