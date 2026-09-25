import { OrdersRealtimeBoard } from "@/components/orders-realtime-board";
import { getOrders } from "@/lib/api";

export default async function OrdersPage() {
  const orders = await getOrders();

  return (
    <section className="space-y-7">
      <div>
        <p className="text-sm font-medium text-zinc-500">Operação</p>
        <div className="mt-1 flex items-end justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Meus pedidos</h1>
            <p className="mt-1 text-sm text-zinc-500">Quadro sincronizado com os demais terminais.</p>
          </div>
        </div>
      </div>
      <OrdersRealtimeBoard initialOrders={orders} />
    </section>
  );
}
