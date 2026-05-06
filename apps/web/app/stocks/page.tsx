"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { ApiError, apiGet } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type { FIFOStock, FishQuality, FishType } from "@/types/api";

function formatDate(value: string): string {
  return new Date(value).toLocaleString("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function qualityClass(quality: FishQuality): string {
  if (quality === "baik") return "bg-emerald-100 text-emerald-800";
  if (quality === "sedang") return "bg-amber-100 text-amber-800";
  return "bg-rose-100 text-rose-800";
}

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Gagal memuat data.";
}

export default function StocksPage() {
  const [stocks, setStocks] = useState<FIFOStock[]>([]);
  const [fishTypes, setFishTypes] = useState<FishType[]>([]);
  const [selectedFishType, setSelectedFishType] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let ignore = false;

    async function loadFishTypes() {
      try {
        const data = await apiGet<FishType[]>("/fish-types");
        if (!ignore) setFishTypes(data);
      } catch (err) {
        if (!ignore) setError(errorMessage(err));
      }
    }

    loadFishTypes();
    return () => {
      ignore = true;
    };
  }, []);

  useEffect(() => {
    let ignore = false;

    async function loadStocks() {
      setLoading(true);
      setError(null);

      try {
        const path = selectedFishType
          ? `/stocks/fifo?fish_type_id=${selectedFishType}`
          : "/stocks/fifo";
        const data = await apiGet<FIFOStock[]>(path);
        if (!ignore) setStocks(data);
      } catch (err) {
        if (!ignore) {
          setStocks([]);
          setError(errorMessage(err));
        }
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    loadStocks();
    return () => {
      ignore = true;
    };
  }, [selectedFishType]);

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Daftar Stok FIFO"
        description="Batch paling lama masuk tampil sebagai prioritas."
        actions={
          <Link
            href="/stocks/new"
            className="rounded-lg bg-emerald-700 px-3 py-2 text-sm font-semibold text-white"
          >
            Tambah
          </Link>
        }
      />

      <div className="mb-4 grid gap-3">
        <select
          value={selectedFishType}
          onChange={(event) => setSelectedFishType(event.target.value)}
          className="h-11 w-full rounded-lg border border-slate-300 bg-white px-3 text-sm text-slate-900 outline-none focus:border-emerald-600"
        >
          <option value="">Semua jenis ikan</option>
          {fishTypes.map((fishType) => (
            <option key={fishType.id} value={fishType.id}>
              {fishType.name}
            </option>
          ))}
        </select>

        <div className="grid grid-cols-2 gap-2">
          <Link
            href="/stocks/new"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-center text-sm font-semibold text-slate-800"
          >
            Stok Masuk
          </Link>
          <Link
            href="/stock-outs/new"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-center text-sm font-semibold text-slate-800"
          >
            Ikan Keluar
          </Link>
        </div>
      </div>

      {error ? (
        <div className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat data stok...
        </div>
      ) : null}

      {!loading && stocks.length === 0 && !error ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center">
          <p className="text-sm font-semibold text-slate-900">
            Tidak ada stok tersedia
          </p>
          <p className="mt-1 text-sm text-slate-600">
            Belum ada batch available untuk filter ini.
          </p>
        </div>
      ) : null}

      {!loading && stocks.length > 0 ? (
        <div className="grid gap-3">
          {stocks.map((stock) => (
            <article
              key={stock.id}
              className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="mb-3 flex items-center gap-3">
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-emerald-700 text-sm font-semibold text-white">
                  {stock.fifo_rank}
                </span>
                <div className="min-w-0 flex-1">
                  <h2 className="truncate text-base font-semibold text-slate-950">
                    {stock.fish_type_name}
                  </h2>
                  <p className="text-xs text-slate-500">
                    {formatDate(stock.entered_at)}
                  </p>
                </div>
                <span
                  className={[
                    "rounded-full px-2.5 py-1 text-xs font-semibold",
                    qualityClass(stock.quality),
                  ].join(" ")}
                >
                  {stock.quality}
                </span>
              </div>

              <dl className="grid gap-2 text-sm">
                <div className="flex justify-between gap-3">
                  <dt className="text-slate-500">Sisa berat</dt>
                  <dd className="font-semibold text-slate-950">
                    {stock.remaining_weight_kg} kg
                  </dd>
                </div>
                <div className="flex justify-between gap-3">
                  <dt className="text-slate-500">Cold storage</dt>
                  <dd className="text-right font-medium text-slate-800">
                    {stock.cold_storage_name}
                  </dd>
                </div>
                <div className="flex justify-between gap-3">
                  <dt className="text-slate-500">Lokasi</dt>
                  <dd className="text-right text-slate-800">
                    {stock.location_label ?? "-"}
                  </dd>
                </div>
              </dl>
            </article>
          ))}
        </div>
      ) : null}
    </div>
  );
}
