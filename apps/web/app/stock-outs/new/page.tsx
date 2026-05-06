"use client";

import Link from "next/link";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { ApiError, apiGet, apiPost } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type {
  ApiFieldError,
  FIFOStock,
  FishType,
  StockOut,
} from "@/types/api";

type StockOutForm = {
  fish_type_id: string;
  total_weight_kg: string;
  destination: string;
  out_at: string;
  notes: string;
};

const initialForm: StockOutForm = {
  fish_type_id: "",
  total_weight_kg: "",
  destination: "",
  out_at: toDatetimeLocalValue(new Date()),
  notes: "",
};

export default function NewStockOutPage() {
  const [fishTypes, setFishTypes] = useState<FishType[]>([]);
  const [fifoStocks, setFifoStocks] = useState<FIFOStock[]>([]);
  const [form, setForm] = useState<StockOutForm>(initialForm);
  const [createdStockOut, setCreatedStockOut] = useState<StockOut | null>(null);
  const [loadingFishTypes, setLoadingFishTypes] = useState(true);
  const [loadingFifo, setLoadingFifo] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fifoError, setFifoError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<ApiFieldError[]>([]);

  useEffect(() => {
    let ignore = false;

    async function loadFishTypes() {
      try {
        const data = await apiGet<FishType[]>("/fish-types");
        if (!ignore) {
          setFishTypes(data);
          setForm((current) => ({
            ...current,
            fish_type_id: current.fish_type_id || data[0]?.id || "",
          }));
        }
      } catch (err) {
        if (!ignore) setError(errorMessage(err));
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
    if (!form.fish_type_id) {
      return;
    }

    let ignore = false;

    async function loadFifoPreview() {
      setLoadingFifo(true);
      setFifoError(null);

      try {
        const data = await apiGet<FIFOStock[]>(
          `/stocks/fifo?fish_type_id=${form.fish_type_id}`,
        );
        if (!ignore) setFifoStocks(data);
      } catch (err) {
        if (!ignore) {
          setFifoStocks([]);
          setFifoError(errorMessage(err));
        }
      } finally {
        if (!ignore) setLoadingFifo(false);
      }
    }

    loadFifoPreview();
    return () => {
      ignore = true;
    };
  }, [form.fish_type_id]);

  const availableWeight = useMemo(
    () =>
      fifoStocks.reduce(
        (total, item) => total + Number(item.remaining_weight_kg),
        0,
      ),
    [fifoStocks],
  );

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setCreatedStockOut(null);
    setFieldErrors([]);

    try {
      const created = await apiPost<StockOut>("/stock-outs", {
        fish_type_id: form.fish_type_id,
        total_weight_kg: Number(form.total_weight_kg),
        destination: form.destination,
        out_at: new Date(form.out_at).toISOString(),
        notes: emptyToNull(form.notes),
      });
      setCreatedStockOut(created);
      setForm((current) => ({
        ...initialForm,
        fish_type_id: current.fish_type_id,
      }));
      await refreshFifo(form.fish_type_id);
    } catch (err) {
      setError(errorMessage(err));
      setFieldErrors(extractFieldErrors(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function refreshFifo(fishTypeID: string) {
    if (!fishTypeID) return;
    const data = await apiGet<FIFOStock[]>(
      `/stocks/fifo?fish_type_id=${fishTypeID}`,
    );
    setFifoStocks(data);
  }

  const isFishTypeEmpty = !loadingFishTypes && fishTypes.length === 0;

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Catat Ikan Keluar"
        description="Sistem akan mengambil batch terlama lebih dulu."
        actions={
          <Link
            href="/stocks"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-semibold text-slate-800"
          >
            FIFO
          </Link>
        }
      />

      {loadingFishTypes ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat jenis ikan...
        </div>
      ) : null}

      {isFishTypeEmpty ? (
        <div className="grid gap-3 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900">
          <p className="font-semibold">Jenis ikan belum tersedia</p>
          <p>Tambahkan jenis ikan sebelum mencatat pengeluaran.</p>
          <Link
            href="/fish-types"
            className="rounded-lg bg-white px-3 py-2 text-center font-semibold text-amber-900"
          >
            Tambah Jenis Ikan
          </Link>
        </div>
      ) : null}

      {!loadingFishTypes && !isFishTypeEmpty ? (
        <>
          <form
            onSubmit={handleSubmit}
            className="mb-5 grid gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
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
                value={form.fish_type_id}
                onChange={(event) => {
                  const fishTypeID = event.target.value;
                  setCreatedStockOut(null);
                  setFifoStocks([]);
                  setForm((current) => ({
                    ...current,
                    fish_type_id: fishTypeID,
                  }));
                }}
                className="h-11 rounded-lg border border-slate-300 bg-white px-3 text-sm outline-none focus:border-emerald-600"
                required
              >
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
                htmlFor="total_weight_kg"
              >
                Berat keluar kg
              </label>
              <input
                id="total_weight_kg"
                type="number"
                inputMode="decimal"
                min="0.01"
                step="0.01"
                value={form.total_weight_kg}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    total_weight_kg: event.target.value,
                  }))
                }
                className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
                placeholder="40"
                required
              />
              <p className="text-xs text-slate-500">
                Stok tersedia untuk jenis ini: {availableWeight.toFixed(2)} kg
              </p>
            </div>

            <div className="grid gap-1.5">
              <label
                className="text-sm font-medium text-slate-700"
                htmlFor="destination"
              >
                Tujuan pengeluaran
              </label>
              <input
                id="destination"
                value={form.destination}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    destination: event.target.value,
                  }))
                }
                className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
                placeholder="Restoran Laut Makassar"
                required
              />
            </div>

            <div className="grid gap-1.5">
              <label
                className="text-sm font-medium text-slate-700"
                htmlFor="out_at"
              >
                Waktu keluar
              </label>
              <input
                id="out_at"
                type="datetime-local"
                value={form.out_at}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    out_at: event.target.value,
                  }))
                }
                className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
                required
              />
            </div>

            <div className="grid gap-1.5">
              <label
                className="text-sm font-medium text-slate-700"
                htmlFor="notes"
              >
                Catatan
              </label>
              <textarea
                id="notes"
                value={form.notes}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    notes: event.target.value,
                  }))
                }
                className="min-h-24 rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-emerald-600"
                placeholder="Pengeluaran untuk pesanan"
              />
            </div>

            {fieldErrors.length > 0 ? (
              <ul className="grid gap-1 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-800">
                {fieldErrors.map((item) => (
                  <li key={`${item.field}-${item.message}`}>
                    {item.field}: {item.message}
                  </li>
                ))}
              </ul>
            ) : null}

            {error ? (
              <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
                {error}
              </div>
            ) : null}

            <button
              type="submit"
              disabled={submitting || fifoStocks.length === 0}
              className="h-11 rounded-lg bg-emerald-700 px-4 text-sm font-semibold text-white disabled:bg-slate-400"
            >
              {submitting ? "Menyimpan..." : "Catat Ikan Keluar"}
            </button>
          </form>

          {createdStockOut ? (
            <section className="mb-5 rounded-lg border border-emerald-200 bg-emerald-50 p-4">
              <h2 className="text-sm font-semibold text-emerald-950">
                Pengeluaran berhasil dicatat
              </h2>
              <p className="mt-1 text-sm text-emerald-900">
                Total {createdStockOut.total_weight_kg} kg ke{" "}
                {createdStockOut.destination}.
              </p>
              <div className="mt-3 grid gap-2">
                {createdStockOut.items.map((item) => (
                  <div
                    key={item.stock_batch_id}
                    className="flex justify-between gap-3 rounded-lg bg-white px-3 py-2 text-sm"
                  >
                    <span className="truncate text-slate-600">
                      Batch {shortID(item.stock_batch_id)}
                    </span>
                    <span className="font-semibold text-slate-950">
                      {item.weight_kg} kg
                    </span>
                  </div>
                ))}
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2">
                <Link
                  href="/stocks"
                  className="rounded-lg bg-white px-3 py-2 text-center text-sm font-semibold text-emerald-900"
                >
                  Lihat FIFO
                </Link>
                <Link
                  href="/dashboard"
                  className="rounded-lg bg-white px-3 py-2 text-center text-sm font-semibold text-emerald-900"
                >
                  Dashboard
                </Link>
              </div>
            </section>
          ) : null}

          <section className="grid gap-3">
            <div className="flex items-center justify-between gap-3">
              <h2 className="text-base font-semibold text-slate-950">
                Preview FIFO
              </h2>
              <span className="text-sm text-slate-500">
                {availableWeight.toFixed(2)} kg
              </span>
            </div>

            {loadingFifo ? (
              <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
                Memuat FIFO...
              </div>
            ) : null}

            {fifoError ? (
              <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
                {fifoError}
              </div>
            ) : null}

            {!loadingFifo && fifoStocks.length === 0 && !fifoError ? (
              <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center">
                <p className="text-sm font-semibold text-slate-900">
                  Tidak ada stok tersedia
                </p>
                <p className="mt-1 text-sm text-slate-600">
                  Tambahkan stok masuk untuk jenis ikan ini terlebih dulu.
                </p>
              </div>
            ) : null}

            {!loadingFifo && fifoStocks.length > 0 ? (
              <div className="grid gap-2">
                {fifoStocks.map((stock) => (
                  <article
                    key={stock.id}
                    className="rounded-lg border border-slate-200 bg-white p-3"
                  >
                    <div className="flex items-center gap-3">
                      <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-emerald-700 text-xs font-semibold text-white">
                        {stock.fifo_rank}
                      </span>
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-semibold text-slate-950">
                          {stock.cold_storage_name}
                        </p>
                        <p className="text-xs text-slate-500">
                          {formatDate(stock.entered_at)}
                        </p>
                      </div>
                      <p className="text-sm font-semibold text-slate-950">
                        {stock.remaining_weight_kg} kg
                      </p>
                    </div>
                  </article>
                ))}
              </div>
            ) : null}
          </section>
        </>
      ) : null}
    </div>
  );
}

function toDatetimeLocalValue(date: Date) {
  const offsetMs = date.getTimezoneOffset() * 60 * 1000;
  return new Date(date.getTime() - offsetMs).toISOString().slice(0, 16);
}

function emptyToNull(value: string) {
  const trimmed = value.trim();
  return trimmed === "" ? null : trimmed;
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

function shortID(value: string) {
  return value.slice(0, 8);
}

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Gagal memproses request.";
}

function extractFieldErrors(error: unknown): ApiFieldError[] {
  if (!(error instanceof ApiError) || !Array.isArray(error.errors)) {
    return [];
  }

  return error.errors.filter(
    (item): item is ApiFieldError =>
      Boolean(item) &&
      typeof item === "object" &&
      "field" in item &&
      "message" in item,
  );
}
