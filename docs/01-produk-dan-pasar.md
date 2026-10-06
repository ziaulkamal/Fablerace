# 01 — Produk & Analisis Pasar

> Bagian dari [Kajian Produk Fablerace](../README.md) · Lanjut: [02 — Model Bisnis](02-model-bisnis.md)

## 1. Visi Produk

**Fablerace** adalah platform manajemen konten (CMS) multi-tenant berbasis API — "WordPress-as-a-Service yang ditulis ulang untuk era modern". Satu instalasi platform melayani banyak pelanggan (tenant); setiap tenant mendapat ruang konten sendiri dengan model data berjenjang ala WordPress (konten hierarkis, taxonomy, meta fleksibel), kontrol akses multi-role, media library berbasis object storage, dan automasi SEO/keyword bawaan.

**Pernyataan posisi:**

> Untuk tim konten dan agensi yang lelah mengurus hosting, plugin, dan patch keamanan WordPress, Fablerace adalah CMS berlangganan $25/bulan yang memberi struktur konten sefleksibel WordPress dengan performa Go, API modern, dan automasi SEO bawaan — tanpa server yang harus dirawat sendiri.

## 2. Masalah yang Diselesaikan

| Masalah di lapangan | Dampak | Jawaban Fablerace |
|---|---|---|
| WordPress self-hosted butuh perawatan terus-menerus (update core/plugin, keamanan, backup) | Biaya tersembunyi jam kerja; risiko peretasan | Fully-managed, satu binary Go, attack surface kecil |
| Plugin SEO (Yoast dkk.) berbayar terpisah dan memberatkan situs | Biaya tambahan, performa turun | Keyword & automasi SEO bawaan platform (sitemap, meta, slug, RSS) |
| Headless CMS modern (Contentful, Storyblok) mahal saat tim membesar | Harga per-seat melonjak ($300+/bulan untuk tim kecil) | Harga flat per lisensi, role tak dibatasi jumlah seat secara agresif |
| Struktur konten headless CMS sering terlalu "datar" | Sulit memodelkan hierarki halaman/kategori bertingkat | Model berjenjang native: konten ber-`parent_id`, term bertingkat, meta EAV |
| Media storage lokal tidak tahan skala | Disk penuh, backup sulit | Object storage (MinIO/S3) sejak hari pertama, kuota per lisensi |

## 3. Target Pasar

**Segmen primer**

1. **Agensi digital & web studio** — mengelola 5–50 situs klien; butuh satu dasbor, role per klien, dan API untuk frontend apa pun. Satu lisensi per situs klien → model reseller alami.
2. **Media kecil / publisher niche** — blog jaringan, media lokal; butuh taxonomy dalam, penjadwalan terbit, dan SEO otomatis.
3. **Tim konten startup/UKM** — butuh CMS untuk situs marketing + blog yang dikonsumsi frontend Next.js/Astro via API.

**Segmen sekunder**: developer indie yang menjual jasa pembuatan situs (white-label), institusi yang butuh CMS on-premise (lisensi self-hosted — peluang revenue kedua, lihat [02 — Model Bisnis](02-model-bisnis.md)).

**Pasar geografis awal**: Indonesia/Asia Tenggara (kompetisi managed-CMS lokal masih tipis, sensitivitas harga cocok dengan titik $25) dengan arsitektur yang sejak awal siap pasar global (API English-first, i18n konten).

## 4. Analisis Kompetitor

| Kompetitor | Harga (per situs/bulan, 2025) | Kekuatan | Celah yang dieksploitasi Fablerace |
|---|---|---|---|
| **WordPress.com** (Business) | ~$25–40 | Ekosistem terbesar, familiar | Berat, teknologi lama, API bukan warga kelas satu, plugin = risiko |
| **Ghost(Pro)** | mulai ~$9–25 (Creator/Team) | Fokus publishing & newsletter, cepat | Model konten datar (tanpa hierarki/custom taxonomy dalam), kurang cocok untuk situs terstruktur & agensi multi-klien |
| **Webflow** (CMS plan) | ~$23–29 | Visual builder kuat | Bukan API-first murni, lock-in desain, item CMS dibatasi kuota ketat |
| **Contentful** | Free → melonjak ~$300 (Team) | Enterprise-grade, ekosistem SDK | Lompatan harga brutal; tim kecil terjepit di antara free dan $300 |
| **Storyblok** | Free → ~$99+ (Entry/Teams) | Visual editor + headless | Sama: gap harga untuk tim kecil; per-seat pricing |
| **Sanity** | Free → ~$15/seat (Growth) | Fleksibilitas skema tinggi | Per-seat; query GROQ = kurva belajar; hierarki harus dirakit sendiri |
| **WordPress self-hosted** | "gratis" + hosting $10–30 | Kontrol penuh | Total cost of ownership tinggi (maintenance, keamanan, performa) |

> Angka harga adalah harga publik yang umum diketahui per 2025 dan bisa berubah; fungsinya di kajian ini sebagai pembanding posisi, bukan acuan kontraktual.

**Kesimpulan positioning**: ada celah nyata di titik **$25 flat/situs** untuk produk yang (a) API-first, (b) struktur konten sedalam WordPress, (c) SEO otomatis bawaan, dan (d) tidak menghukum pertumbuhan tim dengan harga per-seat. Kompetitor besar bermain di bawah $10 (fitur dipangkas) atau melompat ke $99–300 (enterprise).

## 5. Diferensiasi Inti

1. **Model data berjenjang ala WordPress, di atas stack modern** — `contents.parent_id`, taxonomy & term bertingkat, meta EAV; mental model yang sudah dipahami jutaan praktisi, tanpa beban legacy PHP. Detail di [04 — Desain Data](04-desain-data.md).
2. **SEO & keyword automation sebagai fitur platform, bukan plugin** — sitemap, robots, RSS, auto-meta, manajemen keyword per konten. Detail di [06 — Keyword & Automasi](06-keyword-seo-automasi.md).
3. **Ekonomi Go**: satu binary, konsumsi memori kecil → kepadatan tenant per server tinggi → margin sehat di harga $25 (dihitung di [02 — Model Bisnis](02-model-bisnis.md)).
4. **Multi-role berjenjang per tenant** yang cocok untuk alur kerja agensi (owner → admin → editor → author → contributor → viewer). Detail di [05 — RBAC](05-rbac.md).

## 6. Risiko Pasar

| Risiko | Kemungkinan | Mitigasi |
|---|---|---|
| "Kenapa tidak WordPress saja?" — inersia pasar | Tinggi | Fokus ke buyer yang sudah sakit kepala dengan maintenance (agensi); importer konten WordPress (WXR) sebagai fitur onboarding |
| Kompetitor menurunkan harga tier kecil | Sedang | Diferensiasi fitur (hierarki + SEO otomatis) lebih sulit ditiru daripada harga |
| Churn tinggi di segmen UKM | Sedang | Model agensi (1 lisensi = 1 situs klien) membuat biaya berpindah tinggi secara alami |
| Biaya akuisisi lebih besar dari proyeksi | Sedang | Mulai dari kanal murah: komunitas developer, content-led growth, program reseller agensi |

---

Berikutnya: [02 — Model Bisnis & Unit Economics](02-model-bisnis.md)
