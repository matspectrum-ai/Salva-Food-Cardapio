import { NextResponse } from "next/server";
import { API_URL, SESSION_COOKIE } from "@/lib/api";

export async function POST(request: Request) {
  const input = await request.json();
  const upstream = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
    cache: "no-store",
  });

  const payload = await upstream.json().catch(() => ({}));
  if (!upstream.ok) {
    return NextResponse.json(payload, { status: upstream.status });
  }

  const { token, expires_at: expiresAt, principal } = payload as {
    token: string;
    expires_at: string;
    principal: unknown;
  };

  const response = NextResponse.json({ principal });
  response.cookies.set(SESSION_COOKIE, token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    expires: new Date(expiresAt),
  });
  return response;
}
