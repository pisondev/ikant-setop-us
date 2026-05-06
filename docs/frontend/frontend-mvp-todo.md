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
| Frontend foundation | Belum | Belum ada `src/lib/api.ts`, `src/types/api.ts`, app shell, atau bottom nav. |
| `/` | Belum | Masih template default Next.js. |
| `/stocks` | Sebagian | Ada halaman FIFO list, sudah fetch API dan fallback mock, tetapi belum memakai helper/type bersama. |
| Halaman frontend lain | Belum | Dashboard, master data, stock-in, stock-out, dan riwayat belum ada. |

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
| `/` | High | Belum | Tidak ada |
| `/dashboard` | High | Belum | `GET /dashboard/summary`, `GET /dashboard/recent-movements` |
| `/fish-types` | Medium | Belum | `GET /fish-types`, `POST /fish-types` |
| `/cold-storages` | Medium | Belum | `GET /cold-storages`, `POST /cold-storages` |
| `/stocks/new` | High | Belum | `GET /fish-types`, `GET /cold-storages`, `POST /stocks` |
| `/stocks` | High | Sebagian | `GET /stocks/fifo`, `GET /fish-types` |
| `/stock-outs/new` | High | Belum | `GET /fish-types`, `GET /stocks/fifo?fish_type_id={id}`, `POST /stock-outs` |
| `/stock-outs` | Medium | Belum | `GET /stock-outs` |

## Checklist Foundation

```txt
[ ] Buat API helper di apps/web/src/lib/api.ts atau apps/web/lib/api.ts.
[ ] Buat TypeScript type API bersama.
[ ] Setup NEXT_PUBLIC_API_BASE_URL di apps/web/.env.local.
[ ] Buat layout mobile-first reusable.
[ ] Buat bottom navigation.
[ ] Buat page header reusable.
[ ] Ubah / dari template Next.js menjadi redirect/link ke dashboard.
[ ] Hapus fallback mock dari /stocks setelah backend siap dipakai di dev.
[ ] Pastikan loading, error, dan empty state punya pola yang konsisten.
```

## Checklist Per Halaman

### `/`

Acceptance criteria:

```txt
[ ] User membuka / dan langsung diarahkan ke /dashboard atau melihat tombol masuk dashboard.
[ ] Tidak ada konten template Next.js tersisa.
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
[ ] Halaman /fish-types tersedia.
[ ] GET /fish-types terintegrasi.
[ ] POST /fish-types terintegrasi.
[ ] Form field name, image_url, description tersedia.
[ ] Data baru muncul setelah submit.
[ ] Loading state tersedia.
[ ] Error state tersedia.
[ ] Empty state tersedia.
```

### `/cold-storages`

Acceptance criteria:

```txt
[ ] Halaman /cold-storages tersedia.
[ ] GET /cold-storages terintegrasi.
[ ] POST /cold-storages terintegrasi.
[ ] Form field name, location_label, description tersedia.
[ ] Data baru muncul setelah submit.
[ ] Loading state tersedia.
[ ] Error state tersedia.
[ ] Empty state tersedia.
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
[ ] Halaman /stocks/new tersedia.
[ ] GET /fish-types terintegrasi.
[ ] GET /cold-storages terintegrasi.
[ ] POST /stocks terintegrasi.
[ ] User bisa memilih jenis ikan.
[ ] User bisa memilih kualitas ikan.
[ ] User bisa input berat ikan.
[ ] User bisa memilih waktu masuk.
[ ] User bisa memilih cold storage.
[ ] Jika sukses, user diarahkan ke /stocks atau mendapat tombol lihat stok.
[ ] Error validasi backend tampil.
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
[x] Fallback mock tersedia ketika API gagal.
```

Yang masih perlu:

```txt
[ ] Gunakan API helper bersama.
[ ] Gunakan TypeScript type bersama.
[ ] Perbaiki tulisan/karakter rusak pada empty state.
[ ] Tambahkan error state eksplisit, bukan langsung fallback diam-diam.
[ ] Tambahkan tombol menuju /stocks/new dan /stock-outs/new.
[ ] Validasi tampilan mobile setelah layout global tersedia.
[ ] Hapus mock fallback ketika backend menjadi dependency wajib.
```

### `/stock-outs/new`

Acceptance criteria:

```txt
[ ] Halaman /stock-outs/new tersedia.
[ ] GET /fish-types terintegrasi.
[ ] FIFO preview tampil dari GET /stocks/fifo?fish_type_id={id}.
[ ] POST /stock-outs terintegrasi.
[ ] User bisa input fish_type_id, total_weight_kg, destination, out_at, notes.
[ ] Response data.items ditampilkan sebagai summary batch yang dipakai.
[ ] Error insufficient stock dari backend tampil jelas.
[ ] Setelah sukses, user bisa menuju /stocks, /dashboard, atau input lagi.
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

1. Selesaikan frontend foundation: API helper, type bersama, layout mobile.
2. Ganti `/` dari template Next.js.
3. Rapikan `/stocks` agar memakai helper/type dan error state yang jelas.
4. Buat `/fish-types` dan `/cold-storages` supaya data master bisa dibuat dari UI.
5. Buat `/stocks/new`.
6. Buat `/stock-outs/new`.
7. Buat `/dashboard`.
8. Buat `/stock-outs`.
9. Jalankan demo flow end-to-end.

## Definition of Done Frontend MVP

```txt
[ ] /dashboard bisa menampilkan summary dan recent movement dari backend.
[ ] /fish-types bisa menambah jenis ikan.
[ ] /cold-storages bisa menambah cold storage.
[ ] /stocks/new bisa mencatat stok masuk.
[ ] /stocks bisa menampilkan stok FIFO dari backend.
[ ] /stock-outs/new bisa mencatat pengeluaran dan menampilkan batch yang dipakai.
[ ] /stock-outs bisa menampilkan riwayat pengeluaran.
[ ] Error backend bisa dibaca user.
[ ] Semua halaman utama punya loading state.
[ ] Semua halaman utama punya empty state.
[ ] Tampilan nyaman di mobile.
```
