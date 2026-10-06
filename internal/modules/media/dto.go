package media

import (
	"github.com/google/uuid"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

// Whitelist MIME & batas ukuran (docs/07-storage-media.md §3). Deteksi
// magic-byte menyusul di alur confirm — ekstensi/MIME kiriman klien hanyalah
// klaim awal.
var allowedMime = map[string]int64{
	"image/jpeg":      10 << 20,
	"image/png":       10 << 20,
	"image/webp":      10 << 20,
	"image/avif":      10 << 20,
	"image/gif":       10 << 20,
	"application/pdf": 10 << 20,
	"video/mp4":       100 << 20, // khusus Pro+ (enforcement tier menyusul di service)
}

type PresignRequest struct {
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

func (r *PresignRequest) Validate() error {
	if r.FileName == "" {
		return apperror.BadRequest("file_name_required", "file_name wajib diisi")
	}
	max, ok := allowedMime[r.MimeType]
	if !ok {
		return apperror.BadRequest("mime_not_allowed", "tipe file tidak diizinkan")
	}
	if r.SizeBytes <= 0 || r.SizeBytes > max {
		return apperror.BadRequest("size_invalid", "ukuran file di luar batas yang diizinkan")
	}
	return nil
}

type PresignResponse struct {
	MediaID   uuid.UUID `json:"media_id"`
	UploadURL string    `json:"upload_url"`
	ExpiresIn int64     `json:"expires_in"` // detik
}
