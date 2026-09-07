import { LoginForm } from "@/components/login-form";

export default function LoginPage() {
  return (
    <main className="grid min-h-screen place-items-center bg-zinc-100 px-5 py-10">
      <section className="card w-full max-w-md p-7 sm:p-9">
        <div className="mb-8">
          <p className="mb-2 text-sm font-semibold text-zinc-500">Salva Food</p>
          <h1 className="text-3xl font-bold tracking-tight">Entre na sua operação</h1>
          <p className="mt-3 text-sm leading-6 text-zinc-500">
            Acesse pedidos, cardápio e operação do estabelecimento em um único painel.
          </p>
        </div>
        <LoginForm />
      </section>
    </main>
  );
}
