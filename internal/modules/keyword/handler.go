package keyword

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func idParam(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, apperror.BadRequest("id_invalid", "parameter id harus UUID")
	}
	return id, nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, apperror.BadRequest("body_invalid", "body JSON tidak valid"))
		return
	}
	k, err := h.svc.Create(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, k)
}

func (h *Handler) Coverage(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.svc.Coverage(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

// AttachToContent dipasang di /contents/{id}/keywords oleh server — endpoint
// milik modul keyword meski path-nya di bawah contents
// (RELASI: internal/platform/server — pemasangan lintas-path, pola yang sama
// dengan taxonomy.SetContentTerms).
func (h *Handler) AttachToContent(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	var in AttachRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, apperror.BadRequest("body_invalid", "body JSON tidak valid"))
		return
	}
	if err := h.svc.AttachToContent(r.Context(), id, in); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
