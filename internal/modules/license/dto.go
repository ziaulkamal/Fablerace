package license

import (
	"strings"

	"github.com/ziaulkamal/fablerace/internal/shared/apperror"
)

type ActivateRequest struct {
	LicenseKey string `json:"license_key"`
}

func (r *ActivateRequest) Validate() error {
	r.LicenseKey = strings.TrimSpace(r.LicenseKey)
	if r.LicenseKey == "" {
		return apperror.BadRequest("license_key_required", "license_key wajib diisi")
	}
	return nil
}

// UsageResponse: status lisensi + pemakaian kuota (metering docs/02 §4).
type UsageResponse struct {
	Plan         Plan   `json:"plan"`
	Status       Status `json:"status"`
	Seats        int    `json:"seats"`
	SeatsUsed    int    `json:"seats_used"`
	StorageBytes int64  `json:"storage_bytes"`
	StorageQuota int64  `json:"storage_quota"`
	ContentCount int    `json:"content_count"`
}
