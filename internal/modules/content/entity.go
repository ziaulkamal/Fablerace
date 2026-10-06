// Package content adalah modul EXEMPLAR Fablerace: CRUD konten berjenjang ala
// WordPress (parent_id), slug otomatis, meta EAV, dan siklus status terbit.
// Modul lain meniru anatomi file & pola di sini
// (docs/lampiran-b-konvensi-kode.md §2, skema di docs/04-desain-data.md §2.2).
package content

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusReview    Status = "review"
	StatusScheduled Status = "scheduled"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusReview, StatusScheduled, StatusPublished, StatusArchived:
		return true
	}
	return false
}

// Content memetakan tabel contents. Hierarki lewat ParentID (self-FK);
// soft delete lewat deleted_at (status "trashed" versi WordPress).
type Content struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	ContentTypeID uuid.UUID
	ParentID      *uuid.UUID
	AuthorID      uuid.UUID
	Title         string
	Slug          string
	Excerpt       string
	Body          string
	Status        Status
	PublishedAt   *time.Time
	ScheduledAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// MetaItem memetakan content_meta (EAV ala wp_postmeta); Value jsonb agar
// nilai terstruktur tetap bisa di-query (docs/04-desain-data.md §2.2).
type MetaItem struct {
	Key   string
	Value json.RawMessage
}

// RELASI: internal/modules/taxonomy — penautan konten↔term hidup di tabel
// term_relationships milik modul taxonomy; penghapusan konten mengandalkan
// ON DELETE CASCADE di migrasi, bukan kode di modul ini.
// RELASI: internal/modules/keyword — content_keywords menautkan konten ke
// keyword; blok SEO (content_seo) juga dikelola modul keyword.
