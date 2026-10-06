package license

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

func (h *Handler) Usage(w http.ResponseWriter, r *http.Request) {
	u, err := h.svc.Usage(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, u)
}

func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	var in ActivateRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, apperror.BadRequest("body_invalid", "body JSON tidak valid"))
		return
	}
	l, err := h.svc.Activate(r.Context(), in)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.OK(w, l)
}

// NEXT: handler webhook payment gateway di path publik /billing/webhook
// (verifikasi signature SEBELUM parsing — docs/08-keamanan.md §6 SSRF/forgery),
// dipasang server di luar grup tenant.
