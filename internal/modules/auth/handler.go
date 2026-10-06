package auth

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

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.BadRequest("body_invalid", "body JSON tidak valid")
	}
	return nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in RegisterRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	tok, err := h.svc.Register(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, tok)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in LoginRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	tok, err := h.svc.Login(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, tok)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var in RefreshRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	tok, err := h.svc.Refresh(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, tok)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var in RefreshRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.svc.Logout(r.Context(), in.RefreshToken); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
