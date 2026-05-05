"use client"
import { useEffect, useState } from "react"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL

// Tipe data dari API
type FishType = {
  id: string
  name: string
}

type Stock = {
  id: string
  fish_type_name: string
  quality: string
  remaining_weight_kg: number
  entered_at: string
  cold_storage_name: string
  location_label: string
  fifo_rank: number
}

// Mock data sementara (sebelum API Pison siap)
const mockFishTypes: FishType[] = [
  { id: "fish-001", name: "Tuna" },
  { id: "fish-002", name: "Tongkol" },
  { id: "fish-003", name: "Cakalang" },
  { id: "fish-004", name: "Bandeng" },
  { id: "fish-005", name: "Kerapu" },
]

const mockStocks: Stock[] = [
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
    fish_type_name: "Tongkol",
    quality: "baik",
    remaining_weight_kg: 60,
    entered_at: "2026-05-02T08:30:00Z",
    cold_storage_name: "Cold Storage B",
    location_label: "Zona B - Rak 2",
    fifo_rank: 1,
  },
  {
    id: "stock-004",
    fish_type_name: "Cakalang",
    quality: "buruk",
    remaining_weight_kg: 20,
    entered_at: "2026-05-02T09:00:00Z",
    cold_storage_name: "Cold Storage C",
    location_label: "Zona C - Rak 3",
    fifo_rank: 1,
  },
]

// Fungsi bantu: format tanggal jadi lebih mudah dibaca
function formatTanggal(isoString: string): string {
  const date = new Date(isoString)
  return date.toLocaleString("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

// Fungsi bantu: warna badge kualitas
function warnaBadgeKualitas(quality: string): string {
  if (quality === "baik") return "background:#d1fae5;color:#065f46"
  if (quality === "sedang") return "background:#fef9c3;color:#854d0e"
  return "background:#fee2e2;color:#991b1b"
}

export default function StocksPage() {
  const [stocks, setStocks] = useState<Stock[]>([])
  const [fishTypes, setFishTypes] = useState<FishType[]>([])
  const [selectedFishType, setSelectedFishType] = useState("")
  const [loading, setLoading] = useState(true)
  const [useMock, setUseMock] = useState(false)

  // Ambil jenis ikan untuk dropdown
  useEffect(() => {
    fetch(`${API_BASE}/fish-types`)
      .then((res) => res.json())
      .then((json) => setFishTypes(json.data))
      .catch(() => {
        // API belum siap, pakai mock data
        setFishTypes(mockFishTypes)
        setUseMock(true)
      })
  }, [])

  // Ambil data FIFO, diulang kalau filter berubah
  useEffect(() => {
    setLoading(true)

    if (useMock) {
      // Pakai mock data, filter manual
      const filtered = selectedFishType
        ? mockStocks.filter((s) => {
            const fish = mockFishTypes.find((f) => f.id === selectedFishType)
            return fish ? s.fish_type_name === fish.name : true
          })
        : mockStocks
      setStocks(filtered)
      setLoading(false)
      return
    }

    const url = selectedFishType
      ? `${API_BASE}/stocks/fifo?fish_type_id=${selectedFishType}`
      : `${API_BASE}/stocks/fifo`

    fetch(url)
      .then((res) => res.json())
      .then((json) => {
        setStocks(json.data)
        setLoading(false)
      })
      .catch(() => {
        // Fallback ke mock kalau API error
        setStocks(mockStocks)
        setLoading(false)
        setUseMock(true)
      })
  }, [selectedFishType, useMock])

  return (
    <div style={{ maxWidth: 480, margin: "0 auto", padding: "16px" }}>

      {/* Header */}
      <div style={{ marginBottom: 16 }}>
        <h1 style={{ fontSize: 20, fontWeight: "bold", margin: 0 }}>
          Daftar Stok Ikan
        </h1>
        <p style={{ fontSize: 13, color: "#6b7280", margin: "4px 0 0" }}>
          Diurutkan berdasarkan FIFO (masuk paling lama menjadi prioritas pertama)
        </p>
        {useMock && (
          <p style={{ fontSize: 12, color: "#b45309", marginTop: 4 }}>
            Menggunakan data simulasi (API belum tersambung)
          </p>
        )}
      </div>

      {/* Filter jenis ikan */}
      <div style={{ marginBottom: 16 }}>
        <select
          value={selectedFishType}
          onChange={(e) => setSelectedFishType(e.target.value)}
          style={{
            width: "100%",
            padding: "10px 12px",
            borderRadius: 8,
            border: "1px solid #d1d5db",
            fontSize: 14,
          }}
        >
          <option value="">Semua Jenis Ikan</option>
          {fishTypes.map((f) => (
            <option key={f.id} value={f.id}>
              {f.name}
            </option>
          ))}
        </select>
      </div>

      {/* Loading state */}
      {loading && (
        <p style={{ textAlign: "center", color: "#6b7280" }}>
          Memuat data stok...
        </p>
      )}

      {/* Empty state */}
      {!loading && stocks.length === 0 && (
        <div style={{ textAlign: "center", padding: "40px 0", color: "#6b7280" }}>
          <p style={{ fontSize: 32 }}>🐟</p>
          <p style={{ fontWeight: "bold" }}>Tidak ada stok tersedia</p>
          <p style={{ fontSize: 13 }}>
            Belum ada ikan yang masuk untuk jenis ini
          </p>
        </div>
      )}

      {/* Daftar stok */}
      {!loading &&
        stocks.map((stok) => (
          <div
            key={stok.id}
            style={{
              border: "1px solid #e5e7eb",
              borderRadius: 12,
              padding: 16,
              marginBottom: 12,
              background: "#fff",
            }}
          >
            {/* Baris atas: ranking + nama ikan + badge kualitas */}
            <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 8 }}>
              <span
                style={{
                  background: "#1e40af",
                  color: "#fff",
                  borderRadius: 999,
                  width: 28,
                  height: 28,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: 13,
                  fontWeight: "bold",
                  flexShrink: 0,
                }}
              >
                {stok.fifo_rank}
              </span>
              <span style={{ fontWeight: "bold", fontSize: 16, flex: 1 }}>
                {stok.fish_type_name}
              </span>
              <span
                style={{
                  ...Object.fromEntries(
                    warnaBadgeKualitas(stok.quality)
                      .split(";")
                      .map((s) => s.split(":"))
                  ),
                  padding: "2px 10px",
                  borderRadius: 999,
                  fontSize: 12,
                  fontWeight: 600,
                }}
              >
                {stok.quality}
              </span>
            </div>

            {/* Detail stok */}
            <div style={{ fontSize: 13, color: "#374151", display: "flex", flexDirection: "column", gap: 4 }}>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span style={{ color: "#6b7280" }}>Sisa Berat</span>
                <span style={{ fontWeight: "bold" }}>
                  {stok.remaining_weight_kg} kg
                </span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span style={{ color: "#6b7280" }}>Waktu Masuk</span>
                <span>{formatTanggal(stok.entered_at)}</span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span style={{ color: "#6b7280" }}>Lokasi</span>
                <span>{stok.cold_storage_name}</span>
              </div>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <span style={{ color: "#6b7280" }}>Zona</span>
                <span>{stok.location_label}</span>
              </div>
            </div>
          </div>
        ))}
    </div>
  )
}