"use client";

import { useState } from "react";
import type { OrderStatus } from "@/lib/api";

const nextStatus: Partial<Record<OrderStatus, OrderStatus>> = {
  ANALYSIS: "PRODUCTION",
  PRODUCTION: "READY",
  READY: "FINALIZED",
};

const labels: Partial<Record<OrderStatus, string>> = {
  ANALYSIS: "Iniciar preparo",
  PRODUCTION: "Marcar como pronto",
  READY: "Finalizar pedido",
};

export function OrderTransitionButton({ orderId, status }: { orderId: string; status: OrderStatus }) {
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  const target = nextStatus[status];
  if (!target) return null;

  async function transition() {
    setPending(true);
    setError("");
    try {
      const response = await fetch(`/api/orders/${orderId}/transition`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: target }),
      });
      if (!response.ok) {
        const payload = await response.json().catch(() => ({}));
        setError(payload?.error?.message ?? "Falha ao atualizar pedido.");
      }
    } catch {
      setError("Falha de rede ao atualizar pedido.");
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="mt-4 space-y-2">
      <button className="button-primary w-full" disabled={pending} onClick={transition} type="button">
        {pending ? "Atualizando…" : labels[status]}
      </button>
      {error ? <p className="text-xs text-red-600">{error}</p> : null}
    </div>
  );
}
