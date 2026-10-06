// Package cache memegang klien Redis: cache read-through, rate limit, dan
// backend antrian jobs (docs/03-arsitektur.md §5–6). Kunci cache WAJIB
// berprefix "t:{tenant_id}:" untuk data per-tenant (docs/08-keamanan.md §3).
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, addr, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: ping redis: %w", err)
	}
	return rdb, nil
}

// NEXT: worker asynq (scheduled publish, webhook deliver, sitemap, metering —
// docs/03-arsitektur.md §5) memakai koneksi Redis yang sama; inisialisasi
// client asynq ditambahkan di paket ini saat modul job pertama ditulis.
