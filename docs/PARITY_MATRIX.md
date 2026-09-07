# Anota AI -> Salva Food Parity Matrix

Legend:
- `CONFIRMED_PUBLIC`: comportamento confirmado em material público.
- `CONFIRMED_AUTH`: validado em sessão autenticada do alvo.
- `INFERRED`: inferência ainda não validada.
- `UNKNOWN`: precisa de inspeção.
- `TODO`: não implementado.
- `PARITY`: implementado + testado contra evidência.

| Domain | Capability | Evidence | Build |
|---|---|---:|---:|
| Onboarding | cadastro inicial do estabelecimento | CONFIRMED_PUBLIC | TODO |
| Onboarding | importar configuração/cardápio do iFood | CONFIRMED_PUBLIC | TODO |
| Onboarding | configuração manual | CONFIRMED_PUBLIC | TODO |
| Orders | Kanban de pedidos | CONFIRMED_PUBLIC | TODO |
| Orders | análise -> produção -> pronto/entrega -> finalizado | CONFIRMED_PUBLIC | TODO |
| Orders | busca/filtros por canal/tipo | CONFIRMED_PUBLIC | TODO |
| Orders | editar pedido existente via PDV | CONFIRMED_PUBLIC | TODO |
| Orders | reimprimir comanda | CONFIRMED_PUBLIC | TODO |
| Orders | cancelar/finalizar pedido | CONFIRMED_PUBLIC | TODO |
| Catalog | cardápio digital | CONFIRMED_PUBLIC | TODO |
| Catalog | categorias/produtos/complementos/combos | CONFIRMED_PUBLIC | TODO |
| Catalog | QR Code | CONFIRMED_PUBLIC | TODO |
| Catalog | etiquetas/restrições alimentares | CONFIRMED_PUBLIC | TODO |
| Catalog | produtos em destaque / sugestão / upsell | CONFIRMED_PUBLIC | TODO |
| Storefront | repetir pedido | CONFIRMED_PUBLIC | TODO |
| Storefront | agendamento de pedido | CONFIRMED_PUBLIC | TODO |
| Delivery | taxas por bairro | CONFIRMED_PUBLIC | TODO |
| Delivery | taxas por raio | CONFIRMED_PUBLIC | TODO |
| Delivery | regiões desenhadas no mapa | CONFIRMED_PUBLIC | TODO |
| Delivery | importação de bairros por planilha | CONFIRMED_PUBLIC | TODO |
| PDV | balcão | CONFIRMED_PUBLIC | TODO |
| PDV | delivery/telefone | CONFIRMED_PUBLIC | TODO |
| PDV | mesas/comandas | CONFIRMED_PUBLIC | TODO |
| Hall | mapa/estado de mesas | CONFIRMED_PUBLIC | TODO |
| Hall | comanda por mesa/pessoa | CONFIRMED_PUBLIC | TODO |
| Hall | app garçom | CONFIRMED_PUBLIC | TODO |
| Hall | QR de mesa / autoatendimento | CONFIRMED_PUBLIC | TODO |
| Hall | leitura de comanda por código/QR | CONFIRMED_PUBLIC | TODO |
| Hall | balança integrada | CONFIRMED_PUBLIC | TODO |
| Kitchen | KDS | CONFIRMED_PUBLIC | TODO |
| Couriers | cadastro/edição/exclusão de entregadores | CONFIRMED_PUBLIC | TODO |
| Couriers | centrais de entregadores | CONFIRMED_PUBLIC | TODO |
| Couriers | veículo e diária | CONFIRMED_PUBLIC | TODO |
| Delivery AI | atribuição de pedido ao entregador | CONFIRMED_PUBLIC | TODO |
| Delivery AI | mapa/localização em tempo real | CONFIRMED_PUBLIC | TODO |
| Delivery AI | produtividade/relatórios | CONFIRMED_PUBLIC | TODO |
| Logistics | iFood logística sob demanda | CONFIRMED_PUBLIC | TODO |
| Channels | atendente WhatsApp | CONFIRMED_PUBLIC | TODO |
| Channels | Facebook Messenger | CONFIRMED_PUBLIC | TODO |
| Channels | Instagram Direct | CONFIRMED_PUBLIC | TODO |
| Channels | entendimento de áudio | CONFIRMED_PUBLIC | TODO |
| Channels | handoff para humano / pausar robô | CONFIRMED_PUBLIC | TODO |
| Channels | montagem automática de pedido | CONFIRMED_PUBLIC | TODO |
| Marketing | recuperador de vendas | CONFIRMED_PUBLIC | TODO |
| Marketing | cupons | CONFIRMED_PUBLIC | TODO |
| Marketing | cashback | CONFIRMED_PUBLIC | TODO |
| Marketing | cashback por dia/item/categoria/PDV | CONFIRMED_PUBLIC | TODO |
| Marketing | programa de fidelidade | CONFIRMED_PUBLIC | TODO |
| Marketing | compre mais e ganhe mais | CONFIRMED_PUBLIC | TODO |
| Marketing | pesquisa de satisfação | CONFIRMED_PUBLIC | TODO |
| Ads | Meta Pixel | CONFIRMED_PUBLIC | TODO |
| Ads | Google Analytics / tracking | CONFIRMED_PUBLIC | TODO |
| Payments | Pix online | CONFIRMED_PUBLIC | TODO |
| Payments | cartão online | CONFIRMED_PUBLIC | TODO |
| Payments | NuPay/Google Pay/Apple Pay | CONFIRMED_PUBLIC | TODO |
| Payments | antifraude | CONFIRMED_PUBLIC | TODO |
| Cash | abertura/fechamento de caixa | CONFIRMED_PUBLIC | TODO |
| Cash | entradas/saídas/retiradas/despesas | CONFIRMED_PUBLIC | TODO |
| Cash | vendas por forma de pagamento | CONFIRMED_PUBLIC | TODO |
| Inventory | posição de estoque | CONFIRMED_PUBLIC | TODO |
| Inventory | insumos | CONFIRMED_PUBLIC | TODO |
| Inventory | produtos de venda | CONFIRMED_PUBLIC | TODO |
| Inventory | fichas técnicas/receitas | CONFIRMED_PUBLIC | TODO |
| Inventory | entrada/saída/ajuste | CONFIRMED_PUBLIC | TODO |
| Inventory | baixa automática por venda | CONFIRMED_PUBLIC | TODO |
| Inventory | alerta de baixo estoque/esgotado | CONFIRMED_PUBLIC | TODO |
| Purchasing | compras/fornecedores | CONFIRMED_PUBLIC | TODO |
| Purchasing | importação/lançamento de NF de compra | CONFIRMED_PUBLIC | TODO |
| Finance | receitas/despesas/controle financeiro | CONFIRMED_PUBLIC | TODO |
| Finance | CMV/margem/insumos | CONFIRMED_PUBLIC | TODO |
| Fiscal | NFC-e | CONFIRMED_PUBLIC | TODO |
| Fiscal | reemissão/cancelamento/erros fiscais | CONFIRMED_PUBLIC | TODO |
| Fiscal | envio/impressão de nota | CONFIRMED_PUBLIC | TODO |
| Reports | geral/entradas/clientes | CONFIRMED_PUBLIC | TODO |
| Reports | pedidos/produtos | CONFIRMED_PUBLIC | TODO |
| Reports | fidelidade | CONFIRMED_PUBLIC | TODO |
| Reports | estoque | CONFIRMED_PUBLIC | TODO |
| Reports | entregadores | CONFIRMED_PUBLIC | TODO |
| Integrations | pedidos iFood | CONFIRMED_PUBLIC | TODO |
| Integrations | múltiplas lojas iFood | CONFIRMED_PUBLIC | TODO |
| Integrations | integrações ERP/parceiros | CONFIRMED_PUBLIC | TODO |
| Printing | configuração/teste de impressora | CONFIRMED_PUBLIC | TODO |
| Printing | personalização de comanda/mensagem | CONFIRMED_PUBLIC | TODO |
| Settings | abrir/fechar loja temporariamente | CONFIRMED_PUBLIC | TODO |
| Settings | horários e dados do estabelecimento | CONFIRMED_PUBLIC | TODO |
| SaaS | permissões/RBAC exatos | UNKNOWN | TODO |
| SaaS | multiunidade/organização interna exata | UNKNOWN | TODO |
| SaaS | billing/assinatura/admin interno | UNKNOWN | TODO |
| Edge cases | regras exatas de cancelamento/estorno | UNKNOWN | TODO |
| Edge cases | retries/timeouts de canais externos | UNKNOWN | TODO |
| Edge cases | regras exatas de estoque/fiscal/financeiro | UNKNOWN | TODO |
