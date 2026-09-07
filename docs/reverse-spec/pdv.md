# PDV / Pedidos balcão

## PDV-001 — Tela principal
Status: `CONFIRMED_AUTH`
Rota: `/main/pdv-refactor/products`

### Modos observados
- `Delivery e Balcão`.
- `Mesas e Comandas`.

### Controles principais
- Filtros.
- Pesquisa de itens.
- Navegação/seleção de item.
- Rascunhos.
- Telefone e nome do cliente.
- Observação do pedido.
- Entrega.
- Pagamentos.
- CPF/CNPJ.
- Ajuste de valor.
- Geração do pedido.

### Atalhos visíveis
- `D`: Delivery e Balcão.
- `M`: Mesas e Comandas.
- `F`: filtros.
- `P`: pesquisar.
- `A`: próximo.
- `CTRL+X`: rascunhos.
- `O`: observação.
- `E`: entrega.
- `R`: pagamentos.
- `T`: CPF/CNPJ.
- `Y`: ajustar valor.
- `ENTER`: gerar pedido.
### Estados observados
Com pedido vazio, ações dependentes do contexto aparecem desabilitadas, incluindo próximo, entrega, pagamentos, ajuste de valor e geração do pedido.

### Acceptance criteria iniciais
1. O PDV deve suportar operação keyboard-first sem depender do mouse.
2. Rascunhos devem preservar cliente, itens, observação e contexto do pedido.
3. A geração deve ser bloqueada enquanto requisitos mínimos não forem satisfeitos.
4. Totais devem ser recalculados deterministicamente após entrega, descontos e ajustes.
5. A troca entre Delivery/Balcão e Mesas/Comandas deve preservar apenas estado compatível.

### Ainda a validar
- fluxo completo de seleção de produto e adicionais;
- validações de telefone/cliente;
- regras de entrega e endereço;
- formas de pagamento e split;
- aplicação de cupom/cashback no PDV;
- cancelamento/edição após gerar pedido.