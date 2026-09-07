# Gestor de cardápio

## CAT-001 — Navegação do módulo
Status: `CONFIRMED_AUTH`

Subitens observados:
- Gestor
- Imagens do cardápio
- Edição em massa
- Potencializador de cardápio
- Importação inteligente do cardápio
- Importação de Cardápio iFood

## CAT-002 — Gestor
Status: `CONFIRMED_AUTH`
Rota: `/main/menu-v4/manager`

### Estrutura observada
- CTA `Enviar cardápio` com orientação para cadastro via fotos.
- CTA `Nova categoria`.
- Busca de categoria.
- Lista de categorias com drag handle, indicando ordenação manual.
- Categoria selecionada exibe seus itens no painel principal.
- CTA `Adicionar Item` no topo e ao fim da listagem.
- Busca de itens na categoria.
- Toggle de esgotamento no nível da categoria.
- Menu de ações da categoria.
### Item na listagem
Cada item observado possui:
- drag handle para ordenação;
- thumbnail/placeholder de imagem;
- nome;
- toggle `Esgotar`;
- preço em reais;
- ação com ícone de link;
- menu de ações (`...`);
- ação de navegação/edição (`>`).

### Categorias observadas na conta de teste
Foram vistas múltiplas categorias, suficientes para confirmar listagem e seleção de categoria. Os nomes e produtos específicos da conta não são requisitos do Salva Food e não devem virar fixtures de produção.

### Acceptance criteria iniciais
1. Categoria e item devem possuir ordem explícita e reordenável.
2. Esgotamento de categoria e item devem ser operações independentes.
3. Alterações de preço devem usar centavos inteiros no domínio.
4. Busca de categoria e busca de item devem ser escopos distintos.
5. Estado de disponibilidade precisa refletir no storefront sem inconsistência eventual prolongada.

### Pontos ainda desconhecidos
- Campos completos de criação/edição de item.
- Modelagem exata de complementos, variações e combos.
- Regras de preço promocional.
- Comportamento da ação de link.
- Semântica e limites da importação inteligente por fotos.
- Regras de merge/conflito na importação do iFood.
