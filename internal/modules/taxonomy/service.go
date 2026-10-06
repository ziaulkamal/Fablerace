package taxonomy

import (
	"context"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	CreateTaxonomy(ctx context.Context, in CreateTaxonomyRequest) (*Taxonomy, error)
	ListTaxonomies(ctx context.Context) ([]Taxonomy, error)
	CreateTerm(ctx context.Context, taxonomyID uuid.UUID, in CreateTermRequest) (*Term, error)
	ListTerms(ctx context.Context, taxonomyID uuid.UUID) ([]Term, error)
	SetContentTerms(ctx context.Context, contentID uuid.UUID, in SetContentTermsRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// NEXT (docs/10-roadmap.md minggu 3–4):
//  - CreateTerm: tolak parent_id bila taxonomy non-hierarkis; slug unik per
//    taxonomy via shared/slugutil (pola ensureSlug di modul content).
//  - SetContentTerms: verifikasi konten & semua term se-tenant SEBELUM replace
//    (RELASI: internal/modules/content — FindByID dipakai sebagai cek kepemilikan,
//    dipanggil via interface service yang di-inject dari cmd/api/main.go).

func (s *service) CreateTaxonomy(ctx context.Context, in CreateTaxonomyRequest) (*Taxonomy, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return nil, apperror.NotImplemented("taxonomy.create")
}

func (s *service) ListTaxonomies(ctx context.Context) ([]Taxonomy, error) {
	return nil, apperror.NotImplemented("taxonomy.list")
}

func (s *service) CreateTerm(ctx context.Context, taxonomyID uuid.UUID, in CreateTermRequest) (*Term, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return nil, apperror.NotImplemented("taxonomy.create_term")
}

func (s *service) ListTerms(ctx context.Context, taxonomyID uuid.UUID) ([]Term, error) {
	return nil, apperror.NotImplemented("taxonomy.list_terms")
}

func (s *service) SetContentTerms(ctx context.Context, contentID uuid.UUID, in SetContentTermsRequest) error {
	return apperror.NotImplemented("taxonomy.set_content_terms")
}
