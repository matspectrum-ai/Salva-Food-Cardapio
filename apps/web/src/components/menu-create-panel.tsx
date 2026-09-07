"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import type { Category } from "@/lib/api";

function parsePriceToCents(value: string) {
  const normalized = value.trim().replace(/\./g, "").replace(",", ".");
  const amount = Number(normalized);
  if (!Number.isFinite(amount) || amount < 0) return null;
  return Math.round(amount * 100);
}

export function MenuCreatePanel({ categories }: { categories: Category[] }) {
  const router = useRouter();
  const [categoryName, setCategoryName] = useState("");
  const [itemName, setItemName] = useState("");
  const [categoryId, setCategoryId] = useState(categories[0]?.id ?? "");
  const [price, setPrice] = useState("");
  const [soldByWeight, setSoldByWeight] = useState(false);
  const [pending, setPending] = useState<"category" | "item" | "">("");
  const [error, setError] = useState("");

  async function createCategory(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending("category");
    setError("");
    const response = await fetch("/api/menu/categories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: categoryName, sort_order: categories.length }),
    });
    if (!response.ok) {
      const payload = await response.json().catch(() => ({}));
      setError(payload?.error?.message ?? "Falha ao criar categoria.");
      setPending("");
      return;
    }
    setCategoryName("");
    setPending("");
    router.refresh();
  }

  async function createItem(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const priceCents = parsePriceToCents(price);
    if (priceCents === null) {
      setError("Informe um preço válido.");
      return;
    }
    setPending("item");
    setError("");
    const response = await fetch("/api/menu/items", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        category_id: categoryId,
        name: itemName,
        price_cents: priceCents,
        sort_order: 0,
        sold_by_weight: soldByWeight,
      }),
    });
    if (!response.ok) {
      const payload = await response.json().catch(() => ({}));
      setError(payload?.error?.message ?? "Falha ao criar item.");
      setPending("");
      return;
    }
    setItemName("");
    setPrice("");
    setSoldByWeight(false);
    setPending("");
    router.refresh();
  }

  return (
    <div className="grid gap-4 xl:grid-cols-2">
      <form className="card p-5" onSubmit={createCategory}>
        <p className="text-sm font-semibold">Nova categoria</p>
        <div className="mt-4 flex gap-3">
          <input className="field" onChange={(event) => setCategoryName(event.target.value)} placeholder="Ex.: Hambúrgueres" required value={categoryName} />
          <button className="button-primary shrink-0" disabled={pending !== ""} type="submit">{pending === "category" ? "Criando…" : "Criar"}</button>
        </div>
      </form>

      <form className="card p-5" onSubmit={createItem}>
        <p className="text-sm font-semibold">Novo item</p>
        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          <input className="field" onChange={(event) => setItemName(event.target.value)} placeholder="Nome do item" required value={itemName} />
          <select className="field" disabled={!categories.length} onChange={(event) => setCategoryId(event.target.value)} required value={categoryId}>
            {categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
          </select>
          <input className="field" inputMode="decimal" onChange={(event) => setPrice(event.target.value)} placeholder="Preço, ex.: 29,90" required value={price} />
          <label className="flex items-center gap-2 rounded-xl border border-zinc-200 px-3 text-sm">
            <input checked={soldByWeight} onChange={(event) => setSoldByWeight(event.target.checked)} type="checkbox" />
            Vendido por peso
          </label>
        </div>
        <button className="button-primary mt-3" disabled={pending !== "" || !categories.length} type="submit">{pending === "item" ? "Criando…" : "Adicionar item"}</button>
      </form>
      {error ? <p className="text-sm text-red-600 xl:col-span-2">{error}</p> : null}
    </div>
  );
}
