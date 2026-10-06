# Lampiran A — Sketsa Endpoint API

> Bagian dari [Kajian Produk Fablerace](../README.md). Konvensi umum (envelope, pagination, auth) di [03 §4](03-arsitektur.md). Status implementasi scaffold: ✅ fungsional · 🩻 kerangka (stub 501).

Semua endpoint di bawah `/api/v1`, butuh `X-Tenant-ID` + Bearer JWT kecuali ditandai *(publik)* atau *(API key)*.

## auth (🩻)

| Method | Path | Guard | Deskripsi |
|---|---|---|---|
| POST | `/auth/register` | — | Daftar user + tenant pertama (mulai trial) |
| POST | `/auth/login` | rate-limit ketat | Email+password → access + refresh token |
| POST | `/auth/refresh` | — | Rotasi refresh token ([08 §2](08-keamanan.md)) |
| POST | `/auth/logout` | JWT | Revoke refresh token |

## tenant (🩻)

| Method | Path | Guard |
|---|---|---|
| GET | `/tenants` | JWT — daftar tenant milik user |
| POST | `/tenants` | JWT — buat tenant baru (lisensi trial) |
| GET/PATCH | `/tenants/current` | `settings.manage` untuk PATCH |
| GET | `/tenants/current/members` | `user.manage` |
| POST | `/tenants/current/invites` | `user.invite` |

## content (✅ exemplar)

| Method | Path | Guard |
|---|---|---|
| GET | `/contents` | `content.view` — filter `?type=&status=&parent_id=&search=`, pagination |
| POST | `/contents` | `content.create` — slug otomatis bila kosong |
| GET | `/contents/{id}` | `content.view` — termasuk meta & SEO |
| PATCH | `/contents/{id}` | `content.edit_any` / `edit_own` (ownership check di service) |
| DELETE | `/contents/{id}` | `content.delete` — soft delete (trash) |
| POST | `/contents/{id}/publish` | `content.publish` (+ownership utk Author) |
| GET | `/contents/{id}/children` | `content.view` — turunan langsung (hierarki) |
| PUT | `/contents/{id}/meta` | `content.edit_*` — upsert meta EAV |

## taxonomy (🩻)

| Method | Path | Guard |
|---|---|---|
| GET/POST | `/taxonomies` | view: `content.view` · tulis: `taxonomy.manage` |
| GET/POST | `/taxonomies/{id}/terms` | idem — term hierarkis via `parent_id` |
| PUT | `/contents/{id}/terms` | `content.edit_*` — set relasi konten↔term |

## media (🩻)

| Method | Path | Guard |
|---|---|---|
| POST | `/media/presign` | `media.upload` — alur di [07 §2](07-storage-media.md) |
| POST | `/media/{id}/confirm` | `media.upload` |
| GET | `/media` | `content.view` |
| DELETE | `/media/{id}` | `media.manage` |

## keyword (🩻)

| Method | Path | Guard |
|---|---|---|
| GET/POST | `/keywords` | `keyword.manage` (GET: `content.view`) |
| PUT | `/contents/{id}/keywords` | `content.edit_*` — tautkan keyword, set primary |
| GET | `/keywords/{id}/coverage` | `content.view` — konten yang membidik keyword ini |

## license (🩻)

| Method | Path | Guard |
|---|---|---|
| GET | `/license` | `billing.view` — status lisensi & pemakaian kuota |
| POST | `/license/activate` | `billing.manage` — aktivasi license key |
| POST | `/billing/webhook` | signature gateway *(publik, verified)* |

## Delivery publik per tenant *(API key / publik)*

| Method | Path | Catatan |
|---|---|---|
| GET | `/public/contents?slug=` | *(API key)* konten published + blok `seo` siap render |
| GET | `/sitemap.xml`, `/robots.txt`, `/feed.xml` | *(publik, resolve via domain/subdomain)* — [06 §2](06-keyword-seo-automasi.md) |

## Sistem

| Method | Path | Catatan |
|---|---|---|
| GET | `/healthz` | liveness *(publik)* |
| GET | `/readyz` | cek DB/Redis/MinIO *(publik)* |
