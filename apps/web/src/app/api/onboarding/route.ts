import { NextResponse } from "next/server";
import { API_URL, SESSION_COOKIE } from "@/lib/api";

export async function POST(request: Request) {
  const input = await request.json();
  const created = await fetch(`${API_URL}/api/v1/onboarding`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
    cache: "no-store",
  });

  const createdPayload = await created.json().catch(() => ({}));
  if (!created.ok) {
    return NextResponse.json(createdPayload, { status: created.status });
  }

  const login = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: input?.owner?.email,
      password: input?.owner?.password,
    }),
    cache: "no-store",
  });
  const loginPayload = await login.json().catch(() => ({}));
  if (!login.ok) {
    return NextResponse.json(
      {
        error: {
          message: "Conta criada, mas não foi possível iniciar a sessão.",
        },
        onboarding: createdPayload,
      },
      { status: login.status },
    );
  }

  const { token, expires_at: expiresAt, principal } = loginPayload as {
    token: string;
    expires_at: string;
    principal: unknown;
  };
  const response = NextResponse.json({
    onboarding: createdPayload,
    principal,
  });
  response.cookies.set(SESSION_COOKIE, token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    expires: new Date(expiresAt),
  });
  return response;
}
