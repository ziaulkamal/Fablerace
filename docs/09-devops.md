# 09 — DevOps & Deployment

> Bagian dari [Kajian Produk Fablerace](../README.md) · Sebelumnya: [08 — Keamanan](08-keamanan.md) · Lanjut: [10 — Roadmap](10-roadmap.md)

## 1. Containerization

- **Dockerfile multi-stage** ([`Dockerfile`](../Dockerfile)): stage build (golang:1.24-alpine, `CGO_ENABLED=0`) → stage runtime `alpine` minimal, user non-root, binary tunggal ±20 MB. Migrasi SQL ikut ter-embed di binary (golang-migrate `iofs`) sehingga image = aplikasi lengkap.
- **docker-compose** ([`deploy/docker-compose.yml`](../deploy/docker-compose.yml)): `api` + `postgres:16-alpine` + `redis:7-alpine` + `minio` — satu perintah `docker compose up` untuk dev maupun prod single-node; healthcheck antar service, volume persisten, jaringan internal (hanya `api` yang terekspos).

Prod single-node ditambah reverse proxy **Caddy** (TLS otomatis per custom domain tenant — fitur penting karena tiap tenant membawa domain sendiri, [02 §1](02-model-bisnis.md)).

## 2. Konfigurasi

12-factor: semua konfigurasi via env (`.env.example` terdokumentasi, dibaca `internal/platform/config`). Tidak ada file config per environment; tidak ada secret di repo.

## 3. CI/CD (GitHub Actions — pasca commit pertama)

```
push → lint (golangci-lint) → go vet → go test → build image → push registry → deploy
```

- PR: lint + vet + test wajib hijau.
- Tag `v*`: build & push image + deploy ke prod (SSH + `docker compose pull && up -d` — cukup sampai skala menuntut lebih).
- Migrasi dijalankan otomatis oleh aplikasi saat boot (idempoten, bertahan di versi yang sama bila gagal → container lama tetap melayani).

## 4. Observability

| Sinyal | MVP | Skala |
|---|---|---|
| Logging | Structured JSON (`log/slog`) ke stdout, dikumpulkan Docker; `request_id` + `tenant_id` di setiap baris | Loki/CloudWatch |
| Metrics | `GET /healthz` (liveness) + `GET /readyz` (cek DB/Redis/MinIO) | `/metrics` Prometheus + Grafana |
| Tracing | — (belum sepadan) | OpenTelemetry (`NEXT:` di `platform/server`) |
| Alerting | Uptime monitor eksternal (healthz) + alert disk/memori node | Alertmanager |

`tenant_id` di setiap log line penting: debugging "kenapa situs klien X lambat" adalah pertanyaan harian bisnis multi-tenant.

## 5. Jalur Scaling (ditulis sekarang supaya tidak menjadi keputusan panik)

1. **Fase 1 (0–±150 tenant)**: 1 node compose (semua service satu host) — sesuai asumsi unit economics [02 §3](02-model-bisnis.md).
2. **Fase 2**: pisahkan data plane — Postgres & MinIO ke node/managed sendiri; API tetap stateless → tambah replika API di belakang LB (JWT stateless & jobs di Redis membuat ini trivial).
3. **Fase 3**: worker jobs jadi deployment terpisah (binary sama, flag `--worker`); read replica Postgres untuk API publik read-heavy.
4. **Fase 4 (ribuan tenant)**: Kubernetes/Nomad bila armada node >3; arsitektur monolith modular tidak menghalangi — ia justru menunda kebutuhan ini selama mungkin.

## 6. Lingkungan

| Env | Tujuan | Infra |
|---|---|---|
| `dev` | lokal developer | docker compose (hot-reload via `air`, opsional) |
| `staging` | verifikasi pra-rilis | 1 VPS kecil, compose, data sintetis |
| `prod` | pelanggan | VPS 8c/16GB + backup offsite ([08 §7](08-keamanan.md)) |

---

Berikutnya: [10 — Roadmap & Estimasi](10-roadmap.md)
