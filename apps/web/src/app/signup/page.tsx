import Link from "next/link";
import { SignupForm } from "@/components/signup-form";

export default function SignupPage() {
  return (
    <main className="min-h-screen bg-zinc-100 px-5 py-10">
      <section className="card mx-auto w-full max-w-2xl p-7 sm:p-9">
        <div className="mb-8">
          <p className="mb-2 text-sm font-semibold text-zinc-500">Salva Food</p>
          <h1 className="text-3xl font-bold tracking-tight">Crie sua operação</h1>
          <p className="mt-3 max-w-xl text-sm leading-6 text-zinc-500">
            Criaremos sua organização, a primeira unidade e o acesso de proprietário em uma única operação.
          </p>
        </div>
        <SignupForm />
        <p className="mt-6 text-center text-sm text-zinc-500">
          Já possui acesso?{" "}
          <Link className="font-semibold text-zinc-900 underline-offset-4 hover:underline" href="/login">
            Entrar
          </Link>
        </p>
      </section>
    </main>
  );
}
