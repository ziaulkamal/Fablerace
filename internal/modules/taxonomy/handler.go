package taxonomy

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

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.BadRequest("body_invalid", "body JSON tidak valid")
	}
	return nil
}

func idParam(r *http.Request, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		return uuid.Nil, apperror.BadRequest("id_invalid", "parameter "+name+" harus UUID")
	}
	return id, nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListTaxonomies(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CreateTaxonomyRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	t, err := h.svc.CreateTaxonomy(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, t)
}

func (h *Handler) ListTerms(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.svc.ListTerms(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, items)
}

func (h *Handler) CreateTerm(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	var in CreateTermRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	t, err := h.svc.CreateTerm(r.Context(), id, in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, t)
}

// SetContentTerms dipasang di path /contents/{id}/terms oleh server —
// endpoint milik modul taxonomy meski path-nya di bawah contents
// (RELASI: internal/platform/server — pemasangan lintas-path).
func (h *Handler) SetContentTerms(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}
	var in SetContentTermsRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.svc.SetContentTerms(r.Context(), id, in); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
