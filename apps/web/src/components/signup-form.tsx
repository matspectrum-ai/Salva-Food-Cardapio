"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

export function SignupForm() {
  const router = useRouter();
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError("");
    const form = new FormData(event.currentTarget);
    const payload = {
      tenant_name: String(form.get("tenant_name") ?? ""),
      establishment_name: String(form.get("establishment_name") ?? ""),
      timezone: String(form.get("timezone") ?? ""),
      owner: {
        name: String(form.get("owner_name") ?? ""),
        cpf: String(form.get("cpf") ?? ""),
        email: String(form.get("email") ?? ""),
        phone: String(form.get("phone") ?? ""),
        password: String(form.get("password") ?? ""),
      },
    };

    const response = await fetch("/api/onboarding", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      setError(body?.error?.message ?? "Não foi possível criar sua conta.");
      setPending(false);
      return;
    }
    router.replace("/dashboard");
    router.refresh();
  }

  return (
    <form className="space-y-6" onSubmit={onSubmit}>
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="Nome da operação" name="tenant_name" placeholder="Minha operação" />
        <Field label="Nome do estabelecimento" name="establishment_name" placeholder="Unidade principal" />
      </div>
      <label className="block space-y-2">
        <span className="text-sm font-medium text-zinc-700">Fuso horário</span>
        <input className="field" defaultValue="America/Santarem" name="timezone" required />
      </label>

      <div className="border-t border-zinc-200 pt-6">
        <p className="mb-4 text-sm font-semibold text-zinc-900">Dados do proprietário</p>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Nome completo" name="owner_name" autoComplete="name" />
          <Field label="CPF" name="cpf" inputMode="numeric" />
          <Field label="E-mail" name="email" type="email" autoComplete="email" />
          <Field label="Telefone" name="phone" type="tel" autoComplete="tel" />
        </div>
        <label className="mt-4 block space-y-2">
          <span className="text-sm font-medium text-zinc-700">Senha</span>
          <input className="field" minLength={8} name="password" required type="password" autoComplete="new-password" />
          <span className="block text-xs leading-5 text-zinc-500">
            Use pelo menos 8 caracteres, com maiúscula, minúscula, número e símbolo.
          </span>
        </label>
      </div>
      {error ? <p className="rounded-xl bg-red-50 px-4 py-3 text-sm text-red-700">{error}</p> : null}
      <button className="button-primary w-full" disabled={pending} type="submit">
        {pending ? "Criando operação…" : "Criar operação"}
      </button>
    </form>
  );
}

type FieldProps = {
  label: string;
  name: string;
  type?: string;
  placeholder?: string;
  autoComplete?: string;
  inputMode?: "numeric" | "text" | "tel" | "email";
};

function Field({ label, name, type = "text", placeholder, autoComplete, inputMode }: FieldProps) {
  return (
    <label className="block space-y-2">
      <span className="text-sm font-medium text-zinc-700">{label}</span>
      <input
        autoComplete={autoComplete}
        className="field"
        inputMode={inputMode}
        name={name}
        placeholder={placeholder}
        required
        type={type}
      />
    </label>
  );
}
