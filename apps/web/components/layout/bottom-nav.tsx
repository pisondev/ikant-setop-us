"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { href: "/stocks", label: "Stok", active: ["/stocks"] },
  { href: "/stocks/new", label: "Masuk", active: ["/stocks/new"] },
  { href: "/stock-outs/new", label: "Keluar", active: ["/stock-outs"] },
  { href: "/dashboard", label: "Dashboard", active: ["/dashboard"] },
  {
    href: "/fish-types",
    label: "Master",
    active: ["/fish-types", "/cold-storages"],
  },
];

export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="fixed inset-x-0 bottom-0 z-20 border-t border-slate-200 bg-white/95 shadow-sm backdrop-blur">
      <div className="mx-auto grid h-16 max-w-lg grid-cols-5 px-2">
        {navItems.map((item) => {
          const isActive = item.active.some(
            (activePath) =>
              pathname === activePath || pathname.startsWith(`${activePath}/`),
          );

          return (
            <Link
              key={item.href}
              href={item.href}
              className={[
                "flex min-w-0 flex-col items-center justify-center gap-1 px-1 text-xs font-medium",
                isActive ? "text-emerald-700" : "text-slate-500",
              ].join(" ")}
            >
              <span
                className={[
                  "h-1.5 w-1.5 rounded-full",
                  isActive ? "bg-emerald-600" : "bg-transparent",
                ].join(" ")}
              />
              <span className="truncate">{item.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
