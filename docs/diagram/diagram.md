# 1. Context Diagram
```mermaid
flowchart LR
    Baso["Baso<br/>Admin Gudang / Pekerja Lapangan"]
    Daeng["Daeng Syamsul<br/>Pemilik Usaha / Punggawa"]
    System["Sistem Manajemen Inventori Ikan<br/>Mobile-First Web App"]
    DB[("PostgreSQL Database")]
    Cold["Cold Storage"]
    Buyer["Tujuan Pengeluaran<br/>Restoran / Pembeli / Pihak Lain"]

    Baso -->|"Input ikan masuk<br/>jenis, kualitas, berat, waktu masuk, lokasi"| System
    Baso -->|"Input ikan keluar<br/>berat keluar, tujuan, catatan"| System
    Daeng -->|"Melihat dashboard<br/>stok, FIFO, alur masuk/keluar"| System
    System -->|"Menyimpan & membaca data"| DB
    System -->|"Menampilkan lokasi penyimpanan"| Cold
    System -->|"Mencatat tujuan pengeluaran"| Buyer

    DB -->|"Data stok, histori, lokasi, pengeluaran"| System
    System -->|"Update stok real-time"| Daeng
    System -->|"Daftar stok terurut FIFO"| Baso
```

# 2. Usecase Diagram
```plantuml
@startuml
left to right direction

actor "Baso\nAdmin Gudang / Pekerja Lapangan" as Baso
actor "Daeng Syamsul\nPemilik Usaha / Punggawa" as Daeng

rectangle "Sistem Manajemen Inventori Ikan" {
  usecase "Mencatat Ikan Masuk" as UC1
  usecase "Memilih Jenis Ikan" as UC2
  usecase "Menginput Kualitas Ikan" as UC3
  usecase "Menginput Berat Ikan" as UC4
  usecase "Memilih Lokasi Cold Storage" as UC5

  usecase "Melihat Daftar Stok" as UC6
  usecase "Memfilter Stok per Jenis Ikan" as UC7
  usecase "Melihat Urutan FIFO" as UC8

  usecase "Mencatat Ikan Keluar" as UC9
  usecase "Menginput Tujuan Pengeluaran" as UC10
  usecase "Mengurangi Stok" as UC11

  usecase "Melihat Dashboard Monitoring" as UC12
  usecase "Melihat Riwayat Aktivitas Stok" as UC13
  usecase "Memperbarui Kualitas Ikan" as UC14
}

Baso --> UC1
Baso --> UC6
Baso --> UC7
Baso --> UC8
Baso --> UC9
Baso --> UC14

Daeng --> UC6
Daeng --> UC7
Daeng --> UC8
Daeng --> UC12
Daeng --> UC13
Daeng --> UC14

UC1 .> UC2 : include
UC1 .> UC3 : include
UC1 .> UC4 : include
UC1 .> UC5 : include

UC9 .> UC10 : include
UC9 .> UC11 : include

UC6 .> UC8 : include
UC7 .> UC8 : include

@enduml
```

# 3. ER Diagram
```dbml
Project fish_inventory_fifo {
  database_type: "PostgreSQL"
  Note: '''
  Sistem Manajemen Inventori Ikan berbasis FIFO untuk pencatatan stok masuk,
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

  Note: '''
  Menyimpan data pengguna sistem.
  Untuk MVP, role bisa digunakan secara sederhana atau hanya sebagai data persona.
  '''
}

Table fish_types {
  id uuid [pk]
  name varchar(100) [not null, unique]
  image_url text
  description text
  created_at timestamp [not null]
  updated_at timestamp [not null]

  Note: '''
  Master data jenis ikan, misalnya tuna, tongkol, cakalang, kakap.
  '''
}

Table cold_storages {
  id uuid [pk]
  name varchar(100) [not null]
  location_label varchar(150)
  description text
  created_at timestamp [not null]
  updated_at timestamp [not null]

  Note: '''
  Data lokasi penyimpanan ikan.
  Bisa berupa Cold Storage A, Cold Storage B, rak, zona, atau label lokasi sederhana.
  '''
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

  Note: '''
  Satu baris mewakili satu batch stok ikan yang masuk.
  FIFO dihitung berdasarkan entered_at.
  '''
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

  Note: '''
  Header transaksi pengeluaran ikan.
  Detail batch yang dikurangi disimpan pada stock_out_items.
  '''
}

Table stock_out_items {
  id uuid [pk]
  stock_out_id uuid [not null]
  stock_batch_id uuid [not null]
  weight_kg decimal(10,2) [not null]
  created_at timestamp [not null]

  Note: '''
  Detail batch yang terpakai saat pengeluaran ikan.
  Satu pengeluaran bisa mengambil dari beberapa batch sesuai FIFO.
  '''
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

  Note: '''
  Histori aktivitas stok.
  Mencatat stok masuk, stok keluar, perubahan kualitas, perubahan lokasi, dan adjustment.
  '''
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

# 4. Sequence Diagram
```mermaid
sequenceDiagram
    actor User as Baso / Daeng
    participant FE as Next.js Frontend
    participant API as Go Fiber Backend
    participant DB as PostgreSQL

    User->>FE: Buka halaman pengeluaran ikan
    FE->>API: GET /api/fish-types
    API->>DB: Ambil daftar jenis ikan
    DB-->>API: Data jenis ikan
    API-->>FE: Response jenis ikan

    User->>FE: Pilih jenis ikan dan input berat keluar
    FE->>API: POST /api/stock-outs
    API->>DB: Cari batch stok berdasarkan jenis ikan
    DB-->>API: Batch stok terurut FIFO

    API->>API: Validasi stok mencukupi
    API->>DB: Kurangi stok dari batch terlama
    API->>DB: Simpan histori pengeluaran
    DB-->>API: Transaksi berhasil

    API-->>FE: Response stok berhasil dikeluarkan
    FE-->>User: Tampilkan update stok dan dashboard
```
