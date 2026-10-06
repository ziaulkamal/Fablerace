package keyword

import (
	"strings"

	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type CreateRequest struct {
	Phrase string `json:"phrase"`
}

func (r *CreateRequest) Validate() error {
	r.Phrase = strings.TrimSpace(strings.ToLower(r.Phrase))
	if r.Phrase == "" || len(r.Phrase) > 160 {
		return apperror.BadRequest("phrase_invalid", "phrase wajib diisi, maksimal 160 karakter")
	}
	return nil
}

// AttachRequest mengganti seluruh keyword sebuah konten (idempoten);
// PrimaryID harus salah satu anggota KeywordIDs.
type AttachRequest struct {
	KeywordIDs []uuid.UUID `json:"keyword_ids"`
	PrimaryID  *uuid.UUID  `json:"primary_id"`
}

// CoverageItem: satu baris jawaban "konten mana saja membidik keyword ini" —
// deteksi kanibalisasi (docs/06-keyword-seo-automasi.md §1).
type CoverageItem struct {
	ContentID uuid.UUID `json:"content_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	IsPrimary bool      `json:"is_primary"`
}
