# Salão / pedidos locais

## HALL-001 — Pedidos salão
Status: `CONFIRMED_AUTH`
Rota: `/main/local-orders/tables`

### Estrutura observada
- Abas `Mesas` e `Comandas`.
- Busca por mesa.
- Filtro por status.
- Legenda de estados: `Livre`, `Ocupada`, `Fechando conta`.
- CTA `Criar mesa`.
- CTA `Novo pedido`.
- Cada mesa possui nome/número, ação `+ Pedido`, menu secundário e estado visual.
- A conta de teste apresentou múltiplas mesas livres.

## HALL-002 — Onboarding App do Garçom
Status: `CONFIRMED_AUTH`

Ao entrar no módulo foi exibido um modal promocional/onboarding com:
- indicação `App do Garçom` e gratuidade;
- QR Code para abrir o app;
- link público `https://garcom.anota.ai`;
- ação de compartilhamento;
- preview de celular e CTA para começar;
- mensagem de que funciona em celulares sem instalação de app.

### Pontos ainda desconhecidos
- transições completas de estado de mesa;
- estrutura de comanda individual/coletiva;
- divisão de conta e pagamentos;
- transferência/junção de mesas;
- permissões do garçom e vínculo com usuário.
