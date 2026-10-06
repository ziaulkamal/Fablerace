package user

import (
	"strings"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type InviteRequest struct {
	Email    string `json:"email"`
	RoleSlug string `json:"role_slug"` // admin|editor|author|contributor|viewer — owner tidak bisa di-invite
}

func (r *InviteRequest) Validate() error {
	if !strings.Contains(r.Email, "@") {
		return apperror.BadRequest("email_invalid", "email tidak valid")
	}
	switch r.RoleSlug {
	case "admin", "editor", "author", "contributor", "viewer":
		return nil
	}
	return apperror.BadRequest("role_invalid",
		"role_slug harus salah satu: admin, editor, author, contributor, viewer")
}

type UpdateRoleRequest struct {
	RoleSlug string `json:"role_slug"`
}
