import "server-only";

import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import { API_URL, SESSION_COOKIE } from "./api";

export async function proxyToAPI(path: string, init: RequestInit) {
  const token = (await cookies()).get(SESSION_COOKIE)?.value;
  if (!token) {
    return NextResponse.json({ error: { message: "unauthorized" } }, { status: 401 });
  }

  const headers = new Headers(init.headers);
  headers.set("Authorization", `Bearer ${token}`);
  if (init.body && !headers.has("Content-Type")) headers.set("Content-Type", "application/json");

  const upstream = await fetch(`${API_URL}${path}`, { ...init, headers, cache: "no-store" });
  const body = await upstream.text();
  const responseHeaders = new Headers();
  const contentType = upstream.headers.get("content-type");
  if (contentType) responseHeaders.set("content-type", contentType);
  const replay = upstream.headers.get("idempotent-replay");
  if (replay) responseHeaders.set("idempotent-replay", replay);

  return new NextResponse(body || null, { status: upstream.status, headers: responseHeaders });
}
