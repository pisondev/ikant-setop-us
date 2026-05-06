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
| `/stocks` FIFO list | Farhan | Sebagian lanjut | Sudah memakai API helper/type bersama, error state eksplisit, dan tanpa mock fallback. |
| Dashboard frontend | Farhan | Belum | Route `/dashboard` belum ada. |
| Master data frontend | Pannayaka | Belum | Route `/fish-types` dan `/cold-storages` belum ada. |
| Stock-in page | Mikail | Belum | Route `/stocks/new` belum ada. |
| Stock-out page | Mikail | Belum | Route `/stock-outs/new` belum ada. |
| Stock-out history | Adit | Belum | Route `/stock-outs` belum ada. |
| README dan progress docs | Adit + Pison | Sebagian | README sudah disinkronkan, perlu update lagi saat frontend bertambah. |
| Demo dan QA manual | Adit + Farhan | Belum | Menunggu halaman frontend inti selesai. |

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
[ ] /dashboard tersedia.
[ ] /fish-types tersedia.
[ ] /cold-storages tersedia.
[ ] /stocks/new tersedia.
[x] /stocks tersedia.
[~] /stocks terhubung ke API tanpa mock, tetapi halaman tujuan tombol belum tersedia.
[ ] /stock-outs/new tersedia.
[ ] /stock-outs tersedia.
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
[ ] Buat /stocks/new.
[ ] Buat /stock-outs/new.
[ ] Pastikan loading/error/success state ada di form input.
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
[ ] Buat /dashboard.
[ ] Validasi angka dashboard setelah stock in/out.
[ ] Jalankan skenario FIFO manual dari docs/fifo-test-scenarios.md.
```

### Pannayaka - Master Data dan UX Copy

Fokus:

- route `/fish-types`,
- route `/cold-storages`,
- wording aplikasi,
- validasi flow persona Baso dan Daeng.

Checklist:

```txt
[ ] Buat /fish-types.
[ ] Integrasi GET /fish-types.
[ ] Integrasi POST /fish-types.
[ ] Buat /cold-storages.
[ ] Integrasi GET /cold-storages.
[ ] Integrasi POST /cold-storages.
[ ] Review wording halaman stok dan pengeluaran.
[ ] Review empty state agar jelas untuk user lapangan.
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
[ ] Buat /stock-outs.
[ ] Integrasi GET /stock-outs.
[x] README disinkronkan dengan kondisi backend/frontend saat ini.
[ ] Update README lagi setelah frontend lengkap.
[ ] Siapkan demo script final.
[ ] Siapkan QA checklist final.
[ ] Catat bug/mismatch sebelum submission.
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

1. Pannayaka membuat `/fish-types` dan `/cold-storages`.
2. Mikail membuat `/stocks/new`.
3. Mikail membuat `/stock-outs/new`.
4. Farhan membuat `/dashboard`.
5. Adit membuat `/stock-outs`.
6. Adit dan Farhan menjalankan demo flow dan QA manual.

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
[ ] Frontend Next.js menyala dengan halaman MVP lengkap.
[ ] /dashboard bisa menampilkan summary dan recent movement.
[ ] /fish-types bisa menambah jenis ikan.
[ ] /cold-storages bisa menambah cold storage.
[ ] /stocks/new bisa mencatat stok masuk.
[ ] /stocks bisa menampilkan stok FIFO dari backend tanpa mock.
[ ] /stock-outs/new bisa mencatat ikan keluar.
[ ] /stock-outs bisa menampilkan riwayat pengeluaran.
[ ] Demo flow bisa dijalankan tanpa error fatal.
[ ] README final sesuai fitur yang benar-benar selesai.
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
