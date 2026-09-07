# Relatórios

## REP-001 — Relatório Geral
Status: `CONFIRMED_AUTH`
Rota: `/main/reports/general`

Métricas observadas:
- faturamento;
- ticket médio;
- total de pedidos;
- clientes ativos;
- pedidos e entregas;
- total das vendas;
- formas de pagamento;
- distribuição semanal por dia.

Há filtro por intervalo de datas.

## REP-002 — Relatório de Caixas
Status: `CONFIRMED_AUTH`
Rota: `/main/reports/cash`

### Estrutura observada
- filtro por intervalo;
- filtro de status da conciliação;
- acesso ao módulo financeiro;
- estado vazio quando não existem caixas fechados;
- texto vincula caixa conciliado, lançamentos financeiros e despesas.
## REP-003 — Relatório de Entradas
Status: `CONFIRMED_AUTH`
Rota: `/main/reports/entries`

Métricas/controles observados:
- período;
- horários com mais pedidos;
- formas de pagamento com drill-down `Ver detalhes`;
- valor bruto;
- taxas de entrega;
- taxas de serviço;
- valor líquido recebido;
- total de pedidos;
- ticket médio, mínimo e máximo;
- export/download.

A explicação da UI define valor líquido como total após taxas, descontos e cashback.

## REP-004 — Relatório de Pedidos
Status: `CONFIRMED_AUTH`
Rota: `/main/reports/list-orders`

Filtros observados: pagamento, entrega, status e origem, além de busca e período. Ações incluem impressão, download e emissão de nota fiscal em lote.
Colunas observadas na listagem:
- número do pedido;
- status;
- cliente;
- origem;
- data;
- forma de entrega;
- pagamento;
- valor;
- nota fiscal.

Estados de pedido disponíveis no filtro incluem `Em Análise`, `Em Produção`, `Pronto`, `Finalizado` e `Cancelado`.

### Acceptance criteria iniciais
1. Relatórios devem ler de projeções/queries, não alterar o domínio operacional.
2. Métricas financeiras precisam ter fórmula documentada e testada.
3. Filtros por período respeitam timezone do estabelecimento.
4. Exportações devem refletir exatamente os filtros ativos.
5. Valores históricos devem usar snapshots financeiros do pedido, não configuração atual.