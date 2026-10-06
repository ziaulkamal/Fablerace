package user

import (
	"encoding/json"
	"net/http"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListMembers(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	var in InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, apperror.BadRequest("body_invalid", "body JSON tidak valid"))
		return
	}
	if err := h.svc.Invite(r.Context(), in); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
