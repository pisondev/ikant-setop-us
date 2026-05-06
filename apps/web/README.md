# Ikan't Setop Us Web

Next.js frontend untuk MVP inventori ikan FIFO.

## Status

Update terakhir: 2026-05-06.

```txt
[x] Frontend foundation tahap 1
[x] Entry page /
[x] FIFO stock list /stocks
[x] Master fish types /fish-types
[x] Master cold storages /cold-storages
[x] Stock-in form /stocks/new
[ ] Dashboard
[ ] Stock-out form
[ ] Stock-out history
```

## Environment

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8081/api/v1
```

Jika env belum dibuat, helper API memakai default lokal yang sama.

## Development

```bash
npm install
npm run dev
```

App lokal:

```txt
http://localhost:3000
```

## Verification

```bash
npm run lint
npm run build
```

Catatan environment saat ini: Node `22.12.0` menampilkan warning engine dari salah satu dependency yang meminta `22.13.0+`, tetapi lint dan build tetap lolos.
