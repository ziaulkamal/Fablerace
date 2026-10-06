package content

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/pagination"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// Handler hanya menerjemahkan HTTP ↔ service: decode DTO, panggil service,
// tulis envelope. Tanpa SQL, tanpa aturan bisnis (docs/lampiran-b §2).
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

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.BadRequest("body_invalid", "body JSON tidak valid")
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in CreateRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	c, err := h.svc.Create(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, toResponse(c, nil))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	c, meta, err := h.svc.Get(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, toResponse(c, meta))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f ListFilter
	if v, err := uuid.Parse(q.Get("type")); err == nil {
		f.ContentTypeID = &v
	}
	if v := Status(q.Get("status")); v != "" && v.Valid() {
		f.Status = &v
	}
	if v, err := uuid.Parse(q.Get("parent_id")); err == nil {
		f.ParentID = &v
	}
	f.Search = q.Get("search")

	p := pagination.FromQuery(r)
	items, total, err := h.svc.List(r.Context(), f, p)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.List(w, toResponses(items), response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	var in UpdateRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	c, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, toResponse(c, nil))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	c, err := h.svc.Publish(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, toResponse(c, nil))
}

func (h *Handler) Children(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	items, err := h.svc.Children(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, toResponses(items))
}

func (h *Handler) SetMeta(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		response.Error(w, err)
		return
	}
	var in MetaRequest
	if err := decode(r, &in); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.svc.SetMeta(r.Context(), id, in); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}
