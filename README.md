# Fablerace

**CMS multi-tenant berbasis API — "WordPress-as-a-Service" yang ditulis ulang dengan Go.**

Satu platform melayani banyak tenant (situs/workspace). Setiap tenant mendapat struktur konten berjenjang ala WordPress (konten ber-parent, taxonomy & term bertingkat, meta fleksibel), kontrol akses multi-role, media library di object storage (MinIO), dan automasi keyword/SEO bawaan — dijual berlangganan **$25/bulan/lisensi** (1 lisensi = 1 tenant).

| | |
|---|---|
| Backend | Go 1.24+ · chi · pgx |
| Data | PostgreSQL 16 (RLS per tenant) · Redis 7 |
| Storage | MinIO (S3-compatible, prefix per tenant) |
| Runtime | Docker / docker-compose, single binary |

## 📚 Dokumen Kajian

Kajian utuh produk & teknis, dibaca berurutan:

| # | Dokumen | Isi singkat |
|---|---|---|
| 01 | [Produk & Pasar](docs/01-produk-dan-pasar.md) | Visi, target pasar, kompetitor, positioning $25 |
| 02 | [Model Bisnis](docs/02-model-bisnis.md) | Tier & lisensi, unit economics (margin ±77%), break-even, metering |
| 03 | [Arsitektur](docs/03-arsitektur.md) | Stack & alasan, diagram, multi-tenancy (shared DB + RLS), API, jobs |
| 04 | [Desain Data & ERD](docs/04-desain-data.md) | Skema berjenjang ala WordPress: contents, meta EAV, taxonomy, dst. |
| 05 | [RBAC & Multi-Role](docs/05-rbac.md) | 7 role berjenjang + permission matrix + strategi implementasi |
| 06 | [Keyword & Automasi SEO](docs/06-keyword-seo-automasi.md) | Keyword coverage, auto-slug/meta/sitemap/RSS, scheduling, webhooks, FTS, AI |
| 07 | [Storage & Media](docs/07-storage-media.md) | Topologi MinIO, presigned upload, kuota, image processing |
| 08 | [Keamanan](docs/08-keamanan.md) | AuthN, isolasi tenant dua lapis, rate limit, OWASP, audit |
| 09 | [DevOps](docs/09-devops.md) | Docker, compose, CI/CD, observability, jalur scaling |
| 10 | [Roadmap](docs/10-roadmap.md) | MVP 8–10 minggu → v1.0 → v1.5, risiko, metrik gerbang |
| A | [Lampiran: Sketsa API](docs/lampiran-a-api.md) | Endpoint per modul |
| B | [Lampiran: Konvensi Kode](docs/lampiran-b-konvensi-kode.md) | Struktur modul, layer, aturan komentar `NEXT:`/`RELASI:` |

Log pengerjaan: [TASK.md](TASK.md)

## 🏗️ Scaffold Kode

Repo ini sekaligus berisi kerangka proyek Go yang mengikuti [konvensi modul](docs/lampiran-b-konvensi-kode.md): satu folder per modul domain (`internal/modules/{auth,tenant,user,content,taxonomy,media,keyword,license}`), tiap modul berlapis `entity / dto / repository / service / handler / routes`. Modul **`content`** diimplementasi fungsional penuh sebagai exemplar; modul lain berupa kerangka ber-interface dengan komentar `NEXT:` (arah lanjutan) dan `RELASI:` (keterkaitan antar modul).

### Menjalankan (dev)

```bash
cp .env.example .env            # sesuaikan secret
docker compose -f deploy/docker-compose.yml up --build
# API      → http://localhost:8080/healthz
# MinIO UI → http://localhost:9001
```

Migrasi database dijalankan otomatis saat aplikasi boot. Tanpa Docker:

```bash
go build ./... && go vet ./...  # verifikasi
go run ./cmd/api                # butuh Postgres/Redis/MinIO lokal (lihat .env.example)
```

### Peta direktori

```
cmd/api/            entrypoint
internal/modules/   modul domain (1 folder = 1 modul)
internal/platform/  infra lintas-modul: config, database, cache, storage, server, middleware
internal/shared/    util murni: apperror, response, pagination, slugutil
migrations/         SQL (embed, auto-run saat boot)
deploy/             docker-compose
```

### Backlog teknis yang hidup di kode

```bash
grep -rn "NEXT:"   internal/   # arah pengembangan berikutnya
grep -rn "RELASI:" internal/   # peta korelasi antar modul
```

## Lisensi

Proprietary — hak cipta pemilik repo. (Keputusan lisensi final dibahas di [10 — Roadmap](docs/10-roadmap.md).)
