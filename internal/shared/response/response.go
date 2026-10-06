// Package response menyatukan envelope JSON API (docs/03-arsitektur.md §4):
// sukses  → {"success":true,"data":...,"meta":{...}}
// gagal   → {"success":false,"error":{"code":...,"message":...}}
package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type envelope struct {
	Success bool     `json:"success"`
	Data    any      `json:"data,omitempty"`
	Meta    *Meta    `json:"meta,omitempty"`
	Error   *errBody `json:"error,omitempty"`
}

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func write(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func JSON(w http.ResponseWriter, status int, data any) {
	write(w, status, envelope{Success: true, Data: data})
}
func OK(w http.ResponseWriter, data any)      { JSON(w, http.StatusOK, data) }
func Created(w http.ResponseWriter, data any) { JSON(w, http.StatusCreated, data) }
func NoContent(w http.ResponseWriter)         { w.WriteHeader(http.StatusNoContent) }

func List(w http.ResponseWriter, data any, meta Meta) {
	write(w, http.StatusOK, envelope{Success: true, Data: data, Meta: &meta})
}

// Error memetakan *apperror.Error ke HTTP; error lain menjadi 500 generik
// dan penyebab aslinya hanya masuk log — tidak pernah ke klien (docs/08 §6).
func Error(w http.ResponseWriter, err error) {
	var ae *apperror.Error
	if !errors.As(err, &ae) {
		ae = apperror.Internal(err)
	}
	if ae.Status >= http.StatusInternalServerError {
		slog.Error("request failed", "code", ae.Code, "err", ae.Err)
	}
	write(w, ae.Status, envelope{Success: false, Error: &errBody{Code: ae.Code, Message: ae.Message}})
}
