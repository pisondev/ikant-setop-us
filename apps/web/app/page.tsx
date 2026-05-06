import Link from "next/link";
import { PageHeader } from "@/components/layout/page-header";

export default function Home() {
  return (
    <div className="flex flex-1 flex-col">
      <PageHeader
        title="Ikan't Setop Us"
        description="Inventori ikan cold storage berbasis FIFO."
      />

      <section className="grid gap-3">
        <Link
          href="/stocks"
          className="rounded-lg bg-emerald-700 px-4 py-3 text-center text-sm font-semibold text-white"
        >
          Lihat Stok FIFO
        </Link>
        <Link
          href="/stocks/new"
          className="rounded-lg border border-slate-300 bg-white px-4 py-3 text-center text-sm font-semibold text-slate-800"
        >
          Tambah Stok Masuk
        </Link>
        <Link
          href="/stock-outs/new"
          className="rounded-lg border border-slate-300 bg-white px-4 py-3 text-center text-sm font-semibold text-slate-800"
        >
          Catat Ikan Keluar
        </Link>
      </section>
    </div>
  );
}
