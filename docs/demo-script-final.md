# Demo Script Final | Ikan't Setop Us

Update terakhir: 2026-05-06.

Dokumen ini dipakai saat presentasi MVP. Flow utama yang harus ditunjukkan:

```txt
input stok masuk -> data tersimpan -> stok tampil FIFO -> stok keluar dicatat -> dashboard berubah
```

## Persiapan

Jalankan service lokal:

```bash
docker compose up -d
cd apps/api
go run ./cmd/api
```

Di terminal lain:

```bash
cd apps/web
npm run dev
```

URL:

```txt
Frontend: http://localhost:3000
Backend : http://localhost:8081/api/v1
```

## Branch

Branch untuk tahap demo dan QA:

```bash
git switch -c feat/adit-docs-demo
```

Jika branch sudah ada:

```bash
git switch feat/adit-docs-demo
```

## Data Demo

Gunakan nama unik jika database sudah pernah dipakai:

```txt
Jenis ikan       : Tuna Demo <tanggal/jam>
Cold Storage A   : Cold Storage Demo A
Cold Storage B   : Cold Storage Demo B
Batch pertama    : 50 kg, kualitas baik, masuk 08.00
Batch kedua      : 30 kg, kualitas sedang, masuk 09.00
Pengeluaran      : 60 kg ke Restoran Laut Makassar
Ekspektasi akhir : batch pertama habis, batch kedua tersisa 20 kg
```

## Alur Presentasi

1. Buka `/dashboard`.
2. Tunjukkan angka stok awal dan recent movements.
3. Buka `/fish-types`, tambah jenis ikan demo jika belum ada.
4. Buka `/cold-storages`, tambah Cold Storage A dan B jika belum ada.
5. Buka `/stocks/new`, input batch pertama 50 kg.
6. Input batch kedua 30 kg dengan waktu masuk lebih baru.
7. Buka `/stocks`, filter jenis ikan demo.
8. Pastikan FIFO menampilkan batch 50 kg sebagai rank 1 dan batch 30 kg sebagai rank 2.
9. Buka `/stock-outs/new`, catat pengeluaran 60 kg.
10. Pastikan summary batch terpakai menampilkan 50 kg + 10 kg.
11. Buka `/stocks`, pastikan sisa stok tinggal 20 kg.
12. Buka `/stock-outs`, pastikan riwayat pengeluaran muncul.
13. Buka `/dashboard`, pastikan total stok dan recent movements berubah.

## Hasil Validasi API Terakhir

Run ID:

```txt
20260506154516
```

Hasil:

```txt
[x] Master jenis ikan berhasil dibuat.
[x] Cold storage A dan B berhasil dibuat.
[x] Dua batch stok masuk berhasil dibuat.
[x] FIFO sebelum pengeluaran: 1:50kg, 2:30kg.
[x] Stock-out 60 kg memakai dua batch: 50kg + 10kg.
[x] FIFO setelah pengeluaran: 1:20kg.
[x] Dashboard total available berubah net +20 kg.
[x] GET /stock-outs?fish_type_id=... mengembalikan transaksi demo.
```

Catatan:

```txt
Validasi ini dilakukan lewat API lokal. Validasi klik browser dan tampilan mobile masih perlu dijalankan.
```
