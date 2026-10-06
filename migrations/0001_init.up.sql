-- 0001_init: skema penuh Fablerace (docs/04-desain-data.md).
-- Konvensi RLS (docs/08-keamanan.md §3):
--   * Tabel ber-tenant_id yang diakses SETELAH konteks tenant ada → RLS + FORCE.
--   * Tabel resolusi/identitas (tenants, users, roles, permissions,
--     role_permissions, tenant_users, refresh_tokens, licenses, billing_events)
--     SENGAJA tanpa RLS: dibaca middleware sebelum app.tenant_id ter-set;
--     isolasi via filter eksplisit di query.
--   * Tabel penghubung tanpa tenant_id (content_meta, content_seo, revisions,
--     term_relationships, content_keywords, webhook_deliveries) tanpa RLS:
--     isolasi via verifikasi kepemilikan baris induk oleh service + FK CASCADE.

-- ===== Identitas & akses =====================================================

CREATE TABLE tenants (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          text NOT NULL,
    slug          text NOT NULL UNIQUE,
    custom_domain text UNIQUE,
    status        text NOT NULL DEFAULT 'active'
                  CHECK (status IN ('active','suspended','cancelled')),
    settings      jsonb NOT NULL DEFAULT '{}',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email             text NOT NULL,
    password_hash     text NOT NULL,
    name              text NOT NULL DEFAULT '',
    is_platform_admin boolean NOT NULL DEFAULT false,
    status            text NOT NULL DEFAULT 'active'
                      CHECK (status IN ('active','disabled')),
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_unique ON users (lower(email));

-- tenant_id NULL = role sistem (seed di bawah); non-NULL = role kustom tenant (v1.5).
CREATE TABLE roles (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid REFERENCES tenants(id) ON DELETE CASCADE,
    slug        text NOT NULL,
    name        text NOT NULL,
    description text NOT NULL DEFAULT '',
    UNIQUE (tenant_id, slug)
);
CREATE UNIQUE INDEX roles_system_slug_unique ON roles (slug) WHERE tenant_id IS NULL;

CREATE TABLE permissions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        text NOT NULL UNIQUE,
    description text NOT NULL DEFAULT ''
);

CREATE TABLE role_permissions (
    role_id       uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id uuid NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE tenant_users (
    tenant_id  uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    uuid NOT NULL REFERENCES roles(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, user_id)
);
CREATE INDEX tenant_users_user_idx ON tenant_users (user_id);

CREATE TABLE refresh_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX refresh_tokens_user_idx ON refresh_tokens (user_id);

-- ===== Konten berjenjang ala WordPress ======================================

CREATE TABLE content_types (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    slug            text NOT NULL,
    name            text NOT NULL,
    is_hierarchical boolean NOT NULL DEFAULT false, -- page=true, post=false
    created_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, slug)
);

CREATE TABLE contents (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    content_type_id uuid NOT NULL REFERENCES content_types(id),
    parent_id       uuid REFERENCES contents(id) ON DELETE SET NULL, -- hierarki ala post_parent
    author_id       uuid NOT NULL REFERENCES users(id),
    title           text NOT NULL,
    slug            text NOT NULL,
    excerpt         text NOT NULL DEFAULT '',
    body            text NOT NULL DEFAULT '',
    status          text NOT NULL DEFAULT 'draft'
                    CHECK (status IN ('draft','review','scheduled','published','archived')),
    published_at    timestamptz,
    scheduled_at    timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    deleted_at      timestamptz, -- soft delete = "trashed"
    search_vector   tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(excerpt, '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(body, '')), 'C')
    ) STORED
);
-- Resolve URL & keunikan slug per (tenant, type); baris soft-deleted dibebaskan.
CREATE UNIQUE INDEX contents_slug_unique
    ON contents (tenant_id, content_type_id, slug) WHERE deleted_at IS NULL;
-- Pola listing dominan (docs/04-desain-data.md §3).
CREATE INDEX contents_listing_idx
    ON contents (tenant_id, status, published_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX contents_parent_idx ON contents (tenant_id, parent_id) WHERE deleted_at IS NULL;
CREATE INDEX contents_search_idx ON contents USING gin (search_vector);

-- EAV ala wp_postmeta (docs/04-desain-data.md §2.2).
CREATE TABLE content_meta (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id uuid NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    meta_key   text NOT NULL,
    meta_value jsonb NOT NULL DEFAULT 'null',
    UNIQUE (content_id, meta_key)
);
CREATE INDEX content_meta_key_idx ON content_meta (meta_key);
CREATE INDEX content_meta_value_idx ON content_meta USING gin (meta_value);

CREATE TABLE content_seo (
    content_id        uuid PRIMARY KEY REFERENCES contents(id) ON DELETE CASCADE,
    meta_title        text NOT NULL DEFAULT '',
    meta_description  text NOT NULL DEFAULT '',
    canonical_url     text NOT NULL DEFAULT '',
    og_image_media_id uuid, -- FK ke media ditambahkan setelah tabel media dibuat
    noindex           boolean NOT NULL DEFAULT false
);

CREATE TABLE revisions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id    uuid NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    author_id     uuid NOT NULL REFERENCES users(id),
    title         text NOT NULL,
    excerpt       text NOT NULL DEFAULT '',
    body          text NOT NULL DEFAULT '',
    meta_snapshot jsonb NOT NULL DEFAULT '{}',
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX revisions_content_idx ON revisions (content_id, created_at DESC);

-- ===== Taxonomy berjenjang ===================================================

CREATE TABLE taxonomies (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    slug            text NOT NULL,
    name            text NOT NULL,
    is_hierarchical boolean NOT NULL DEFAULT false, -- category=true, tag=false
    created_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, slug)
);

CREATE TABLE terms (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    taxonomy_id uuid NOT NULL REFERENCES taxonomies(id) ON DELETE CASCADE,
    parent_id   uuid REFERENCES terms(id) ON DELETE SET NULL, -- term berjenjang
    name        text NOT NULL,
    slug        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (taxonomy_id, slug)
);
CREATE INDEX terms_taxonomy_idx ON terms (tenant_id, taxonomy_id);

CREATE TABLE term_relationships (
    content_id uuid NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    term_id    uuid NOT NULL REFERENCES terms(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, term_id)
);
CREATE INDEX term_relationships_term_idx ON term_relationships (term_id);

-- ===== Media =================================================================

CREATE TABLE media (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    uploader_id uuid NOT NULL REFERENCES users(id),
    file_name   text NOT NULL,
    object_key  text NOT NULL UNIQUE, -- path MinIO: t/{tenant_id}/media/...
    mime_type   text NOT NULL,
    size_bytes  bigint NOT NULL DEFAULT 0,
    width       integer,
    height      integer,
    alt_text    text NOT NULL DEFAULT '',
    variants    jsonb NOT NULL DEFAULT '{}',
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX media_tenant_idx ON media (tenant_id, created_at DESC);

ALTER TABLE content_seo
    ADD CONSTRAINT content_seo_og_image_fk
    FOREIGN KEY (og_image_media_id) REFERENCES media(id) ON DELETE SET NULL;

-- ===== Keyword & SEO =========================================================

CREATE TABLE keywords (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    phrase     text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, phrase)
);

CREATE TABLE content_keywords (
    content_id uuid NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    keyword_id uuid NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    is_primary boolean NOT NULL DEFAULT false,
    PRIMARY KEY (content_id, keyword_id)
);
-- Maks satu keyword utama per konten (docs/06-keyword-seo-automasi.md §1).
CREATE UNIQUE INDEX content_keywords_primary_unique
    ON content_keywords (content_id) WHERE is_primary;
CREATE INDEX content_keywords_keyword_idx ON content_keywords (keyword_id);

-- ===== Lisensi, integrasi & audit ===========================================

CREATE TABLE licenses (
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            uuid NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
    license_key          text NOT NULL UNIQUE,
    plan                 text NOT NULL DEFAULT 'pro'
                         CHECK (plan IN ('starter','pro','agency')),
    seats                integer NOT NULL DEFAULT 15,
    status               text NOT NULL DEFAULT 'trialing'
                         CHECK (status IN ('trialing','active','past_due','grace','suspended','cancelled')),
    current_period_start timestamptz,
    current_period_end   timestamptz,
    provider             text NOT NULL DEFAULT 'manual', -- manual (beta) | paddle | midtrans
    provider_ref         text NOT NULL DEFAULT '',
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE billing_events (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    license_id uuid NOT NULL REFERENCES licenses(id) ON DELETE CASCADE,
    type       text NOT NULL,
    payload    jsonb NOT NULL DEFAULT '{}',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX billing_events_license_idx ON billing_events (license_id, created_at DESC);

CREATE TABLE api_keys (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name         text NOT NULL,
    key_hash     text NOT NULL UNIQUE, -- SHA-256; plaintext hanya tampil sekali
    scopes       text[] NOT NULL DEFAULT '{content.view}',
    expires_at   timestamptz,
    last_used_at timestamptz,
    revoked_at   timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX api_keys_tenant_idx ON api_keys (tenant_id);

CREATE TABLE webhooks (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    url        text NOT NULL,
    events     text[] NOT NULL DEFAULT '{}',
    secret     text NOT NULL, -- HMAC signature payload keluar
    is_active  boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE webhook_deliveries (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id      uuid NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event           text NOT NULL,
    payload         jsonb NOT NULL DEFAULT '{}',
    status          text NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending','delivered','failed')),
    attempts        integer NOT NULL DEFAULT 0,
    last_attempt_at timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX webhook_deliveries_pending_idx
    ON webhook_deliveries (webhook_id, created_at) WHERE status = 'pending';

CREATE TABLE audit_logs (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     uuid REFERENCES users(id) ON DELETE SET NULL,
    action      text NOT NULL,
    entity_type text NOT NULL,
    entity_id   uuid,
    changes     jsonb NOT NULL DEFAULT '{}',
    ip          text NOT NULL DEFAULT '',
    user_agent  text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_tenant_idx ON audit_logs (tenant_id, created_at DESC);

-- ===== Row-Level Security ====================================================
-- FORCE agar policy berlaku juga untuk pemilik tabel (user aplikasi di
-- docker-compose adalah owner DB). app.tenant_id di-set per transaksi oleh
-- database.WithTenant.

DO $$
DECLARE t text;
BEGIN
    FOREACH t IN ARRAY ARRAY[
        'content_types','contents','taxonomies','terms','media','keywords',
        'api_keys','webhooks','audit_logs'
    ] LOOP
        EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
        EXECUTE format('ALTER TABLE %I FORCE ROW LEVEL SECURITY', t);
        EXECUTE format(
            'CREATE POLICY tenant_isolation ON %I
                USING (tenant_id = current_setting(''app.tenant_id'', true)::uuid)
                WITH CHECK (tenant_id = current_setting(''app.tenant_id'', true)::uuid)', t);
    END LOOP;
END $$;

-- ===== Seed: role sistem & permission (docs/05-rbac.md §2) ==================

INSERT INTO permissions (slug, description) VALUES
    ('content.view',     'Melihat konten'),
    ('content.create',   'Membuat konten'),
    ('content.edit_own', 'Mengubah konten milik sendiri'),
    ('content.edit_any', 'Mengubah konten siapa pun'),
    ('content.publish',  'Menerbitkan konten'),
    ('content.delete',   'Menghapus (trash) konten'),
    ('taxonomy.manage',  'Mengelola taxonomy & term'),
    ('media.upload',     'Mengunggah media'),
    ('media.manage',     'Mengelola/menghapus media'),
    ('keyword.manage',   'Mengelola keyword & SEO'),
    ('user.invite',      'Mengundang anggota'),
    ('user.manage',      'Mengelola anggota & role'),
    ('apikey.manage',    'Mengelola API key'),
    ('webhook.manage',   'Mengelola webhook'),
    ('settings.manage',  'Mengelola pengaturan tenant'),
    ('billing.view',     'Melihat lisensi & tagihan'),
    ('billing.manage',   'Mengelola lisensi & tagihan');

INSERT INTO roles (tenant_id, slug, name, description) VALUES
    (NULL, 'owner',       'Owner',       'Pemilik tenant, pemegang billing'),
    (NULL, 'admin',       'Admin',       'Kelola tenant tanpa billing'),
    (NULL, 'editor',      'Editor',      'Kelola semua konten'),
    (NULL, 'author',      'Author',      'Tulis & terbitkan konten sendiri'),
    (NULL, 'contributor', 'Contributor', 'Tulis draft tanpa publish'),
    (NULL, 'viewer',      'Viewer',      'Baca saja');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON (
       (r.slug = 'owner')
    OR (r.slug = 'admin'       AND p.slug NOT IN ('billing.view','billing.manage'))
    OR (r.slug = 'editor'      AND p.slug IN ('content.view','content.create','content.edit_own',
                                              'content.edit_any','content.publish','content.delete',
                                              'taxonomy.manage','media.upload','media.manage','keyword.manage'))
    OR (r.slug = 'author'      AND p.slug IN ('content.view','content.create','content.edit_own',
                                              'content.publish','media.upload'))
    OR (r.slug = 'contributor' AND p.slug IN ('content.view','content.create','content.edit_own','media.upload'))
    OR (r.slug = 'viewer'      AND p.slug = 'content.view')
)
WHERE r.tenant_id IS NULL;

-- Catatan: content_types & taxonomies default (post, page / category, tag)
-- di-provision per tenant oleh tenant.Service.Create, bukan di-seed global
-- (NEXT: internal/modules/tenant/repository.go).
