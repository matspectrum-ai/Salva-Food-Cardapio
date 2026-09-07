# Architecture — Salva Food

## Baseline

Arquitetura inicial: **modular monolith**.

### Apps
- `apps/api`: API HTTP, domínio, realtime e workers leves em Go.
- `apps/web`: backoffice, PDV, KDS, salão e storefront web.
- `packages/contracts`: OpenAPI/JSON Schema/event contracts compartilhados.

### Runtime dependencies
- PostgreSQL: source of truth.
- Redis: cache, rate-limit, ephemeral presence, short-lived locks.
- Object storage S3-compatible: imagens do cardápio, documentos fiscais, exports.

### Core modules
- identity / RBAC
- tenants / establishments / branches
- catalog / menu / modifiers / combos
- customers / addresses
- orders / order state machine
- storefront / cart / checkout
- pdv
- tables / tabs / waiter
- kitchen / KDS
- delivery / couriers / dispatch
- promotions / coupons / cashback / loyalty
- channels / conversations / AI attendant
- payments
- cash register
- inventory / ingredients / recipes / purchasing
- finance
- fiscal
- reports / analytics
- integrations / webhooks
- audit / observability

## Cross-cutting invariants

- `tenant_id` obrigatório em todas as tabelas de negócio.
- RBAC por tenant/estabelecimento.
- IDs UUID/ULID, nunca sequência exposta como boundary principal.
- Valores monetários em centavos (`int64`).
- Order lifecycle é uma state machine explícita.
- Outbox transacional para webhooks, mensagens e integrações.
- Idempotency keys em criação de pedidos, pagamentos e callbacks externos.
- Audit log append-only para ações administrativas e financeiras.

## Why not microservices now
O alvo possui muitos módulos, mas isso não implica que o Salva Food precise começar distribuído. Um modular monolith reduz custo de coordenação, simplifica transações e acelera a obtenção de paridade. Boundaries serão mantidos para extração futura apenas se throughput, isolamento ou ownership justificarem.

## Web session boundary
- `apps/web` usa Next.js App Router + TypeScript.
- A autenticação de domínio continua pertencendo à API Go.
- O BFF do Next existe apenas para transformar o token opaco em cookie `HttpOnly` e encaminhar chamadas server-side.
- O token de sessão não deve ser persistido em `localStorage`/`sessionStorage`.
- Páginas autenticadas derivam tenant e permissões da sessão; o browser não escolhe `tenant_id` em headers de negócio.
