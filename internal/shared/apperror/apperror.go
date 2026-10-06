// Package apperror mendefinisikan error domain ber-kode stabil yang dipetakan
// handler ke HTTP lewat shared/response — detail internal tidak pernah bocor ke klien.
package apperror

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int    // kode HTTP
	Code    string // kode mesin stabil, mis. "content_not_found"
	Message string // pesan aman untuk klien
	Err     error  // penyebab internal (untuk log, bukan respons)
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

func (e *Error) Unwrap() error { return e.Err }

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func (e *Error) Wrap(err error) *Error {
	return &Error{Status: e.Status, Code: e.Code, Message: e.Message, Err: err}
}

func NotFound(code, message string) *Error   { return New(http.StatusNotFound, code, message) }
func BadRequest(code, message string) *Error { return New(http.StatusBadRequest, code, message) }
func Conflict(code, message string) *Error   { return New(http.StatusConflict, code, message) }
func Forbidden(message string) *Error        { return New(http.StatusForbidden, "forbidden", message) }
func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, "unauthorized", message)
}

func Internal(err error) *Error {
	return &Error{Status: http.StatusInternalServerError, Code: "internal_error",
		Message: "terjadi kesalahan internal", Err: err}
}

// NotImplemented dipakai handler/service modul kerangka (stub) agar jalur
// wiring tetap teruji end-to-end sebelum logika nyata ditulis.
func NotImplemented(feature string) *Error {
	return New(http.StatusNotImplemented, "not_implemented",
		fmt.Sprintf("fitur %q belum diimplementasi; lihat marker NEXT: pada modul terkait", feature))
}
