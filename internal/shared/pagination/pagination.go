// Package pagination membaca ?page=&per_page= dengan batas aman
// (default 20, maks 100 — docs/03-arsitektur.md §4).
package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type Params struct {
	Page    int
	PerPage int
}

func FromQuery(r *http.Request) Params {
	p := Params{Page: 1, PerPage: DefaultPerPage}
	if v, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && v > 0 {
		p.Page = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && v > 0 {
		p.PerPage = min(v, MaxPerPage)
	}
	return p
}

func (p Params) Limit() int  { return p.PerPage }
func (p Params) Offset() int { return (p.Page - 1) * p.PerPage }

// NEXT: cursor-based pagination untuk feed publik bervolume besar
// (docs/03-arsitektur.md §4) — ditambah di paket ini agar envelope konsisten.
