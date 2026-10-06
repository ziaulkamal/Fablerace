// Package auth: registrasi, login, rotasi refresh token (docs/08-keamanan.md §2).
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package auth

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken memetakan tabel refresh_tokens; yang disimpan hash-nya,
// bukan token mentah (docs/04-desain-data.md §2.1).
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
