import Link from "next/link";
import type { ReactNode } from "react";
import { getPrincipal } from "@/lib/api";
import { LogoutButton } from "./logout-button";

const navigation = [
  ["Dashboard", "/dashboard"],
  ["Pedidos", "/orders"],
  ["Cardápio", "/menu"],
] as const;

export async function AppShell({ children }: { children: ReactNode }) {
  const principal = await getPrincipal();

  return (
    <div className="min-h-screen bg-zinc-50 text-zinc-950">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-zinc-200 bg-white p-5 lg:block">
        <Link className="mb-8 block text-xl font-bold tracking-tight" href="/dashboard">Salva Food</Link>
        <nav className="space-y-1">
          {navigation.map(([label, href]) => (
            <Link className="nav-link" href={href} key={href}>{label}</Link>
          ))}
        </nav>
        <div className="absolute bottom-5 left-5 right-5 space-y-3 border-t border-zinc-200 pt-4">
          <div>
            <p className="truncate text-sm font-medium">{principal.user.name}</p>
            <p className="truncate text-xs text-zinc-500">{principal.membership.title}</p>
          </div>
          <LogoutButton />
        </div>
      </aside>
      <div className="lg:pl-64">
        <header className="flex min-h-16 items-center justify-between border-b border-zinc-200 bg-white px-5 lg:px-8">
          <Link className="font-bold lg:hidden" href="/dashboard">Salva Food</Link>
          <p className="ml-auto text-sm text-zinc-500">{principal.user.email}</p>
        </header>
        <main className="p-5 lg:p-8">{children}</main>
      </div>
    </div>
  );
}
