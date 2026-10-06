package auth

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Nama tenant pertama dibuat sekaligus saat register (trial 14 hari,
	// docs/02-model-bisnis.md §1). RELASI: internal/modules/tenant — pembuatan
	// tenant + keanggotaan owner; internal/modules/license — lisensi trialing.
	TenantName string `json:"tenant_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // detik
}
