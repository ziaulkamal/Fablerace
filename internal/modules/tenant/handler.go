package tenant

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListMine(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CreateRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	t, err := h.svc.Create(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, t)
}

func (h *Handler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	t, err := h.svc.GetCurrent(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, t)
}

func (h *Handler) UpdateCurrent(w http.ResponseWriter, r *http.Request) {
	var in UpdateSettingsRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	t, err := h.svc.UpdateCurrent(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, t)
}
