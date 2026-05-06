# Diagram Sistem | Ikan't Setop Us

Dokumen ini adalah sumber diagram sistem. File PNG di folder ini adalah artefak visual; jika diagram berubah, update sumber di file ini terlebih dahulu.

Update terakhir: 2026-05-06.

## 1. Context Diagram

```mermaid
flowchart LR
    Baso["Baso<br/>Admin Gudang / Pekerja Lapangan"]
    Daeng["Daeng Syamsul<br/>Pemilik Usaha"]
    FE["Next.js Frontend<br/>Mobile-first web app"]
    API["Go Fiber API<br/>/api/v1"]
    DB[("PostgreSQL Database")]
    Cold["Cold Storage"]
    Buyer["Tujuan Pengeluaran<br/>Restoran / Pembeli"]

    Baso -->|"input stok masuk dan stok keluar"| FE
    Daeng -->|"pantau dashboard dan FIFO"| FE
    FE -->|"HTTP JSON"| API
    API -->|"read/write data"| DB
    API -->|"referensi lokasi"| Cold
    API -->|"catat tujuan pengeluaran"| Buyer
    DB -->|"stok, movement, master data"| API
    API -->|"response API"| FE
```

## 2. Backend Module Diagram

Backend saat ini dipetakan per modul domain di `apps/api/internal/modules`.

```mermaid
flowchart TD
    Main["cmd/api/main.go"]
    V1["Fiber router /api/v1"]
    Fish["modules/fish<br/>GET/POST /fish-types"]
    Storage["modules/storage<br/>GET/POST /cold-storages"]
    Stock["modules/stock<br/>stocks, FIFO, quality, location"]
    Stockout["modules/stockout<br/>stock-out FIFO transaction"]
    Dashboard["modules/dashboard<br/>summary, recent movements"]
    DB[("Database")]

    Main --> V1
    V1 --> Fish
    V1 --> Storage
    V1 --> Stock
    V1 --> Stockout
    V1 --> Dashboard
    Fish --> DB
    Storage --> DB
    Stock --> DB
    Stockout --> DB
    Dashboard --> DB
```

## 3. Use Case Diagram

```plantuml
@startuml
left to right direction

actor "Baso\nAdmin Gudang / Pekerja Lapangan" as Baso
actor "Daeng Syamsul\nPemilik Usaha" as Daeng

rectangle "Sistem Manajemen Inventori Ikan" {
  usecase "Kelola Jenis Ikan" as UC1
  usecase "Kelola Cold Storage" as UC2
  usecase "Mencatat Ikan Masuk" as UC3
  usecase "Melihat Stok FIFO" as UC4
  usecase "Memfilter FIFO per Jenis Ikan" as UC5
  usecase "Mencatat Ikan Keluar" as UC6
  usecase "Mengurangi Stok Berdasarkan FIFO" as UC7
  usecase "Melihat Dashboard" as UC8
  usecase "Melihat Recent Movement" as UC9
  usecase "Melihat Riwayat Pengeluaran" as UC10
  usecase "Memperbarui Kualitas Stok" as UC11
  usecase "Memperbarui Lokasi Stok" as UC12
}

Baso --> UC1
Baso --> UC2
Baso --> UC3
Baso --> UC4
Baso --> UC5
Baso --> UC6
Baso --> UC10
Baso --> UC11
Baso --> UC12

Daeng --> UC4
Daeng --> UC5
Daeng --> UC8
Daeng --> UC9
Daeng --> UC10

UC6 .> UC7 : include
UC5 .> UC4 : include
UC8 .> UC9 : include

@enduml
```

## 4. ER Diagram

Sumber DBML terpisah juga tersedia di `docs/diagram/diagram_erd.dbml`.

```dbml
Project fish_inventory_fifo {
  database_type: "PostgreSQL"
  Note: '''
  Sistem manajemen inventori ikan berbasis FIFO untuk stok masuk,
  stok keluar, lokasi cold storage, kualitas ikan, dan histori aktivitas stok.
  '''
}

Enum user_role {
  owner
  warehouse_admin
}

Enum fish_quality {
  baik
  sedang
  buruk
}

Enum stock_status {
  available
  depleted
}

Enum stock_movement_type {
  in
  out
  quality_update
  location_update
  adjustment
}

Table users {
  id uuid [pk]
  name varchar(120) [not null]
  role user_role [not null]
  email varchar(120)
  password_hash text
  created_at timestamp [not null]
  updated_at timestamp [not null]
}

Table fish_types {
  id uuid [pk]
  name varchar(100) [not null, unique]
  image_url text
  description text
  created_at timestamp [not null]
  updated_at timestamp [not null]
}

Table cold_storages {
  id uuid [pk]
  name varchar(100) [not null]
  location_label varchar(150)
  description text
  created_at timestamp [not null]
  updated_at timestamp [not null]
}

Table stock_batches {
  id uuid [pk]
  fish_type_id uuid [not null]
  cold_storage_id uuid [not null]
  quality fish_quality [not null]
  initial_weight_kg decimal(10,2) [not null]
  remaining_weight_kg decimal(10,2) [not null]
  entered_at timestamp [not null]
  status stock_status [not null]
  notes text
  created_by uuid
  created_at timestamp [not null]
  updated_at timestamp [not null]
}

Table stock_outs {
  id uuid [pk]
  destination varchar(150) [not null]
  total_weight_kg decimal(10,2) [not null]
  out_at timestamp [not null]
  notes text
  created_by uuid
  created_at timestamp [not null]
  updated_at timestamp [not null]
}

Table stock_out_items {
  id uuid [pk]
  stock_out_id uuid [not null]
  stock_batch_id uuid [not null]
  weight_kg decimal(10,2) [not null]
  created_at timestamp [not null]
}

Table stock_movements {
  id uuid [pk]
  stock_batch_id uuid [not null]
  movement_type stock_movement_type [not null]
  weight_kg decimal(10,2)
  previous_quality fish_quality
  new_quality fish_quality
  previous_cold_storage_id uuid
  new_cold_storage_id uuid
  description text
  created_by uuid
  created_at timestamp [not null]
}

Ref: stock_batches.fish_type_id > fish_types.id
Ref: stock_batches.cold_storage_id > cold_storages.id
Ref: stock_batches.created_by > users.id
Ref: stock_outs.created_by > users.id
Ref: stock_out_items.stock_out_id > stock_outs.id
Ref: stock_out_items.stock_batch_id > stock_batches.id
Ref: stock_movements.stock_batch_id > stock_batches.id
Ref: stock_movements.previous_cold_storage_id > cold_storages.id
Ref: stock_movements.new_cold_storage_id > cold_storages.id
Ref: stock_movements.created_by > users.id
```

## 5. Sequence Diagram - Stock In

```mermaid
sequenceDiagram
    actor User as Baso
    participant FE as Next.js Frontend
    participant API as Go Fiber API /api/v1
    participant DB as Database

    User->>FE: Buka /stocks/new
    FE->>API: GET /fish-types
    API->>DB: Ambil master jenis ikan
    DB-->>API: Data jenis ikan
    API-->>FE: Fish types
    FE->>API: GET /cold-storages
    API->>DB: Ambil master cold storage
    DB-->>API: Data cold storage
    API-->>FE: Cold storages
    User->>FE: Isi form stok masuk
    FE->>API: POST /stocks
    API->>DB: Insert stock_batches
    API->>DB: Insert stock_movements type in
    DB-->>API: Transaksi berhasil
    API-->>FE: Stock batch created
    FE-->>User: Tampilkan sukses / arahkan ke /stocks
```

## 6. Sequence Diagram - Stock Out FIFO

```mermaid
sequenceDiagram
    actor User as Baso
    participant FE as Next.js Frontend
    participant API as Go Fiber API /api/v1
    participant DB as Database

    User->>FE: Buka /stock-outs/new
    FE->>API: GET /fish-types
    API-->>FE: Fish types
    User->>FE: Pilih jenis ikan
    FE->>API: GET /stocks/fifo?fish_type_id={id}
    API->>DB: Ambil batch available terurut entered_at
    DB-->>API: FIFO batches
    API-->>FE: FIFO preview
    User->>FE: Isi berat keluar dan tujuan
    FE->>API: POST /stock-outs
    API->>DB: Begin transaction
    API->>DB: Lock batch available FOR UPDATE
    API->>API: Validasi stok cukup
    API->>DB: Insert stock_outs
    API->>DB: Insert stock_out_items
    API->>DB: Update remaining_weight_kg dan status batch
    API->>DB: Insert stock_movements type out
    API->>DB: Commit transaction
    API-->>FE: Stock out created dengan items
    FE-->>User: Tampilkan batch yang dipakai
```

## 7. Catatan Sinkronisasi

- Endpoint sequence memakai base path `/api/v1`.
- Update kualitas dan lokasi sudah ada di backend, tetapi bukan prioritas halaman frontend MVP.
- Auth dan role user belum diimplementasikan di aplikasi MVP walaupun tabel `users` sudah ada di schema.
