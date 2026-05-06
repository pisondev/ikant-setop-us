# Jobdesc dan Progress Tim | Ikan't Setop Us

Dokumen ini adalah tracker owner dan progress. Detail checklist teknis frontend ada di `docs/frontend/frontend-mvp-todo.md`; jangan menduplikasi checklist panjang di sini.

Update terakhir: 2026-05-06.

## Ringkasan Status

| Area | Owner | Status | Catatan |
|---|---|---|---|
| Backend API dan database | Pison | Selesai MVP | Endpoint kontrak API sudah tersedia dan dites. |
| API contract dan Postman | Pison | Selesai | Ada di `docs/api`. |
| Frontend foundation | Mikail | Selesai tahap 1 | API helper, types, app shell, bottom nav, dan page header tersedia. |
| `/` entry page | Mikail | Selesai sementara | Template Next.js sudah diganti menjadi entry point aplikasi. |
| `/stocks` FIFO list | Farhan | Selesai implementasi | Sudah memakai API helper/type bersama, error state eksplisit, dan tanpa mock fallback. |
| Dashboard frontend | Farhan | Selesai | Route `/dashboard` sudah summary dan recent movements. |
| Master data frontend | Pannayaka | Selesai | Route `/fish-types` dan `/cold-storages` sudah GET/POST API. |
| Stock-in page | Mikail | Selesai | Route `/stocks/new` sudah GET master data dan POST `/stocks`. |
| Stock-out page | Mikail | Selesai | Route `/stock-outs/new` sudah preview FIFO dan POST `/stock-outs`. |
| Stock-out history | Adit | Selesai | Route `/stock-outs` sudah GET `/stock-outs`, filter dasar, dan detail batch FIFO. |
| README dan progress docs | Adit + Pison | Selesai sementara | README dan tracker sudah disinkronkan dengan frontend MVP saat ini. |
| Demo dan QA manual | Adit + Farhan | Sebagian | Demo flow API run 20260506154516 lolos; validasi browser/mobile masih perlu dijalankan. |

## Progress Teknis Aktual

### Backend

```txt
[x] Migration schema awal tersedia.
[x] Backend Go Fiber tersambung ke database.
[x] Endpoint /api/v1/health tersedia.
[x] Modul fish: GET/POST /fish-types.
[x] Modul storage: GET/POST /cold-storages.
[x] Modul stock: GET/POST /stocks, GET /stocks/fifo, GET /stocks/{id}.
[x] Modul stock: PATCH /stocks/{id}/quality dan /location.
[x] Modul stockout: GET/POST /stock-outs.
[x] Modul dashboard: GET /dashboard/summary dan /recent-movements.
[x] Logic FIFO stock-out transaksional tersedia.
[x] Unit test backend tersedia.
[x] go test ./... -count=1 lolos.
[x] go vet ./... lolos.
[x] go test -race ./... -count=1 lolos.
```

### Frontend

```txt
[x] Frontend foundation tahap 1 selesai.
[x] API helper tersedia.
[x] TypeScript types API tersedia.
[x] Mobile app shell tersedia.
[x] Bottom navigation tersedia.
[x] / tidak lagi memakai template Next.js.
[x] /dashboard tersedia.
[x] /fish-types tersedia.
[x] /cold-storages tersedia.
[x] /stocks/new tersedia.
[x] /stocks tersedia.
[x] /stocks terhubung ke API tanpa mock.
[x] /stock-outs/new tersedia.
[x] /stock-outs tersedia.
```

Keterangan:

```txt
[x] selesai
[~] sebagian
[ ] belum
```

## Pembagian Tugas

### Mikail - Frontend Lead dan Core Input Flow

Fokus:

- frontend foundation,
- layout mobile-first,
- API helper dan types,
- route `/`,
- route `/stocks/new`,
- route `/stock-outs/new`.

Checklist:

```txt
[x] Buat API helper.
[x] Buat type API bersama.
[x] Buat layout mobile-first.
[x] Buat bottom navigation.
[x] Ubah / dari template Next.js.
[x] Buat /stocks/new.
[x] Buat /stock-outs/new.
[x] Pastikan loading/error/success state ada di form input.
```

### Farhan - FIFO, Dashboard, dan QA Data Flow

Fokus:

- route `/stocks`,
- integrasi dashboard,
- validasi shape response API,
- skenario QA FIFO.

Checklist:

```txt
[x] Buat halaman awal /stocks.
[x] Integrasi awal GET /stocks/fifo.
[x] Integrasi awal GET /fish-types.
[x] Rapikan /stocks memakai API helper dan shared types.
[x] Tambahkan error state eksplisit di /stocks.
[x] Tambahkan tombol navigasi ke input stok dan input keluar.
[x] Buat /dashboard.
[x] Validasi angka dashboard setelah stock in/out via API demo run 20260506154516.
[~] Jalankan skenario FIFO manual dari docs/fifo-test-scenarios.md: flow utama sudah PASS via API, skenario lengkap belum semua dijalankan.
```

### Pannayaka - Master Data dan UX Copy

Fokus:

- route `/fish-types`,
- route `/cold-storages`,
- wording aplikasi,
- validasi flow persona Baso dan Daeng.

Checklist:

```txt
[x] Buat /fish-types.
[x] Integrasi GET /fish-types.
[x] Integrasi POST /fish-types.
[x] Buat /cold-storages.
[x] Integrasi GET /cold-storages.
[x] Integrasi POST /cold-storages.
[x] Review wording halaman master data.
[x] Review empty state halaman master data agar jelas untuk user lapangan.
```

### Adit - Dokumentasi, Demo Flow, QA Manual

Fokus:

- route `/stock-outs`,
- README,
- demo script,
- manual QA checklist,
- bahan presentasi.

Checklist:

```txt
[x] Buat /stock-outs.
[x] Integrasi GET /stock-outs.
[x] README disinkronkan dengan kondisi backend/frontend saat ini.
[x] Update README lagi setelah frontend lengkap.
[x] Siapkan demo script final.
[x] Siapkan QA checklist final.
[~] Catat bug/mismatch sebelum submission: belum ada bug aplikasi dari API demo; validasi browser/mobile pending.
```

### Pison - Backend Lead dan Integrator

Fokus:

- database,
- backend API,
- API contract,
- Postman,
- integrasi dan debugging backend/frontend.

Checklist:

```txt
[x] Database migration tersedia.
[x] Backend modular per domain tersedia.
[x] API contract terpenuhi.
[x] Postman collection tersedia.
[x] Postman environment tersedia.
[x] Unit test backend tersedia.
[x] Race test backend lolos.
[ ] Bantu frontend saat ada mismatch response atau CORS.
```

## Prioritas Berikutnya

Urutan kerja yang paling aman:

1. Adit dan Farhan menjalankan demo flow dan QA manual.
2. Validasi visual mobile lewat browser.
3. Catat bug/mismatch sebelum submission.

## Demo Flow Target

```txt
1. Buka dashboard.
2. Tambah jenis ikan Tuna jika belum ada.
3. Tambah Cold Storage A dan B jika belum ada.
4. Input Tuna 50 kg kualitas baik ke Cold Storage A pada jam 08.00.
5. Input Tuna 30 kg kualitas sedang ke Cold Storage B pada jam 09.00.
6. Buka /stocks dan filter Tuna.
7. Pastikan batch 50 kg menjadi FIFO rank #1.
8. Catat pengeluaran Tuna 60 kg.
9. Pastikan batch pertama habis.
10. Pastikan batch kedua tersisa 20 kg.
11. Buka dashboard.
12. Pastikan total stok dan recent movements berubah.
```

## Definition of Done Tim

```txt
[x] Backend Go Fiber menyala.
[x] Database schema tersedia.
[x] API contract tersedia.
[x] Postman collection dan environment tersedia.
[x] Frontend Next.js menyala dengan halaman MVP lengkap.
[x] /dashboard bisa menampilkan summary dan recent movement.
[x] /fish-types bisa menambah jenis ikan.
[x] /cold-storages bisa menambah cold storage.
[x] /stocks/new bisa mencatat stok masuk.
[x] /stocks bisa menampilkan stok FIFO dari backend tanpa mock.
[x] /stock-outs/new bisa mencatat ikan keluar.
[x] /stock-outs bisa menampilkan riwayat pengeluaran.
[~] Demo flow bisa dijalankan tanpa error fatal: API flow PASS, browser/mobile pending.
[x] README final sesuai fitur yang benar-benar selesai.
```

## Branch dan Commit

Branch integrasi yang disarankan:

```txt
main
dev
feat/mikail-core-frontend
feat/farhan-fifo-dashboard
feat/pannayaka-master-data
feat/adit-docs-demo
```

Sintaks checkout tahap berikutnya:

```bash
git switch -c feat/adit-docs-demo
```

Jika branch sudah ada:

```bash
git switch feat/adit-docs-demo
```

Commit convention:

```txt
feat(scope): message
fix(scope): message
docs(scope): message
test(scope): message
refactor(scope): message
chore(scope): message
```

Contoh:

```bash
feat(stocks): add fifo stock list
feat(stock-out): add stock out form
docs(readme): update setup guide
test(fifo): add manual fifo scenario
```
