# Manual QA Checklist | Ikan't Setop Us

Update terakhir: 2026-05-06.

Checklist ini dipakai setelah semua halaman MVP tersedia.

## Environment

```txt
[ ] Database Docker hidup.
[ ] Backend hidup di http://localhost:8081/api/v1/health.
[ ] Frontend hidup di http://localhost:3000.
[ ] Browser tidak menampilkan console error fatal.
```

## Flow Utama

```txt
[ ] /dashboard terbuka.
[ ] /fish-types bisa menambah jenis ikan demo.
[ ] /cold-storages bisa menambah dua cold storage demo.
[ ] /stocks/new bisa input batch 50 kg.
[ ] /stocks/new bisa input batch 30 kg dengan waktu lebih baru.
[ ] /stocks menampilkan FIFO rank 1 untuk batch 50 kg.
[ ] /stock-outs/new bisa mencatat pengeluaran 60 kg.
[ ] Summary pengeluaran menampilkan dua item batch: 50 kg dan 10 kg.
[ ] /stocks setelah pengeluaran menampilkan sisa 20 kg.
[ ] /stock-outs menampilkan riwayat pengeluaran demo.
[ ] /dashboard total stok dan recent movements berubah.
```

## State Halaman

```txt
[ ] Semua halaman utama punya loading state.
[ ] Semua halaman utama punya error state saat backend mati.
[ ] Semua halaman utama punya empty state saat data kosong.
[ ] Error insufficient stock tampil jelas di /stock-outs/new.
[ ] Filter jenis ikan di /stocks bekerja.
[ ] Filter jenis ikan, tujuan, date_from, dan date_to di /stock-outs bekerja.
```

## Validasi Mobile

Viewport minimum yang perlu dicek:

```txt
[ ] 360 x 800.
[ ] 390 x 844.
[ ] 430 x 932.
```

Checklist visual:

```txt
[ ] Bottom nav tidak menutup tombol utama.
[ ] Teks tombol tidak overflow.
[ ] Kartu stok dan riwayat tidak melebar keluar layar.
[ ] Form input nyaman dipakai dengan keyboard mobile.
[ ] Halaman bisa discroll sampai konten terakhir.
```

## Hasil QA

```txt
Tanggal:
Tester:
Browser:
Device/viewport:

Hasil akhir: PASS / FAIL
Catatan bug:
-
```
