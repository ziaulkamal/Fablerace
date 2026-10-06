package content

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/platform/database"
	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
)

// Repository: satu-satunya tempat SQL modul content. Semua operasi berjalan di
// dalam database.WithTenant sehingga RLS aktif; filter tenant_id eksplisit di
// query adalah lapis pertahanan pertama, RLS lapis kedua (docs/08-keamanan.md §3).
type Repository interface {
	Create(ctx context.Context, c *Content) error
	FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Content, error)
	List(ctx context.Context, tenantID uuid.UUID, f ListFilter, p pagination.Params) ([]Content, int, error)
	Update(ctx context.Context, c *Content) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	Children(ctx context.Context, tenantID, parentID uuid.UUID) ([]Content, error)
	SlugExists(ctx context.Context, tenantID, contentTypeID uuid.UUID, slug string) (bool, error)
	ReplaceMeta(ctx context.Context, tenantID, contentID uuid.UUID, items []MetaItem) error
	ListMeta(ctx context.Context, tenantID, contentID uuid.UUID) ([]MetaItem, error)
}

var ErrNotFound = apperror.NotFound("content_not_found", "konten tidak ditemukan")

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

const contentColumns = `id, tenant_id, content_type_id, parent_id, author_id,
	title, slug, excerpt, body, status, published_at, scheduled_at, created_at, updated_at`

func scanContent(row pgx.Row) (*Content, error) {
	var c Content
	err := row.Scan(&c.ID, &c.TenantID, &c.ContentTypeID, &c.ParentID, &c.AuthorID,
		&c.Title, &c.Slug, &c.Excerpt, &c.Body, &c.Status,
		&c.PublishedAt, &c.ScheduledAt, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("content: scan: %w", err)
	}
	return &c, nil
}

func (r *pgxRepository) Create(ctx context.Context, c *Content) error {
	return database.WithTenant(ctx, r.pool, c.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			INSERT INTO contents (tenant_id, content_type_id, parent_id, author_id,
				title, slug, excerpt, body, status, scheduled_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			RETURNING id, created_at, updated_at`,
			c.TenantID, c.ContentTypeID, c.ParentID, c.AuthorID,
			c.Title, c.Slug, c.Excerpt, c.Body, c.Status, c.ScheduledAt,
		).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	})
}

func (r *pgxRepository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Content, error) {
	var out *Content
	err := database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		var ferr error
		out, ferr = scanContent(tx.QueryRow(ctx, `
			SELECT `+contentColumns+` FROM contents
			WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, tenantID, id))
		return ferr
	})
	return out, err
}

func (r *pgxRepository) List(ctx context.Context, tenantID uuid.UUID, f ListFilter, p pagination.Params) ([]Content, int, error) {
	where := []string{"tenant_id = $1", "deleted_at IS NULL"}
	args := []any{tenantID}

	add := func(cond string, val any) {
		args = append(args, val)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if f.ContentTypeID != nil {
		add("content_type_id = $%d", *f.ContentTypeID)
	}
	if f.Status != nil {
		add("status = $%d", *f.Status)
	}
	if f.ParentID != nil {
		add("parent_id = $%d", *f.ParentID)
	}
	if f.Search != "" {
		// FTS lewat kolom generated search_vector (docs/06-keyword-seo-automasi.md §4).
		// NEXT: Meilisearch bila typo-tolerance/ranking Postgres terbukti kurang.
		add("search_vector @@ websearch_to_tsquery('simple', $%d)", f.Search)
	}
	cond := strings.Join(where, " AND ")

	var (
		items []Content
		total int
	)
	err := database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			"SELECT count(*) FROM contents WHERE "+cond, args...).Scan(&total); err != nil {
			return fmt.Errorf("content: count: %w", err)
		}
		rows, err := tx.Query(ctx, fmt.Sprintf(
			"SELECT %s FROM contents WHERE %s ORDER BY created_at DESC LIMIT %d OFFSET %d",
			contentColumns, cond, p.Limit(), p.Offset()), args...)
		if err != nil {
			return fmt.Errorf("content: list: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			c, err := scanContent(rows)
			if err != nil {
				return err
			}
			items = append(items, *c)
		}
		return rows.Err()
	})
	return items, total, err
}

func (r *pgxRepository) Update(ctx context.Context, c *Content) error {
	return database.WithTenant(ctx, r.pool, c.TenantID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE contents SET title=$3, slug=$4, excerpt=$5, body=$6, parent_id=$7,
				status=$8, published_at=$9, scheduled_at=$10, updated_at=now()
			WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL`,
			c.TenantID, c.ID, c.Title, c.Slug, c.Excerpt, c.Body, c.ParentID,
			c.Status, c.PublishedAt, c.ScheduledAt)
		if err != nil {
			return fmt.Errorf("content: update %s: %w", c.ID, err)
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	})
	// NEXT: tulis snapshot ke tabel revisions pada setiap update
	// (docs/06-keyword-seo-automasi.md §3 "Auto-revisi") — dalam transaksi yang sama.
}

func (r *pgxRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	return database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE contents SET deleted_at = now()
			WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL`, tenantID, id)
		if err != nil {
			return fmt.Errorf("content: soft delete %s: %w", id, err)
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func (r *pgxRepository) Children(ctx context.Context, tenantID, parentID uuid.UUID) ([]Content, error) {
	var items []Content
	err := database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT `+contentColumns+` FROM contents
			WHERE tenant_id=$1 AND parent_id=$2 AND deleted_at IS NULL
			ORDER BY created_at ASC`, tenantID, parentID)
		if err != nil {
			return fmt.Errorf("content: children: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			c, err := scanContent(rows)
			if err != nil {
				return err
			}
			items = append(items, *c)
		}
		return rows.Err()
	})
	return items, err
	// NEXT: query seluruh subtree via WITH RECURSIVE (kedalaman maks 10) untuk
	// breadcrumb/sitemap — docs/04-desain-data.md §2.2.
}

func (r *pgxRepository) SlugExists(ctx context.Context, tenantID, contentTypeID uuid.UUID, slug string) (bool, error) {
	var exists bool
	err := database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM contents
				WHERE tenant_id=$1 AND content_type_id=$2 AND slug=$3 AND deleted_at IS NULL)`,
			tenantID, contentTypeID, slug).Scan(&exists)
	})
	return exists, err
}

// ReplaceMeta meng-upsert seluruh pasangan key→value (EAV). content_meta tidak
// ber-tenant_id; isolasi datang dari verifikasi kepemilikan konten oleh service
// SEBELUM memanggil ini + FK cascade (docs/04-desain-data.md §2.2).
func (r *pgxRepository) ReplaceMeta(ctx context.Context, tenantID, contentID uuid.UUID, items []MetaItem) error {
	return database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		for _, m := range items {
			if _, err := tx.Exec(ctx, `
				INSERT INTO content_meta (content_id, meta_key, meta_value)
				VALUES ($1,$2,$3)
				ON CONFLICT (content_id, meta_key) DO UPDATE SET meta_value = EXCLUDED.meta_value`,
				contentID, m.Key, m.Value); err != nil {
				return fmt.Errorf("content: upsert meta %q: %w", m.Key, err)
			}
		}
		return nil
	})
}

func (r *pgxRepository) ListMeta(ctx context.Context, tenantID, contentID uuid.UUID) ([]MetaItem, error) {
	var items []MetaItem
	err := database.WithTenant(ctx, r.pool, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT meta_key, meta_value FROM content_meta
			WHERE content_id = $1 ORDER BY meta_key`, contentID)
		if err != nil {
			return fmt.Errorf("content: list meta: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var m MetaItem
			if err := rows.Scan(&m.Key, &m.Value); err != nil {
				return fmt.Errorf("content: scan meta: %w", err)
			}
			items = append(items, m)
		}
		return rows.Err()
	})
	return items, err
}
