# Pembagian Tugas Tim — Ikan't Setop Us

Dokumen ini menjadi acuan pembagian tugas tim **Ikan't Setop Us** untuk pengerjaan MVP aplikasi **FishFlow: FIFO Fish Inventory Management System**.

Pembagian tugas ini disusun berdasarkan rencana frontend MVP yang sudah disepakati, yaitu aplikasi web **Next.js mobile-first** yang membuktikan alur utama:

```txt
Input stok ikan masuk
→ data tersimpan di database
→ stok tampil berdasarkan FIFO
→ stok keluar dapat dicatat
→ dashboard ikut berubah
```

---

## 1. Prinsip Pembagian Tugas

Pembagian tugas dibuat agar setiap anggota memiliki tanggung jawab yang jelas, bisa dikerjakan paralel, dan tetap mengarah ke MVP yang sama.

### 1.1 Fokus Utama MVP

MVP tidak mengejar semua fitur besar. Fokus utama adalah memastikan fitur inti berjalan end-to-end:

1. User dapat membuat master jenis ikan.
2. User dapat membuat master cold storage.
3. User dapat mencatat stok ikan masuk.
4. User dapat melihat stok berdasarkan FIFO.
5. User dapat mencatat ikan keluar.
6. Sistem dapat mengurangi stok berdasarkan FIFO.
7. Dashboard dapat menampilkan perubahan stok.
8. Riwayat pengeluaran dan recent movement dapat ditampilkan.

### 1.2 Batasan MVP

Fitur berikut tidak menjadi prioritas MVP:

- login/register,
- role-based access control,
- payment,
- marketplace,
- integrasi alat timbang,
- sensor cold storage,
- notifikasi realtime,
- grafik kompleks,
- halaman detail stok yang kompleks.

### 1.3 Peran Pison

Pison berperan sebagai **backend lead dan system integrator**.

Tanggung jawab utama Pison:

- merancang database PostgreSQL,
- membuat migration,
- membuat backend API menggunakan Go Fiber,
- memastikan endpoint sesuai API contract,
- memastikan logic FIFO berjalan,
- membantu integrasi frontend-backend,
- review pull request yang menyentuh API contract,
- membantu debugging integrasi jika ada mismatch data.

Karena Pison fokus pada backend, anggota lain difokuskan pada frontend, testing, dokumentasi, dan demo.

---

# 2. Pembagian Tugas per Anggota

---

## 2.1 Mikail — Frontend Lead & Core UI Implementer

### Peran Utama

Mikail bertanggung jawab sebagai **Frontend Lead** yang mengerjakan fondasi frontend dan halaman input utama.

Mikail memegang bagian yang paling banyak berhubungan dengan UI, layout, dan interaksi user. Karena fitur utama aplikasi harus nyaman digunakan oleh pekerja lapangan seperti Baso, Mikail perlu memastikan tampilan mobile-first, tombol mudah ditekan, dan form tidak terlalu panjang.

---

### A. Frontend Foundation

Mikail mengerjakan fondasi awal frontend:

```txt
src/lib/api.ts
src/types/api.ts
src/components/layout/app-shell.tsx
src/components/layout/bottom-nav.tsx
src/components/layout/page-header.tsx
```

Detail tugas:

- setup penggunaan `NEXT_PUBLIC_API_BASE_URL`,
- membuat helper API:
  - `apiGet(path)`,
  - `apiPost(path, body)`,
  - `apiPatch(path, body)`,
- membuat types awal berdasarkan API contract,
- membuat layout utama mobile-first,
- membuat bottom navigation,
- membuat page header reusable,
- membuat halaman `/` yang redirect atau mengarahkan user ke `/dashboard`.

Output yang diharapkan:

```txt
[ ] API helper tersedia
[ ] TypeScript types awal tersedia
[ ] Layout mobile-first tersedia
[ ] Bottom navigation tersedia
[ ] Halaman awal `/` mengarah ke dashboard
```

---

### B. Halaman `/stocks/new`

**Nama halaman:** Input Stok Ikan Masuk

**Tujuan:** Digunakan oleh Baso atau admin gudang untuk mencatat ikan yang baru masuk setelah proses pembongkaran, pemilahan, dan penimbangan.

**API yang digunakan:**

```http
GET /fish-types
GET /cold-storages
POST /stocks
```

**Field form wajib:**

- `fish_type_id`
- `quality`
- `initial_weight_kg`
- `entered_at`
- `cold_storage_id`

**Field opsional:**

- `notes`

**Behavior:**

- Frontend mengambil daftar jenis ikan dari `GET /fish-types`.
- Frontend mengambil daftar lokasi cold storage dari `GET /cold-storages`.
- User mengisi form stok masuk.
- Frontend mengirim request ke `POST /stocks`.
- Jika berhasil, user diarahkan ke `/stocks`.
- Jika gagal, error dari backend ditampilkan.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /stocks/new
[ ] User bisa memilih jenis ikan
[ ] User bisa memilih kualitas ikan
[ ] User bisa input berat ikan
[ ] User bisa memilih waktu masuk
[ ] User bisa memilih cold storage
[ ] User bisa submit stok masuk
[ ] Data stok baru muncul di daftar FIFO
[ ] Loading state tersedia
[ ] Error state tersedia
```

---

### C. Halaman `/stock-outs/new`

**Nama halaman:** Input Ikan Keluar

**Tujuan:** Mencatat pengeluaran ikan dari cold storage dan mengurangi stok berdasarkan FIFO.

**API yang digunakan:**

```http
GET /fish-types
GET /stocks/fifo?fish_type_id={fish_type_id}
POST /stock-outs
```

**Field form wajib:**

- `fish_type_id`
- `total_weight_kg`
- `destination`
- `out_at`

**Field opsional:**

- `notes`

**Behavior:**

- User memilih jenis ikan.
- Frontend mengambil preview FIFO berdasarkan jenis ikan.
- User mengisi berat ikan keluar.
- User mengisi tujuan pengeluaran.
- Frontend mengirim request ke `POST /stock-outs`.
- Backend menentukan batch yang dikurangi berdasarkan FIFO.
- Frontend menampilkan summary batch yang dipakai dari response `items`.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /stock-outs/new
[ ] User bisa memilih jenis ikan
[ ] FIFO preview tampil setelah jenis ikan dipilih
[ ] User bisa input berat keluar
[ ] User bisa input tujuan pengeluaran
[ ] User bisa submit pengeluaran
[ ] Response items ditampilkan sebagai summary
[ ] Jika stok tidak cukup, error backend ditampilkan
[ ] Setelah berhasil, dashboard dan daftar stok bisa berubah
```

---

### Deliverable Mikail

```txt
[ ] src/lib/api.ts
[ ] src/types/api.ts
[ ] Layout mobile-first
[ ] Bottom navigation
[ ] Page header
[ ] Redirect halaman /
[ ] /stocks/new
[ ] /stock-outs/new
[ ] Loading state pada halaman input
[ ] Error state pada halaman input
[ ] Success state setelah submit
```

### Branch yang Disarankan

```bash
feat/frontend-foundation
feat/stock-in-page
feat/stock-out-page
```

---

## 2.2 Farhan — Data Flow, API Integration, FIFO QA

### Peran Utama

Farhan bertanggung jawab sebagai **Data Flow & FIFO QA**.

Karena fitur utama aplikasi sangat bergantung pada validitas data dan urutan FIFO, Farhan fokus memastikan data yang tampil di frontend sesuai dengan response backend dan logic FIFO berjalan benar.

---

### A. Halaman `/stocks`

**Nama halaman:** Daftar Stok FIFO

**Tujuan:** Menampilkan stok ikan yang tersedia dan membantu user melihat urutan stok berdasarkan FIFO.

**API yang digunakan:**

```http
GET /stocks/fifo
GET /fish-types
```

**Query FIFO:**

FIFO keseluruhan:

```http
GET /stocks/fifo
```

FIFO per jenis ikan:

```http
GET /stocks/fifo?fish_type_id={fish_type_id}
```

**Informasi yang ditampilkan:**

- ranking FIFO,
- jenis ikan,
- kualitas ikan,
- berat tersisa,
- waktu masuk,
- lokasi cold storage,
- status stok.

**Fitur:**

- menampilkan daftar stok berdasarkan FIFO keseluruhan,
- filter berdasarkan jenis ikan,
- menampilkan FIFO per jenis ikan,
- tombol menuju `/stocks/new`,
- tombol menuju `/stock-outs/new`,
- tampilan card mobile-friendly.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /stocks
[ ] Stok tampil berdasarkan FIFO
[ ] FIFO rank tampil jelas
[ ] Filter jenis ikan tersedia
[ ] Saat filter jenis ikan dipilih, data berubah sesuai query
[ ] Stok dengan remaining_weight_kg = 0 tidak diprioritaskan
[ ] Empty state tersedia
[ ] Loading state tersedia
[ ] Error state tersedia
```

---

### B. Dashboard Data Integration

Farhan membantu integrasi data pada halaman:

```txt
/dashboard
```

**API yang digunakan:**

```http
GET /dashboard/summary
GET /dashboard/recent-movements
```

**Informasi yang divalidasi:**

- total berat stok tersedia,
- total batch stok,
- total batch available,
- total batch depleted,
- total ikan masuk hari ini,
- total ikan keluar hari ini,
- ringkasan stok berdasarkan jenis ikan,
- ringkasan stok berdasarkan cold storage,
- recent movements.

**Acceptance criteria:**

```txt
[ ] Dashboard menampilkan summary dari backend
[ ] Recent movements tampil
[ ] Angka dashboard berubah setelah stok masuk
[ ] Angka dashboard berubah setelah stok keluar
[ ] Jika data kosong, empty state tetap rapi
```

---

### C. Validasi API Contract

Farhan bertugas memastikan response backend cocok dengan tipe data frontend.

Data yang perlu dicek:

```txt
FishType
ColdStorage
FIFOStock
DashboardSummary
RecentMovement
StockOut
```

Checklist validasi:

```txt
[ ] fish_type_name tersedia pada FIFO stock
[ ] fifo_rank tersedia pada FIFO stock
[ ] remaining_weight_kg berupa number
[ ] entered_at berupa ISO string
[ ] cold_storage_name tersedia
[ ] location_label tersedia jika ada
[ ] dashboard summary sesuai format
[ ] recent movement sesuai format
[ ] error insufficient stock bisa dibaca frontend
```

Jika ada mismatch antara backend dan frontend, Farhan mencatatnya dan menginfokan ke Pison.

---

### D. Skenario Testing FIFO

Farhan membuat dan menjalankan skenario testing FIFO.

**Skenario utama:**

```txt
1. Tambah jenis ikan Tuna.
2. Tambah Cold Storage A.
3. Tambah Cold Storage B.
4. Input Tuna 50 kg kualitas baik ke Cold Storage A pada jam 08.00.
5. Input Tuna 30 kg kualitas sedang ke Cold Storage B pada jam 09.00.
6. Buka /stocks dan filter Tuna.
7. Pastikan Tuna 50 kg menjadi FIFO rank #1.
8. Catat pengeluaran Tuna 60 kg.
9. Pastikan batch pertama habis.
10. Pastikan batch kedua tersisa 20 kg.
11. Buka dashboard.
12. Pastikan total stok dan recent movement berubah.
```

**Acceptance criteria testing:**

```txt
[ ] FIFO per jenis ikan berjalan benar
[ ] Pengeluaran stok mengambil batch paling lama dulu
[ ] Jika batch pertama tidak cukup, sistem mengambil batch berikutnya
[ ] Stok habis berubah menjadi depleted
[ ] Dashboard berubah setelah stock out
```

---

### Deliverable Farhan

```txt
[ ] /stocks
[ ] Filter FIFO per jenis ikan
[ ] FIFO rank tampil
[ ] Dashboard data integration
[ ] Checklist validasi API contract
[ ] Skenario testing FIFO
[ ] Catatan bug/mismatch API
```

### Branch yang Disarankan

```bash
feat/fifo-stock-list
feat/dashboard-integration
test/fifo-demo-scenario
```

---

## 2.3 Pannayaka — Product Flow, Master Data Pages, UX Copy

### Peran Utama

Pannayaka bertanggung jawab pada **master data, product flow, dan UX copy**.

Halaman master data relatif lebih sederhana dibanding halaman FIFO dan stock out, tetapi tetap penting karena data master digunakan oleh halaman input stok masuk dan input ikan keluar.

Selain itu, Pannayaka memastikan bahasa di aplikasi mudah dipahami oleh persona seperti Baso dan Daeng Syamsul.

---

### A. Halaman `/fish-types`

**Nama halaman:** Master Jenis Ikan

**Tujuan:** Mengelola daftar jenis ikan yang akan digunakan pada input stok.

**API yang digunakan:**

```http
GET /fish-types
POST /fish-types
```

**Field form:**

- `name`
- `image_url`
- `description`

**Fitur:**

- menampilkan daftar jenis ikan,
- menambahkan jenis ikan baru,
- menampilkan empty state jika belum ada data,
- menampilkan loading state,
- menampilkan error state.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /fish-types
[ ] Daftar jenis ikan tampil
[ ] User bisa menambah jenis ikan
[ ] Data baru muncul setelah submit
[ ] Jenis ikan baru bisa digunakan di /stocks/new
[ ] Jenis ikan baru bisa digunakan di /stock-outs/new
```

---

### B. Halaman `/cold-storages`

**Nama halaman:** Master Cold Storage

**Tujuan:** Mengelola daftar lokasi penyimpanan ikan.

**API yang digunakan:**

```http
GET /cold-storages
POST /cold-storages
```

**Field form:**

- `name`
- `location_label`
- `description`

**Fitur:**

- menampilkan daftar cold storage,
- menambahkan cold storage baru,
- menampilkan empty state jika belum ada data,
- menampilkan loading state,
- menampilkan error state.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /cold-storages
[ ] Daftar cold storage tampil
[ ] User bisa menambah cold storage
[ ] Data baru muncul setelah submit
[ ] Cold storage baru bisa digunakan di /stocks/new
```

---

### C. UX Copy dan Wording Aplikasi

Pannayaka mengecek semua teks pada aplikasi agar sesuai konteks studi kasus.

**Contoh wording yang disarankan:**

```txt
Tambah Stok Ikan
Catat Ikan Keluar
Prioritas FIFO
Stok Masuk Hari Ini
Stok Keluar Hari Ini
Sisa Berat
Lokasi Penyimpanan
Kualitas Ikan
Daftar Stok
Riwayat Pengeluaran
Cold Storage
```

**Hindari istilah yang terlalu teknis:**

```txt
Create Stock Batch
Submit Stock Movement
Transaction Header
Inventory Mutation
```

---

### D. Validasi Flow Persona

Pannayaka memastikan alur aplikasi sesuai kebutuhan Baso dan Daeng Syamsul.

**Flow Baso:**

```txt
Input stok masuk
→ cek daftar FIFO
→ catat ikan keluar
```

**Flow Daeng Syamsul:**

```txt
Buka dashboard
→ lihat ringkasan stok
→ lihat prioritas FIFO
→ pantau recent movement
```

**Acceptance criteria:**

```txt
[ ] Label halaman mudah dimengerti
[ ] Tombol utama mudah ditemukan
[ ] Alur Baso tidak terlalu panjang
[ ] Alur Daeng mudah dipahami dari dashboard
[ ] Empty state menggunakan bahasa yang jelas
```

---

### Deliverable Pannayaka

```txt
[ ] /fish-types
[ ] /cold-storages
[ ] UX wording aplikasi
[ ] Validasi user flow Baso
[ ] Validasi user flow Daeng
[ ] Empty state wording
```

### Branch yang Disarankan

```bash
feat/master-fish-types
feat/master-cold-storages
chore/ux-copy
```

---

## 2.4 Adit — Documentation, Demo Flow, QA, Presentation Support

### Peran Utama

Adit bertanggung jawab pada **dokumentasi, demo flow, QA manual, dan support presentasi**.

Adit juga mengerjakan halaman riwayat pengeluaran yang relatif lebih sederhana dibanding halaman input stok masuk dan input ikan keluar.

---

### A. Halaman `/stock-outs`

**Nama halaman:** Riwayat Pengeluaran Ikan

**Tujuan:** Menampilkan riwayat ikan keluar dari cold storage.

**API yang digunakan:**

```http
GET /stock-outs
```

**Informasi yang ditampilkan:**

- tujuan pengeluaran,
- total berat keluar,
- waktu keluar,
- catatan,
- batch yang dikurangi jika tersedia,
- jenis ikan pada item pengeluaran.

**Fitur:**

- menampilkan daftar pengeluaran,
- link ke halaman input ikan keluar,
- empty state,
- loading state,
- error state.

**Acceptance criteria:**

```txt
[ ] User bisa membuka halaman /stock-outs
[ ] Riwayat pengeluaran tampil
[ ] Data pengeluaran yang baru dibuat muncul
[ ] User dapat melihat tujuan pengeluaran
[ ] User dapat melihat total berat keluar
[ ] User dapat melihat waktu pengeluaran
[ ] Empty state tersedia
[ ] Loading state tersedia
[ ] Error state tersedia
```

---

### B. README Project

Adit mengurus `README.md` agar repository mudah dipahami.

Isi minimal README:

```txt
1. Nama project
2. Deskripsi singkat aplikasi
3. Studi kasus yang diselesaikan
4. Tech stack
5. Struktur folder
6. Cara menjalankan database
7. Cara menjalankan backend
8. Cara menjalankan frontend
9. Environment variable yang dibutuhkan
10. Daftar fitur MVP
11. Demo flow
12. Link repository dan anggota tim
```

**Acceptance criteria:**

```txt
[ ] README menjelaskan project secara singkat
[ ] README menjelaskan cara run backend
[ ] README menjelaskan cara run frontend
[ ] README menjelaskan cara run database
[ ] README menyebut fitur MVP
[ ] README bisa dipahami anggota lain
```

---

### C. Demo Script

Adit menyiapkan script demo untuk video pitching atau presentasi.

**Demo flow utama:**

```txt
1. Buka dashboard.
2. Tambah jenis ikan jika belum ada.
3. Tambah cold storage jika belum ada.
4. Input stok Tuna 50 kg kualitas baik ke Cold Storage A.
5. Input stok Tuna 30 kg kualitas sedang ke Cold Storage B.
6. Buka daftar stok FIFO.
7. Pilih filter Tuna.
8. Tampilkan batch Tuna terurut FIFO.
9. Catat pengeluaran Tuna 40 kg ke Restoran Laut Makassar.
10. Sistem mengurangi stok dari batch terlama.
11. Buka dashboard.
12. Dashboard menunjukkan stok dan recent movement terbaru.
```

**Acceptance criteria:**

```txt
[ ] Demo flow tertulis jelas
[ ] Setiap step punya halaman yang harus dibuka
[ ] Setiap step punya data input yang jelas
[ ] Demo bisa dilakukan dalam waktu singkat
[ ] Demo menunjukkan masalah FIFO terjawab
```

---

### D. QA Checklist

Adit membuat checklist manual sebelum submission.

**Checklist teknis:**

```txt
[ ] Backend menyala
[ ] Frontend menyala
[ ] Database menyala
[ ] Environment variable sudah benar
[ ] API base URL frontend benar
[ ] Tidak ada error fatal di console
```

**Checklist fitur:**

```txt
[ ] /dashboard bisa dibuka
[ ] /fish-types bisa tambah data
[ ] /cold-storages bisa tambah data
[ ] /stocks/new bisa input stok
[ ] /stocks menampilkan FIFO
[ ] /stock-outs/new bisa catat pengeluaran
[ ] /stock-outs menampilkan riwayat
[ ] Dashboard berubah setelah stok masuk
[ ] Dashboard berubah setelah stok keluar
```

**Checklist UI:**

```txt
[ ] Tampilan mobile tidak rusak
[ ] Tombol utama mudah ditemukan
[ ] Loading state tampil
[ ] Error backend tampil di frontend
[ ] Empty state tampil saat data kosong
```

---

### Deliverable Adit

```txt
[ ] /stock-outs
[ ] README.md
[ ] Demo script
[ ] QA checklist
[ ] Screenshot atau bahan video demo
[ ] Catatan bug sebelum final submission
```

### Branch yang Disarankan

```bash
feat/stock-out-history
docs/readme
docs/demo-script
test/manual-qa
```

---

# 3. Ringkasan Pembagian Halaman

| Area / Halaman | Owner Utama | Support |
|---|---|---|
| Frontend foundation | Mikail | Farhan |
| `src/lib/api.ts` | Mikail | Farhan |
| `src/types/api.ts` | Mikail | Farhan |
| Layout mobile + bottom nav | Mikail | Pannayaka |
| `/` redirect | Mikail | - |
| `/dashboard` | Farhan | Mikail |
| `/fish-types` | Pannayaka | Mikail |
| `/cold-storages` | Pannayaka | Mikail |
| `/stocks/new` | Mikail | Farhan |
| `/stocks` | Farhan | Mikail |
| `/stock-outs/new` | Mikail | Farhan |
| `/stock-outs` | Adit | Farhan |
| UX copy | Pannayaka | Adit |
| README & docs | Adit | Pison |
| Demo script | Adit | Pannayaka |
| Manual QA | Adit + Farhan | Semua |

---

# 4. Urutan Kerja

## Tahap 1 — Foundation dan Persiapan

### Mikail

```txt
Setup layout, API helper, types, dan bottom nav.
```

### Pannayaka

```txt
Mulai /fish-types dan /cold-storages.
```

### Farhan

```txt
Siapkan skenario data FIFO dan cek format API response.
```

### Adit

```txt
Siapkan README, demo script, dan QA checklist.
```

---

## Tahap 2 — Master Data dan Input Stok

### Pannayaka

```txt
Integrasi /fish-types dan /cold-storages ke API asli.
```

### Mikail

```txt
Mulai /stocks/new karena butuh fish types dan cold storages.
```

### Farhan

```txt
Cek data stok masuk muncul benar di database dan API.
```

### Adit

```txt
Update README dengan cara setup terbaru.
```

---

## Tahap 3 — FIFO dan Stock Out

### Farhan

```txt
Kerjakan /stocks dan validasi FIFO.
```

### Mikail

```txt
Kerjakan /stock-outs/new.
```

### Adit

```txt
Mulai /stock-outs dan update demo script sesuai UI terbaru.
```

### Pannayaka

```txt
Review UX wording pada halaman stok dan pengeluaran.
```

---

## Tahap 4 — Dashboard, QA, dan Polish

### Farhan

```txt
Integrasi dashboard dan validasi angka stok.
```

### Mikail

```txt
Polish UI halaman input utama.
```

### Adit

```txt
Tes demo flow full dari awal sampai akhir.
```

### Pannayaka

```txt
Cek apakah wording dan flow mudah dipahami untuk Baso dan Daeng.
```

### Semua Anggota

```txt
Bug fixing dan final review.
```

---

# 5. Prioritas Berdasarkan Deadline

Jika waktu sangat terbatas, gunakan prioritas berikut.

## Mikail

```txt
[HIGH] Layout mobile
[HIGH] /stocks/new
[HIGH] /stock-outs/new
```

## Farhan

```txt
[HIGH] /stocks
[HIGH] Validasi FIFO
[HIGH] Dashboard integration minimal
```

## Pannayaka

```txt
[MEDIUM] /fish-types
[MEDIUM] /cold-storages
[HIGH] UX wording dan flow check
```

## Adit

```txt
[MEDIUM] /stock-outs
[HIGH] README
[HIGH] Demo script
[HIGH] QA checklist
```

---

# 6. GitHub Issues yang Disarankan

## Mikail

```txt
[FE] Setup frontend API helper, types, and mobile layout
[FE] Build stock-in page /stocks/new
[FE] Build stock-out form /stock-outs/new
```

## Farhan

```txt
[FE] Build FIFO stock list /stocks
[FE] Integrate dashboard summary and recent movements
[QA] Validate FIFO stock-out scenario
```

## Pannayaka

```txt
[FE] Build fish type master page /fish-types
[FE] Build cold storage master page /cold-storages
[UX] Review copywriting and Baso-Daeng user flow
```

## Adit

```txt
[FE] Build stock-out history page /stock-outs
[DOCS] Write README and setup instructions
[QA] Prepare demo script and manual test checklist
```

---

# 7. Branch Strategy

Agar rapi, gunakan branch `dev` sebagai tempat integrasi sebelum merge ke `main`.

Struktur branch:

```txt
main
└── dev
    ├── feat/mikail-core-frontend
    ├── feat/farhan-fifo-dashboard
    ├── feat/pannayaka-master-data
    └── feat/adit-docs-demo
```

## Branch per Anggota

### Mikail

```bash
git checkout -b feat/mikail-core-frontend
```

### Farhan

```bash
git checkout -b feat/farhan-fifo-dashboard
```

### Pannayaka

```bash
git checkout -b feat/pannayaka-master-data
```

### Adit

```bash
git checkout -b feat/adit-docs-demo
```

---

# 8. Commit Convention

Gunakan format commit sederhana:

```txt
type(scope): message
```

Contoh:

```bash
feat(frontend): add mobile app shell
feat(stocks): add stock in form
feat(fifo): add fifo stock list
feat(master-data): add fish type page
docs(readme): add local setup guide
test(manual): add fifo demo checklist
fix(stock-out): handle insufficient stock error
```

Tipe commit yang digunakan:

```txt
feat     = fitur baru
fix      = perbaikan bug
docs     = dokumentasi
style    = perubahan tampilan tanpa logic besar
refactor = perapian kode
test     = testing/checklist
chore    = setup/tools/config
```

---

# 9. Checklist Integrasi Harian

Sebelum push atau merge ke `dev`, setiap anggota wajib cek:

```txt
[ ] App frontend masih bisa run
[ ] Tidak ada error TypeScript besar
[ ] Tidak ada file .env ikut ter-commit
[ ] Halaman yang dikerjakan bisa dibuka
[ ] API endpoint yang dipakai sesuai contract
[ ] Loading state tersedia
[ ] Error state minimal tersedia
[ ] Tampilan mobile tidak rusak total
```

---

# 10. Definition of Done Tim

MVP tim dianggap selesai jika:

```txt
[ ] Backend Go Fiber menyala
[ ] PostgreSQL menyala
[ ] Frontend Next.js menyala
[ ] /dashboard bisa menampilkan summary dan recent movement
[ ] /fish-types bisa menambah jenis ikan
[ ] /cold-storages bisa menambah cold storage
[ ] /stocks/new bisa mencatat stok masuk
[ ] /stocks bisa menampilkan stok FIFO
[ ] /stock-outs/new bisa mencatat ikan keluar
[ ] /stock-outs bisa menampilkan riwayat pengeluaran
[ ] Stok keluar mengurangi batch berdasarkan FIFO
[ ] Dashboard berubah setelah stok masuk dan keluar
[ ] Demo flow bisa dijalankan tanpa error fatal
[ ] README cukup jelas untuk menjalankan project
[ ] Video/presentasi bisa menunjukkan alur input → FIFO → output
```

---

# 11. Catatan Koordinasi dengan Backend

Agar frontend tidak terlalu lama menunggu backend, Pison perlu memprioritaskan endpoint berikut:

```http
GET  /fish-types
POST /fish-types

GET  /cold-storages
POST /cold-storages

POST /stocks
GET  /stocks/fifo

POST /stock-outs
GET  /stock-outs

GET  /dashboard/summary
GET  /dashboard/recent-movements
```

Jika endpoint belum selesai, anggota frontend boleh menggunakan mock data sementara, tetapi harus diberi komentar jelas agar nanti mudah diganti ke API asli.

Contoh komentar:

```ts
// TODO: replace mock data with GET /stocks/fifo
```

---

# 12. Risiko dan Mitigasi

## Risiko 1 — Backend belum siap saat frontend dikerjakan

Mitigasi:

```txt
Gunakan mock data sementara.
Tetapkan response shape sesuai api-contract.
Farhan bertugas mengganti mock ke API asli setelah backend siap.
```

## Risiko 2 — Banyak konflik Git

Mitigasi:

```txt
Gunakan branch per anggota.
Jangan mengedit file yang sama tanpa koordinasi.
Merge ke dev secara bertahap.
```

## Risiko 3 — Tampilan mobile tidak rapi

Mitigasi:

```txt
Mikail membuat layout dasar reusable.
Semua halaman mengikuti komponen layout yang sama.
Pannayaka dan Adit membantu review mobile flow.
```

## Risiko 4 — FIFO tidak terlihat saat demo

Mitigasi:

```txt
Farhan menyiapkan data demo dengan dua batch ikan yang sama.
Demo harus menunjukkan batch lama berkurang lebih dulu.
```

## Risiko 5 — Waktu tidak cukup

Mitigasi:

```txt
Fokus pada halaman minimum demo:
1. /stocks/new
2. /stocks
3. /stock-outs/new
4. /dashboard
```

Halaman master data dan riwayat pengeluaran tetap dikerjakan jika waktu cukup, atau menggunakan data seed dari backend jika perlu.

---

# 13. Rangkuman Singkat

## Mikail

Fokus pada:

```txt
frontend foundation, layout mobile, /stocks/new, /stock-outs/new
```

## Farhan

Fokus pada:

```txt
/stocks, FIFO validation, dashboard integration, API response checking
```

## Pannayaka

Fokus pada:

```txt
/fish-types, /cold-storages, UX copy, user flow validation
```

## Adit

Fokus pada:

```txt
/stock-outs, README, demo script, manual QA, presentation support
```

## Pison

Fokus pada:

```txt
backend, database, API contract, FIFO logic, integration support
```
