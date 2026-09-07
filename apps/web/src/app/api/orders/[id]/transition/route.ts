import { proxyToAPI } from "@/lib/bff";

export async function POST(request: Request, context: { params: Promise<{ id: string }> }) {
  const { id } = await context.params;
  const body = await request.text();
  return proxyToAPI(`/api/v1/orders/${encodeURIComponent(id)}/transitions`, {
    method: "POST",
    body,
  });
}
