# Entregas

## DEL-001 — Navegação do módulo
Status: `CONFIRMED_AUTH`

Subitens confirmados:
- Cadastro entregadores.
- Relatório entregadores.
- Áreas de entrega.

## DEL-002 — Relatório de áreas de entrega
Status: `CONFIRMED_AUTH`
Rota: `/main/deliveryman/regions`

### Estrutura observada
- filtro por período;
- breadcrumb `Entregas`;
- estado vazio quando não há vendas no intervalo;
- foco analítico por área/região de entrega.

### Acceptance criteria iniciais
1. Regiões devem ser entidades versionadas por estabelecimento.
2. Pedidos devem registrar snapshot da região/taxa aplicadas no momento da compra.
3. Alterar uma região não pode reescrever pedidos históricos.
4. Relatórios devem agregar pelo snapshot efetivo do pedido.

### Ainda a validar
- CRUD de regiões/áreas;
- taxa por bairro, raio e polígono;
- cadastro de entregador;
- dispatch e tracking;
- relatórios individuais de entregador;
- configuração do módulo Entregadores.