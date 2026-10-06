package auth

import (
	"context"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Service interface {
	Register(ctx context.Context, in RegisterRequest) (*TokenResponse, error)
	Login(ctx context.Context, in LoginRequest) (*TokenResponse, error)
	Refresh(ctx context.Context, in RefreshRequest) (*TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

// NEXT (urutan implementasi — docs/10-roadmap.md minggu 1–2):
//  1. Register: hash Argon2id (parameter OWASP), buat user + tenant pertama +
//     lisensi trialing dalam SATU transaksi.
//     RELASI: internal/modules/tenant.Service.Create & internal/modules/license —
//     dipanggil via interface service, di-inject dari cmd/api/main.go.
//  2. Login: verifikasi Argon2id, terbitkan JWT HS256 (klaim minimal sub/exp/iat)
//     + refresh token opaque (hash SHA-256 ke tabel).
//     RELASI: internal/platform/middleware/auth.go — verifikator token yang
//     diterbitkan di sini; secret & TTL sama-sama dari platform/config.
//  3. Refresh: ROTASI — token lama revoke, terbit pasangan baru; pemakaian ulang
//     token ter-revoke = indikasi pencurian → revoke seluruh keluarga token user
//     (docs/08-keamanan.md §2).
//  4. Rate limit login 5 gagal/menit per email+IP
//     (RELASI: internal/platform/middleware/ratelimit.go — limiter khusus auth).

func (s *service) Register(ctx context.Context, in RegisterRequest) (*TokenResponse, error) {
	return nil, apperror.NotImplemented("auth.register")
}

func (s *service) Login(ctx context.Context, in LoginRequest) (*TokenResponse, error) {
	return nil, apperror.NotImplemented("auth.login")
}

func (s *service) Refresh(ctx context.Context, in RefreshRequest) (*TokenResponse, error) {
	return nil, apperror.NotImplemented("auth.refresh")
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return apperror.NotImplemented("auth.logout")
}
