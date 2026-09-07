# Phase 1 Core — Catalog + Orders

Status: `IMPLEMENTATION_SPEC`
Evidence: `docs/reverse-spec/catalog.md`, `orders.md`, `pdv.md`.

## Scope
Primeiro vertical slice executável do Salva Food:
- categorias;
- itens de cardápio;
- ordenação;
- esgotamento;
- preço em centavos;
- criação de pedido;
- lifecycle explícito do pedido;
- isolamento obrigatório por tenant.

## Catalog acceptance criteria
1. Categoria exige `tenant_id`, ID, nome e ordem não negativa.
2. Item exige tenant, ID, categoria, nome e preço não negativo.
3. Preço é `int64` em centavos.
4. Categoria e item possuem `sort_order` explícito.
5. Categoria e item podem ser esgotados independentemente.
6. Item nunca pode referenciar categoria de outro tenant.
7. Alterações inválidas retornam erro de domínio, não panic.
## Order acceptance criteria
1. Pedido exige tenant, ID e origem.
2. Estado inicial é `ANALYSIS`.
3. Transição normal: `ANALYSIS -> PRODUCTION -> READY -> FINALIZED`.
4. `CANCELLED` é terminal e pode ser alcançado enquanto o pedido não está finalizado/cancelado.
5. `FINALIZED` é terminal.
6. Pular etapas retorna erro de domínio.
7. Repetir a mesma transição terminal não altera o pedido.
8. Total do pedido é soma determinística dos snapshots dos itens.
9. Snapshot do item guarda nome, quantidade e preço unitário no momento da compra.
10. Nenhuma consulta/alteração pode cruzar tenant.

## HTTP slice
Endpoints iniciais internos do Salva Food:
- `POST /api/v1/catalog/categories`;
- `POST /api/v1/catalog/items`;
- `GET /api/v1/catalog/categories`;
- `POST /api/v1/orders`;
- `GET /api/v1/orders`;
- `POST /api/v1/orders/{id}/transitions`.

O header `X-Tenant-ID` é obrigatório nesta fase. Auth/RBAC substituirá este boundary simplificado posteriormente sem alterar o domínio.