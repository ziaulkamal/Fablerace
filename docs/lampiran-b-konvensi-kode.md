# Lampiran B — Konvensi Modul & Clean Code

> Bagian dari [Kajian Produk Fablerace](../README.md). Konvensi ini **mengikat** untuk seluruh kontribusi kode; scaffold di repo ini adalah contoh hidupnya.

## 1. Prinsip

1. **Nama modul = nama folder.** Satu domain bisnis = satu folder di `internal/modules/`. Tidak ada file domain yang hidup di luar folder modulnya.
2. **Layer tetap, nama file tetap.** Setiap modul berisi file dengan peran yang sama persis — pembaca yang paham satu modul otomatis paham semuanya.
3. **Dependency satu arah.** `handler → service → repository`. Antar-modul hanya boleh bergantung pada **interface service** modul lain — tidak pernah ke repository modul lain.
4. **Komentar adalah penunjuk arah, bukan narasi.** Dua marker wajib (§4); selain itu komentar hanya untuk constraint yang tidak terlihat dari kode.

## 2. Anatomi Modul

```
internal/modules/content/          ← nama folder = nama modul
├── entity.go       # Struct domain murni (tanpa tag JSON API, tanpa logika DB)
├── dto.go          # Request/response API + validasi input (mass-assignment guard, 08 §6)
├── repository.go   # Interface Repository + implementasi pgx; SATU-SATUNYA tempat SQL modul ini
├── service.go      # Interface Service + implementasi; seluruh aturan bisnis & ownership check
├── handler.go      # HTTP: parse/validasi DTO → panggil service → tulis response envelope
└── routes.go       # Registrasi route modul + guard RBAC per route
```

Peran per layer:

| File | Boleh | Tidak boleh |
|---|---|---|
| `entity.go` | struct domain, konstanta status | import pgx/chi, logika |
| `dto.go` | struct req/res, fungsi `Validate()`, mapper entity↔dto | query DB |
| `repository.go` | SQL parametrik, scan → entity | aturan bisnis, HTTP |
| `service.go` | orkestrasi, transaksi, panggil service modul lain via interface | SQL mentah, baca http.Request |
| `handler.go` | decode, validasi, status code | SQL, aturan bisnis |
| `routes.go` | wiring path→handler + `RequirePermission(...)` | logika apa pun |

## 3. Peta Direktori Lengkap

```
cmd/api/main.go            # komposisi: config → infra → modul → router → serve
internal/
├── modules/               # domain (lihat §2): auth, tenant, user, content,
│                          # taxonomy, media, keyword, license
├── platform/              # infrastruktur lintas-modul — TIDAK berisi aturan bisnis
│   ├── config/            # env → struct typed
│   ├── database/          # pool pgx, migrasi embed, helper transaksi RLS
│   ├── cache/             # klien Redis
│   ├── storage/           # klien MinIO, konstruksi object key per tenant
│   ├── server/            # router chi, mounting routes, graceful shutdown
│   └── middleware/        # auth, tenant, rbac, ratelimit, license
└── shared/                # util MURNI (tanpa I/O): apperror, response,
                           # pagination, slugutil
migrations/                # SQL golang-migrate, embed ke binary
deploy/docker-compose.yml
```

Aturan import antar lapisan: `modules → platform → shared` (searah). `shared` tidak meng-import apa pun dari proyek; `platform` tidak meng-import `modules`; modul tidak saling import kecuali interface service (di-inject dari `main.go`).

## 4. Konvensi Komentar `NEXT:` dan `RELASI:`

Dua marker ini adalah **kontrak tim** — permintaan eksplisit desain proyek ini: setiap titik yang menyentuh modul lain atau menunda keputusan harus meninggalkan jejak yang bisa di-grep.

```go
// RELASI: internal/modules/taxonomy — term_relationships menautkan konten ini
// ke term; penghapusan konten harus membersihkan relasi (lihat service.Delete).
//
// NEXT: tambah tabel slug_redirects + redirect 301 saat slug konten published
// berubah (docs/06-keyword-seo-automasi.md §2).
```

| Marker | Artinya | Kapan wajib |
|---|---|---|
| `RELASI:` | fungsi/file ini berkorelasi dengan fungsi di modul lain — sebut path modul + apa keterkaitannya | setiap kali service/handler menyentuh data atau interface modul lain |
| `NEXT:` | arah pengembangan berikutnya yang disengaja ditunda — sebut apa & rujuk dokumen kajian | setiap keputusan "cukup untuk sekarang" |

Cara pakai harian:

```bash
grep -rn "NEXT:"   internal/   # backlog teknis yang hidup di kode
grep -rn "RELASI:" internal/   # peta keterkaitan antar modul
```

`TODO` tanpa marker dilarang masuk `main` — pakai `NEXT:` agar seragam dan ter-grep.

## 5. Gaya Kode

- `gofmt` + `go vet` wajib bersih; `golangci-lint` di CI ([09 §3](09-devops.md)).
- Error dibungkus konteks: `fmt.Errorf("content: find %s: %w", id, err)` — prefix nama modul.
- Error domain memakai `shared/apperror` (kode stabil, mis. `content_not_found`) → handler memetakan ke HTTP + envelope ([03 §4](03-arsitektur.md)) tanpa membocorkan detail internal.
- Context selalu parameter pertama; `tenant_id`/`user_id` dibaca hanya lewat helper context middleware — tidak pernah dari header mentah di dalam modul.
- Tanpa variabel global kecuali di `main.go` (komposisi). Dependensi di-inject lewat konstruktor `NewService(repo, ...)`.
