# 06 — Keyword, SEO & Stack Automasi

> Bagian dari [Kajian Produk Fablerace](../README.md) · Sebelumnya: [05 — RBAC](05-rbac.md) · Lanjut: [07 — Storage & Media](07-storage-media.md)

Automasi SEO adalah diferensiasi utama ([01 §5](01-produk-dan-pasar.md)): hal-hal yang di WordPress membutuhkan 3–4 plugin menjadi perilaku bawaan platform.

## 1. Manajemen Keyword

Model data: `keywords` + `content_keywords` dengan flag `is_primary` ([04 §2.5](04-desain-data.md)).

- Setiap konten dapat ditautkan ke banyak keyword; **satu** keyword utama (`is_primary`) per konten — dipakai untuk analisis kepadatan & saran internal link.
- Keyword unik per tenant (`UNIQUE(tenant_id, phrase)`) → satu keyword bisa dipantau lintas banyak konten: dasbor "keyword coverage" (konten apa saja yang membidik frasa ini) — mencegah **keyword cannibalization** (dua konten bersaing untuk frasa yang sama), masalah nyata yang Yoast-pun tidak tangani gratis.
- Analisis on-save (service `keyword`): keberadaan keyword di judul, slug, meta description, heading pertama, dan kepadatan di body → skor sederhana 0–100 dikembalikan di respons API, frontend bebas menampilkan.

## 2. Automasi SEO Bawaan

| Automasi | Mekanisme | Letak di scaffold |
|---|---|---|
| **Auto-slug** | Judul → slug ternormalisasi (transliterasi dasar, lowercase, dash), unik per (tenant, type); bentrok → sufiks `-2` | `internal/shared/slugutil` + dipakai `content/service.go` |
| **Auto meta description** | Bila `content_seo.meta_description` kosong: potong cerdas dari excerpt/body (batas 155 char, potong di batas kata) | service `keyword` (`NEXT:` di scaffold) |
| **sitemap.xml** | Dibangun ulang via job `seo:rebuild_sitemap` setiap publish/unpublish; hasil di-cache Redis + disalin ke MinIO; endpoint publik `GET /sitemap.xml` per tenant/domain | job worker ([03 §5](03-arsitektur.md)) |
| **robots.txt** | Dihasilkan dari setting tenant + flag `noindex` per konten | endpoint publik |
| **RSS/Atom feed** | `GET /feed.xml` — konten published terbaru per tenant | endpoint publik |
| **Canonical & OG tags** | `content_seo` menyimpan canonical_url + og_image; API publik mengembalikan blok `seo` siap-render untuk frontend | modul `content` + `keyword` |
| **Redirect slug** | Mengubah slug konten published otomatis mencatat redirect 301 (tabel `slug_redirects`) | pasca-MVP, `NEXT:` di `content/service.go` |

## 3. Automasi Alur Kerja

| Automasi | Mekanisme |
|---|---|
| **Scheduled publishing** | Konten status `scheduled` + `scheduled_at`; job cron per menit menerbitkan yang jatuh tempo ([03 §5](03-arsitektur.md)) |
| **Webhooks keluar** | Event domain (`content.published`, `content.updated`, `media.uploaded`, …) → `webhook:deliver` dengan HMAC signature & retry eksponensial; log di `webhook_deliveries` — inilah jembatan "stack automasi lain": Zapier/n8n/Make tinggal menangkap webhook |
| **Build hook frontend** | Kasus khusus webhook: URL deploy-hook Vercel/Netlify → situs statis tenant rebuild otomatis saat konten terbit |
| **Auto-revisi** | Setiap update konten menyimpan snapshot ke `revisions` (maks N revisi, default 20, konfigurasi per tenant) |
| **Auto audit** | Semua mutasi penting tercatat `audit_logs` via helper service |

## 4. Pencarian (Full-Text Search)

**Fase 1 (MVP): PostgreSQL FTS.** Kolom generated `search_vector` (judul berbobot A, excerpt B, body C) + index GIN; bahasa config per tenant (`simple` default — aman untuk Indonesia, `english` opsional). Nol infrastruktur tambahan.

**Fase 2 (bila kebutuhan terbukti): Meilisearch** — typo-tolerance & ranking lebih baik; disinkron via event yang sama dengan webhook. Keputusan ini ditandai `NEXT:` di `content/repository.go`. Jangan mulai dengan Meilisearch: tambahan satu stateful service = biaya operasional permanen yang belum dibayar kebutuhan.

## 5. AI Assist (Add-on $10/bln — [02 §1](02-model-bisnis.md))

Fitur tier atas, memanggil LLM API (Claude) dari service `keyword`:

- Saran keyword dari draft konten.
- Generate meta title/description (user tetap konfirmasi — bukan auto-apply).
- Saran internal link ke konten lama yang relevan.

Arsitektur: panggilan LLM selalu lewat job queue (bukan request path), hasil di-cache per revisi konten agar biaya token terkendali; batas pemakaian bulanan per tenant dihitung di metering. Margin add-on tetap sehat karena biaya inferensi per artikel (ratusan token) jauh di bawah $10/bulan pada volume wajar — tetap dipantau lewat metering yang sama.

---

Berikutnya: [07 — Storage & Media (MinIO)](07-storage-media.md)
