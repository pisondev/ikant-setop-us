# Skenario Testing FIFO

## Tools yang Dipakai untuk Testing

Ada dua pilihan:

### Pilihan A: Postman
1. Download Postman di https://www.postman.com/downloads/
2. Buka Postman → New Request
3. Masukkan URL, pilih method (GET/POST), isi body jika diperlukan
4. Klik Send, lihat response di bawah

### Pilihan B: Terminal dengan `curl`
1. Buka terminal di VS Code → `Ctrl + `` ` ``
2. Ketik command `curl` yang ada di setiap skenario
3. Lihat output JSON-nya langsung

### Setup Awal Postman
1. Buka Postman → Environments → New Environment
2. Tambah variable:
   - `base_url` = `http://localhost:8081/api/v1`
3. Sekarang bisa pakai `{{base_url}}` di setiap request

---

## Cara Baca Dokumen Ini

Setiap skenario punya:
- **Setup**: data awal yang perlu dimasukkan
- **Command**: sintaks curl + contoh Postman
- **Expected Result**: hasil yang harusnya muncul
- **Checklist Validasi**: apa yang perlu dicek dari response JSON
- **Hasil**: (bisa diisi PASS atau FAIL)

---

## Setup Master Data

Jalankan ini sekali di awal sebelum semua skenario.

### Buat Jenis Ikan

**Postman:**
```
Method : POST
URL    : {{base_url}}/fish-types
Body   : raw → JSON
```

**curl Tuna:**
```bash
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tuna",
    "image_url": "/images/tuna.png",
    "description": "Ikan tuna segar"
  }'
```

**curl Tongkol:**
```bash
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tongkol",
    "image_url": "/images/tongkol.png",
    "description": "Ikan tongkol segar"
  }'
```

**curl Cakalang:**
```bash
curl -X POST http://localhost:8081/api/v1/fish-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cakalang",
    "image_url": "/images/cakalang.png",
    "description": "Ikan cakalang segar"
  }'
```

**Expected response:**
```json
{
  "success": true,
  "data": {
    "id": "xxxx-xxxx-xxxx",
    "name": "Tuna"
  }
}
```

**Cek semua jenis ikan berhasil dibuat:**
```bash
curl http://localhost:8081/api/v1/fish-types
```

---

### Buat Cold Storage

**curl Cold Storage A:**
```bash
curl -X POST http://localhost:8081/api/v1/cold-storages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cold Storage A",
    "location_label": "Zona A - Rak 1",
    "description": "Penyimpanan utama"
  }'
```

**curl Cold Storage B:**
```bash
curl -X POST http://localhost:8081/api/v1/cold-storages \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cold Storage B",
    "location_label": "Zona B - Rak 2",
    "description": "Penyimpanan sekunder"
  }'
```

**Cek berhasil:**
```bash
curl http://localhost:8081/api/v1/cold-storages
```

> Simpan semua ID yang muncul di response. Contoh:
> ```
> TUNA_ID     = "7f9d8c5a-4a3b-4f2e-8a21-111111111111"
> TONGKOL_ID  = "7f9d8c5a-4a3b-4f2e-8a21-222222222222"
> CS_A_ID     = "16dfdc88-831f-4e28-b2c7-333333333333"
> CS_B_ID     = "16dfdc88-831f-4e28-b2c7-444444444444"
> ```

---

## Skenario Dasar FIFO

---

### Skenario FIFO Ranking Tampil Benar

**Tujuan:** Stok yang masuk lebih dulu dapat ranking FIFO lebih tinggi.

#### Input 3 batch Tuna berurutan

**Batch A (08.00):**
```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "GANTI_DENGAN_TUNA_ID",
    "cold_storage_id": "GANTI_DENGAN_CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 50,
    "entered_at": "2026-05-01T08:00:00Z",
    "notes": "Batch A - tangkapan pagi"
  }'
```

**Batch B (09.00):**
```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "GANTI_DENGAN_TUNA_ID",
    "cold_storage_id": "GANTI_DENGAN_CS_B_ID",
    "quality": "baik",
    "initial_weight_kg": 30,
    "entered_at": "2026-05-01T09:00:00Z",
    "notes": "Batch B - tangkapan pagi"
  }'
```

**Batch C (10.00):**
```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "GANTI_DENGAN_TUNA_ID",
    "cold_storage_id": "GANTI_DENGAN_CS_A_ID",
    "quality": "sedang",
    "initial_weight_kg": 20,
    "entered_at": "2026-05-01T10:00:00Z",
    "notes": "Batch C - tangkapan siang"
  }'
```

#### Cek FIFO

```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=GANTI_DENGAN_TUNA_ID"
```

#### Expected Response
```json
{
  "success": true,
  "data": [
    {
      "fifo_rank": 1,
      "remaining_weight_kg": 50,
      "entered_at": "2026-05-01T08:00:00Z",
      "fish_type_name": "Tuna"
    },
    {
      "fifo_rank": 2,
      "remaining_weight_kg": 30,
      "entered_at": "2026-05-01T09:00:00Z"
    },
    {
      "fifo_rank": 3,
      "remaining_weight_kg": 20,
      "entered_at": "2026-05-01T10:00:00Z"
    }
  ]
}
```

#### Checklist Validasi
```
[ ] response.success = true
[ ] response.data.length = 3
[ ] data[0].fifo_rank = 1
[ ] data[0].entered_at = "...08:00:00Z" (paling lama)
[ ] data[1].fifo_rank = 2
[ ] data[2].fifo_rank = 3
[ ] Tidak ada data dari jenis ikan lain
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Filter Jenis Ikan Tidak Mencampur Data

**Tujuan:** Filter Tuna hanya tampil Tuna, filter Tongkol hanya tampil Tongkol.

#### Tambah 1 batch Tongkol (masuk jam 07.00, lebih awal dari Tuna)

```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "GANTI_DENGAN_TONGKOL_ID",
    "cold_storage_id": "GANTI_DENGAN_CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 40,
    "entered_at": "2026-05-01T07:00:00Z",
    "notes": "Tongkol pagi sekali"
  }'
```

#### Cek filter Tuna
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=GANTI_DENGAN_TUNA_ID"
```

#### Cek filter Tongkol
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=GANTI_DENGAN_TONGKOL_ID"
```

#### Checklist Validasi
```
Filter Tuna:
[ ] response.data.length = 3
[ ] Semua item fish_type_name = "Tuna"
[ ] Tidak ada item dengan fish_type_name = "Tongkol"

Filter Tongkol:
[ ] response.data.length = 1
[ ] data[0].fish_type_name = "Tongkol"
[ ] data[0].fifo_rank = 1
[ ] data[0].remaining_weight_kg = 40
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Tanpa Filter Tampilkan Semua

```bash
curl http://localhost:8081/api/v1/stocks/fifo
```

#### Checklist Validasi
```
[ ] response.data.length = 4 (3 Tuna + 1 Tongkol)
[ ] data[0] adalah Tongkol (entered_at 07.00, paling lama)
[ ] data[1] adalah Tuna 08.00
[ ] data[2] adalah Tuna 09.00
[ ] data[3] adalah Tuna 10.00
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

## Skenario Pengeluaran Stok

---

### Skenario Keluar Pas Satu Batch

**Input 1 batch:**
```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 50,
    "entered_at": "2026-05-01T08:00:00Z",
    "notes": "Satu-satunya batch"
  }'
```

**Keluarkan tepat 50 kg:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 50,
    "destination": "Restoran Laut Makassar",
    "out_at": "2026-05-01T12:00:00Z",
    "notes": "Test skenario 2.1"
  }'
```

**Cek stok setelah keluar:**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
Response stock-out:
[ ] success = true
[ ] data.items.length = 1
[ ] data.items[0].weight_kg = 50

Verifikasi stok:
[ ] data = [] (kosong) ATAU data[0].status = "depleted"
[ ] data[0].remaining_weight_kg = 0 (jika masih tampil)
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Keluar Sebagian Batch

**Setup:** Input Tuna 50 kg

**Keluarkan 20 kg saja:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 20,
    "destination": "Warung Seafood Pak Baso",
    "out_at": "2026-05-01T12:00:00Z",
    "notes": "Test skenario 2.2"
  }'
```

**Verifikasi:**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
[ ] success = true
[ ] data.items[0].weight_kg = 20

Verifikasi stok:
[ ] data[0].remaining_weight_kg = 30   <- 50 - 20 = 30
[ ] data[0].status = "available"        <- belum habis
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Keluar Melewati Dua Batch

**Input 2 batch:**
```bash
# Batch A — masuk lebih dulu
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 50,
    "entered_at": "2026-05-01T08:00:00Z",
    "notes": "Batch A"
  }'

# Batch B — masuk belakangan
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_B_ID",
    "quality": "baik",
    "initial_weight_kg": 30,
    "entered_at": "2026-05-01T09:00:00Z",
    "notes": "Batch B"
  }'
```

**Keluarkan 60 kg (lebih dari batch A):**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 60,
    "destination": "Restoran Daeng Syamsul",
    "out_at": "2026-05-01T14:00:00Z",
    "notes": "Test skenario 2.3 - melewati dua batch"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "total_weight_kg": 60,
    "items": [
      { "weight_kg": 50 },
      { "weight_kg": 10 }
    ]
  }
}
```

**Verifikasi:**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
Response stock-out:
[ ] success = true
[ ] data.items.length = 2          <- kena 2 batch
[ ] data.items[0].weight_kg = 50   <- batch A habis semua
[ ] data.items[1].weight_kg = 10   <- batch B dikurangi 10

Verifikasi stok:
[ ] hanya 1 data yang muncul (batch A sudah depleted)
[ ] data[0].remaining_weight_kg = 20   <- 30 - 10 = 20
[ ] data[0].status = "available"
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Keluar Melewati Tiga Batch

**Input 3 batch:**
```bash
# Batch A
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_A_ID","quality":"baik","initial_weight_kg":20,"entered_at":"2026-05-01T07:00:00Z","notes":"Batch A"}'

# Batch B
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_A_ID","quality":"baik","initial_weight_kg":30,"entered_at":"2026-05-01T08:00:00Z","notes":"Batch B"}'

# Batch C
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_B_ID","quality":"sedang","initial_weight_kg":40,"entered_at":"2026-05-01T09:00:00Z","notes":"Batch C"}'
```

**Keluarkan 60 kg:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 60,
    "destination": "Test Tiga Batch",
    "out_at": "2026-05-01T15:00:00Z",
    "notes": "Test skenario 2.4"
  }'
```

**Checklist Validasi:**
```
[ ] data.items.length = 3
[ ] items[0].weight_kg = 20   <- Batch A habis
[ ] items[1].weight_kg = 30   <- Batch B habis
[ ] items[2].weight_kg = 10   <- Batch C dikurangi 10

Verifikasi stok:
[ ] hanya Batch C yang muncul
[ ] Batch C remaining_weight_kg = 30   <- 40 - 10 = 30
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

## Edge Case

---

### Skenario Stok Tidak Cukup

**Setup:** Ada Tuna total 50 kg.

**Minta 999 kg:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 999,
    "destination": "Test Insufficient",
    "out_at": "2026-05-01T16:00:00Z",
    "notes": "Ini harusnya gagal"
  }'
```

**Expected Response:**
```json
{
  "success": false,
  "message": "Insufficient stock",
  "errors": [
    {
      "field": "total_weight_kg",
      "message": "Requested 999 kg, but only 50 kg is available"
    }
  ]
}
```

**Verifikasi (stok seharusnya tidak berubah):**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
Response:
[ ] success = false
[ ] message mengandung kata "Insufficient" atau "stock"
[ ] ada field errors
[ ] HTTP status code = 400 (cek di Postman bagian atas response)

Verifikasi stok (tidak berubah):
[ ] total remaining masih 50 kg
[ ] status semua batch masih "available"
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Empty State Saat Stok Habis

**Setup:** Habiskan semua stok Tuna, lalu jalankan:

```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
API:
[ ] success = true
[ ] data = []  (array kosong, bukan null, bukan error)

Frontend (buka /stocks filter Tuna):
[ ] Halaman tidak crash / tidak error
[ ] Ada teks empty state (misal: "Tidak ada stok tersedia")
[ ] Tidak ada layar putih kosong tanpa keterangan
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Dua Jenis Ikan Tidak Saling Mempengaruhi

**Setup:**
```bash
# Input Tuna 50 kg
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_A_ID","quality":"baik","initial_weight_kg":50,"entered_at":"2026-05-01T08:00:00Z","notes":"Tuna"}'

# Input Tongkol 40 kg
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TONGKOL_ID","cold_storage_id":"CS_A_ID","quality":"baik","initial_weight_kg":40,"entered_at":"2026-05-01T08:00:00Z","notes":"Tongkol"}'
```

**Keluarkan semua Tuna:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 50,
    "destination": "Test Isolasi Jenis",
    "out_at": "2026-05-01T12:00:00Z",
    "notes": "Tuna habis, Tongkol harusnya aman"
  }'
```

**Verifikasi (ikan lain harus tetap utuh):**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TONGKOL_ID"
```

**Checklist Validasi:**
```
[ ] Tongkol remaining_weight_kg = 40  <- tidak berubah
[ ] Tongkol status = "available"
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario FIFO Berdasarkan Waktu, Bukan Lokasi Cold Storage

**Setup:**
```bash
# Batch ini masuk lebih dulu (08.00) tapi di CS A
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_A_ID","quality":"baik","initial_weight_kg":50,"entered_at":"2026-05-01T08:00:00Z","notes":"Masuk lebih dulu"}'

# Batch ini masuk belakangan (09.00) tapi di CS B
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","cold_storage_id":"CS_B_ID","quality":"baik","initial_weight_kg":30,"entered_at":"2026-05-01T09:00:00Z","notes":"Masuk belakangan"}'
```

**Aksi:**
```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

**Checklist Validasi:**
```
[ ] data[0].fifo_rank = 1 -> yang masuk jam 08.00 (CS A)
[ ] data[1].fifo_rank = 2 -> yang masuk jam 09.00 (CS B)
[ ] FIFO tidak diurutkan berdasarkan nama cold storage
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

## Dashboard

---

### Skenario Dashboard Update Setelah Stok Masuk

**Catat angka dashboard sebelum:**
```bash
curl http://localhost:8081/api/v1/dashboard/summary
```
> Tulis di sini → `today_stock_in_kg` sebelum: ______

**Input stok baru:**
```bash
curl -X POST http://localhost:8081/api/v1/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "cold_storage_id": "CS_A_ID",
    "quality": "baik",
    "initial_weight_kg": 50,
    "entered_at": "2026-05-02T08:00:00Z",
    "notes": "Test dashboard"
  }'
```

**Cek dashboard sesudah:**
```bash
curl http://localhost:8081/api/v1/dashboard/summary
```

**Checklist Validasi:**
```
[ ] today_stock_in_kg bertambah 50
[ ] total_available_weight_kg bertambah 50
[ ] total_stock_batches bertambah 1
[ ] fish_type_summary Tuna: available_weight_kg bertambah 50
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Dashboard Update Setelah Stok Keluar

**Setup:** Ada stok Tuna 50 kg.

**Catat dashboard sebelum:**
```bash
curl http://localhost:8081/api/v1/dashboard/summary
```
> `today_stock_out_kg` sebelum: ______

**Keluarkan 30 kg:**
```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{
    "fish_type_id": "TUNA_ID",
    "total_weight_kg": 30,
    "destination": "Test Dashboard",
    "out_at": "2026-05-02T12:00:00Z",
    "notes": "Test dashboard keluar"
  }'
```

**Cek dashboard sesudah:**
```bash
curl http://localhost:8081/api/v1/dashboard/summary
```

**Checklist Validasi:**
```
[ ] today_stock_out_kg bertambah 30
[ ] total_available_weight_kg berkurang 30
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

### Skenario Recent Movements Tercatat

```bash
curl "http://localhost:8081/api/v1/dashboard/recent-movements?limit=5"
```

**Checklist Validasi:**
```
[ ] Ada entry movement_type = "in" untuk setiap stok masuk
[ ] Ada entry movement_type = "out" untuk setiap stok keluar
[ ] Entry terbaru ada di index 0 (paling atas)
[ ] Field fish_type_name terisi, bukan null
[ ] Field weight_kg berupa number, bukan string
```

**Hasil:** PASS / FAIL  
**Catatan:**

---

## Validasi API Contract

### Cek Response `/stocks/fifo`

```bash
curl "http://localhost:8081/api/v1/stocks/fifo?fish_type_id=TUNA_ID"
```

Lihat response, lalu centang:

| Field | Tipe Harusnya | Contoh Benar | Contoh Salah | Hasil |
|---|---|---|---|---|
| `id` | string UUID | `"5cf59487-..."` | `123` | ✅/❌ |
| `fish_type_name` | string | `"Tuna"` | `null` | ✅/❌ |
| `quality` | string enum | `"baik"` | `"BAIK"` atau `1` | ✅/❌ |
| `remaining_weight_kg` | number | `50` | `"50"` | ✅/❌ |
| `entered_at` | ISO string | `"2026-05-01T08:00:00Z"` | `"01-05-2026"` | ✅/❌ |
| `cold_storage_name` | string | `"Cold Storage A"` | `null` | ✅/❌ |
| `location_label` | string | `"Zona A - Rak 1"` | `null` | ✅/❌ |
| `fifo_rank` | number | `1` | `"1"` | ✅/❌ |

### Cek Response `POST /stock-outs`

```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","total_weight_kg":10,"destination":"Test","out_at":"2026-05-02T10:00:00Z"}'
```

| Field | Harusnya | Hasil |
|---|---|---|
| `success` | boolean `true` | ✅/❌ |
| `data.items` | array | ✅/❌ |
| `data.items[].stock_batch_id` | string UUID | ✅/❌ |
| `data.items[].weight_kg` | number | ✅/❌ |

### Cek Response Error Insufficient Stock

```bash
curl -X POST http://localhost:8081/api/v1/stock-outs \
  -H "Content-Type: application/json" \
  -d '{"fish_type_id":"TUNA_ID","total_weight_kg":99999,"destination":"Test Error","out_at":"2026-05-02T10:00:00Z"}'
```

| Field | Harusnya | Hasil |
|---|---|---|
| `success` | boolean `false` | ✅/❌ |
| `message` | string (ada kata "stock") | ✅/❌ |
| `errors` | array | ✅/❌ |
| `errors[0].field` | `"total_weight_kg"` | ✅/❌ |
| HTTP status code | `400` (cek di Postman) | ✅/❌ |

---

## Ringkasan Hasil Testing

| Skenario | Prioritas | Hasil |
|---|---|---|
| 1.1 FIFO ranking tampil benar | HIGH | ✅/❌ |
| 1.2 Filter jenis ikan bekerja | HIGH | ✅/❌ |
| 1.3 Tanpa filter tampil semua | MEDIUM | ✅/❌ |
| 2.1 Keluar pas satu batch | HIGH | ✅/❌ |
| 2.2 Keluar sebagian batch | HIGH | ✅/❌ |
| 2.3 Keluar melewati dua batch | HIGH | ✅/❌ |
| 2.4 Keluar melewati tiga batch | MEDIUM | ✅/❌ |
| 3.1 Insufficient stock ditolak | HIGH | ✅/❌ |
| 3.2 Empty state saat stok habis | MEDIUM | ✅/❌ |
| 3.3 Dua jenis ikan tidak saling pengaruh | HIGH | ✅/❌ |
| 3.4 FIFO berdasarkan waktu bukan lokasi | MEDIUM | ✅/❌ |
| 4.1 Dashboard update stok masuk | MEDIUM | ✅/❌ |
| 4.2 Dashboard update stok keluar | MEDIUM | ✅/❌ |
| 4.3 Recent movements tercatat | LOW | ✅/❌ |
| 5.x Validasi API contract | HIGH | ✅/❌ |

---