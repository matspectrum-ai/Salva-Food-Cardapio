import "server-only";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export const API_URL = (process.env.SALVA_API_URL ?? "http://127.0.0.1:8080").replace(/\/$/, "");
export const SESSION_COOKIE = process.env.SALVA_SESSION_COOKIE ?? "salva_food_session";

export type Permission = string;
export type Principal = {
  user: { id: string; name: string; cpf: string; email: string; phone: string; image_url?: string };
  membership: { tenant_id: string; user_id: string; title: string; status: "ACTIVE" | "INACTIVE"; permissions: Permission[] };
};

export type Category = {
  tenant_id: string; id: string; name: string; sort_order: number; sold_out: boolean;
};

export type Item = {
  tenant_id: string; id: string; category_id: string; name: string;
  price_cents: number; sort_order: number; sold_out: boolean; sold_by_weight: boolean;
};

export type OrderStatus = "ANALYSIS" | "PRODUCTION" | "READY" | "FINALIZED" | "CANCELLED";
export type Order = {
  tenant_id: string; id: string; source: string; status: OrderStatus;
  items: Array<{ item_id: string; name: string; quantity: number; unit_price_cents: number }>;
  total_cents: number;
};

async function sessionToken() {
  return (await cookies()).get(SESSION_COOKIE)?.value ?? "";
}

export async function apiFetch(path: string, init: RequestInit = {}) {
  const token = await sessionToken();
  if (!token) redirect("/login");

  const headers = new Headers(init.headers);
  headers.set("Authorization", `Bearer ${token}`);
  if (init.body && !headers.has("Content-Type")) headers.set("Content-Type", "application/json");

  const response = await fetch(`${API_URL}${path}`, { ...init, headers, cache: "no-store" });
  if (response.status === 401) redirect("/login");
  if (!response.ok) {
    const body = await response.text();
    throw new Error(`API ${response.status}: ${body || response.statusText}`);
  }
  return response;
}

export async function getPrincipal(): Promise<Principal> {
  return (await apiFetch("/api/v1/me")).json();
}

export async function getOrders(): Promise<Order[]> {
  const payload = await (await apiFetch("/api/v1/orders")).json() as { data: Order[] };
  return payload.data;
}

export async function getCategories(): Promise<Category[]> {
  const payload = await (await apiFetch("/api/v1/catalog/categories")).json() as { data: Category[] };
  return payload.data;
}

export async function getItems(): Promise<Item[]> {
  const payload = await (await apiFetch("/api/v1/catalog/items")).json() as { data: Item[] };
  return payload.data;
}
