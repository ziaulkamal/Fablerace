# 10 — Roadmap, Estimasi & Risiko

> Bagian dari [Kajian Produk Fablerace](../README.md) · Sebelumnya: [09 — DevOps](09-devops.md) · Lampiran: [A — API](lampiran-a-api.md) · [B — Konvensi Kode](lampiran-b-konvensi-kode.md)

## 1. Prinsip Prioritas

Sebuah tenant membayar $25/bulan hanya bila dalam 30 menit pertama ia bisa: membuat konten berjenjang → mengunggah media → menerbitkan → mengambil lewat API dengan API key → melihat sitemap/feed-nya hidup. Semua yang tidak melayani loop itu keluar dari MVP.

## 2. Fase

### MVP (target: 8–10 minggu, 2 engineer)

| Minggu | Fokus | Acuan kajian |
|---|---|---|
| 1–2 | Fondasi: scaffold (sudah ada di repo ini), auth JWT+refresh, tenant+keanggotaan, RBAC middleware | [03](03-arsitektur.md), [05](05-rbac.md) |
| 3–4 | Modul content penuh (hierarki, meta, revisi, status) + taxonomy | [04](04-desain-data.md) |
| 5 | Media: presign → confirm → process; kuota dasar | [07](07-storage-media.md) |
| 6 | Keyword + automasi inti: auto-slug, sitemap, RSS, scheduled publishing | [06](06-keyword-seo-automasi.md) |
| 7 | API key publik + rate limit + license enforcement (manual billing dulu: aktivasi lisensi oleh admin) | [08](08-keamanan.md) |
| 8 | Hardening: audit log, backup, staging, uji beban kepadatan tenant (validasi asumsi [02 §3](02-model-bisnis.md)) | [08](08-keamanan.md), [09](09-devops.md) |
| 9–10 | Private beta 10–20 tenant desain-partner (gratis/diskon seumur hidup) → iterasi | — |

**Sengaja TIDAK di MVP**: billing otomatis (aktivasi manual dulu — 20 pelanggan pertama tidak butuh self-serve), admin dashboard UI (API-first; beta memakai Postman/collection + frontend contoh), AI assist, Meilisearch, 2FA, role kustom.

### v1.0 — Komersial (bulan 4–5)

- Integrasi billing Paddle + Midtrans, siklus dunning otomatis ([02 §1](02-model-bisnis.md)).
- Admin dashboard web (frontend terpisah, konsumsi API yang sama).
- Webhook keluar + build hook; importer WordPress (WXR) sebagai senjata onboarding ([01 §6](01-produk-dan-pasar.md)).
- Custom domain + TLS otomatis per tenant (Caddy on-demand).

### v1.5 — Ekspansi (bulan 6–9)

- AI assist add-on ($10/bln), role kustom, 2FA, Meilisearch (bila terbukti perlu), tier Agency + manajemen multi-lisensi satu dasbor, lisensi self-hosted tahunan.

## 3. Estimasi Biaya Pra-Revenue

| Pos | Estimasi |
|---|---|
| 2 engineer × 3 bulan (MVP → v1.0 awal) | dominan — angka tergantung struktur tim (founder/karyawan) |
| Infra dev+staging+prod awal | <$150/bln |
| Tooling (registry, monitor, email transaksional) | <$100/bln |

Titik kritis finansial bukan infra (sangat murah berkat kepadatan Go) melainkan **waktu engineering** — alasan kuat menjaga scope MVP kejam.

## 4. Risiko & Mitigasi (teknis — risiko pasar di [01 §6](01-produk-dan-pasar.md))

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Asumsi kepadatan 150 tenant/node salah | Margin tergerus | Uji beban minggu-8 dengan profil trafik realistis; arsitektur sudah siap pisah node ([09 §5](09-devops.md)) |
| Scope creep fitur CMS ("tambah page builder dong") | MVP molor | Garis produk tegas: Fablerace = struktur konten + API + SEO; presentasi = urusan frontend tenant |
| Kebocoran lintas tenant | Fatal | RLS dua lapis + test otomatis isolasi per endpoint ([08 §3](08-keamanan.md)) — test isolasi adalah bagian definition-of-done setiap modul |
| Ketergantungan 1–2 engineer | Bus factor | Kajian ini + konvensi kode ([lampiran-b](lampiran-b-konvensi-kode.md)) + komentar `NEXT:`/`RELASI:` menjadikan codebase self-documenting |
| Biaya support membengkak di $25 | Margin tergerus | Dokumentasi publik + onboarding otomatis; target <0,2 jam support/tenant/bln dipantau sejak beta |

## 5. Metrik Keberhasilan per Fase

| Fase | Metrik gerbang untuk lanjut |
|---|---|
| Beta | ≥10 tenant aktif mingguan; time-to-first-publish <30 menit; 0 insiden isolasi |
| v1.0 | ≥50 lisensi berbayar dalam 90 hari komersial; churn bulanan <5% |
| v1.5 | ±315 lisensi (break-even, [02 §3](02-model-bisnis.md)); margin kotor terukur ≥70% |

---

Lampiran: [A — Sketsa API](lampiran-a-api.md) · [B — Konvensi Kode](lampiran-b-konvensi-kode.md)
