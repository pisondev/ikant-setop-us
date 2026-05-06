# Ikan't Setop Us

Mobile-first web app untuk mencatat stok ikan masuk, stok FIFO, pengeluaran ikan, dan ringkasan dashboard cold storage.

Project ini dibuat sebagai MVP untuk membuktikan alur:

```txt
input stok masuk -> data tersimpan -> stok tampil FIFO -> stok keluar dicatat -> dashboard berubah
```

## Status Project

Update terakhir: 2026-05-06.

| Area | Status | Catatan |
|---|---|---|
| Database schema | Selesai | Migration awal tersedia di `apps/api/migrations`. |
| Backend API | Selesai untuk kontrak MVP | Endpoint sudah dipetakan per modul dan punya unit test. |
| Postman docs | Selesai | Collection dan environment tersedia di `docs/api`. |
| Frontend | Sebagian | Foundation, master data, `/`, dan `/stocks` sudah tersedia; stock-in/out dan dashboard belum. |
| Dokumentasi progress | Disinkronkan | Detail status ada di `docs/frontend/frontend-mvp-todo.md` dan `docs/jobdesc_progress.md`. |

## Fitur MVP

- Master jenis ikan.
- Master cold storage.
- Input stok ikan masuk.
- Daftar stok FIFO.
- Input ikan keluar berdasarkan FIFO.
- Riwayat pengeluaran.
- Dashboard summary dan recent movements.

## Tech Stack

- Frontend: Next.js, React, TypeScript.
- Backend: Go Fiber.
- Database: PostgreSQL 16.
- API docs: Postman collection dan API contract markdown.

## Struktur Folder

```txt
apps/
  api/
    cmd/api/                  Entry point backend
    internal/modules/          Modul backend per domain
    migrations/                SQL migration
  web/
    app/                       Next.js app router
    components/layout/          App shell, bottom nav, page header
    lib/                        API helper
    types/                      Shared API types
docs/
  api/                         API contract, Postman collection, environment
  diagram/                     Diagram sistem dan DBML
  frontend/                    Frontend MVP source of truth
  jobdesc_progress.md          Progress dan pembagian tugas tim
docker-compose.yaml            Database lokal
Makefile                       Helper migration Windows
```

## Backend API

Base URL lokal:

```txt
http://localhost:8081/api/v1
```

Endpoint MVP:

```txt
GET    /health
GET    /fish-types
POST   /fish-types
GET    /cold-storages
POST   /cold-storages
GET    /stocks
POST   /stocks
GET    /stocks/fifo
GET    /stocks/{id}
PATCH  /stocks/{id}/quality
PATCH  /stocks/{id}/location
GET    /stock-outs
POST   /stock-outs
GET    /dashboard/summary
GET    /dashboard/recent-movements
```

Dokumen API:

- `docs/api/api-contract-v1.md`
- `docs/api/ikant-setop-us-api.postman_collection.json`
- `docs/api/ikant-setop-us-local.postman_environment.json`

## Environment

Root `.env` / `apps/api/.env` untuk backend:

```env
APP_ENV=development
APP_NAME=ikant-setop-us-api
APP_PORT=8081
APP_VERSION=v1
CORS_ALLOWED_ORIGINS=http://localhost:3000
DB_HOST=localhost
DB_PORT=5438
DB_USER=ikant_user
DB_PASSWORD=ikant_pass
DB_NAME=ikant_setop_us_db
DB_SSLMODE=disable
```

Frontend `apps/web/.env.local`:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api/v1
```

## Cara Menjalankan

Jalankan database:

```bash
docker compose up -d
```

Jalankan migration dari root repo:

```bash
make migrate-up-windows
```

Jalankan backend:

```bash
cd apps/api
go run ./cmd/api
```

Jalankan frontend:

```bash
cd apps/web
npm install
npm run dev
```

Frontend lokal:

```txt
http://localhost:3000
```

## Test

Backend:

```bash
cd apps/api
go test ./... -count=1
go vet ./...
go test -race ./... -count=1
```

Frontend:

```bash
cd apps/web
npm run lint
npm run build
```

## Demo Flow Target

1. Buka dashboard.
2. Tambah jenis ikan jika belum ada.
3. Tambah cold storage jika belum ada.
4. Input stok Tuna 50 kg kualitas baik ke Cold Storage A.
5. Input stok Tuna 30 kg kualitas sedang ke Cold Storage B.
6. Buka daftar stok FIFO dan filter Tuna.
7. Catat pengeluaran Tuna 40 kg.
8. Sistem mengurangi batch paling lama lebih dulu.
9. Dashboard dan recent movements berubah.

## Dokumen Utama

- `docs/frontend/frontend-mvp-todo.md` adalah single source of truth status frontend MVP.
- `docs/jobdesc_progress.md` adalah tracker progress per owner.
- `docs/diagram/diagram.md` adalah sumber diagram sistem.
- `docs/api/api-contract-v1.md` adalah kontrak API backend.

## Catatan

Fitur login, role-based access, notifikasi realtime, grafik kompleks, integrasi alat timbang, dan sensor cold storage belum masuk scope MVP.
