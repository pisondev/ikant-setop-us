"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import { ApiError, apiGet, apiPost } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type {
  ApiFieldError,
  ColdStorage,
  FishQuality,
  FishType,
  StockBatch,
} from "@/types/api";

type StockInForm = {
  fish_type_id: string;
  cold_storage_id: string;
  quality: FishQuality;
  initial_weight_kg: string;
  entered_at: string;
  notes: string;
};

const initialForm: StockInForm = {
  fish_type_id: "",
  cold_storage_id: "",
  quality: "baik",
  initial_weight_kg: "",
  entered_at: toDatetimeLocalValue(new Date()),
  notes: "",
};

export default function NewStockPage() {
  const router = useRouter();
  const [fishTypes, setFishTypes] = useState<FishType[]>([]);
  const [coldStorages, setColdStorages] = useState<ColdStorage[]>([]);
  const [form, setForm] = useState<StockInForm>(initialForm);
  const [loadingMasters, setLoadingMasters] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<ApiFieldError[]>([]);

  useEffect(() => {
    let ignore = false;

    async function loadMasters() {
      try {
        const [fishTypeData, storageData] = await Promise.all([
          apiGet<FishType[]>("/fish-types"),
          apiGet<ColdStorage[]>("/cold-storages"),
        ]);

        if (!ignore) {
          setFishTypes(fishTypeData);
          setColdStorages(storageData);
          setForm((current) => ({
            ...current,
            fish_type_id: current.fish_type_id || fishTypeData[0]?.id || "",
            cold_storage_id:
              current.cold_storage_id || storageData[0]?.id || "",
          }));
        }
      } catch (err) {
        if (!ignore) setError(errorMessage(err));
      } finally {
        if (!ignore) setLoadingMasters(false);
      }
    }

    loadMasters();
    return () => {
      ignore = true;
    };
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setFieldErrors([]);

    try {
      await apiPost<StockBatch>("/stocks", {
        fish_type_id: form.fish_type_id,
        cold_storage_id: form.cold_storage_id,
        quality: form.quality,
        initial_weight_kg: Number(form.initial_weight_kg),
        entered_at: new Date(form.entered_at).toISOString(),
        notes: emptyToNull(form.notes),
      });

      router.push("/stocks");
      router.refresh();
    } catch (err) {
      setError(errorMessage(err));
      setFieldErrors(extractFieldErrors(err));
    } finally {
      setSubmitting(false);
    }
  }

  const isMasterDataEmpty = !loadingMasters && (fishTypes.length === 0 || coldStorages.length === 0);

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Tambah Stok Masuk"
        description="Catat batch ikan baru agar masuk ke urutan FIFO."
        actions={
          <Link
            href="/stocks"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-semibold text-slate-800"
          >
            FIFO
          </Link>
        }
      />

      {loadingMasters ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat data master...
        </div>
      ) : null}

      {isMasterDataEmpty ? (
        <div className="grid gap-3 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900">
          <p className="font-semibold">Data master belum lengkap</p>
          <p>
            Tambahkan minimal satu jenis ikan dan satu cold storage sebelum
            mencatat stok masuk.
          </p>
          <div className="grid grid-cols-2 gap-2">
            <Link
              href="/fish-types"
              className="rounded-lg bg-white px-3 py-2 text-center font-semibold text-amber-900"
            >
              Jenis Ikan
            </Link>
            <Link
              href="/cold-storages"
              className="rounded-lg bg-white px-3 py-2 text-center font-semibold text-amber-900"
            >
              Storage
            </Link>
          </div>
        </div>
      ) : null}

      {!loadingMasters && !isMasterDataEmpty ? (
        <form
          onSubmit={handleSubmit}
          className="grid gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
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
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  fish_type_id: event.target.value,
                }))
              }
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
              htmlFor="cold_storage_id"
            >
              Cold storage
            </label>
            <select
              id="cold_storage_id"
              value={form.cold_storage_id}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  cold_storage_id: event.target.value,
                }))
              }
              className="h-11 rounded-lg border border-slate-300 bg-white px-3 text-sm outline-none focus:border-emerald-600"
              required
            >
              {coldStorages.map((storage) => (
                <option key={storage.id} value={storage.id}>
                  {storage.name}
                  {storage.location_label ? ` - ${storage.location_label}` : ""}
                </option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="grid gap-1.5">
              <label
                className="text-sm font-medium text-slate-700"
                htmlFor="quality"
              >
                Kualitas
              </label>
              <select
                id="quality"
                value={form.quality}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    quality: event.target.value as FishQuality,
                  }))
                }
                className="h-11 rounded-lg border border-slate-300 bg-white px-3 text-sm outline-none focus:border-emerald-600"
              >
                <option value="baik">Baik</option>
                <option value="sedang">Sedang</option>
                <option value="buruk">Buruk</option>
              </select>
            </div>

            <div className="grid gap-1.5">
              <label
                className="text-sm font-medium text-slate-700"
                htmlFor="initial_weight_kg"
              >
                Berat kg
              </label>
              <input
                id="initial_weight_kg"
                type="number"
                inputMode="decimal"
                min="0.01"
                step="0.01"
                value={form.initial_weight_kg}
                onChange={(event) =>
                  setForm((current) => ({
                    ...current,
                    initial_weight_kg: event.target.value,
                  }))
                }
                className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
                placeholder="50"
                required
              />
            </div>
          </div>

          <div className="grid gap-1.5">
            <label
              className="text-sm font-medium text-slate-700"
              htmlFor="entered_at"
            >
              Waktu masuk
            </label>
            <input
              id="entered_at"
              type="datetime-local"
              value={form.entered_at}
              onChange={(event) =>
                setForm((current) => ({
                  ...current,
                  entered_at: event.target.value,
                }))
              }
              className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
              required
            />
          </div>

          <div className="grid gap-1.5">
            <label className="text-sm font-medium text-slate-700" htmlFor="notes">
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
              placeholder="Tangkapan pagi"
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
            disabled={submitting}
            className="h-11 rounded-lg bg-emerald-700 px-4 text-sm font-semibold text-white disabled:bg-slate-400"
          >
            {submitting ? "Menyimpan..." : "Simpan Stok Masuk"}
          </button>
        </form>
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
