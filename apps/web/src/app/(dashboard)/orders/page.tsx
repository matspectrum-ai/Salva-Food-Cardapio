import { OrderTransitionButton } from "@/components/order-transition-button";
import { getOrders, type Order, type OrderStatus } from "@/lib/api";

const money = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const columns: Array<{ status: OrderStatus; title: string }> = [
  { status: "ANALYSIS", title: "Em análise" },
  { status: "PRODUCTION", title: "Em produção" },
  { status: "READY", title: "Prontos" },
];

function OrderCard({ order }: { order: Order }) {
  return (
    <article className="rounded-xl border border-zinc-200 bg-white p-4 shadow-sm">
      <div className="flex items-center justify-between gap-3">
        <span className="text-xs font-semibold uppercase tracking-wide text-zinc-500">{order.source}</span>
        <span className="text-sm font-bold">{money.format(order.total_cents / 100)}</span>
      </div>
      <p className="mt-3 font-mono text-xs text-zinc-500">#{order.id.slice(0, 8)}</p>
      <ul className="mt-3 space-y-1 text-sm text-zinc-700">
        {order.items.map((item, index) => (
          <li className="flex justify-between gap-3" key={`${item.item_id}-${index}`}>
            <span>{item.quantity}× {item.name}</span>
            <span>{money.format((item.quantity * item.unit_price_cents) / 100)}</span>
          </li>
        ))}
      </ul>
      <OrderTransitionButton orderId={order.id} status={order.status} />
    </article>
  );
}

export default async function OrdersPage() {
  const orders = await getOrders();

  return (
    <section className="space-y-7">
      <div>
        <p className="text-sm font-medium text-zinc-500">Operação</p>
        <h1 className="mt-1 text-3xl font-bold tracking-tight">Meus pedidos</h1>
      </div>
      <div className="grid gap-4 xl:grid-cols-3">
        {columns.map((column) => {
          const entries = orders.filter((order) => order.status === column.status);
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
    </section>
  );
}
