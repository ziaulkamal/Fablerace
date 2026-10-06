// Package user: profil user & manajemen anggota per tenant (undang, ubah role).
// Status: KERANGKA — wiring lengkap, logika menyusul (lihat marker NEXT:).
package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Email           string
	Name            string
	IsPlatformAdmin bool
	CreatedAt       time.Time
}

// Member = user + role-nya di tenant aktif (join tenant_users → roles).
// RELASI: internal/modules/tenant — tabel keanggotaan; modul ini mengelola
// isi keanggotaan (role), modul tenant mengelola keberadaan tenant-nya.
type Member struct {
	UserID   uuid.UUID
	Email    string
	Name     string
	RoleSlug string
}
