# Dummy Data untuk Testing

## Cara Pakai Dokumen Ini

Dummy data ini digunakan untuk tahap development dan testing.

Ada dua format di dokumen ini:
- **Format tabel**: untuk dibaca dan dipahami
- **Format curl/JSON**: langsung dijalankan ke API

---

## Master Data

### 1. Jenis Ikan (Fish Types)

| # | Nama | Deskripsi |
|---|---|---|
| 1 | Tuna | Ikan tuna sirip kuning, tangkapan laut dalam |
| 2 | Tongkol | Ikan tongkol abu-abu, tangkapan pesisir |
| 3 | Cakalang | Ikan cakalang, ukuran sedang |
| 4 | Bandeng | Ikan bandeng tambak, kualitas premium |
| 5 | Kerapu | Ikan kerapu merah, harga tinggi |

**curl — input semua jenis ikan:**
```bash
# Tuna
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Tuna","image_url":"/images/tuna.png","description":"Ikan tuna sirip kuning, tangkapan laut dalam"}'

# Tongkol
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Tongkol","image_url":"/images/tongkol.png","description":"Ikan tongkol abu-abu, tangkapan pesisir"}'

# Cakalang
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Cakalang","image_url":"/images/cakalang.png","description":"Ikan cakalang, ukuran sedang"}'

# Bandeng
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Bandeng","image_url":"/images/bandeng.png","description":"Ikan bandeng tambak, kualitas premium"}'

# Kerapu
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{"name":"Kerapu","image_url":"/images/kerapu.png","description":"Ikan kerapu merah, harga tinggi"}'
```

**Verifikasi:**
```bash
curl http://localhost:8081/api/v1/fish-types
# Harusnya ada 5 item di response.data
```

---

### 2. Cold Storage (Lokasi Penyimpanan)

| # | Nama | Label Lokasi | Keterangan |
|---|---|---|---|
| 1 | Cold Storage A | Zona A - Rak 1 | Penyimpanan utama, kualitas baik |
| 2 | Cold Storage B | Zona B - Rak 2 | Penyimpanan sekunder |
| 3 | Cold Storage C | Zona C - Rak 3 | Penyimpanan tambahan, kualitas sedang |

**curl — input semua cold storage:**
```bash
# Cold Storage A
curl -X POST http://localhost:8081/api/v1/cold-storages \
  -H "Content-Type: application/json" \
  -d '{"name":"Cold Storage A","location_label":"Zona A - Rak 1","description":"Penyimpanan utama, untuk ikan kualitas baik"}'

# Cold Storage B
curl -X POST http://localhost:8081/api/v1/cold-storages \
  -H "Content-Type: application/json" \
  -d '{"name":"Cold Storage B","location_label":"Zona B - Rak 2","description":"Penyimpanan sekunder"}'

# Cold Storage C
curl -X POST http://localhost:8081/api/v1/cold-storages \
  -H "Content-Type: application/json" \
  -d '{"name":"Cold Storage C","location_label":"Zona C - Rak 3","description":"Penyimpanan tambahan, kualitas sedang"}'
```

**Verifikasi:**
```bash
curl http://localhost:8081/api/v1/cold-storages
# Harusnya ada 3 item di response.data
```

---

## Dummy Stok Masuk

Simulasi tiga hari operasional di pelabuhan Makassar.
Setelah input master data, catat ID yang muncul lalu ganti placeholder di bawah.

```
TUNA_ID     = "(isi dari response POST /fish-types)"
TONGKOL_ID  = "(isi dari response POST /fish-types)"
CAKALANG_ID = "(isi dari response POST /fish-types)"
BANDENG_ID  = "(isi dari response POST /fish-types)"
KERAPU_ID   = "(isi dari response POST /fish-types)"
CS_A_ID     = "(isi dari response POST /cold-storages)"
CS_B_ID     = "(isi dari response POST /cold-storages)"
CS_C_ID     = "(isi dari response POST /cold-storages)"
```

---

### Hari 1

Kapal baru bersandar pagi hari. Bongkar muat dimulai jam 07.00.
Ada 4 batch masuk dari berbagai jenis ikan.

| # | Jenis | Kualitas | Berat | Waktu Masuk | Cold Storage | Catatan |
|---|---|---|---|---|---|---|
| S01 | Tuna | baik | 80 kg | 07.30 | CS A | Tangkapan terbaik pagi |
| S02 | Tongkol | baik | 60 kg | 08.00 | CS A | Ukuran besar, segar |
| S03 | Tuna | sedang | 45 kg | 09.15 | CS B | Ada beberapa yang lecet |
| S04 | Cakalang | baik | 35 kg | 10.00 | CS C | Tangkapan siang |

**curl Hari 1:**
```bash
# S01 - Tuna 80kg jam 07.30
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 80,
    "entered_at": "2026-05-01T07:30:00Z",
    "notes": "Tangkapan terbaik pagi, kapal pertama"
  }'

# S02 - Tongkol 60kg jam 08.00
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TONGKOL_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 60,
    "entered_at": "2026-05-01T08:00:00Z",
    "notes": "Ukuran besar, masih sangat segar"
  }'

# S03 - Tuna 45kg jam 09.15
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_B_ID",
    "quality": "sedang",
    "initial_weight_kg": 45,
    "entered_at": "2026-05-01T09:15:00Z",
    "notes": "Ada beberapa yang lecet saat bongkar muat"
  }'

# S04 - Cakalang 35kg jam 10.00
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "CAKALANG_ID",
    "cold_storage_id": "CS_C_ID",
    "quality": "baik",
    "initial_weight_kg": 35,
    "entered_at": "2026-05-01T10:00:00Z",
    "notes": "Tangkapan siang dari kapal kedua"
  }'
```

---

### Hari 2 (Hari Sibuk)

Dua kapal bersandar bersamaan. Volume lebih banyak dari hari sebelumnya.
Masuk 5 batch dengan berbagai kualitas.

| # | Jenis | Kualitas | Berat | Waktu Masuk | Cold Storage | Catatan |
|---|---|---|---|---|---|---|
| S05 | Bandeng | baik | 50 kg | 06.45 | CS A | Kapal paling awal |
| S06 | Tuna | baik | 100 kg | 07.00 | CS A | Batch besar, kualitas premium |
| S07 | Kerapu | baik | 25 kg | 07.30 | CS B | Harga tinggi, prioritas |
| S08 | Tongkol | sedang | 70 kg | 08.30 | CS B | Sedikit lebih tua dari kemarin |
| S09 | Cakalang | buruk | 20 kg | 09.00 | CS C | Kualitas menurun, perlu segera keluar |

**curl Hari 2:**
```bash
# S05 - Bandeng 50kg jam 06.45
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "BANDENG_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 50,
    "entered_at": "2026-05-02T06:45:00Z",
    "notes": "Kapal paling awal, bandeng tambak segar"
  }'

# S06 - Tuna 100kg jam 07.00
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 100,
    "entered_at": "2026-05-02T07:00:00Z",
    "notes": "Batch besar kualitas premium dari kapal besar"
  }'

# S07 - Kerapu 25kg jam 07.30
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "KERAPU_ID",
    "cold_storage_id": "CS_B_ID",
    "quality": "baik",
    "initial_weight_kg": 25,
    "entered_at": "2026-05-02T07:30:00Z",
    "notes": "Harga tinggi, prioritaskan pengeluaran cepat"
  }'

# S08 - Tongkol 70kg jam 08.30
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TONGKOL_ID",
    "cold_storage_id": "CS_B_ID",
    "quality": "sedang",
    "initial_weight_kg": 70,
    "entered_at": "2026-05-02T08:30:00Z",
    "notes": "Sedikit lebih tua dari batch tongkol kemarin"
  }'

# S09 - Cakalang 20kg jam 09.00
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "CAKALANG_ID",
    "cold_storage_id": "CS_C_ID",
    "quality": "buruk",
    "initial_weight_kg": 20,
    "entered_at": "2026-05-02T09:00:00Z",
    "notes": "Kualitas menurun, perlu segera dikeluarkan"
  }'
```

---

### Hari 3

Volume normal. Tiga batch masuk dari kapal reguler.

| # | Jenis | Kualitas | Berat | Waktu Masuk | Cold Storage | Catatan |
|---|---|---|---|---|---|---|
| S10 | Tuna | baik | 60 kg | 07.00 | CS A | Tangkapan rutin |
| S11 | Bandeng | sedang | 40 kg | 08.00 | CS B | Kualitas biasa |
| S12 | Tongkol | baik | 55 kg | 09.30 | CS A | Segar dari kapal pagi |

**curl Hari 3:**
```bash
# S10 - Tuna 60kg
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 60,
    "entered_at": "2026-05-03T07:00:00Z",
    "notes": "Tangkapan rutin hari Rabu"
  }'

# S11 - Bandeng 40kg
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "BANDENG_ID",
    "cold_storage_id": "CS_B_ID",
    "quality": "sedang",
    "initial_weight_kg": 40,
    "entered_at": "2026-05-03T08:00:00Z",
    "notes": "Kualitas biasa, dari tambak reguler"
  }'

# S12 - Tongkol 55kg
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TONGKOL_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 55,
    "entered_at": "2026-05-03T09:30:00Z",
    "notes": "Masih segar dari kapal pagi"
  }'
```

---

## Dummy Stok Keluar

Simulasikan pengeluaran ikan ke berbagai tujuan.
Jalankan ini setelah skenario semua stok masuk dijalankan.

### Pengeluaran Hari 1

| # | Jenis | Berat Keluar | Tujuan | Waktu | Catatan |
|---|---|---|---|---|---|
| O01 | Tuna | 30 kg | Restoran Laut Makassar | 12.00 | Pesanan makan siang |
| O02 | Tongkol | 20 kg | Warung Daeng Baso | 13.00 | Pesanan reguler harian |

**curl Pengeluaran Hari 1:**
```bash
# O01 - Tuna 30kg ke restoran
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 30,
    "destination": "Restoran Laut Makassar",
    "out_at": "2026-05-01T12:00:00Z",
    "notes": "Pesanan makan siang, butuh cepat"
  }'

# O02 - Tongkol 20kg ke warung
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TONGKOL_ID",
    "total_weight_kg": 20,
    "destination": "Warung Daeng Baso",
    "out_at": "2026-05-01T13:00:00Z",
    "notes": "Pesanan reguler harian"
  }'
```

---

### Pengeluaran Hari 2

| # | Jenis | Berat Keluar | Tujuan | Waktu | Catatan |
|---|---|---|---|---|---|
| O03 | Tuna | 90 kg | Pabrik Pengolahan Ikan A | 10.00 | Pesanan besar |
| O04 | Cakalang | 35 kg | Pasar Sentral Makassar | 11.00 | Semua cakalang hari 1 habis |
| O05 | Kerapu | 15 kg | Hotel Grand Makassar | 14.00 | Pesanan hotel premium |

**curl Pengeluaran Hari 2:**
```bash
# O03 - Tuna 90kg ke pabrik
# Catatan: ini akan kena S01 (80kg) + sebagian S03 (10kg dari 45kg)
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 90,
    "destination": "Pabrik Pengolahan Ikan A",
    "out_at": "2026-05-02T10:00:00Z",
    "notes": "Pesanan besar untuk diolah menjadi produk kaleng"
  }'

# O04 - Cakalang 35kg ke pasar
# Catatan: ini akan menghabiskan S04 (35kg hari 1)
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "CAKALANG_ID",
    "total_weight_kg": 35,
    "destination": "Pasar Sentral Makassar",
    "out_at": "2026-05-02T11:00:00Z",
    "notes": "Semua cakalang hari pertama habis untuk pasar"
  }'

# O05 - Kerapu 15kg ke hotel
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "KERAPU_ID",
    "total_weight_kg": 15,
    "destination": "Hotel Grand Makassar",
    "out_at": "2026-05-02T14:00:00Z",
    "notes": "Pesanan premium untuk restoran hotel"
  }'
```

---

### Pengeluaran Hari 3

| # | Jenis | Berat Keluar | Tujuan | Waktu | Catatan |
|---|---|---|---|---|---|
| O06 | Tongkol | 50 kg | Distributor Sulawesi | 09.00 | Pengiriman antar kota |
| O07 | Bandeng | 30 kg | Warung Makan Bu Rahma | 11.00 | Langganan tetap |

**curl Pengeluaran Hari 3:**
```bash
# O06 - Tongkol 50kg ke distributor
# Catatan: sisa Tongkol hari 1 = 40kg (sudah keluar 20kg dari O02)
# Jadi akan kena S02 sisa 40kg + sebagian S08 10kg
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TONGKOL_ID",
    "total_weight_kg": 50,
    "destination": "Distributor Sulawesi Selatan",
    "out_at": "2026-05-03T09:00:00Z",
    "notes": "Pengiriman rutin antar kota"
  }'

# O07 - Bandeng 30kg ke warung
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "BANDENG_ID",
    "total_weight_kg": 30,
    "destination": "Warung Makan Bu Rahma",
    "out_at": "2026-05-03T11:00:00Z",
    "notes": "Langganan tetap setiap hari Rabu"
  }'
```

---

## Kalkulasi Stok yang Diharapkan

Setelah semua input di atas dijalankan, ini kondisi stok yang harusnya ada di sistem.
Gunakan ini untuk validasi manual setelah testing.

### Tuna
| Batch | Masuk | Awal | Keluar | Sisa | Status |
|---|---|---|---|---|---|
| S01 | 01 Mei 07.30 | 80 kg | 80 kg (O03) | 0 kg | depleted |
| S03 | 01 Mei 09.15 | 45 kg | 10 kg (O03) | 35 kg | available |
| S06 | 02 Mei 07.00 | 100 kg | 0 kg | 100 kg | available |
| S10 | 03 Mei 07.00 | 60 kg | 0 kg | 60 kg | available |

**Total Tuna tersisa: 195 kg**
**FIFO order: S03 → S06 → S10**

### Tongkol
| Batch | Masuk | Awal | Keluar | Sisa | Status |
|---|---|---|---|---|---|
| S02 | 01 Mei 08.00 | 60 kg | 60 kg (O02+O06) | 0 kg | depleted |
| S08 | 02 Mei 08.30 | 70 kg | 10 kg (O06) | 60 kg | available |
| S12 | 03 Mei 09.30 | 55 kg | 0 kg | 55 kg | available |

**Total Tongkol tersisa: 115 kg**
**FIFO order: S08 → S12**

### Cakalang
| Batch | Masuk | Awal | Keluar | Sisa | Status |
|---|---|---|---|---|---|
| S04 | 01 Mei 10.00 | 35 kg | 35 kg (O04) | 0 kg | depleted |
| S09 | 02 Mei 09.00 | 20 kg | 0 kg | 20 kg | available |

**Total Cakalang tersisa: 20 kg**
**FIFO order: S09 (satu-satunya)**

### Bandeng
| Batch | Masuk | Awal | Keluar | Sisa | Status |
|---|---|---|---|---|---|
| S05 | 02 Mei 06.45 | 50 kg | 30 kg (O07) | 20 kg | available |
| S11 | 03 Mei 08.00 | 40 kg | 0 kg | 40 kg | available |

**Total Bandeng tersisa: 60 kg**
**FIFO order: S05 → S11**

### Kerapu
| Batch | Masuk | Awal | Keluar | Sisa | Status |
|---|---|---|---|---|---|
| S07 | 02 Mei 07.30 | 25 kg | 15 kg (O05) | 10 kg | available |

**Total Kerapu tersisa: 10 kg**

---

### Ringkasan Dashboard yang Diharapkan

Setelah semua data diinput, cek endpoint ini:
```bash
curl http://localhost:8081/api/v1/dashboard/summary
```

Yang harusnya muncul:
```json
{
  "data": {
    "total_available_weight_kg": 400,
    "total_stock_batches": 12,
    "total_available_batches": 8,
    "total_depleted_batches": 4,
    "fish_type_summary": [
      { "fish_type_name": "Tuna",    "available_weight_kg": 195 },
      { "fish_type_name": "Tongkol", "available_weight_kg": 115 },
      { "fish_type_name": "Bandeng", "available_weight_kg": 60  },
      { "fish_type_name": "Cakalang","available_weight_kg": 20  },
      { "fish_type_name": "Kerapu",  "available_weight_kg": 10  }
    ]
  }
}
```

---

## Mock Data untuk Frontend (Tanpa API)

Kalau API belum siap, data ini bisa langsung digunakan di kode frontend
sebagai placeholder sementara.

```typescript
// mock-data.ts — taruh di src/lib/mock-data.ts

export const mockFishTypes = [
  { id: "fish-001", name: "Tuna",    image_url: "/images/tuna.png" },
  { id: "fish-002", name: "Tongkol", image_url: "/images/tongkol.png" },
  { id: "fish-003", name: "Cakalang",image_url: "/images/cakalang.png" },
  { id: "fish-004", name: "Bandeng", image_url: "/images/bandeng.png" },
  { id: "fish-005", name: "Kerapu",  image_url: "/images/kerapu.png" },
]

export const mockColdStorages = [
  { id: "cs-001", name: "Cold Storage A", location_label: "Zona A - Rak 1" },
  { id: "cs-002", name: "Cold Storage B", location_label: "Zona B - Rak 2" },
  { id: "cs-003", name: "Cold Storage C", location_label: "Zona C - Rak 3" },
]

export const mockFifoStocks = [
  {
    id: "stock-001",
    fish_type_name: "Tuna",
    quality: "sedang",
    remaining_weight_kg: 35,
    entered_at: "2026-05-01T09:15:00Z",
    cold_storage_name: "Cold Storage B",
    location_label: "Zona B - Rak 2",
    fifo_rank: 1,
  },
  {
    id: "stock-002",
    fish_type_name: "Tuna",
    quality: "baik",
    remaining_weight_kg: 100,
    entered_at: "2026-05-02T07:00:00Z",
    cold_storage_name: "Cold Storage A",
    location_label: "Zona A - Rak 1",
    fifo_rank: 2,
  },
  {
    id: "stock-003",
    fish_type_name: "Tuna",
    quality: "baik",
    remaining_weight_kg: 60,
    entered_at: "2026-05-03T07:00:00Z",
    cold_storage_name: "Cold Storage A",
    location_label: "Zona A - Rak 1",
    fifo_rank: 3,
  },
  {
    id: "stock-004",
    fish_type_name: "Tongkol",
    quality: "sedang",
    remaining_weight_kg: 60,
    entered_at: "2026-05-02T08:30:00Z",
    cold_storage_name: "Cold Storage B",
    location_label: "Zona B - Rak 2",
    fifo_rank: 1,
  },
  {
    id: "stock-005",
    fish_type_name: "Tongkol",
    quality: "baik",
    remaining_weight_kg: 55,
    entered_at: "2026-05-03T09:30:00Z",
    cold_storage_name: "Cold Storage A",
    location_label: "Zona A - Rak 1",
    fifo_rank: 2,
  },
  {
    id: "stock-006",
    fish_type_name: "Cakalang",
    quality: "buruk",
    remaining_weight_kg: 20,
    entered_at: "2026-05-02T09:00:00Z",
    cold_storage_name: "Cold Storage C",
    location_label: "Zona C - Rak 3",
    fifo_rank: 1,
  },
  {
    id: "stock-007",
    fish_type_name: "Bandeng",
    quality: "baik",
    remaining_weight_kg: 20,
    entered_at: "2026-05-02T06:45:00Z",
    cold_storage_name: "Cold Storage A",
    location_label: "Zona A - Rak 1",
    fifo_rank: 1,
  },
  {
    id: "stock-008",
    fish_type_name: "Bandeng",
    quality: "sedang",
    remaining_weight_kg: 40,
    entered_at: "2026-05-03T08:00:00Z",
    cold_storage_name: "Cold Storage B",
    location_label: "Zona B - Rak 2",
    fifo_rank: 2,
  },
  {
    id: "stock-009",
    fish_type_name: "Kerapu",
    quality: "baik",
    remaining_weight_kg: 10,
    entered_at: "2026-05-02T07:30:00Z",
    cold_storage_name: "Cold Storage B",
    location_label: "Zona B - Rak 2",
    fifo_rank: 1,
  },
]

export const mockDashboardSummary = {
  total_available_weight_kg: 400,
  total_stock_batches: 12,
  total_available_batches: 8,
  total_depleted_batches: 4,
  today_stock_in_kg: 155,
  today_stock_out_kg: 80,
  fish_type_summary: [
    { fish_type_name: "Tuna",     available_weight_kg: 195 },
    { fish_type_name: "Tongkol",  available_weight_kg: 115 },
    { fish_type_name: "Bandeng",  available_weight_kg: 60  },
    { fish_type_name: "Cakalang", available_weight_kg: 20  },
    { fish_type_name: "Kerapu",   available_weight_kg: 10  },
  ],
}

export const mockRecentMovements = [
  {
    id: "mov-001",
    movement_type: "out",
    fish_type_name: "Bandeng",
    weight_kg: 30,
    description: "Bandeng 30 kg keluar ke Warung Makan Bu Rahma",
    created_at: "2026-05-03T11:00:00Z",
  },
  {
    id: "mov-002",
    movement_type: "out",
    fish_type_name: "Tongkol",
    weight_kg: 50,
    description: "Tongkol 50 kg keluar ke Distributor Sulawesi Selatan",
    created_at: "2026-05-03T09:00:00Z",
  },
  {
    id: "mov-003",
    movement_type: "in",
    fish_type_name: "Tongkol",
    weight_kg: 55,
    description: "Tongkol 55 kg masuk ke Cold Storage A",
    created_at: "2026-05-03T09:30:00Z",
  },
  {
    id: "mov-004",
    movement_type: "in",
    fish_type_name: "Tuna",
    weight_kg: 60,
    description: "Tuna 60 kg masuk ke Cold Storage A",
    created_at: "2026-05-03T07:00:00Z",
  },
]
```