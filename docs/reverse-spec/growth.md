# Growth / Venda mais

## GROW-001 — Recuperador de vendas
Status: `CONFIRMED_AUTH`
Rota: `/main/trigger`

Modos confirmados:
- automático, com análise de padrões de compra;
- manual, por gatilhos/disparos configuráveis.

No modo manual, a UI permite múltiplos disparos e criação via `Novo gatilho` quando o automático está desativado.

## GROW-002 — Cashback
Status: `CONFIRMED_AUTH`
Rota: `/main/cashback`

Regras observadas:
- preset rápido de 5% com validade de 15 dias;
- cálculo apenas após pedido concluído;
- base = valor do pedido após cupom, excluindo frete;
- resgate usa todo o saldo disponível;
- não combina cashback e cupom no mesmo pedido;
- créditos expiram;
- novo cashback gerado renova o prazo do saldo;
- ao desativar, para de gerar novos créditos, mas saldo existente continua resgatável até expirar.
## GROW-003 — Cupons
Status: `CONFIRMED_AUTH`
Rota: `/main/coupon`

Campos/atributos confirmados na listagem:
- status;
- código do cupom;
- retirados;
- desconto fixo ou percentual;
- alvo do desconto: produto ou frete;
- pedido mínimo;
- validade;
- visibilidade;
- ativação individual;
- edição;
- ativação em massa;
- acesso a desempenho.

A tela também oferece sugestões de cupons inteligentes.

## GROW-004 — Compre + Ganhe +
Status: `CONFIRMED_AUTH`
Rota: `/main/promotion`

Módulo de promoções com criação e gerenciamento. Os campos completos de criação ainda não foram exercitados.

### Acceptance criteria iniciais
1. Promoções devem ser avaliadas por engine de regras determinística.
2. Compatibilidade entre cupom/cashback/promoção deve ser explícita e testada.
3. Cálculo financeiro deve gerar breakdown auditável.
4. Validade deve respeitar timezone do estabelecimento.
5. Recuperação automática precisa de consentimento, rate limiting e trilha de envio.