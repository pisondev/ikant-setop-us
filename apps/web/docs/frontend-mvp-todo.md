# Frontend MVP Todo | "Ikan't Setop Us"

Dokumen ini menjadi acuan pengerjaan frontend MVP untuk aplikasi FishFlow: FIFO Fish Inventory Management System.

Frontend dibuat menggunakan **Next.js** dengan pendekatan **mobile-first UI**. Fokus MVP adalah membuktikan alur utama aplikasi berjalan end-to-end:

```txt
Input stok ikan masuk
→ data tersimpan di database
→ stok tampil berdasarkan FIFO
→ stok keluar dapat dicatat
→ dashboard ikut berubah
```

---

## 1. Prinsip Frontend MVP

### 1.1 Mobile-first

Aplikasi harus nyaman digunakan dari layar HP karena user utama di lapangan adalah pekerja seperti Baso.

Prioritas desain:

- tombol besar,
- form pendek,
- minim typing,
- mudah digunakan satu tangan,
- informasi stok mudah dibaca,
- layout tetap rapi di layar kecil.

### 1.2 MVP Scope

Frontend hanya fokus pada fitur inti:

- dashboard monitoring,
- master jenis ikan,
- master cold storage,
- input stok ikan masuk,
- daftar stok FIFO,
- input ikan keluar,
- riwayat pengeluaran sederhana,
- recent movement.

Fitur yang tidak masuk MVP:

- login/register,
- role-based access control,
- payment,
- marketplace,
- integrasi alat timbang,
- sensor cold storage,
- notifikasi realtime,
- grafik kompleks,
- halaman detail stok yang kompleks.

---

## 2. Environment

Frontend mengambil API dari backend Go Fiber.

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api/v1
```

Semua request API sebaiknya dibuat melalui helper:

```txt
src/lib/api.ts
```

---

## 3. Struktur Route MVP

Struktur route yang disarankan:

```txt
src/app/
├── page.tsx
├── dashboard/
│   └── page.tsx
├── fish-types/
│   └── page.tsx
├── cold-storages/
│   └── page.tsx
├── stocks/
│   ├── page.tsx
│   └── new/
│       └── page.tsx
└── stock-outs/
    ├── page.tsx
    └── new/
        └── page.tsx
```

Halaman yang sengaja tidak diwajibkan untuk MVP:

```txt
/stocks/[id]
```

Alasannya, update kualitas dan update lokasi bisa ditunda agar MVP tidak terlalu kompleks. Jika waktu masih cukup, halaman detail stok bisa ditambahkan sebagai fitur tambahan.

---

# 4. Page Todo

---

## 4.1 `/`

### Nama Halaman

Home / Redirect Page

### Tujuan

Halaman awal aplikasi. Untuk MVP, halaman ini boleh langsung redirect ke dashboard.

### Fitur

- Menampilkan nama aplikasi secara singkat, atau
- Redirect otomatis ke `/dashboard`.

### API yang Dibutuhkan

Tidak ada.

### MVP Priority

High

### Acceptance Criteria

- User membuka `/`.
- User langsung diarahkan ke `/dashboard`, atau melihat tombol masuk ke dashboard.

---

## 4.2 `/dashboard`

### Nama Halaman

Dashboard Monitoring

### Tujuan

Memberikan ringkasan stok ikan kepada Daeng Syamsul atau pengelola gudang.

### Informasi yang Ditampilkan

- Total berat stok tersedia.
- Total batch stok.
- Total batch available.
- Total batch depleted.
- Total ikan masuk hari ini.
- Total ikan keluar hari ini.
- Ringkasan stok berdasarkan jenis ikan.
- Ringkasan stok berdasarkan cold storage.
- Recent movements.

### Komponen UI

- Summary card total stok.
- Summary card ikan masuk hari ini.
- Summary card ikan keluar hari ini.
- List stok per jenis ikan.
- List stok per cold storage.
- Recent movement list.
- Tombol cepat:
  - Tambah Stok Masuk
  - Catat Ikan Keluar
  - Lihat FIFO

### API yang Dibutuhkan

```http
GET /dashboard/summary
GET /dashboard/recent-movements
```

### State

- loading
- error
- success
- empty recent movement

### MVP Priority

High

### Acceptance Criteria

- User dapat melihat ringkasan stok.
- Data berubah setelah stok masuk atau stok keluar dibuat.
- Recent movement menampilkan aktivitas terakhir.

---

## 4.3 `/fish-types`

### Nama Halaman

Master Jenis Ikan

### Tujuan

Mengelola daftar jenis ikan yang akan digunakan pada input stok.

### Informasi yang Ditampilkan

- Nama jenis ikan.
- Gambar ikan jika ada.
- Deskripsi singkat jika ada.

### Fitur

- Melihat daftar jenis ikan.
- Menambahkan jenis ikan baru.

### Form Tambah Jenis Ikan

Field:

- `name`
- `image_url`
- `description`

### API yang Dibutuhkan

```http
GET /fish-types
POST /fish-types
```

### State

- loading list
- submit loading
- error
- empty state

### MVP Priority

Medium

### Acceptance Criteria

- User dapat melihat daftar jenis ikan.
- User dapat menambah jenis ikan.
- Jenis ikan yang baru ditambahkan bisa dipakai di halaman input stok masuk dan input ikan keluar.

---

## 4.4 `/cold-storages`

### Nama Halaman

Master Cold Storage

### Tujuan

Mengelola daftar lokasi penyimpanan ikan.

### Informasi yang Ditampilkan

- Nama cold storage.
- Label lokasi.
- Deskripsi.

### Fitur

- Melihat daftar cold storage.
- Menambahkan lokasi cold storage baru.

### Form Tambah Cold Storage

Field:

- `name`
- `location_label`
- `description`

### API yang Dibutuhkan

```http
GET /cold-storages
POST /cold-storages
```

### State

- loading list
- submit loading
- error
- empty state

### MVP Priority

Medium

### Acceptance Criteria

- User dapat melihat daftar cold storage.
- User dapat menambah cold storage.
- Cold storage yang baru ditambahkan bisa dipakai di halaman input stok masuk.

---

## 4.5 `/stocks/new`

### Nama Halaman

Input Stok Ikan Masuk

### Tujuan

Digunakan Baso atau admin gudang untuk mencatat ikan yang baru masuk setelah proses pembongkaran, pemilahan, dan penimbangan.

### Prinsip UI

- Mobile-first.
- Form pendek.
- Input cepat.
- Tombol besar.
- Dropdown/select untuk data master.

### Form Field

Wajib:

- `fish_type_id`
- `quality`
- `initial_weight_kg`
- `entered_at`
- `cold_storage_id`

Opsional:

- `notes`

### Input Behavior

- `remaining_weight_kg` tidak diinput dari frontend.
- Backend otomatis mengisi `remaining_weight_kg = initial_weight_kg`.
- Backend otomatis mengisi `status = available`.
- Backend membuat histori movement type `in`.

### API yang Dibutuhkan

```http
GET /fish-types
GET /cold-storages
POST /stocks
```

### Request Body

```json
{
  "fish_type_id": "uuid",
  "cold_storage_id": "uuid",
  "quality": "baik",
  "initial_weight_kg": 50,
  "entered_at": "2026-05-01T08:00:00Z",
  "notes": "Tangkapan pagi"
}
```

### Setelah Submit Berhasil

Redirect ke:

```txt
/stocks
```

atau tampilkan tombol:

```txt
Lihat Daftar Stok
```

### State

- loading master data
- submit loading
- validation error
- success
- failed submit

### MVP Priority

High

### Acceptance Criteria

- User dapat mencatat stok ikan masuk.
- Data stok tersimpan di database.
- Stok baru muncul di daftar FIFO.
- Dashboard berubah setelah stok masuk.

---

## 4.6 `/stocks`

### Nama Halaman

Daftar Stok FIFO

### Tujuan

Menampilkan stok ikan yang tersedia dan membantu user melihat urutan FIFO.

### Informasi yang Ditampilkan

- Ranking FIFO.
- Jenis ikan.
- Kualitas ikan.
- Berat tersisa.
- Waktu masuk.
- Lokasi cold storage.
- Status stok.

### Fitur

- Melihat stok berdasarkan FIFO keseluruhan.
- Filter berdasarkan jenis ikan untuk melihat FIFO per jenis ikan.
- Tombol menuju input stok masuk.
- Tombol menuju input ikan keluar.
- Tampilan stok dalam bentuk card agar nyaman di mobile.

### API yang Dibutuhkan

```http
GET /stocks/fifo
GET /fish-types
```

### Query yang Digunakan

Untuk FIFO keseluruhan:

```http
GET /stocks/fifo
```

Untuk FIFO per jenis ikan:

```http
GET /stocks/fifo?fish_type_id={fish_type_id}
```

### State

- loading
- error
- empty stock
- filter active
- success

### MVP Priority

High

### Acceptance Criteria

- User dapat melihat stok terurut FIFO.
- User dapat memilih jenis ikan dan melihat FIFO per jenis ikan.
- Stok dengan `remaining_weight_kg = 0` tidak diprioritaskan sebagai stok available.
- User dapat memahami stok mana yang sebaiknya dikeluarkan lebih dulu.

---

## 4.7 `/stock-outs/new`

### Nama Halaman

Input Ikan Keluar

### Tujuan

Mencatat pengeluaran ikan dari cold storage dan mengurangi stok berdasarkan FIFO.

### Prinsip UI

- User memilih jenis ikan.
- User memasukkan berat keluar.
- User mengisi tujuan pengeluaran.
- Backend menentukan batch yang dikurangi berdasarkan FIFO.
- Frontend menampilkan hasil batch yang dipakai.

### Form Field

Wajib:

- `fish_type_id`
- `total_weight_kg`
- `destination`
- `out_at`

Opsional:

- `notes`

### API yang Dibutuhkan

```http
GET /fish-types
GET /stocks/fifo?fish_type_id={fish_type_id}
POST /stock-outs
```

### Request Body

```json
{
  "fish_type_id": "uuid",
  "total_weight_kg": 40,
  "destination": "Restoran Laut Makassar",
  "out_at": "2026-05-01T12:00:00Z",
  "notes": "Pengeluaran untuk pesanan makan siang"
}
```

### Response Penting

Frontend perlu menampilkan `items` dari response.

```json
{
  "id": "uuid",
  "fish_type_id": "uuid",
  "destination": "Restoran Laut Makassar",
  "total_weight_kg": 40,
  "out_at": "2026-05-01T12:00:00Z",
  "items": [
    {
      "stock_batch_id": "uuid",
      "weight_kg": 25
    },
    {
      "stock_batch_id": "uuid",
      "weight_kg": 15
    }
  ],
  "created_at": "2026-05-01T12:01:00Z"
}
```

### Preview FIFO Sebelum Submit

Setelah user memilih jenis ikan, frontend menampilkan daftar stok available untuk jenis ikan tersebut:

```http
GET /stocks/fifo?fish_type_id={fish_type_id}
```

Tujuannya agar user tahu batch mana yang akan diprioritaskan.

### Setelah Submit Berhasil

Tampilkan summary:

```txt
Pengeluaran berhasil dicatat.
Sistem mengambil stok dari beberapa batch berdasarkan FIFO.
```

Lalu sediakan tombol:

- Lihat Daftar Stok
- Lihat Dashboard
- Catat Pengeluaran Lagi

### State

- loading fish types
- loading FIFO preview
- submit loading
- insufficient stock error
- success
- failed submit

### MVP Priority

High

### Acceptance Criteria

- User dapat mencatat ikan keluar.
- Backend mengurangi stok berdasarkan FIFO.
- Jika stok tidak cukup, frontend menampilkan error dari backend.
- Setelah pengeluaran berhasil, dashboard dan daftar stok berubah.

---

## 4.8 `/stock-outs`

### Nama Halaman

Riwayat Pengeluaran Ikan

### Tujuan

Menampilkan riwayat ikan keluar dari cold storage.

### Informasi yang Ditampilkan

- Tujuan pengeluaran.
- Total berat keluar.
- Waktu keluar.
- Catatan.
- Batch yang dikurangi jika tersedia.
- Jenis ikan pada item pengeluaran.

### Fitur

- Melihat daftar pengeluaran.
- Link ke halaman input ikan keluar.

### API yang Dibutuhkan

```http
GET /stock-outs
```

### Optional Query

```http
GET /stock-outs?date_from=2026-05-01&date_to=2026-05-02
```

### State

- loading
- error
- empty state
- success

### MVP Priority

Medium

### Acceptance Criteria

- User dapat melihat riwayat ikan keluar.
- Data pengeluaran yang baru dibuat muncul di daftar.
- User dapat memahami tujuan pengeluaran dan jumlah berat keluar.

---

# 5. Komponen Frontend yang Disarankan

## 5.1 Layout Components

```txt
src/components/layout/
├── app-shell.tsx
├── bottom-nav.tsx
└── page-header.tsx
```

### Fungsi

- `AppShell`: wrapper layout mobile.
- `BottomNav`: navigasi bawah untuk mobile.
- `PageHeader`: judul halaman dan action button.

Menu bottom nav MVP:

```txt
Dashboard
Stok
Tambah
Keluar
Master
```

---

## 5.2 Dashboard Components

```txt
src/components/dashboard/
├── summary-card.tsx
├── fish-type-summary-list.tsx
├── storage-summary-list.tsx
└── recent-movement-list.tsx
```

---

## 5.3 Stock Components

```txt
src/components/stocks/
├── stock-card.tsx
├── fifo-rank-badge.tsx
├── quality-badge.tsx
└── stock-empty-state.tsx
```

---

## 5.4 Form Components

```txt
src/components/forms/
├── fish-type-select.tsx
├── cold-storage-select.tsx
├── quality-select.tsx
├── weight-input.tsx
└── submit-button.tsx
```

---

# 6. API Helper Todo

Buat file:

```txt
src/lib/api.ts
```

Fungsi dasar yang perlu ada:

```txt
apiGet(path)
apiPost(path, body)
apiPatch(path, body)
```

Base URL:

```txt
process.env.NEXT_PUBLIC_API_BASE_URL
```

Response mengikuti format backend:

```ts
export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  errors?: unknown;
  meta?: unknown;
};
```

---

# 7. TypeScript Types Todo

Buat file:

```txt
src/types/api.ts
```

Types minimum:

```ts
export type FishQuality = "baik" | "sedang" | "buruk";

export type StockStatus = "available" | "depleted";

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  errors?: unknown;
  meta?: unknown;
};

export type FishType = {
  id: string;
  name: string;
  image_url?: string | null;
  description?: string | null;
  created_at: string;
  updated_at: string;
};

export type ColdStorage = {
  id: string;
  name: string;
  location_label?: string | null;
  description?: string | null;
  created_at: string;
  updated_at: string;
};

export type FIFOStock = {
  id: string;
  fish_type_name: string;
  quality: FishQuality;
  remaining_weight_kg: number;
  entered_at: string;
  cold_storage_name: string;
  location_label?: string | null;
  fifo_rank: number;
};

export type DashboardSummary = {
  total_available_weight_kg: number;
  total_stock_batches: number;
  total_available_batches: number;
  total_depleted_batches: number;
  today_stock_in_kg: number;
  today_stock_out_kg: number;
  fish_type_summary: {
    fish_type_id: string;
    fish_type_name: string;
    available_weight_kg: number;
    available_batches: number;
  }[];
  cold_storage_summary: {
    cold_storage_id: string;
    cold_storage_name: string;
    available_weight_kg: number;
    available_batches: number;
  }[];
};

export type RecentMovement = {
  id: string;
  stock_batch_id: string;
  movement_type: "in" | "out" | "quality_update" | "location_update" | "adjustment";
  fish_type_name: string;
  weight_kg?: number | null;
  description: string;
  created_at: string;
};

export type StockOut = {
  id: string;
  destination: string;
  total_weight_kg: number;
  out_at: string;
  notes?: string | null;
  items: {
    stock_batch_id: string;
    fish_type_name?: string;
    weight_kg: number;
  }[];
  created_at: string;
};
```

---

# 8. Timeline Pengerjaan Frontend

## Tahap 1 : Foundation

- Setup `NEXT_PUBLIC_API_BASE_URL`.
- Buat `src/lib/api.ts`.
- Buat `src/types/api.ts`.
- Buat layout mobile sederhana.
- Buat bottom navigation.

## Tahap 2 : Master Data

- Buat `/fish-types`.
- Buat `/cold-storages`.

Tujuannya agar data master tersedia untuk form stok.

## Tahap 3 : Stok Masuk

- Buat `/stocks/new`.
- Integrasi `GET /fish-types`.
- Integrasi `GET /cold-storages`.
- Integrasi `POST /stocks`.

Ini fitur input utama untuk Baso.

## Tahap 4 : FIFO List

- Buat `/stocks`.
- Integrasi `GET /stocks/fifo`.
- Tambahkan filter per jenis ikan.
- Tampilkan urutan FIFO.

Ini fitur utama untuk menjawab masalah FIFO.

## Tahap 5 : Stok Keluar

- Buat `/stock-outs/new`.
- Integrasi `GET /fish-types`.
- Integrasi `GET /stocks/fifo?fish_type_id={id}`.
- Integrasi `POST /stock-outs`.

Ini membuktikan proses end-to-end.

## Tahap 6 : Dashboard

- Buat `/dashboard`.
- Integrasi `GET /dashboard/summary`.
- Integrasi `GET /dashboard/recent-movements`.

## Tahap 7 : Riwayat Pengeluaran

- Buat `/stock-outs`.
- Integrasi `GET /stock-outs`.

## Tahap 8 : Polish MVP

- Tambahkan loading state.
- Tambahkan error state.
- Tambahkan empty state.
- Rapikan mobile spacing.
- Pastikan tombol utama mudah ditekan.
- Pastikan demo flow lancar.

---

# 9. Demo Flow

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

---

# 10. Prioritas Halaman (Plan B)

Jika waktu sangat terbatas, kerjakan halaman dengan urutan ini:

```txt
1. /stocks/new
2. /stocks
3. /stock-outs/new
4. /dashboard
5. /fish-types
6. /cold-storages
7. /stock-outs
```

Minimum demo:

```txt
/stocks/new
/stocks
/stock-outs/new
/dashboard
```

Dengan 4 halaman itu, MVP sudah bisa menunjukkan:

```txt
input → proses FIFO → output dashboard
```

---

# 11. Definition of Done MVP Frontend

Frontend MVP dianggap selesai jika:

- `/stocks/new` bisa membuat stok masuk.
- `/stocks` bisa menampilkan stok terurut FIFO.
- `/stock-outs/new` bisa mencatat ikan keluar dan menampilkan hasil pengurangan batch.
- `/dashboard` bisa menampilkan summary dan recent movement.
- Minimal data master `fish-types` dan `cold-storages` bisa dibuat dari UI atau sudah tersedia dari seed/backend.
- Semua halaman utama nyaman dibuka di layar mobile.
- Error dari backend bisa ditampilkan ke user.
- Loading state tersedia pada request utama.
