# Salva Food

Plataforma SaaS para restaurantes com objetivo de paridade funcional verificável com os fluxos observáveis do Anota AI, sob identidade própria **Salva Food**.

## Princípios

- Paridade funcional é tratada como especificação testável, não como promessa visual.
- Nenhum comportamento é considerado replicado sem evidência + acceptance test.
- Arquitetura inicial: modular monolith, evitando microserviços prematuros.
- Multi-tenant por design.
- Backend: Go.
- Frontend: TypeScript/React (Next.js quando o scaffold de UI for iniciado).
- PostgreSQL como source of truth; Redis apenas para estado efêmero/cache/coordenação.
- Eventos de domínio + transactional outbox para integrações assíncronas.

## Status

Fase 0 — Reverse-spec / mapa de paridade: **iniciada em 2026-09-06**.

Veja `docs/PARITY_MATRIX.md`, `docs/ARCHITECTURE.md` e `docs/ROADMAP.md`.
