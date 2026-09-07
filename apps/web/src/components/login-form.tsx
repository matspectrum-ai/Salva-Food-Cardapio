"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

export function LoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true);
    setError("");
    const response = await fetch("/api/session/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!response.ok) {
      const payload = await response.json().catch(() => ({}));
      setError(payload?.error?.message ?? "Não foi possível entrar.");
      setPending(false);
      return;
    }
    router.replace("/dashboard");
    router.refresh();
  }

  return (
    <form className="space-y-5" onSubmit={onSubmit}>
      <label className="block space-y-2">
        <span className="text-sm font-medium text-zinc-700">E-mail</span>
        <input
          autoComplete="email"
          className="field"
          onChange={(event) => setEmail(event.target.value)}
          required
          type="email"
          value={email}
        />
      </label>
      <label className="block space-y-2">
        <span className="text-sm font-medium text-zinc-700">Senha</span>
        <input
          autoComplete="current-password"
          className="field"
          onChange={(event) => setPassword(event.target.value)}
          required
          type="password"
          value={password}
        />
      </label>
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <button className="button-primary w-full" disabled={pending} type="submit">
        {pending ? "Entrando…" : "Entrar"}
      </button>
    </form>
  );
}
