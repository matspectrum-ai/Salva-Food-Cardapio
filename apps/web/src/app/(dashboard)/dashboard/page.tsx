import { getCategories, getItems, getOrders } from "@/lib/api";

const money = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });

export default async function DashboardPage() {
  const [orders, categories, items] = await Promise.all([getOrders(), getCategories(), getItems()]);
  const openOrders = orders.filter((order) => !["FINALIZED", "CANCELLED"].includes(order.status));
  const movedCents = orders.reduce((total, order) => total + order.total_cents, 0);
  const soldOut = items.filter((item) => item.sold_out).length;

  const metrics = [
    ["Pedidos abertos", String(openOrders.length)],
    ["Movimentado", money.format(movedCents / 100)],
    ["Categorias", String(categories.length)],
    ["Itens esgotados", String(soldOut)],
  ];

  return (
    <section className="space-y-7">
      <div>
        <p className="text-sm font-medium text-zinc-500">Visão operacional</p>
        <h1 className="mt-1 text-3xl font-bold tracking-tight">Dashboard</h1>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {metrics.map(([label, value]) => (
          <article className="card p-5" key={label}>
            <p className="text-sm text-zinc-500">{label}</p>
            <p className="metric mt-2">{value}</p>
          </article>
        ))}
      </div>
      <article className="card p-5">
        <div className="flex items-end justify-between gap-4">
          <div>
            <p className="text-sm text-zinc-500">Operação agora</p>
            <h2 className="mt-1 text-lg font-bold">Fila de pedidos</h2>
          </div>
          <span className="rounded-full bg-zinc-100 px-3 py-1 text-xs font-semibold text-zinc-600">
            {openOrders.length} ativos
          </span>
        </div>
        <div className="mt-5 grid gap-3 sm:grid-cols-3">
          {(["ANALYSIS", "PRODUCTION", "READY"] as const).map((status) => (
            <div className="rounded-xl border border-zinc-200 p-4" key={status}>
              <p className="text-xs font-semibold uppercase tracking-wide text-zinc-500">{status}</p>
              <p className="mt-2 text-2xl font-bold">{orders.filter((order) => order.status === status).length}</p>
            </div>
          ))}
        </div>
      </article>
    </section>
  );
}
