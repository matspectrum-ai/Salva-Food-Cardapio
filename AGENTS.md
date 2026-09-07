# Salva Food — Repository Instructions

## Goal
Construir o Salva Food com paridade funcional verificável em relação aos comportamentos observáveis do produto-alvo, sem copiar código-fonte, segredos, marcas ou assets proprietários.

## Engineering loop
Observe -> Specify -> Plan -> Implement -> Verify -> Record evidence.

## Rules
1. Não implementar feature sem acceptance criteria observáveis.
2. Diferenciar `CONFIRMED_PUBLIC`, `CONFIRMED_AUTH`, `INFERRED` e `UNKNOWN` na matriz de paridade.
3. Todo endpoint público deve possuir contrato versionado em `packages/contracts`.
4. Todo módulo deve ter testes de domínio antes de integração externa.
5. Fluxos críticos precisam de testes E2E antes de serem marcados como `PARITY`.
6. Nunca acoplar domínio diretamente a WhatsApp/iFood/gateway fiscal/pagamento; usar ports/adapters.
7. Toda integração assíncrona deve ser idempotente e auditável.
8. Multi-tenancy obrigatório em todas as entidades de negócio.
9. Dinheiro sempre em inteiros de centavos; nunca float.
10. Horários persistidos com timezone explícito do estabelecimento.
11. Não introduzir microserviços sem evidência de necessidade operacional.
12. Cada divergência descoberta no black-box deve gerar fixture/teste de regressão.

## Definition of Done for parity
Uma feature só recebe status `PARITY` quando há:
- evidência de comportamento do alvo;
- spec/acceptance criteria;
- implementação;
- testes automatizados;
- comparação black-box do fluxo principal e edge cases relevantes.
