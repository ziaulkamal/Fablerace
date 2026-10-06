// Package slugutil menormalkan judul menjadi slug URL
// (docs/06-keyword-seo-automasi.md §2 "Auto-slug").
package slugutil

import (
	"fmt"
	"strings"
	"unicode"
)

const maxLen = 160

// Make menurunkan judul ke slug: lowercase, non-alfanumerik → '-', tanpa
// '-' ganda/tepi. Keunikan per (tenant, content_type) ditegakkan pemanggil
// lewat repository (RELASI: internal/modules/content — service.generateSlug).
func Make(s string) string {
	var b strings.Builder
	prevDash := true // cegah dash di awal
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case unicode.IsLetter(r) && r < 128, unicode.IsDigit(r):
			b.WriteRune(r)
			prevDash = false
		case unicode.IsLetter(r):
			// NEXT: tabel transliterasi (é→e, ü→u, aksara non-Latin) agar slug
			// tetap informatif untuk judul ber-diakritik; sementara dilewati.
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		default:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > maxLen {
		out = strings.Trim(out[:maxLen], "-")
	}
	if out == "" {
		out = "untitled"
	}
	return out
}

// WithSuffix menambah sufiks angka saat bentrok: "judul" → "judul-2".
func WithSuffix(base string, n int) string {
	return fmt.Sprintf("%s-%d", base, n)
}
