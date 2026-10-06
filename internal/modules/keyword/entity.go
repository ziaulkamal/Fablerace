// Package keyword: manajemen keyword per konten + automasi SEO
// (docs/06-keyword-seo-automasi.md). Diferensiasi utama produk.
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package keyword

import (
	"time"

	"github.com/google/uuid"
)

type Keyword struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Phrase    string // unik per tenant → dasar analisis coverage/kanibalisasi
	CreatedAt time.Time
}

// ContentKeyword memetakan content_keywords; IsPrimary maksimal satu per konten.
// RELASI: internal/modules/content — konten yang dianalisis; event publish dari
// content.Service.Publish memicu rebuild sitemap & webhook di modul ini.
type ContentKeyword struct {
	ContentID uuid.UUID
	KeywordID uuid.UUID
	IsPrimary bool
}

// SEOMeta memetakan content_seo (dipisah dari meta EAV karena selalu dibaca
// bersama konten — docs/04-desain-data.md §2.2).
type SEOMeta struct {
	ContentID       uuid.UUID
	MetaTitle       string
	MetaDescription string
	CanonicalURL    string
	OGImageMediaID  *uuid.UUID
	NoIndex         bool
}
