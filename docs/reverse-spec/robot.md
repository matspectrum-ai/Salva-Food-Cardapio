# Robô / Atendimento automatizado

## BOT-001 — Personalização de mensagens
Status: `CONFIRMED_AUTH`
Rota: `/main/bot/bot-messages`

### Estrutura observada
- abas `Mensagens` e `Personalidade`;
- filtros por edição, redes sociais, favoritos e status;
- busca de mensagens;
- reset global para padrão;
- reset por mensagem;
- edição individual com preview de conversa;
- placeholders `Nome do cliente`, `Link do cardápio`, `Separar mensagens` e `Saudação`.

### Intents confirmadas
Existem mensagens configuráveis para saudação, cardápio, promoção, pedido, pagamento, tempo de produção, horários, saiu para entrega, pedido concluído/em produção, agradecimento, ajuda, reserva, atendente humano, fallback, retirada, feedback, emprego, observação, despedida, cancelamento, áudio não entendido, risada, nota fiscal, ambiguidade de intenção, indisponibilidade, entrega/retirada, cupom e loja aberta/fechada.

O conjunto observado é maior que esta síntese; a implementação deve tratar mensagens como catálogo configurável de intents e não como campos hardcoded.

## BOT-002 — Personalidades
Status: `CONFIRMED_AUTH`

Personalidades observadas:
- Padrão;
- Descontraído e divertido;
- Direto ao ponto;
- Fantasia;
- Geek;
- Clássico e retrô.
## BOT-003 — Configurações
Status: `CONFIRMED_AUTH`
Rota base: `/main/bot/configuration`

Subseções confirmadas:
1. Central do robô.
2. Interações robô.
3. Notificações robô.
4. Sons das notificações.
5. Imagens cardápio.
6. Imagens promoção.

A central exibe status e tipo do robô e permite ativação. Na conta observada, o tipo exibido era `Anota Nuvem`.

### Interações configuráveis observadas
- perguntar quando não entender;
- avisar estabelecimento fechado;
- enviar link para jogos;
- enviar CSAT;
- lembrar carrinho não finalizado;
- avisar abertura em 30 minutos;
- desativar aviso de fechado na saudação;
- desativar mensagens de conversão na saudação;
- envio automático de nota fiscal via WhatsApp.

### Notificações administrativas
- observação do pedido;
- pedido incorreto;
- chamada de atendente humano;
- solicitação de nota fiscal.
### Sons
Há configuração independente para `Chamados atendente` e `Solicitações de clientes`, com opções Som 1, Som 2 ou sem notificação e ação de teste antes de salvar.

### Imagens
- `Imagem cardápio`: fallback enviado quando o cliente tem problemas de conexão.
- `Imagem promoção`: configurável por dia da semana, com upload e descrição da imagem.

### Acceptance criteria iniciais
1. Mensagens devem ser versionáveis por intent e tenant.
2. Personalidade deve ser uma transformação configurável, não substituir a intenção semântica.
3. Handoff humano deve ser evento de domínio auditável.
4. Notificações devem ter preferência por categoria e canal.
5. Media assets devem ficar em object storage, nunca no banco relacional.
6. Automação deve possuir evals para intent routing, tool calling e regressão de respostas.
7. Alterações de prompt/personalidade precisam de rollback.