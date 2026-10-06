package content

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ziaulkamal/fablerace/internal/platform/database"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
)

// Integration test modul exemplar terhadap Postgres nyata, termasuk uji
// ISOLASI RLS lintas-tenant — bagian definition-of-done tiap modul
// (docs/10-roadmap.md §4 "test isolasi").
//
// Jalankan:
//   TEST_DATABASE_URL=postgres://fablerace:fablerace@localhost:5432/fablerace?sslmode=disable \
//     go test ./internal/modules/content/ -v
// Tanpa TEST_DATABASE_URL test di-skip (CI tanpa DB tetap hijau).
func TestRepositoryIntegration(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL tidak di-set; lewati integration test")
	}
	ctx := context.Background()

	if err := database.Migrate(dbURL); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	pool, err := database.NewPool(ctx, dbURL)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	defer pool.Close()

	// --- Fixture: 2 tenant + 1 user + content_type per tenant ---------------
	suffix := uuid.NewString()[:8]
	var tenantA, tenantB, userID uuid.UUID
	if err := pool.QueryRow(ctx,
		`INSERT INTO tenants (name, slug) VALUES ('Tenant A','ta-`+suffix+`') RETURNING id`).Scan(&tenantA); err != nil {
		t.Fatalf("seed tenant A: %v", err)
	}
	if err := pool.QueryRow(ctx,
		`INSERT INTO tenants (name, slug) VALUES ('Tenant B','tb-`+suffix+`') RETURNING id`).Scan(&tenantB); err != nil {
		t.Fatalf("seed tenant B: %v", err)
	}
	if err := pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name)
		 VALUES ('test-`+suffix+`@example.com','x','Tester') RETURNING id`).Scan(&userID); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// content_types ber-RLS FORCE → insert harus lewat WithTenant: sekaligus
	// menguji jalur tulis RLS.
	var typeA uuid.UUID
	if err := database.WithTenant(ctx, pool, tenantA, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO content_types (tenant_id, slug, name, is_hierarchical)
			 VALUES ($1,'page','Page',true) RETURNING id`, tenantA).Scan(&typeA)
	}); err != nil {
		t.Fatalf("seed content_type: %v", err)
	}

	repo := NewRepository(pool)

	// --- Create + hierarki ---------------------------------------------------
	parent := &Content{
		TenantID: tenantA, ContentTypeID: typeA, AuthorID: userID,
		Title: "Halaman Induk", Slug: "halaman-induk", Status: StatusDraft,
	}
	if err := repo.Create(ctx, parent); err != nil {
		t.Fatalf("create parent: %v", err)
	}
	child := &Content{
		TenantID: tenantA, ContentTypeID: typeA, ParentID: &parent.ID, AuthorID: userID,
		Title: "Halaman Anak", Slug: "halaman-anak", Status: StatusPublished,
	}
	if err := repo.Create(ctx, child); err != nil {
		t.Fatalf("create child: %v", err)
	}

	children, err := repo.Children(ctx, tenantA, parent.ID)
	if err != nil || len(children) != 1 || children[0].ID != child.ID {
		t.Fatalf("children: err=%v len=%d", err, len(children))
	}

	// --- SlugExists + List filter -------------------------------------------
	if ok, _ := repo.SlugExists(ctx, tenantA, typeA, "halaman-induk"); !ok {
		t.Fatal("SlugExists: slug seharusnya terdeteksi")
	}
	st := StatusPublished
	items, total, err := repo.List(ctx, tenantA, ListFilter{Status: &st}, pagination.Params{Page: 1, PerPage: 10})
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("list published: err=%v total=%d", err, total)
	}

	// --- Meta EAV -------------------------------------------------------------
	if err := repo.ReplaceMeta(ctx, tenantA, parent.ID, []MetaItem{
		{Key: "hero_color", Value: json.RawMessage(`"#112233"`)},
	}); err != nil {
		t.Fatalf("replace meta: %v", err)
	}
	meta, err := repo.ListMeta(ctx, tenantA, parent.ID)
	if err != nil || len(meta) != 1 || meta[0].Key != "hero_color" {
		t.Fatalf("list meta: err=%v len=%d", err, len(meta))
	}

	// --- UJI ISOLASI RLS: tenant B tidak boleh melihat data tenant A ----------
	if _, err := repo.FindByID(ctx, tenantB, parent.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("RLS BOCOR: tenant B membaca konten tenant A (err=%v)", err)
	}
	if _, totalB, _ := repo.List(ctx, tenantB, ListFilter{}, pagination.Params{Page: 1, PerPage: 10}); totalB != 0 {
		t.Fatalf("RLS BOCOR: list tenant B melihat %d baris tenant A", totalB)
	}

	// --- Soft delete -----------------------------------------------------------
	if err := repo.SoftDelete(ctx, tenantA, child.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, tenantA, child.ID); !errors.Is(err, ErrNotFound) {
		t.Fatal("soft delete: konten terhapus masih terbaca")
	}
}
