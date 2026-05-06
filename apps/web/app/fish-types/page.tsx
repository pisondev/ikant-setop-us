"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { ApiError, apiGet, apiPost } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type { ApiFieldError, FishType } from "@/types/api";

type FishTypeForm = {
  name: string;
  image_url: string;
  description: string;
};

const initialForm: FishTypeForm = {
  name: "",
  image_url: "",
  description: "",
};

export default function FishTypesPage() {
  const [items, setItems] = useState<FishType[]>([]);
  const [form, setForm] = useState<FishTypeForm>(initialForm);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<ApiFieldError[]>([]);

  async function loadItems() {
    setLoading(true);
    setError(null);

    try {
      const data = await apiGet<FishType[]>("/fish-types");
      setItems(data);
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    let ignore = false;

    async function loadInitialItems() {
      try {
        const data = await apiGet<FishType[]>("/fish-types");
        if (!ignore) setItems(data);
      } catch (err) {
        if (!ignore) setError(errorMessage(err));
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    loadInitialItems();
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
      await apiPost<FishType>("/fish-types", {
        name: form.name,
        image_url: emptyToNull(form.image_url),
        description: emptyToNull(form.description),
      });
      setForm(initialForm);
      await loadItems();
    } catch (err) {
      setError(errorMessage(err));
      setFieldErrors(extractFieldErrors(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Jenis Ikan"
        description="Master data untuk pilihan stok masuk dan ikan keluar."
        actions={
          <Link
            href="/cold-storages"
            className="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-semibold text-slate-800"
          >
            Storage
          </Link>
        }
      />

      <form
        onSubmit={handleSubmit}
        className="mb-5 grid gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
      >
        <div className="grid gap-1.5">
          <label className="text-sm font-medium text-slate-700" htmlFor="name">
            Nama jenis ikan
          </label>
          <input
            id="name"
            value={form.name}
            onChange={(event) =>
              setForm((current) => ({ ...current, name: event.target.value }))
            }
            className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
            placeholder="Tuna"
            required
          />
        </div>

        <div className="grid gap-1.5">
          <label
            className="text-sm font-medium text-slate-700"
            htmlFor="image_url"
          >
            URL gambar
          </label>
          <input
            id="image_url"
            value={form.image_url}
            onChange={(event) =>
              setForm((current) => ({
                ...current,
                image_url: event.target.value,
              }))
            }
            className="h-11 rounded-lg border border-slate-300 px-3 text-sm outline-none focus:border-emerald-600"
            placeholder="/images/fish/tuna.png"
          />
        </div>

        <div className="grid gap-1.5">
          <label
            className="text-sm font-medium text-slate-700"
            htmlFor="description"
          >
            Deskripsi
          </label>
          <textarea
            id="description"
            value={form.description}
            onChange={(event) =>
              setForm((current) => ({
                ...current,
                description: event.target.value,
              }))
            }
            className="min-h-24 rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-emerald-600"
            placeholder="Ikan tuna"
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

        <button
          type="submit"
          disabled={submitting}
          className="h-11 rounded-lg bg-emerald-700 px-4 text-sm font-semibold text-white disabled:bg-slate-400"
        >
          {submitting ? "Menyimpan..." : "Tambah Jenis Ikan"}
        </button>
      </form>

      {error ? (
        <div className="mb-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat jenis ikan...
        </div>
      ) : null}

      {!loading && items.length === 0 && !error ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center">
          <p className="text-sm font-semibold text-slate-900">
            Belum ada jenis ikan
          </p>
          <p className="mt-1 text-sm text-slate-600">
            Tambahkan jenis ikan pertama untuk dipakai pada stok masuk.
          </p>
        </div>
      ) : null}

      {!loading && items.length > 0 ? (
        <div className="grid gap-3">
          {items.map((item) => (
            <article
              key={item.id}
              className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
            >
              <h2 className="text-base font-semibold text-slate-950">
                {item.name}
              </h2>
              {item.description ? (
                <p className="mt-1 text-sm text-slate-600">
                  {item.description}
                </p>
              ) : null}
              {item.image_url ? (
                <p className="mt-2 break-all text-xs text-slate-500">
                  {item.image_url}
                </p>
              ) : null}
            </article>
          ))}
        </div>
      ) : null}
    </div>
  );
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
