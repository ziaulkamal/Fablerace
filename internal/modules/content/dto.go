package content

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// DTO eksplisit per endpoint — tidak pernah bind/echo struct DB mentah
// (guard mass-assignment, docs/08-keamanan.md §6).

type CreateRequest struct {
	ContentTypeID uuid.UUID  `json:"content_type_id"`
	ParentID      *uuid.UUID `json:"parent_id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"` // kosong → dibuat otomatis dari judul
	Excerpt       string     `json:"excerpt"`
	Body          string     `json:"body"`
	Status        Status     `json:"status"` // kosong → draft
	ScheduledAt   *time.Time `json:"scheduled_at"`
}

func (r *CreateRequest) Validate() error {
	if r.ContentTypeID == uuid.Nil {
		return apperror.BadRequest("content_type_required", "content_type_id wajib diisi")
	}
	if r.Title == "" || len(r.Title) > 300 {
		return apperror.BadRequest("title_invalid", "title wajib diisi, maksimal 300 karakter")
	}
	if r.Status == "" {
		r.Status = StatusDraft
	}
	if !r.Status.Valid() {
		return apperror.BadRequest("status_invalid", "status tidak dikenal")
	}
	if r.Status == StatusScheduled && r.ScheduledAt == nil {
		return apperror.BadRequest("scheduled_at_required", "status scheduled butuh scheduled_at")
	}
	return nil
}

// UpdateRequest: field pointer = partial update (nil berarti tidak diubah).
type UpdateRequest struct {
	Title       *string    `json:"title"`
	Slug        *string    `json:"slug"`
	Excerpt     *string    `json:"excerpt"`
	Body        *string    `json:"body"`
	ParentID    *uuid.UUID `json:"parent_id"` // NEXT: dukung unset parent (perlu sentinel eksplisit, bukan nil)
	Status      *Status    `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

func (r *UpdateRequest) Validate() error {
	if r.Title != nil && (*r.Title == "" || len(*r.Title) > 300) {
		return apperror.BadRequest("title_invalid", "title tidak boleh kosong, maksimal 300 karakter")
	}
	if r.Status != nil && !r.Status.Valid() {
		return apperror.BadRequest("status_invalid", "status tidak dikenal")
	}
	return nil
}

// MetaRequest: upsert penuh meta EAV — key → nilai JSON bebas.
type MetaRequest map[string]json.RawMessage

type ListFilter struct {
	ContentTypeID *uuid.UUID
	Status        *Status
	ParentID      *uuid.UUID
	Search        string
}

type Response struct {
	ID            uuid.UUID      `json:"id"`
	ContentTypeID uuid.UUID      `json:"content_type_id"`
	ParentID      *uuid.UUID     `json:"parent_id,omitempty"`
	AuthorID      uuid.UUID      `json:"author_id"`
	Title         string         `json:"title"`
	Slug          string         `json:"slug"`
	Excerpt       string         `json:"excerpt,omitempty"`
	Body          string         `json:"body,omitempty"`
	Status        Status         `json:"status"`
	PublishedAt   *time.Time     `json:"published_at,omitempty"`
	ScheduledAt   *time.Time     `json:"scheduled_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Meta          []MetaResponse `json:"meta,omitempty"`
}

type MetaResponse struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func toResponse(c *Content, meta []MetaItem) Response {
	res := Response{
		ID: c.ID, ContentTypeID: c.ContentTypeID, ParentID: c.ParentID,
		AuthorID: c.AuthorID, Title: c.Title, Slug: c.Slug, Excerpt: c.Excerpt,
		Body: c.Body, Status: c.Status, PublishedAt: c.PublishedAt,
		ScheduledAt: c.ScheduledAt, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
	for _, m := range meta {
		res.Meta = append(res.Meta, MetaResponse{Key: m.Key, Value: m.Value})
	}
	return res
}

func toResponses(cs []Content) []Response {
	out := make([]Response, 0, len(cs))
	for i := range cs {
		out = append(out, toResponse(&cs[i], nil))
	}
	return out
}
