# Frontend MVP Todo | Ikan't Setop Us

Dokumen ini adalah single source of truth untuk status frontend MVP. Dokumen pembagian tugas di `docs/jobdesc_progress.md` harus merujuk ke checklist ini agar progress tidak dobel dan tidak saling bertentangan.

Update terakhir: 2026-05-06.

## Target MVP

Alur utama yang harus bisa didemokan:

```txt
input stok masuk -> data tersimpan -> stok tampil FIFO -> stok keluar dicatat -> dashboard berubah
```

## Status Aktual Repo

| Area | Status | Catatan |
|---|---|---|
| Backend API | Selesai | Semua endpoint MVP tersedia di `/api/v1` dan sudah dites. |
| API docs | Selesai | Contract, Postman collection, dan environment ada di `docs/api`. |
| Frontend foundation | Selesai tahap 1 | API helper, shared types, app shell, bottom nav, dan page header sudah ada. |
| `/` | Selesai sementara | Template Next.js sudah diganti menjadi entry point aplikasi ke route yang tersedia. |
| `/stocks` | Sebagian lanjut | Sudah memakai API helper/type bersama dan error state eksplisit. |
| Master data frontend | Selesai | `/fish-types` dan `/cold-storages` sudah terintegrasi GET/POST API. |
| Stock-in frontend | Selesai | `/stocks/new` sudah terintegrasi GET master data dan POST `/stocks`. |
| Stock-out form frontend | Selesai | `/stock-outs/new` sudah punya FIFO preview dan POST `/stock-outs`. |
| Halaman frontend lain | Belum | Dashboard dan riwayat belum ada. |

## API Base URL

Frontend harus mengambil backend dari:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api/v1
```

API contract:

```txt
docs/api/api-contract-v1.md
```

## Struktur Route Target

| Route | Prioritas | Status | API |
|---|---:|---|---|
| `/` | High | Selesai sementara | Tidak ada |
| `/dashboard` | High | Belum | `GET /dashboard/summary`, `GET /dashboard/recent-movements` |
| `/fish-types` | Medium | Selesai | `GET /fish-types`, `POST /fish-types` |
| `/cold-storages` | Medium | Selesai | `GET /cold-storages`, `POST /cold-storages` |
| `/stocks/new` | High | Selesai | `GET /fish-types`, `GET /cold-storages`, `POST /stocks` |
| `/stocks` | High | Sebagian | `GET /stocks/fifo`, `GET /fish-types` |
| `/stock-outs/new` | High | Selesai | `GET /fish-types`, `GET /stocks/fifo?fish_type_id={id}`, `POST /stock-outs` |
| `/stock-outs` | Medium | Belum | `GET /stock-outs` |

## Checklist Foundation

```txt
[x] Buat API helper di apps/web/lib/api.ts.
[x] Buat TypeScript type API bersama di apps/web/types/api.ts.
[~] Setup NEXT_PUBLIC_API_BASE_URL: helper punya default lokal, .env.local tetap disarankan untuk dev.
[x] Buat layout mobile-first reusable.
[x] Buat bottom navigation.
[x] Buat page header reusable.
[x] Ubah / dari template Next.js menjadi entry point aplikasi.
[x] Hapus fallback mock dari /stocks setelah backend siap dipakai di dev.
[~] Pastikan loading, error, dan empty state punya pola yang konsisten: sudah diterapkan di /stocks, perlu diulang pada halaman baru.
```

## Checklist Per Halaman

### `/`

Acceptance criteria:

```txt
[x] User membuka / dan melihat entry point aplikasi.
[x] Tidak ada konten template Next.js tersisa.
```

### `/dashboard`

Informasi yang ditampilkan:

- total berat stok tersedia,
- total batch stok,
- total batch available,
- total batch depleted,
- stok masuk hari ini,
- stok keluar hari ini,
- ringkasan stok per jenis ikan,
- ringkasan stok per cold storage,
- recent movements.

Acceptance criteria:

```txt
[ ] Halaman /dashboard tersedia.
[ ] GET /dashboard/summary terintegrasi.
[ ] GET /dashboard/recent-movements terintegrasi.
[ ] Loading state tersedia.
[ ] Error state tersedia.
[ ] Empty recent movement tetap rapi.
[ ] Angka berubah setelah stok masuk.
[ ] Angka berubah setelah stok keluar.
```

### `/fish-types`

Acceptance criteria:

```txt
[x] Halaman /fish-types tersedia.
[x] GET /fish-types terintegrasi.
[x] POST /fish-types terintegrasi.
[x] Form field name, image_url, description tersedia.
[x] Data baru muncul setelah submit.
[x] Loading state tersedia.
[x] Error state tersedia.
[x] Empty state tersedia.
```

### `/cold-storages`

Acceptance criteria:

```txt
[x] Halaman /cold-storages tersedia.
[x] GET /cold-storages terintegrasi.
[x] POST /cold-storages terintegrasi.
[x] Form field name, location_label, description tersedia.
[x] Data baru muncul setelah submit.
[x] Loading state tersedia.
[x] Error state tersedia.
[x] Empty state tersedia.
```

### `/stocks/new`

Request body:

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

Acceptance criteria:

```txt
[x] Halaman /stocks/new tersedia.
[x] GET /fish-types terintegrasi.
[x] GET /cold-storages terintegrasi.
[x] POST /stocks terintegrasi.
[x] User bisa memilih jenis ikan.
[x] User bisa memilih kualitas ikan.
[x] User bisa input berat ikan.
[x] User bisa memilih waktu masuk.
[x] User bisa memilih cold storage.
[x] Jika sukses, user diarahkan ke /stocks.
[x] Error validasi backend tampil.
```

### `/stocks`

Status saat ini: sebagian selesai di `apps/web/app/stocks/page.tsx`.

Yang sudah ada:

```txt
[x] Halaman /stocks tersedia.
[x] Fetch GET /fish-types.
[x] Fetch GET /stocks/fifo.
[x] Filter jenis ikan memanggil /stocks/fifo?fish_type_id={id}.
[x] Loading state tersedia.
[x] Empty state tersedia.
[x] Error state tampil ketika API gagal.
```

Yang masih perlu:

```txt
[x] Gunakan API helper bersama.
[x] Gunakan TypeScript type bersama.
[x] Perbaiki tulisan/karakter rusak pada empty state.
[x] Tambahkan error state eksplisit, bukan langsung fallback diam-diam.
[x] Tambahkan tombol menuju /stocks/new dan /stock-outs/new.
[x] Validasi build setelah layout global tersedia.
[x] Hapus mock fallback ketika backend menjadi dependency wajib.
[ ] Validasi visual mobile lewat browser setelah halaman lain tersedia.
```

### `/stock-outs/new`

Acceptance criteria:

```txt
[x] Halaman /stock-outs/new tersedia.
[x] GET /fish-types terintegrasi.
[x] FIFO preview tampil dari GET /stocks/fifo?fish_type_id={id}.
[x] POST /stock-outs terintegrasi.
[x] User bisa input fish_type_id, total_weight_kg, destination, out_at, notes.
[x] Response data.items ditampilkan sebagai summary batch yang dipakai.
[x] Error insufficient stock dari backend tampil jelas.
[x] Setelah sukses, user bisa menuju /stocks atau /dashboard.
```

### `/stock-outs`

Acceptance criteria:

```txt
[ ] Halaman /stock-outs tersedia.
[ ] GET /stock-outs terintegrasi.
[ ] Riwayat pengeluaran tampil.
[ ] Items batch yang dikurangi tampil jika tersedia.
[ ] Filter date_from/date_to boleh ditambahkan setelah list dasar stabil.
[ ] Loading state tersedia.
[ ] Error state tersedia.
[ ] Empty state tersedia.
```

## TypeScript Types Minimum

```ts
export type FishQuality = "baik" | "sedang" | "buruk";
export type StockStatus = "available" | "depleted";

export type ApiResponse<T> = {
  success: boolean;
  message: string;
  data?: T;
  errors?: unknown;
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
```

## Urutan Kerja Berikutnya

1. Buat `/dashboard`.
2. Buat `/stock-outs`.
3. Jalankan demo flow end-to-end.

## Definition of Done Frontend MVP

```txt
[ ] /dashboard bisa menampilkan summary dan recent movement dari backend.
[x] /fish-types bisa menambah jenis ikan.
[x] /cold-storages bisa menambah cold storage.
[x] /stocks/new bisa mencatat stok masuk.
[ ] /stocks bisa menampilkan stok FIFO dari backend.
[x] /stock-outs/new bisa mencatat pengeluaran dan menampilkan batch yang dipakai.
[ ] /stock-outs bisa menampilkan riwayat pengeluaran.
[ ] Error backend bisa dibaca user.
[ ] Semua halaman utama punya loading state.
[ ] Semua halaman utama punya empty state.
[ ] Tampilan nyaman di mobile.
```
