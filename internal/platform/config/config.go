// Package config memuat konfigurasi 12-factor dari environment variable
// (docs/09-devops.md §2). Daftar variabel terdokumentasi di .env.example.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv   string // dev | staging | prod
	HTTPAddr string

	DatabaseURL string

	RedisAddr     string
	RedisPassword string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool

	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:   getenv("APP_ENV", "dev"),
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),

		DatabaseURL: getenv("DATABASE_URL",
			"postgres://fablerace:fablerace@localhost:5432/fablerace?sslmode=disable"),

		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		MinioEndpoint:  getenv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getenv("MINIO_ACCESS_KEY", "fablerace"),
		MinioSecretKey: getenv("MINIO_SECRET_KEY", "fablerace-secret"),
		MinioBucket:    getenv("MINIO_BUCKET", "fablerace"),
		MinioUseSSL:    getenvBool("MINIO_USE_SSL", false),

		JWTSecret:       getenv("JWT_SECRET", ""),
		AccessTokenTTL:  getenvDur("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getenvDur("REFRESH_TOKEN_TTL", 30*24*time.Hour),
	}

	if cfg.JWTSecret == "" {
		if cfg.AppEnv != "dev" {
			return nil, fmt.Errorf("config: JWT_SECRET wajib diisi di luar env dev")
		}
		// Hanya untuk dev lokal — token tidak akan valid lintas restart pun tak masalah.
		cfg.JWTSecret = "dev-only-insecure-secret"
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func getenvDur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
