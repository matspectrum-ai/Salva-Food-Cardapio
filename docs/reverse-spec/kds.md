# Cozinha / KDS

## KDS-001 — Mapa de telas
Status: `CONFIRMED_AUTH`
Rota: `/main/kitchen/screens/list`

### Estrutura observada
- título `Cozinha (KDS)`;
- breadcrumb para o módulo;
- seção `Mapa de telas KDS`;
- ação `Adicionar tela KDS`;
- ao menos uma tela chamada `Cozinha` na conta de teste.

### Acceptance criteria iniciais
1. O estabelecimento pode possuir múltiplas telas KDS.
2. Cada tela deve possuir identidade/configuração própria.
3. A criação de tela não deve exigir duplicação do catálogo.
4. Pedidos elegíveis devem ser roteáveis para telas configuradas.

### Ainda a validar
- critérios de roteamento por categoria/item;
- estados do pedido dentro do KDS;
- tempos/SLA e alertas;
- impressão simultânea;
- permissões por tela;
- comportamento offline/reconexão.