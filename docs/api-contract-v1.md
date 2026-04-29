# API Contract V1 | "Ikan't Setop Us"

API contract untuk MVP FishFlow: FIFO Fish Inventory Management System.

Sistem ini digunakan untuk mencatat stok ikan masuk, mengelola lokasi cold storage, menampilkan stok berdasarkan FIFO, mencatat pengeluaran ikan, dan menampilkan dashboard monitoring.

---

## 1. Base URL

### Local Development

```txt
http://localhost:8080/api/v1
```

### Health Check

```http
GET /health
```

Response:

```json
{
  "success": true,
  "message": "API is running",
  "data": {
    "service": "ikant-setop-us-api",
    "version": "v1"
  }
}
```

---

## 2. Response Format

### Success Response

```json
{
  "success": true,
  "message": "Success message",
  "data": {}
}
```

### Error Response

```json
{
  "success": false,
  "message": "Error message",
  "errors": [
    {
      "field": "field_name",
      "message": "Validation error message"
    }
  ]
}
```

---

## 3. Common Rules

- ID menggunakan UUID.
- Waktu menggunakan format ISO 8601.
- Berat ikan menggunakan satuan kilogram.
- `entered_at` adalah waktu ikan pertama kali masuk ke sistem/cold storage.
- FIFO dihitung berdasarkan `entered_at` paling lama.
- Kualitas ikan diperbarui manual oleh admin, tidak otomatis turun berdasarkan waktu.
- Untuk MVP, authentication dapat dibuat opsional atau belum diimplementasikan.

---

## 4. Enums

### Fish Quality

```txt
baik
sedang
buruk
```

### Stock Status

```txt
available
depleted
```

### Stock Movement Type

```txt
in
out
quality_update
location_update
adjustment
```

---

# 5. Fish Types

Master data jenis ikan.

---

## 5.1 Get Fish Types

```http
GET /fish-types
```

Response:

```json
{
  "success": true,
  "message": "Fish types retrieved successfully",
  "data": [
    {
      "id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
      "name": "Tuna",
      "image_url": "/images/fish/tuna.png",
      "description": "Ikan tuna",
      "created_at": "2026-05-01T08:00:00Z",
      "updated_at": "2026-05-01T08:00:00Z"
    }
  ]
}
```

---

## 5.2 Create Fish Type

```http
POST /fish-types
```

Request:

```json
{
  "name": "Tuna",
  "image_url": "/images/fish/tuna.png",
  "description": "Ikan tuna"
}
```

Response:

```json
{
  "success": true,
  "message": "Fish type created successfully",
  "data": {
    "id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
    "name": "Tuna",
    "image_url": "/images/fish/tuna.png",
    "description": "Ikan tuna",
    "created_at": "2026-05-01T08:00:00Z",
    "updated_at": "2026-05-01T08:00:00Z"
  }
}
```

---

# 6. Cold Storages

Master data lokasi penyimpanan ikan.

---

## 6.1 Get Cold Storages

```http
GET /cold-storages
```

Response:

```json
{
  "success": true,
  "message": "Cold storages retrieved successfully",
  "data": [
    {
      "id": "16dfdc88-831f-4e28-b2c7-222222222222",
      "name": "Cold Storage A",
      "location_label": "Zona A - Rak 1",
      "description": "Penyimpanan utama untuk ikan kualitas baik",
      "created_at": "2026-05-01T08:00:00Z",
      "updated_at": "2026-05-01T08:00:00Z"
    }
  ]
}
```

---

## 6.2 Create Cold Storage

```http
POST /cold-storages
```

Request:

```json
{
  "name": "Cold Storage A",
  "location_label": "Zona A - Rak 1",
  "description": "Penyimpanan utama untuk ikan kualitas baik"
}
```

Response:

```json
{
  "success": true,
  "message": "Cold storage created successfully",
  "data": {
    "id": "16dfdc88-831f-4e28-b2c7-222222222222",
    "name": "Cold Storage A",
    "location_label": "Zona A - Rak 1",
    "description": "Penyimpanan utama untuk ikan kualitas baik",
    "created_at": "2026-05-01T08:00:00Z",
    "updated_at": "2026-05-01T08:00:00Z"
  }
}
```

---

# 7. Stocks

Data batch stok ikan yang masuk ke cold storage.

---

## 7.1 Get Stocks

```http
GET /stocks
```

Query params:

| Name | Type | Required | Description |
|---|---|---:|---|
| `fish_type_id` | UUID | No | Filter berdasarkan jenis ikan |
| `quality` | string | No | Filter kualitas: `baik`, `sedang`, `buruk` |
| `cold_storage_id` | UUID | No | Filter lokasi cold storage |
| `status` | string | No | Filter status: `available`, `depleted` |
| `sort` | string | No | Default: `fifo` |

Example:

```http
GET /stocks?fish_type_id=7f9d8c5a-4a3b-4f2e-8a21-111111111111&sort=fifo
```

Response:

```json
{
  "success": true,
  "message": "Stocks retrieved successfully",
  "data": [
    {
      "id": "5cf59487-8233-4dbf-a27e-333333333333",
      "fish_type": {
        "id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
        "name": "Tuna"
      },
      "cold_storage": {
        "id": "16dfdc88-831f-4e28-b2c7-222222222222",
        "name": "Cold Storage A",
        "location_label": "Zona A - Rak 1"
      },
      "quality": "baik",
      "initial_weight_kg": 50,
      "remaining_weight_kg": 50,
      "entered_at": "2026-05-01T08:00:00Z",
      "status": "available",
      "notes": "Tangkapan pagi",
      "created_at": "2026-05-01T08:05:00Z",
      "updated_at": "2026-05-01T08:05:00Z"
    }
  ]
}
```

---

## 7.2 Create Stock Batch

```http
POST /stocks
```

Request:

```json
{
  "fish_type_id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
  "cold_storage_id": "16dfdc88-831f-4e28-b2c7-222222222222",
  "quality": "baik",
  "initial_weight_kg": 50,
  "entered_at": "2026-05-01T08:00:00Z",
  "notes": "Tangkapan pagi"
}
```

Backend behavior:

- Membuat batch stok baru.
- Mengisi `remaining_weight_kg` sama dengan `initial_weight_kg`.
- Mengisi `status` sebagai `available`.
- Membuat histori `stock_movement` dengan type `in`.

Response:

```json
{
  "success": true,
  "message": "Stock batch created successfully",
  "data": {
    "id": "5cf59487-8233-4dbf-a27e-333333333333",
    "fish_type_id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
    "cold_storage_id": "16dfdc88-831f-4e28-b2c7-222222222222",
    "quality": "baik",
    "initial_weight_kg": 50,
    "remaining_weight_kg": 50,
    "entered_at": "2026-05-01T08:00:00Z",
    "status": "available",
    "notes": "Tangkapan pagi",
    "created_at": "2026-05-01T08:05:00Z",
    "updated_at": "2026-05-01T08:05:00Z"
  }
}
```

---

## 7.3 Get Stock Detail

```http
GET /stocks/{id}
```

Response:

```json
{
  "success": true,
  "message": "Stock detail retrieved successfully",
  "data": {
    "id": "5cf59487-8233-4dbf-a27e-333333333333",
    "fish_type": {
      "id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
      "name": "Tuna"
    },
    "cold_storage": {
      "id": "16dfdc88-831f-4e28-b2c7-222222222222",
      "name": "Cold Storage A",
      "location_label": "Zona A - Rak 1"
    },
    "quality": "baik",
    "initial_weight_kg": 50,
    "remaining_weight_kg": 50,
    "entered_at": "2026-05-01T08:00:00Z",
    "status": "available",
    "notes": "Tangkapan pagi",
    "created_at": "2026-05-01T08:05:00Z",
    "updated_at": "2026-05-01T08:05:00Z"
  }
}
```

---

## 7.4 Update Stock Quality

```http
PATCH /stocks/{id}/quality
```

Request:

```json
{
  "quality": "sedang",
  "notes": "Kualitas mata ikan mulai menurun"
}
```

Backend behavior:

- Mengubah kualitas ikan secara manual.
- Membuat histori `stock_movement` dengan type `quality_update`.

Response:

```json
{
  "success": true,
  "message": "Stock quality updated successfully",
  "data": {
    "id": "5cf59487-8233-4dbf-a27e-333333333333",
    "previous_quality": "baik",
    "new_quality": "sedang"
  }
}
```

---

## 7.5 Update Stock Location

```http
PATCH /stocks/{id}/location
```

Request:

```json
{
  "cold_storage_id": "9122572d-8a8b-4898-8e35-444444444444",
  "notes": "Dipindahkan ke Cold Storage B"
}
```

Backend behavior:

- Mengubah lokasi penyimpanan stok.
- Membuat histori `stock_movement` dengan type `location_update`.

Response:

```json
{
  "success": true,
  "message": "Stock location updated successfully",
  "data": {
    "id": "5cf59487-8233-4dbf-a27e-333333333333",
    "previous_cold_storage_id": "16dfdc88-831f-4e28-b2c7-222222222222",
    "new_cold_storage_id": "9122572d-8a8b-4898-8e35-444444444444"
  }
}
```

---

# 8. FIFO Stocks

Endpoint khusus untuk menampilkan stok berdasarkan urutan FIFO.

---

## 8.1 Get FIFO Stocks

```http
GET /stocks/fifo
```

Query params:

| Name | Type | Required | Description |
|---|---|---:|---|
| `fish_type_id` | UUID | No | Jika diisi, FIFO ditampilkan per jenis ikan |
| `limit` | number | No | Jumlah data |
| `offset` | number | No | Offset pagination |

Example:

```http
GET /stocks/fifo?fish_type_id=7f9d8c5a-4a3b-4f2e-8a21-111111111111
```

Response:

```json
{
  "success": true,
  "message": "FIFO stocks retrieved successfully",
  "data": [
    {
      "id": "5cf59487-8233-4dbf-a27e-333333333333",
      "fish_type_name": "Tuna",
      "quality": "baik",
      "remaining_weight_kg": 50,
      "entered_at": "2026-05-01T08:00:00Z",
      "cold_storage_name": "Cold Storage A",
      "location_label": "Zona A - Rak 1",
      "fifo_rank": 1
    }
  ]
}
```

---

# 9. Stock Outs

Pengeluaran ikan dari cold storage.

---

## 9.1 Create Stock Out

```http
POST /stock-outs
```

Request:

```json
{
  "fish_type_id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
  "total_weight_kg": 40,
  "destination": "Restoran Laut Makassar",
  "out_at": "2026-05-01T12:00:00Z",
  "notes": "Pengeluaran untuk pesanan makan siang"
}
```

Backend behavior:

- Mencari stok ikan berdasarkan `fish_type_id`.
- Mengambil batch dengan `remaining_weight_kg > 0`.
- Mengurutkan batch berdasarkan `entered_at` paling lama.
- Mengurangi stok dari batch terlama terlebih dahulu.
- Jika batch pertama tidak cukup, sistem lanjut ke batch berikutnya.
- Menyimpan transaksi ke `stock_outs`.
- Menyimpan detail batch yang dikurangi ke `stock_out_items`.
- Membuat histori `stock_movement` dengan type `out`.
- Jika `remaining_weight_kg` batch menjadi 0, status menjadi `depleted`.

Example FIFO case:

```txt
Request keluar 40 kg Tuna.

Batch A:
- remaining_weight_kg: 25
- entered_at lebih lama

Batch B:
- remaining_weight_kg: 50
- entered_at lebih baru

Result:
- Batch A dikurangi 25 kg menjadi 0 kg
- Batch B dikurangi 15 kg menjadi 35 kg
```

Response:

```json
{
  "success": true,
  "message": "Stock out created successfully",
  "data": {
    "id": "0182c0f1-c188-4d1d-b428-555555555555",
    "fish_type_id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
    "destination": "Restoran Laut Makassar",
    "total_weight_kg": 40,
    "out_at": "2026-05-01T12:00:00Z",
    "items": [
      {
        "stock_batch_id": "5cf59487-8233-4dbf-a27e-333333333333",
        "weight_kg": 25
      },
      {
        "stock_batch_id": "b0bb6184-6313-4949-a736-666666666666",
        "weight_kg": 15
      }
    ],
    "created_at": "2026-05-01T12:01:00Z"
  }
}
```

Error if stock insufficient:

```json
{
  "success": false,
  "message": "Insufficient stock",
  "errors": [
    {
      "field": "total_weight_kg",
      "message": "Requested 100 kg, but only 75 kg is available"
    }
  ]
}
```

---

## 9.2 Get Stock Outs

```http
GET /stock-outs
```

Query params:

| Name | Type | Required | Description |
|---|---|---:|---|
| `fish_type_id` | UUID | No | Filter berdasarkan jenis ikan |
| `destination` | string | No | Filter tujuan pengeluaran |
| `date_from` | string | No | Format YYYY-MM-DD |
| `date_to` | string | No | Format YYYY-MM-DD |

Response:

```json
{
  "success": true,
  "message": "Stock outs retrieved successfully",
  "data": [
    {
      "id": "0182c0f1-c188-4d1d-b428-555555555555",
      "destination": "Restoran Laut Makassar",
      "total_weight_kg": 40,
      "out_at": "2026-05-01T12:00:00Z",
      "notes": "Pengeluaran untuk pesanan makan siang",
      "items": [
        {
          "stock_batch_id": "5cf59487-8233-4dbf-a27e-333333333333",
          "fish_type_name": "Tuna",
          "weight_kg": 25
        }
      ],
      "created_at": "2026-05-01T12:01:00Z"
    }
  ]
}
```

---

# 10. Dashboard

Data ringkas untuk monitoring stok.

---

## 10.1 Get Dashboard Summary

```http
GET /dashboard/summary
```

Response:

```json
{
  "success": true,
  "message": "Dashboard summary retrieved successfully",
  "data": {
    "total_available_weight_kg": 250,
    "total_stock_batches": 12,
    "total_available_batches": 10,
    "total_depleted_batches": 2,
    "today_stock_in_kg": 120,
    "today_stock_out_kg": 40,
    "fish_type_summary": [
      {
        "fish_type_id": "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
        "fish_type_name": "Tuna",
        "available_weight_kg": 80,
        "available_batches": 3
      }
    ],
    "cold_storage_summary": [
      {
        "cold_storage_id": "16dfdc88-831f-4e28-b2c7-222222222222",
        "cold_storage_name": "Cold Storage A",
        "available_weight_kg": 150,
        "available_batches": 7
      }
    ]
  }
}
```

---

## 10.2 Get Recent Movements

```http
GET /dashboard/recent-movements
```

Query params:

| Name | Type | Required | Description |
|---|---|---:|---|
| `limit` | number | No | Default 10 |

Response:

```json
{
  "success": true,
  "message": "Recent movements retrieved successfully",
  "data": [
    {
      "id": "c0ec32ef-10aa-4fbb-8b50-777777777777",
      "stock_batch_id": "5cf59487-8233-4dbf-a27e-333333333333",
      "movement_type": "in",
      "fish_type_name": "Tuna",
      "weight_kg": 50,
      "description": "Stok Tuna 50 kg masuk ke Cold Storage A",
      "created_at": "2026-05-01T08:05:00Z"
    }
  ]
}
```

---

# 11. Suggested MVP Endpoint Priority

Endpoint minimum untuk demo:

```txt
GET    /health

GET    /fish-types
POST   /fish-types

GET    /cold-storages
POST   /cold-storages

POST   /stocks
GET    /stocks
GET    /stocks/fifo

POST   /stock-outs

GET    /dashboard/summary
GET    /dashboard/recent-movements
```

---

# 12. MVP End-to-End Flow

```txt
1. User membuat master jenis ikan.
2. User membuat lokasi cold storage.
3. User mencatat ikan masuk.
4. Sistem menyimpan batch stok ke database.
5. User membuka daftar stok FIFO.
6. Sistem menampilkan stok dari yang paling lama masuk.
7. User mencatat pengeluaran ikan.
8. Sistem mengurangi stok berdasarkan FIFO.
9. Dashboard diperbarui.
```