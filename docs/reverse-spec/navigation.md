# Navegação autenticada

Status: `CONFIRMED_AUTH`

## Shell global observado
- Sidebar fixa à esquerda com busca "Procurando por algo?".
- Banner superior de trial/oferta.
- Área superior direita com estados operacionais e notificações.
- Identificação do estabelecimento na base da sidebar com estado `ABERTO`.
- CTA para download do app desktop.

## Seção "Seu dia a dia"
- Meus pedidos
- Pedidos balcão (PDV)
- Pedidos salão
- Pedidos agendados
- Gestor de cardápio
- Entregas
- Meu Desempenho
- Cozinha - KDS
- Notas fiscais
- Pagamento online
- Robô

## Seção "Meu Salão"
- Gestão de salão
- Configurações Salão

## Seção "Venda mais"
- Recuperador de vendas
- Cashback
- Cupom
- Compre + Ganhe +

## Seção "Análises"
- Relatórios
- Satisfação

## Seção "Configurações"
- Entregadores
- Minha conta
- Configurações

## Seções auxiliares
- Benefícios
- Instruções de ajuda
- Sugestões
- Termos e Políticas

## Rotas confirmadas até agora
- `/main/orders` — Meus pedidos.
- `/main/order-schedule` — Pedidos agendados.
- `/main/menu-v4/manager` — Gestor de cardápio > Gestor.

A presença de um item nesta lista confirma somente que ele existe na navegação autenticada; não confirma seus subfluxos internos.

## Submenus autenticados confirmados

### Gestor de cardápio
Gestor; Imagens do cardápio; Edição em massa; Potencializador; Importação inteligente; Importação iFood.

### Entregas
Cadastro entregadores; Relatório entregadores; Áreas de entrega.

### Notas fiscais
Relatório de NFC-e; Configuração de NFC-e; Inutilização de NFC-e.

### Robô
Chamado atendentes; Feedback de clientes; Personalização; Configurações.

### Configurações Salão
Meu Salão; Meus Garçons; App Garçom; Comandas; Pedidos Balcão (PDV); Taxa de serviço; Cardápio QR Code; Impressoras; Balanças.

### Relatórios
Geral; Caixas; Clientes; Entradas; Pedidos; Funil de conversão; Mesas e comandas; Cupons; Itens; Entregadores; Garçons; Área de Entrega; Satisfação; Cashback.

### Entregadores
Cadastro; Relatórios; Configurações.

### Minha conta
Geral; Informações pessoais; Formas de pagamento; Fatura; Planos; Colaboradores.
### Configurações
Meus Clientes; Meus Pedidos; Impressora; Frente de caixa; Integrações; Cardápio Digital; Redes Sociais; Entregadores; Robô; Estabelecimento; Pedidos agendados; Integração de anúncios.

### Benefícios
Parceiros; Indique e Ganhe.

## Rotas autenticadas adicionais
- `/main/local-orders/tables` — Pedidos salão.
- `/main/pdv-refactor/products` — PDV.
- `/main/deliveryman/regions` — Áreas de entrega/relatório.
- `/main/kitchen/screens/list` — KDS.
- `/main/online-payment/account-and-plans` — Pagamento online.
- `/main/bot/bot-messages` — Personalização do robô.
- `/main/bot/configuration/*` — Configurações do robô.
- `/main/trigger` — Recuperador de vendas.
- `/main/cashback` — Cashback.
- `/main/coupon` — Cupons.
- `/main/promotion` — Compre + Ganhe +.
- `/main/general-configuration/multi-store-integration` — Integrações.
- `/main/general-configuration/digital-menu/*` — Cardápio Digital.
- `/main/general-configuration/social-networks/*` — Redes Sociais.
- `/main/general-configuration/establishment/*` — Estabelecimento.