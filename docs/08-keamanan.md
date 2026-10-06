# 08 — Keamanan

> Bagian dari [Kajian Produk Fablerace](../README.md) · Sebelumnya: [07 — Storage & Media](07-storage-media.md) · Lanjut: [09 — DevOps](09-devops.md)

## 1. Model Ancaman Utama

| Ancaman | Prioritas | Pertahanan |
|---|---|---|
| Kebocoran data lintas tenant | **Kritis** — satu insiden membunuh produk SaaS | Dua lapis: filter `tenant_id` di repository + RLS Postgres (§3) |
| Pengambilalihan akun | Tinggi | Hash Argon2id, rotasi refresh token, rate limit login (§2) |
| Penyalahgunaan API key bocor | Tinggi | Key di-hash, scope minimal, expiry, revoke instan (§4) |
| Upload berbahaya | Sedang | Magic-byte check, whitelist MIME, presigned terbatas ([07 §3](07-storage-media.md)) |
| Abuse/DoS tenant tunggal | Sedang | Rate limit per tenant & per key (§5) |

## 2. Autentikasi

- **Password**: Argon2id (parameter OWASP terkini), tidak pernah log/plaintext.
- **Sesi dashboard**: JWT access token pendek (15 menit, HMAC-SHA256, secret dari env) + **refresh token** opaque di tabel `refresh_tokens` (hash SHA-256, umur 30 hari, **rotasi setiap pakai** — pemakaian ulang token lama = indikasi pencurian → revoke seluruh keluarga token user tersebut).
- Klaim JWT minimal: `sub` (user_id), `exp`, `iat` — **tanpa** role/permission di dalam token (selalu dicek segar dari DB/cache per request, sehingga pencabutan akses berlaku seketika; trade-off satu query ter-cache 60 detik, [05 §3](05-rbac.md)).
- Login rate-limited: 5 percobaan gagal/menit per email+IP, backoff eksponensial.
- 2FA TOTP: pasca-MVP (`NEXT:` di modul `auth`).

## 3. Isolasi Tenant (dua lapis)

**Lapis 1 — aplikasi**: semua query repository wajib menyertakan `tenant_id` dari context (bukan dari payload klien); konvensi ditegakkan review + helper `database.WithTenant`.

**Lapis 2 — PostgreSQL RLS** ([04 §4](04-desain-data.md)): policy `tenant_id = current_setting('app.tenant_id')::uuid` pada semua tabel ber-tenant; koneksi aplikasi memakai role **non-BYPASSRLS**. Bug aplikasi yang lolos lapis 1 menghasilkan **0 baris**, bukan kebocoran.

Aturan pendukung:
- `tenant_id` di context diisi **hanya** oleh middleware tenant setelah memverifikasi keanggotaan user / kepemilikan API key — tidak pernah dipercaya dari header mentah tanpa verifikasi.
- Kunci cache Redis selalu berprefix `t:{tenant_id}:` ([03 §6](03-arsitektur.md)).
- Object key MinIO selalu dikonstruksi server-side ([07 §1](07-storage-media.md)).

## 4. API Key

- Format: `fbr_live_<32 byte random base62>`; hanya ditampilkan sekali saat dibuat; disimpan **SHA-256 hash** (`api_keys.key_hash`).
- Scope eksplisit (subset permission, default read-only), `expires_at` opsional, revoke = soft delete seketika (cache permission key ikut di-invalidate).
- `last_used_at` diperbarui async (tanpa menambah latensi request) untuk audit key tak terpakai.

## 5. Rate Limiting & Kuota

- Token bucket di Redis: per API key dan per tenant (60/300 req-menit sesuai tier, [02 §1](02-model-bisnis.md)); respons `429` + header `Retry-After`.
- Endpoint sensitif (login, presign upload, webhook test) punya limiter terpisah yang lebih ketat.
- **License enforcement middleware**: status lisensi non-`active` menurunkan kemampuan (mis. `suspended` → API read-only) sesuai siklus di [02 §1](02-model-bisnis.md); dicek dari cache dengan TTL 60 detik.

## 6. Checklist OWASP API Top-10 (dipetakan)

| Risiko | Mitigasi di Fablerace |
|---|---|
| BOLA (akses objek orang lain) | RLS + ownership check service (edit_own vs edit_any) |
| Broken authentication | §2 (Argon2id, rotasi refresh, rate limit) |
| Property-level authorization | DTO eksplisit per endpoint — tidak pernah bind/echo struct DB mentah |
| Unrestricted resource consumption | §5 + batas ukuran upload + pagination maks 100 |
| SSRF (URL webhook) | Validasi URL webhook: tolak IP privat/link-local, resolve-then-connect check |
| Injection | Query parametrik pgx 100% — tanpa string concat SQL |
| Security misconfiguration | Image container non-root, secrets hanya via env, TLS wajib di proxy ([09](09-devops.md)) |
| Mass assignment | DTO + validator eksplisit (`internal/shared/validator`) |

## 7. Audit, Backup & Tanggap Insiden

- `audit_logs` ([04 §2.6](04-desain-data.md)) mencatat mutasi penting: siapa, apa, kapan, dari IP mana, diff `changes` — dapat diekspor per tenant (nilai jual ke agensi).
- Backup: `pg_dump` harian + WAL archiving (point-in-time recovery), mirror MinIO ([07 §5](07-storage-media.md)); **uji restore bulanan** dijadwalkan — backup yang tak pernah diuji dianggap tidak ada.
- Data pribadi minimal (email, nama) → eksposur GDPR kecil; endpoint ekspor & hapus akun disiapkan pasca-MVP (`NEXT:` modul `user`).

---

Berikutnya: [09 — DevOps & Deployment](09-devops.md)
