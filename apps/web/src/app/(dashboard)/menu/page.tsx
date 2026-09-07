import { getCategories, getItems } from "@/lib/api";

const money = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });

export default async function MenuPage() {
  const [categories, items] = await Promise.all([getCategories(), getItems()]);

  return (
    <section className="space-y-7">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-zinc-500">Catálogo</p>
          <h1 className="mt-1 text-3xl font-bold tracking-tight">Gestor de cardápio</h1>
        </div>
        <p className="text-sm text-zinc-500">{categories.length} categorias · {items.length} itens</p>
      </div>
      <div className="space-y-5">
        {categories.map((category) => {
          const categoryItems = items.filter((item) => item.category_id === category.id);
          return (
            <article className="card overflow-hidden" key={category.id}>
              <header className="flex items-center justify-between gap-4 border-b border-zinc-200 px-5 py-4">
                <div>
                  <h2 className="font-bold">{category.name}</h2>
                  <p className="mt-1 text-xs text-zinc-500">Ordem {category.sort_order} · {categoryItems.length} itens</p>
                </div>
                {category.sold_out ? <span className="rounded-full bg-red-50 px-3 py-1 text-xs font-semibold text-red-700">Categoria esgotada</span> : null}
              </header>
              <div className="divide-y divide-zinc-100">
                {categoryItems.map((item) => (
                  <div className="flex flex-wrap items-center justify-between gap-4 px-5 py-4" key={item.id}>
                    <div>
                      <p className="font-semibold">{item.name}</p>
                      <p className="mt-1 text-xs text-zinc-500">Ordem {item.sort_order}{item.sold_by_weight ? " · vendido por peso" : ""}</p>
                    </div>
                    <div className="flex items-center gap-3">
                      {item.sold_out ? <span className="rounded-full bg-amber-50 px-2.5 py-1 text-xs font-semibold text-amber-700">Esgotado</span> : null}
                      <span className="min-w-24 text-right font-bold">{money.format(item.price_cents / 100)}</span>
                    </div>
                  </div>
                ))}
                {!categoryItems.length ? <p className="px-5 py-6 text-sm text-zinc-500">Nenhum item nesta categoria.</p> : null}
              </div>
            </article>
          );
        })}
        {!categories.length ? <div className="card p-8 text-center text-sm text-zinc-500">O cardápio ainda não possui categorias.</div> : null}
      </div>
    </section>
  );
}
