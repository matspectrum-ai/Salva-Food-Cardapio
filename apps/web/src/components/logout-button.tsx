"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export function LogoutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);

  async function logout() {
    setPending(true);
    await fetch("/api/session/logout", { method: "POST" });
    router.replace("/login");
    router.refresh();
  }

  return (
    <button className="button-secondary w-full" disabled={pending} onClick={logout} type="button">
      {pending ? "Saindo…" : "Sair"}
    </button>
  );
}
