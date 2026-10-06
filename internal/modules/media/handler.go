package media

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Presign(w http.ResponseWriter, r *http.Request) {
	var in PresignRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, apperror.BadRequest("body_invalid", "body JSON tidak valid"))
		return
	}
	res, err := h.svc.Presign(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, res)
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apperror.BadRequest("id_invalid", "parameter id harus UUID"))
		return
	}
	m, err := h.svc.Confirm(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, m)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	p := pagination.FromQuery(r)
	items, total, err := h.svc.List(r.Context(), p)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.List(w, items, response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apperror.BadRequest("id_invalid", "parameter id harus UUID"))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
