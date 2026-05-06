import type { ReactNode } from "react";
import { BottomNav } from "@/components/layout/bottom-nav";

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-dvh bg-slate-50 text-slate-950">
      <main className="mx-auto flex min-h-dvh w-full max-w-lg flex-col px-4 pb-24 pt-5">
        {children}
      </main>
      <BottomNav />
    </div>
  );
}
