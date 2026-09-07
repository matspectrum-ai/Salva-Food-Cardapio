# Pagamento online

## PAY-001 — Conta e planos
Status: `CONFIRMED_AUTH`
Rota: `/main/online-payment/account-and-plans`

### Onboarding Pix observado
- cadastro de documento da empresa/responsável;
- opção para informar ausência de conta bancária no CNPJ;
- CPF do responsável legal;
- seleção do tipo de chave Pix;
- chave Pix vinculada ao documento da conta;
- número de celular;
- CTA de cadastro permanece bloqueado sem requisitos válidos.

### Formas de pagamento exibidas
- Pix.
- Cartão de crédito e carteiras digitais.

A tela também referencia condições comerciais e integração com o Portal iFood Pago.

### Acceptance criteria iniciais
1. Conta de recebimento e titularidade documental devem ser validadas antes da ativação.
2. Credenciais de PSP nunca entram no domínio; ficam em adapters/secret store.
3. Pagamentos devem usar idempotency keys e ledger auditável.
4. Estado de ativação deve ser separado por método de pagamento.
5. Dados bancários/documentais sensíveis não devem aparecer em logs.