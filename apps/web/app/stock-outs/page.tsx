"use client";

import Link from "next/link";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { ApiError, apiGet } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type { FishType, StockOut } from "@/types/api";

type StockOutFilters = {
  fish_type_id: string;
  destination: string;
  date_from: string;
  date_to: string;
};

const initialFilters: StockOutFilters = {
  fish_type_id: "",
  destination: "",
  date_from: "",
  date_to: "",
};

export default function StockOutHistoryPage() {
  const [stockOuts, setStockOuts] = useState<StockOut[]>([]);
  const [fishTypes, setFishTypes] = useState<FishType[]>([]);
  const [filters, setFilters] = useState<StockOutFilters>(initialFilters);
  const [appliedFilters, setAppliedFilters] =
    useState<StockOutFilters>(initialFilters);
  const [loading, setLoading] = useState(true);
  const [loadingFishTypes, setLoadingFishTypes] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [fishTypeError, setFishTypeError] = useState<string | null>(null);

  useEffect(() => {
    let ignore = false;

    async function loadFishTypes() {
      try {
        const data = await apiGet<FishType[]>("/fish-types");
        if (!ignore) setFishTypes(data);
      } catch (err) {
        if (!ignore) setFishTypeError(errorMessage(err));
      } finally {
        if (!ignore) setLoadingFishTypes(false);
      }
    }

    loadFishTypes();
    return () => {
      ignore = true;
    };
  }, []);

  useEffect(() => {
    let ignore = false;

    async function loadStockOuts() {
      setLoading(true);
      setError(null);

      try {
        const data = await apiGet<StockOut[]>(
          `/stock-outs${buildQuery(appliedFilters)}`,
        );
        if (!ignore) setStockOuts(data);
      } catch (err) {
        if (!ignore) {
          setStockOuts([]);
          setError(errorMessage(err));
        }
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    loadStockOuts();
    return () => {
      ignore = true;
    };
  }, [appliedFilters]);

  const totalWeight = useMemo(
    () =>
      stockOuts.reduce(
        (total, stockOut) => total + Number(stockOut.total_weight_kg),
        0,
      ),
    [stockOuts],
  );

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setAppliedFilters(normalizeFilters(filters));
  }

  function handleReset() {
    setFilters(initialFilters);
    setAppliedFilters(initialFilters);
  }

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Riwayat Pengeluaran"
        description="Daftar ikan keluar beserta batch FIFO yang terpakai."
        actions={
          <Link
            href="/stock-outs/new"
            className="rounded-lg bg-emerald-700 px-3 py-2 text-sm font-semibold text-white"
          >
            Catat
          </Link>
        }
      />

      <form
        onSubmit={handleSubmit}
        className="mb-4 grid gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
      >
        <div className="grid gap-1.5">
          <label
            className="text-sm font-medium text-slate-700"
            htmlFor="fish_type_id"
          >
            Jenis ikan
          </label>
          <select
            id="fish_type_id"
            value={filters.fish_type_id}
            onChange={(event) =>
              setFilters((current) => ({
                ...current,
                fish_type_id: event.target.value,
              }))
            }
            className="h-11 rounded-lg border border-slate-300 bg-white px-3 text-sm text-slate-900 outline-none focus:border-emerald-600"
            disabled={loadingFishTypes}
          >
            <option value="">Semua jenis ikan</option>
            {fishTypes.map((fishType) => (
              <option key={fishType.id} value={fishType.id}>
                {fishType.name}
              </option>
            ))}
          </select>
        </div>

        <div className="grid gap-1.5">
          <label
            className="text-sm font-medium text-slate-700"
            htmlFor="destination"
          >
            Tujuan
          </label>
          <input
            id="destination"
            value={filters.destination}
            onChange={(event) =>
              setFilters((current) => ({
                ...current,
                destination: event.target.value,
              }))
            }
            className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
            placeholder="Restoran Laut Makassar"
          />
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="grid gap-1.5">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="date_from"
            >
              Dari
            </label>
            <input
              id="date_from"
              type="date"
              value={filters.date_from}
              onChange={(event) =>
                setFilters((current) => ({
                  ...current,
                  date_from: event.target.value,
                }))
              }
              className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
            />
          </div>
          <div className="grid gap-1.5">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="date_to"
            >
              Sampai
            </label>
            <input
              id="date_to"
              type="date"
              value={filters.date_to}
              onChange={(event) =>
                setFilters((current) => ({
                  ...current,
                  date_to: event.target.value,
                }))
              }
              className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
            />
          </div>
        </div>

        {fishTypeError ? (
          <div className="rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
            {fishTypeError}
          </div>
        ) : null}

        <div className="grid grid-cols-2 gap-2">
          <button
            type="submit"
            className="h-11 rounded-lg bg-emerald-700 px-4 text-sm font-semibold text-white disabled:bg-slate-400"
            disabled={loading}
          >
            {loading ? "Memuat..." : "Terapkan"}
          </button>
          <button
            type="button"
            onClick={handleReset}
            className="h-11 rounded-lg border border-slate-300 bg-white px-4 text-sm font-semibold text-slate-800"
          >
            Reset
          </button>
        </div>
      </form>

      <section className="mb-4 grid grid-cols-2 gap-3">
        <SummaryCard label="Transaksi" value={stockOuts.length.toString()} />
        <SummaryCard
          label="Total keluar"
          value={`${formatWeight(totalWeight)} kg`}
        />
      </section>

      {error ? (
        <div className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat riwayat pengeluaran...
        </div>
      ) : null}

      {!loading && stockOuts.length === 0 && !error ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center">
          <p className="text-sm font-semibold text-slate-900">
            Belum ada pengeluaran
          </p>
          <p className="mt-1 text-sm text-slate-600">
            Catat ikan keluar untuk melihat riwayat FIFO di sini.
          </p>
          <Link
            href="/stock-outs/new"
            className="mt-4 inline-flex rounded-lg bg-emerald-700 px-3 py-2 text-sm font-semibold text-white"
          >
            Catat Ikan Keluar
          </Link>
        </div>
      ) : null}

      {!loading && stockOuts.length > 0 ? (
        <div className="grid gap-3">
          {stockOuts.map((stockOut) => (
            <article
              key={stockOut.id}
              className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="mb-3 flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <h2 className="truncate text-base font-semibold text-slate-950">
                    {stockOut.destination}
                  </h2>
                  <p className="mt-1 text-xs text-slate-500">
                    {formatDate(stockOut.out_at)}
                  </p>
                </div>
                <p className="shrink-0 rounded-full bg-emerald-100 px-2.5 py-1 text-sm font-semibold text-emerald-800">
                  {formatWeight(stockOut.total_weight_kg)} kg
                </p>
              </div>

              {stockOut.notes ? (
                <p className="mb-3 rounded-lg bg-slate-50 px-3 py-2 text-sm text-slate-600">
                  {stockOut.notes}
                </p>
              ) : null}

              <div className="grid gap-2">
                <p className="text-xs font-semibold uppercase text-slate-500">
                  Batch terpakai
                </p>
                {stockOut.items.length === 0 ? (
                  <div className="rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-600">
                    Detail batch tidak tersedia.
                  </div>
                ) : (
                  stockOut.items.map((item) => (
                    <div
                      key={`${stockOut.id}-${item.stock_batch_id}`}
                      className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 px-3 py-2 text-sm"
                    >
                      <div className="min-w-0">
                        <p className="truncate font-medium text-slate-900">
                          {item.fish_type_name ?? "Batch FIFO"}
                        </p>
                        <p className="text-xs text-slate-500">
                          Batch {shortID(item.stock_batch_id)}
                        </p>
                      </div>
                      <p className="shrink-0 font-semibold text-slate-950">
                        {formatWeight(item.weight_kg)} kg
                      </p>
                    </div>
                  ))
                )}
              </div>
            </article>
          ))}
        </div>
      ) : null}
    </div>
  );
}

function SummaryCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
      <p className="text-xs font-medium text-slate-500">{label}</p>
      <p className="mt-1 text-lg font-semibold text-slate-950">{value}</p>
    </div>
  );
}

function buildQuery(filters: StockOutFilters): string {
  const params = new URLSearchParams();
  const normalizedFilters = normalizeFilters(filters);

  Object.entries(normalizedFilters).forEach(([key, value]) => {
    if (value) params.set(key, value);
  });

  const query = params.toString();
  return query ? `?${query}` : "";
}

function normalizeFilters(filters: StockOutFilters): StockOutFilters {
  return {
    fish_type_id: filters.fish_type_id,
    destination: filters.destination.trim(),
    date_from: filters.date_from,
    date_to: filters.date_to,
  };
}

function formatWeight(value: number): string {
  return new Intl.NumberFormat("id-ID", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(value);
}

function formatDate(value: string): string {
  return new Date(value).toLocaleString("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function shortID(value: string): string {
  return value.slice(0, 8);
}

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Gagal memuat riwayat pengeluaran.";
}
