package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
	"github.com/ziaulkamal/fablerace/internal/shared/response"
)

// RateLimit: fixed-window sederhana per kunci (IP) dalam memori proses —
// cukup untuk melindungi endpoint sensitif (login/presign) di deploy
// single-node fase awal (docs/08-keamanan.md §5).
//
// NEXT: pindah ke token bucket di Redis saat API punya >1 replika
// (docs/09-devops.md §5 fase 2) — state limiter harus bersama; sekaligus
// limit per-tenant & per-API-key sesuai tier (60/300 rpm, docs/02 §1).
func RateLimit(maxPerMinute int) func(http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		window  = time.Now().Truncate(time.Minute)
		counter = make(map[string]int)
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			now := time.Now().Truncate(time.Minute)
			if !now.Equal(window) {
				window = now
				counter = make(map[string]int)
			}
			counter[r.RemoteAddr]++
			over := counter[r.RemoteAddr] > maxPerMinute
			mu.Unlock()

			if over {
				w.Header().Set("Retry-After", "60")
				response.Error(w, apperror.New(http.StatusTooManyRequests,
					"rate_limited", "terlalu banyak permintaan, coba lagi sebentar"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
