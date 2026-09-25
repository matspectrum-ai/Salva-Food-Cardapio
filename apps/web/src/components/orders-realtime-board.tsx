"use client";

import { useEffect, useMemo, useState } from "react";
import { OrderTransitionButton } from "@/components/order-transition-button";
import type { Order, OrderStatus } from "@/lib/api";

const money = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const columns: Array<{ status: OrderStatus; title: string }> = [
  { status: "ANALYSIS", title: "Em análise" },
  { status: "PRODUCTION", title: "Em produção" },
  { status: "READY", title: "Prontos" },
];

function upsertOrder(current: Order[], next: Order): Order[] {
  const index = current.findIndex((order) => order.id === next.id);
  if (index < 0) return [next, ...current];
  const copy = [...current];
  copy[index] = next;
  return copy;
}

function OrderCard({ order }: { order: Order }) {
  const fulfillment = { DELIVERY: "Delivery", PICKUP: "Retirada", DINE_IN: "Salão" }[order.fulfillment_type];
  return (
    <article className="rounded-xl border border-zinc-200 bg-white p-4 shadow-sm">
      <div className="flex items-center justify-between gap-3">
        <span className="text-xs font-semibold uppercase tracking-wide text-zinc-500">{order.source}</span>
        <span className="text-sm font-bold">{money.format(order.total_cents / 100)}</span>
      </div>
      <p className="mt-3 font-mono text-xs text-zinc-500">#{order.id.slice(0, 8)}</p>
      {order.customer ? <p className="mt-2 text-sm font-semibold text-zinc-900">{order.customer.name}</p> : null}
      <div className="mt-1 flex flex-wrap gap-2 text-xs text-zinc-500">
        <span>{fulfillment}</span>
        {order.scheduled_at ? <span>Agendado {new Date(order.scheduled_at).toLocaleString("pt-BR")}</span> : null}
      </div>
      <ul className="mt-3 space-y-1 text-sm text-zinc-700">
        {order.items.map((item, index) => (
          <li className="flex justify-between gap-3" key={`${item.item_id}-${index}`}>
            <span>{item.quantity}× {item.name}</span>
            <span>{money.format((item.quantity * item.unit_price_cents) / 100)}</span>
          </li>
        ))}
      </ul>
      {order.notes ? <p className="mt-3 rounded-lg bg-zinc-50 p-2 text-xs text-zinc-600">Obs.: {order.notes}</p> : null}
      <OrderTransitionButton orderId={order.id} status={order.status} />
    </article>
  );
}

export function OrdersRealtimeBoard({ initialOrders }: { initialOrders: Order[] }) {
  const [orders, setOrders] = useState(initialOrders);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const source = new EventSource("/api/orders/events");
    const applySnapshot = (event: MessageEvent<string>) => {
      try { setOrders(JSON.parse(event.data) as Order[]); } catch { /* reconnect keeps stream alive */ }
    };
    const applyOrder = (event: MessageEvent<string>) => {
      try { setOrders((current) => upsertOrder(current, JSON.parse(event.data) as Order)); } catch { /* ignore malformed event */ }
    };
    source.addEventListener("open", () => setConnected(true));
    source.addEventListener("error", () => setConnected(false));
    source.addEventListener("orders.snapshot", applySnapshot as EventListener);
    source.addEventListener("order.created", applyOrder as EventListener);
    source.addEventListener("order.updated", applyOrder as EventListener);
    return () => source.close();
  }, []);

  const visible = useMemo(() => orders.filter((order) => ["ANALYSIS", "PRODUCTION", "READY"].includes(order.status)), [orders]);

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-xs text-zinc-500">
        <span className={`h-2 w-2 rounded-full ${connected ? "bg-emerald-500" : "bg-zinc-300"}`} />
        {connected ? "Atualização em tempo real ativa" : "Conectando ao tempo real…"}
      </div>
      <div className="grid gap-4 xl:grid-cols-3">
        {columns.map((column) => {
          const entries = visible.filter((order) => order.status === column.status);
          return (
            <section className="rounded-2xl bg-zinc-100 p-3" key={column.status}>
              <header className="mb-3 flex items-center justify-between px-1 py-1">
                <h2 className="font-bold">{column.title}</h2>
                <span className="rounded-full bg-white px-2.5 py-1 text-xs font-semibold text-zinc-600">{entries.length}</span>
              </header>
              <div className="space-y-3">
                {entries.length ? entries.map((order) => <OrderCard key={order.id} order={order} />) : (
                  <div className="rounded-xl border border-dashed border-zinc-300 p-6 text-center text-sm text-zinc-500">Nenhum pedido nesta etapa.</div>
                )}
              </div>
            </section>
          );
        })}
      </div>
    </div>
  );
}

