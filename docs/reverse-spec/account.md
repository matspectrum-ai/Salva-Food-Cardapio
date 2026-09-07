# Minha conta / colaboradores

## ACCOUNT-001 — Navegação da conta
Status: `CONFIRMED_AUTH`

Subitens observados em `Minha conta`:
- Geral
- Informações pessoais
- Formas de pagamento
- Fatura
- Planos
- Colaboradores

## ACCOUNT-002 — Colaboradores
Status: `CONFIRMED_AUTH`
Rota: `/main/account/employees`

Estrutura observada:
- CTA `Novo colaborador`;
- busca por colaborador;
- paginação;
- colunas `Status`, `Funcionário`, `Cargo`, `Telefone`;
- aviso explícito de que desativar um colaborador não encerra sessões já abertas.
## ACCOUNT-003 — Cadastro de colaborador
Status: `CONFIRMED_AUTH`

Campos observados:
- nome do colaborador, obrigatório;
- CPF, obrigatório;
- permissões, obrigatório;
- e-mail, obrigatório;
- senha, obrigatório;
- telefone, obrigatório;
- cargo, obrigatório e apresentado como texto livre;
- imagem opcional por upload/drag-and-drop.

Upload observado:
- PNG, JPG/JPEG, WEBP e HEIC;
- máximo de 20 MB;
- resolução mínima indicada de 200 px.

Ações: `Salvar` e `Cancelar`.
### Catálogo de permissões observado
- Relatórios
- Assistente virtual
- Aumente suas vendas
- Assinatura
- Configurações
- Ponto de venda
- Estabelecimento
- Cupom
- Pesquisa de satisfação
- Itens do menu
- Configurações - Região de atendimento
- Aba - Compre + Ganhe +
- Aba - Cashback
- Aba - Recuperador de vendas
- Aba - Mesas e garçom
- Gestor de cardápio - edição
- Frente de caixa - aberto
- Edição de pedidos prontos
- Emissão de nota fiscal
- Emissão de nota fiscal (NFC-e)
- Financeiro
- Compras
- Controle de estoque
- Dashboard
- Aba - Pagamento Online
## ACCOUNT-004 — Política de senha visível
Status: `CONFIRMED_AUTH`

O cadastro apresenta os seguintes critérios:
- mínimo de 8 caracteres;
- letras maiúsculas e minúsculas;
- pelo menos 1 número;
- símbolo ou caractere especial (exemplos visíveis: `$`, `#`, `@`).

## Implicações para o Salva Food
- `Cargo` não deve ser confundido com autorização: é metadata descritiva.
- autorização deve usar permission set explícito por vínculo do usuário com o tenant.
- desativação e revogação de sessões devem ser operações distintas; o alvo explicitamente não desloga automaticamente ao desativar.
- a semântica exata de enforcement de cada permissão ainda precisa ser exercitada em contas com diferentes colaboradores.