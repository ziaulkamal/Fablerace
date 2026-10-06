// Package media: media library di MinIO — presign upload, confirm, varian
// gambar (docs/07-storage-media.md).
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package media

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	UploaderID uuid.UUID
	FileName   string
	ObjectKey  string // path MinIO — SELALU hasil storage.ObjectKey, tak pernah dari klien
	MimeType   string
	SizeBytes  int64
	Width      *int
	Height     *int
	AltText    string
	Variants   json.RawMessage // {"thumb_640": "…", "thumb_1280": "…"} hasil job media:process
	CreatedAt  time.Time
}
