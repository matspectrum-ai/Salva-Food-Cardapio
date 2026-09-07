import { proxyToAPI } from "@/lib/bff";

export async function POST(request: Request) {
  const body = await request.text();
  return proxyToAPI("/api/v1/catalog/categories", {
    method: "POST",
    body,
  });
}
