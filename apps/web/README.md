# Salva Food Web

Next.js backoffice for the Salva Food platform.

## Runtime

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS 4
- Go API through a server-side BFF boundary

The browser never stores the opaque API session token in `localStorage` or `sessionStorage`. The Next.js BFF stores it in an `HttpOnly` cookie and server-rendered pages call the Go API with Bearer authentication.

## Local development

```bash
cp .env.example .env.local
pnpm install
pnpm dev
```

The Go API is expected at `SALVA_API_URL`, defaulting to `http://127.0.0.1:8080`.
