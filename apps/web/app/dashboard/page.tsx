"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { ApiError, apiGet } from "@/lib/api";
import { PageHeader } from "@/components/layout/page-header";
import type { DashboardSummary, RecentMovement } from "@/types/api";

export default function DashboardPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [movements, setMovements] = useState<RecentMovement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let ignore = false;

    async function loadDashboard() {
      try {
        const [summaryData, movementData] = await Promise.all([
          apiGet<DashboardSummary>("/dashboard/summary"),
          apiGet<RecentMovement[]>("/dashboard/recent-movements?limit=10"),
        ]);

        if (!ignore) {
          setSummary(summaryData);
          setMovements(movementData);
        }
      } catch (err) {
        if (!ignore) setError(errorMessage(err));
      } finally {
        if (!ignore) setLoading(false);
      }
    }

    loadDashboard();
    return () => {
      ignore = true;
    };
  }, []);

  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Dashboard"
        description="Ringkasan stok ikan dan aktivitas terbaru."
        actions={
          <Link
            href="/stocks/new"
            className="rounded-lg bg-emerald-700 px-3 py-2 text-sm font-semibold text-white"
          >
            Tambah
          </Link>
        }
      />

      {loading ? (
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-8 text-center text-sm text-slate-600">
          Memuat dashboard...
        </div>
      ) : null}

      {error ? (
        <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-800">
          {error}
        </div>
      ) : null}

      {!loading && summary ? (
        <div className="grid gap-5">
          <section className="grid gap-3">
            <SummaryCard
              label="Total stok tersedia"
              value={`${formatWeight(summary.total_available_weight_kg)} kg`}
              strong
            />
            <div className="grid grid-cols-2 gap-3">
              <SummaryCard
                label="Masuk hari ini"
                value={`${formatWeight(summary.today_stock_in_kg)} kg`}
              />
              <SummaryCard
                label="Keluar hari ini"
                value={`${formatWeight(summary.today_stock_out_kg)} kg`}
              />
            </div>
            <div className="grid grid-cols-3 gap-3">
              <SummaryCard
                label="Batch"
                value={summary.total_stock_batches.toString()}
              />
              <SummaryCard
                label="Available"
                value={summary.total_available_batches.toString()}
              />
              <SummaryCard
                label="Depleted"
                value={summary.total_depleted_batches.toString()}
              />
            </div>
          </section>

          <section className="grid gap-3">
            <SectionTitle title="Stok per Jenis Ikan" />
            {summary.fish_type_summary.length === 0 ? (
              <EmptyBlock message="Belum ada stok available per jenis ikan." />
            ) : (
              <div className="grid gap-2">
                {summary.fish_type_summary.map((item) => (
                  <ListRow
                    key={item.fish_type_id}
                    label={item.fish_type_name}
                    meta={`${item.available_batches} batch`}
                    value={`${formatWeight(item.available_weight_kg)} kg`}
                  />
                ))}
              </div>
            )}
          </section>

          <section className="grid gap-3">
            <SectionTitle title="Stok per Cold Storage" />
            {summary.cold_storage_summary.length === 0 ? (
              <EmptyBlock message="Belum ada stok available per storage." />
            ) : (
              <div className="grid gap-2">
                {summary.cold_storage_summary.map((item) => (
                  <ListRow
                    key={item.cold_storage_id}
                    label={item.cold_storage_name}
                    meta={`${item.available_batches} batch`}
                    value={`${formatWeight(item.available_weight_kg)} kg`}
                  />
                ))}
              </div>
            )}
          </section>

          <section className="grid gap-3">
            <div className="flex items-center justify-between gap-3">
              <SectionTitle title="Recent Movements" />
              <Link
                href="/stocks"
                className="text-sm font-semibold text-emerald-700"
              >
                FIFO
              </Link>
            </div>
            {movements.length === 0 ? (
              <EmptyBlock message="Belum ada aktivitas stok." />
            ) : (
              <div className="grid gap-2">
                {movements.map((movement) => (
                  <article
                    key={movement.id}
                    className="rounded-lg border border-slate-200 bg-white p-3 shadow-sm"
                  >
                    <div className="mb-1 flex items-center justify-between gap-3">
                      <span className="rounded-full bg-slate-100 px-2.5 py-1 text-xs font-semibold text-slate-700">
                        {movementLabel(movement.movement_type)}
                      </span>
                      <span className="text-xs text-slate-500">
                        {formatDate(movement.created_at)}
                      </span>
                    </div>
                    <p className="text-sm font-semibold text-slate-950">
                      {movement.fish_type_name}
                      {movement.weight_kg ? ` - ${movement.weight_kg} kg` : ""}
                    </p>
                    <p className="mt-1 text-sm text-slate-600">
                      {movement.description || "Aktivitas stok tercatat."}
                    </p>
                  </article>
                ))}
              </div>
            )}
          </section>

          <section className="grid grid-cols-2 gap-3">
            <Link
              href="/stocks/new"
              className="rounded-lg border border-slate-300 bg-white px-3 py-3 text-center text-sm font-semibold text-slate-800"
            >
              Stok Masuk
            </Link>
            <Link
              href="/stock-outs/new"
              className="rounded-lg border border-slate-300 bg-white px-3 py-3 text-center text-sm font-semibold text-slate-800"
            >
              Ikan Keluar
            </Link>
          </section>
        </div>
      ) : null}
    </div>
  );
}

function SummaryCard({
  label,
  value,
  strong,
}: {
  label: string;
  value: string;
  strong?: boolean;
}) {
  return (
    <div className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
      <p className="text-xs font-medium text-slate-500">{label}</p>
      <p
        className={[
          "mt-1 font-semibold text-slate-950",
          strong ? "text-2xl" : "text-lg",
        ].join(" ")}
      >
        {value}
      </p>
    </div>
  );
}

function SectionTitle({ title }: { title: string }) {
  return <h2 className="text-base font-semibold text-slate-950">{title}</h2>;
}

function ListRow({
  label,
  meta,
  value,
}: {
  label: string;
  meta: string;
  value: string;
}) {
  return (
    <article className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
      <div className="min-w-0">
        <p className="truncate text-sm font-semibold text-slate-950">{label}</p>
        <p className="text-xs text-slate-500">{meta}</p>
      </div>
      <p className="shrink-0 text-sm font-semibold text-slate-950">{value}</p>
    </article>
  );
}

function EmptyBlock({ message }: { message: string }) {
  return (
    <div className="rounded-lg border border-slate-200 bg-white px-4 py-6 text-center text-sm text-slate-600">
      {message}
    </div>
  );
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
    hour: "2-digit",
    minute: "2-digit",
  });
}

function movementLabel(value: string): string {
  if (value === "in") return "Masuk";
  if (value === "out") return "Keluar";
  if (value === "quality_update") return "Kualitas";
  if (value === "location_update") return "Lokasi";
  return "Adjustment";
}

function errorMessage(error: unknown): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Gagal memuat dashboard.";
}
