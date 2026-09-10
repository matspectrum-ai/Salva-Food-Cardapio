# Pedidos

## ORD-001 — Quadro de pedidos em tempo real
Status: `CONFIRMED_AUTH`
Rota: `/main/orders`

### Estrutura observada
- Filtro principal `Todos` e filtros por canal/tipo representados por ícones.
- Busca por cliente ou número do pedido.
- CTA `Novo pedido`.
- Ação de configuração no cabeçalho.
- Quadro com três estados visíveis: `Em análise`, `Em produção`, `Prontos para entrega`.
- Contador por coluna.
- Toggle `Aceitar os pedidos automaticamente`.
- Configuração separada de tempos para `Balcão` e `Delivery`.

### Card de pedido pronto
- Identificação do pedido.
- Hora do pedido.
- Cliente e telefone.
- Modalidade de fulfillment; foi observado `Retirada no local`.
- Valor total.
- Ação `NF` presente, podendo estar desabilitada.
- Ação `Finalizar pedido`.

### Acceptance criteria iniciais do Salva Food
1. O quadro deve refletir a contagem por estado sem reload completo.
2. Busca deve aceitar nome do cliente e identificador do pedido.
3. Pedidos devem avançar entre estados por transições explícitas e auditáveis.
4. A finalização deve ser uma ação distinta da mudança para `Prontos para entrega`.
5. A configuração de aceitação automática deve ser persistente por estabelecimento.
6. Estados vazios devem possuir mensagens próprias por coluna.

## ORD-002 — Pedidos agendados
Status: `CONFIRMED_AUTH`
Rota: `/main/order-schedule`

### Estrutura observada
- Filtro por número do pedido.
- Filtro de data (`dd/mm/aaaa`).
- Busca por cliente.
- Abas `Pendentes` e `Aceitos`.
- Lista de pedidos à esquerda.
- Painel de detalhe à direita.
- Estado vazio do detalhe: `Selecione um pedido ao lado`.

### Pontos ainda desconhecidos
- Regra de criação do agendamento.
- Janela mínima/máxima de antecedência.
- Transição exata `Pendente -> Aceito`.
- Cancelamento, rejeição, edição e conflitos de horário.
- Relação entre pedido agendado e o Kanban em tempo real.

## ORD-003 — Contexto de cliente no pedido
Status: `INFERRED`

Implementação inicial do domínio mantém o pedido como snapshot histórico do cliente e endereço selecionados. O pedido pode carregar `customer_id`, `address_id`, dados copiados do cliente/endereço, modalidade `DELIVERY`/`PICKUP`/`DINE_IN`, `scheduled_at` e observação.

Regras implementadas enquanto hipótese de domínio:
- `customer_id` e `address_id` são resolvidos dentro do tenant ativo.
- O endereço precisa pertencer ao cliente selecionado.
- `DELIVERY` exige endereço.
- Alterações futuras no cadastro não devem alterar o snapshot já gravado no pedido.
- A chave de idempotência inclui todos os campos de entrada, portanto trocar cliente, endereço, modalidade, agendamento ou observação produz conflito em uma chave já utilizada.

Essas regras são de implementação/inferência e não devem ser tratadas como comportamento confirmado do Anota AI até nova evidência autenticada.
