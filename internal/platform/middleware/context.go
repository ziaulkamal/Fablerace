// Package middleware berisi rantai guard HTTP: auth → tenant → license → rbac
// (urutan di docs/05-rbac.md §3). File ini memegang kunci context bersama —
// modul domain membaca identitas HANYA lewat helper di sini, tidak pernah dari
// header mentah (docs/lampiran-b-konvensi-kode.md §5).
package middleware

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxTenantID
	ctxPermissions
)

func withUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxUserID, id)
}

func withTenantID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxTenantID, id)
}

func withPermissions(ctx context.Context, perms map[string]struct{}) context.Context {
	return context.WithValue(ctx, ctxPermissions, perms)
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxUserID).(uuid.UUID)
	return id, ok
}

func TenantID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxTenantID).(uuid.UUID)
	return id, ok
}

// HasPermission dipakai service untuk aturan gabungan permission+ownership,
// mis. content.edit_own vs content.edit_any (RELASI: internal/modules/content
// — service.authorizeMutation).
func HasPermission(ctx context.Context, slug string) bool {
	perms, ok := ctx.Value(ctxPermissions).(map[string]struct{})
	if !ok {
		return false
	}
	_, ok = perms[slug]
	return ok
}
