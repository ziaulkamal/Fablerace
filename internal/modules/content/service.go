package content

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/platform/middleware"
	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
	"github.com/ziaulkamal/fablerace/internal/shared/slugutil"
)

// Service memegang seluruh aturan bisnis konten. Identitas (tenant/user/izin)
// dibaca dari context yang diisi middleware — bukan dari parameter klien.
type Service interface {
	Create(ctx context.Context, in CreateRequest) (*Content, error)
	Get(ctx context.Context, id uuid.UUID) (*Content, []MetaItem, error)
	List(ctx context.Context, f ListFilter, p pagination.Params) ([]Content, int, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateRequest) (*Content, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) (*Content, error)
	Children(ctx context.Context, id uuid.UUID) ([]Content, error)
	SetMeta(ctx context.Context, id uuid.UUID, meta MetaRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func identity(ctx context.Context) (tenantID, userID uuid.UUID, err error) {
	tenantID, ok := middleware.TenantID(ctx)
	if !ok {
		return uuid.Nil, uuid.Nil, apperror.Forbidden("tenant tidak ter-resolve")
	}
	userID, ok = middleware.UserID(ctx)
	if !ok {
		return uuid.Nil, uuid.Nil, apperror.Unauthorized("butuh autentikasi")
	}
	return tenantID, userID, nil
}

// authorizeMutation: aturan gabungan RBAC+ownership (docs/05-rbac.md §2) —
// edit_any boleh menyentuh semua; edit_own hanya konten milik sendiri.
// Guard permission dasarnya ada di routes.go; nuansa own/any di sini karena
// butuh data baris (RELASI: internal/platform/middleware/rbac.go).
func authorizeMutation(ctx context.Context, c *Content, userID uuid.UUID) error {
	if middleware.HasPermission(ctx, "content.edit_any") {
		return nil
	}
	if middleware.HasPermission(ctx, "content.edit_own") && c.AuthorID == userID {
		return nil
	}
	return apperror.Forbidden("hanya boleh mengubah konten milik sendiri")
}

func (s *service) Create(ctx context.Context, in CreateRequest) (*Content, error) {
	tenantID, userID, err := identity(ctx)
	if err != nil {
		return nil, err
	}
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if in.ParentID != nil {
		// Parent wajib ada & se-tenant; sekaligus mencegah menautkan ke konten tenant lain.
		if _, err := s.repo.FindByID(ctx, tenantID, *in.ParentID); err != nil {
			return nil, apperror.BadRequest("parent_not_found", "parent_id tidak ditemukan")
		}
		// NEXT: batasi kedalaman hierarki maks 10 level via WITH RECURSIVE
		// (docs/04-desain-data.md §2.2) sebelum fitur page-tree dalam dipakai masif.
	}

	slug, err := s.ensureSlug(ctx, tenantID, in.ContentTypeID, in.Slug, in.Title)
	if err != nil {
		return nil, err
	}

	c := &Content{
		TenantID:      tenantID,
		ContentTypeID: in.ContentTypeID,
		ParentID:      in.ParentID,
		AuthorID:      userID,
		Title:         in.Title,
		Slug:          slug,
		Excerpt:       in.Excerpt,
		Body:          in.Body,
		Status:        in.Status,
		ScheduledAt:   in.ScheduledAt,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, apperror.Internal(err)
	}
	return c, nil
	// RELASI: internal/modules/keyword — analisis skor keyword on-save dipicu
	// dari sini (event) saat modul keyword terimplementasi (docs/06 §1).
}

// ensureSlug: pakai slug kiriman bila ada, selain itu turunkan dari judul;
// keduanya diuji unik per (tenant, content_type), bentrok → sufiks -2, -3, …
func (s *service) ensureSlug(ctx context.Context, tenantID, typeID uuid.UUID, requested, title string) (string, error) {
	base := slugutil.Make(requested)
	if requested == "" {
		base = slugutil.Make(title)
	}
	candidate := base
	for i := 2; ; i++ {
		exists, err := s.repo.SlugExists(ctx, tenantID, typeID, candidate)
		if err != nil {
			return "", apperror.Internal(err)
		}
		if !exists {
			return candidate, nil
		}
		if i > 50 {
			return "", apperror.Conflict("slug_exhausted", "tidak menemukan slug unik; tentukan slug manual")
		}
		candidate = slugutil.WithSuffix(base, i)
	}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Content, []MetaItem, error) {
	tenantID, _, err := identity(ctx)
	if err != nil {
		return nil, nil, err
	}
	c, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, nil, err
	}
	meta, err := s.repo.ListMeta(ctx, tenantID, id)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}
	return c, meta, nil
}

func (s *service) List(ctx context.Context, f ListFilter, p pagination.Params) ([]Content, int, error) {
	tenantID, _, err := identity(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, total, err := s.repo.List(ctx, tenantID, f, p)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return items, total, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, in UpdateRequest) (*Content, error) {
	tenantID, userID, err := identity(ctx)
	if err != nil {
		return nil, err
	}
	if err := in.Validate(); err != nil {
		return nil, err
	}
	c, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := authorizeMutation(ctx, c, userID); err != nil {
		return nil, err
	}

	if in.Title != nil {
		c.Title = *in.Title
	}
	if in.Slug != nil && *in.Slug != c.Slug {
		// NEXT: catat redirect 301 di slug_redirects saat slug konten published
		// berubah (docs/06-keyword-seo-automasi.md §2).
		slug, err := s.ensureSlug(ctx, tenantID, c.ContentTypeID, *in.Slug, c.Title)
		if err != nil {
			return nil, err
		}
		c.Slug = slug
	}
	if in.Excerpt != nil {
		c.Excerpt = *in.Excerpt
	}
	if in.Body != nil {
		c.Body = *in.Body
	}
	if in.ParentID != nil {
		if _, err := s.repo.FindByID(ctx, tenantID, *in.ParentID); err != nil {
			return nil, apperror.BadRequest("parent_not_found", "parent_id tidak ditemukan")
		}
		c.ParentID = in.ParentID
	}
	if in.Status != nil {
		c.Status = *in.Status
		if c.Status == StatusPublished && c.PublishedAt == nil {
			now := time.Now()
			c.PublishedAt = &now
		}
	}
	if in.ScheduledAt != nil {
		c.ScheduledAt = in.ScheduledAt
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, _, err := identity(ctx)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, tenantID, id)
}

func (s *service) Publish(ctx context.Context, id uuid.UUID) (*Content, error) {
	tenantID, userID, err := identity(ctx)
	if err != nil {
		return nil, err
	}
	c, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	// Author boleh publish konten sendiri; selain itu butuh edit_any
	// (pola WordPress publish_posts vs edit_others_posts, docs/05-rbac.md §2).
	if err := authorizeMutation(ctx, c, userID); err != nil {
		return nil, err
	}

	now := time.Now()
	c.Status = StatusPublished
	c.PublishedAt = &now
	c.ScheduledAt = nil
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
	// RELASI: internal/modules/keyword — event publish memicu job
	// seo:rebuild_sitemap & webhook content.published (docs/03-arsitektur.md §5)
	// saat worker asynq diimplementasi (NEXT: internal/platform/cache).
}

func (s *service) Children(ctx context.Context, id uuid.UUID) ([]Content, error) {
	tenantID, _, err := identity(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.FindByID(ctx, tenantID, id); err != nil {
		return nil, err
	}
	items, err := s.repo.Children(ctx, tenantID, id)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return items, nil
}

func (s *service) SetMeta(ctx context.Context, id uuid.UUID, meta MetaRequest) error {
	tenantID, userID, err := identity(ctx)
	if err != nil {
		return err
	}
	c, err := s.repo.FindByID(ctx, tenantID, id) // verifikasi kepemilikan tenant SEBELUM sentuh content_meta
	if err != nil {
		return err
	}
	if err := authorizeMutation(ctx, c, userID); err != nil {
		return err
	}
	items := make([]MetaItem, 0, len(meta))
	for k, v := range meta {
		if k == "" || len(k) > 120 {
			return apperror.BadRequest("meta_key_invalid", "meta key wajib diisi, maksimal 120 karakter")
		}
		items = append(items, MetaItem{Key: k, Value: v})
	}
	if err := s.repo.ReplaceMeta(ctx, tenantID, id, items); err != nil {
		return apperror.Internal(err)
	}
	return nil
}
