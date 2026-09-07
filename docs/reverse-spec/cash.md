# Frente de caixa

## CASH-001 — Configurações
Status: `CONFIRMED_AUTH`
Rota: `/main/general-configuration/backoffice-settings`

Permissões/opções observadas:
- abertura de caixa automática: ao chegar pedido com caixa fechado, abre caixa com saldo inicial zero;
- impressão automática do demonstrativo no fechamento;
- impressão automática de comprovante de retirada;
- impressão automática de comprovante de suprimento.

### Implicações de domínio
- caixa possui ciclo aberto/fechado;
- abertura pode ser explícita ou automática;
- retirada e suprimento são tipos distintos de movimentação;
- fechamento produz demonstrativo;
- impressão é efeito colateral configurável, não parte da transação financeira.

### Acceptance criteria iniciais
1. Movimentações de caixa são append-only/auditáveis.
2. Cada pedido financeiro referencia o caixa efetivo no momento da operação.
3. Abertura automática deve ser idempotente sob pedidos concorrentes.
4. Impressão falha não pode invalidar uma movimentação já persistida.
5. Fechamento deve registrar snapshot dos totais por forma de pagamento.