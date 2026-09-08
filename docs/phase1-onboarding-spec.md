# Phase 1 — Onboarding inicial

Evidence baseline: cadastro inicial do estabelecimento (`CONFIRMED_PUBLIC`).

## Escopo implementado
- criação de tenant/organização;
- criação da primeira unidade;
- criação do primeiro proprietário;
- permission set completo para o proprietário;
- login imediatamente após o cadastro pelo frontend;
- timezone IANA explícito da unidade.

## Invariantes
1. Tenant, unidade, usuário e membership são criados atomicamente.
2. Falha em qualquer etapa faz rollback integral no PostgreSQL.
3. E-mail duplicado não pode deixar tenant ou unidade órfãos.
4. Senha segue a mesma política dos colaboradores e é persistida somente como hash.
5. A resposta pública do onboarding não expõe CPF nem password hash.
6. O owner recebe cargo descritivo `Proprietário` e todas as permissões conhecidas do Salva Food.
7. O frontend troca o token opaco por cookie `HttpOnly`; não usa localStorage/sessionStorage.

## Fora deste slice
- importação de configuração/cardápio do iFood;
- wizard completo de configuração manual;
- billing/assinatura;
- criação de unidades adicionais;
- regras exatas de multiunidade do produto-alvo.
