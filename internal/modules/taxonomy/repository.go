package taxonomy

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// Repository: tabel taxonomies, terms (ber-RLS), term_relationships (isolasi
// via verifikasi kepemilikan konten & term oleh service — pola yang sama dengan
// content_meta di modul content).
type Repository interface {
	CreateTaxonomy(ctx context.Context, t *Taxonomy) error
	ListTaxonomies(ctx context.Context, tenantID uuid.UUID) ([]Taxonomy, error)
	CreateTerm(ctx context.Context, t *Term) error
	ListTerms(ctx context.Context, tenantID, taxonomyID uuid.UUID) ([]Term, error)
	ReplaceContentTerms(ctx context.Context, tenantID, contentID uuid.UUID, termIDs []uuid.UUID) error
}

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) CreateTaxonomy(ctx context.Context, t *Taxonomy) error {
	return apperror.NotImplemented("taxonomy.create")
}

func (r *pgxRepository) ListTaxonomies(ctx context.Context, tenantID uuid.UUID) ([]Taxonomy, error) {
	return nil, apperror.NotImplemented("taxonomy.list")
}

func (r *pgxRepository) CreateTerm(ctx context.Context, t *Term) error {
	return apperror.NotImplemented("taxonomy.create_term")
}

func (r *pgxRepository) ListTerms(ctx context.Context, tenantID, taxonomyID uuid.UUID) ([]Term, error) {
	return nil, apperror.NotImplemented("taxonomy.list_terms")
	// NEXT: kembalikan pohon term (WITH RECURSIVE, pola yang sama dengan
	// contents.parent_id — RELASI: internal/modules/content/repository.go Children).
}

func (r *pgxRepository) ReplaceContentTerms(ctx context.Context, tenantID, contentID uuid.UUID, termIDs []uuid.UUID) error {
	return apperror.NotImplemented("taxonomy.replace_content_terms")
}
