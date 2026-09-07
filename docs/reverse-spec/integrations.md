# Integrações e canais

## INT-001 — Catálogo de integrações
Status: `CONFIRMED_AUTH`
Rota: `/main/general-configuration/multi-store-integration`

### Estrutura observada
- CTA `Catálogo de integrações`;
- CTA `Integrar`;
- credenciais técnicas por loja exibidas na UI;
- seção `Minhas integrações`;
- abas `iFood` e `Outras`;
- habilitação/desabilitação de integrações existentes;
- paginação da listagem.

Nenhum identificador/token observado deve ser persistido na reverse-spec ou em fixtures.

### Acceptance criteria iniciais
1. Cada integração deve ser configurada por tenant/estabelecimento.
2. Segredos ficam fora do banco de domínio e fora de logs.
3. Integrações externas usam ports/adapters.
4. Webhooks e comandos assíncronos devem ser idempotentes e auditáveis.
5. Estado de conexão e erro deve ser observável sem expor credenciais.

## INT-002 — Redes sociais
Status: `CONFIRMED_AUTH`
Rota base: `/main/general-configuration/social-networks`

Subseções confirmadas: `WhatsApp` e `Facebook/Instagram`.
### WhatsApp
A UI confirma:
- estado conectado/desconectado;
- ação de conectar WhatsApp;
- status independente do robô;
- componente auxiliar `Anota AI Responde`;
- ativação vinculada a WhatsApp Web;
- distribuição por aplicativo/extensão auxiliar;
- credencial/token técnico visível na tela de configuração.

A credencial observada foi deliberadamente ignorada e não integra o projeto.

### Implicações de arquitetura
O Salva Food não deve acoplar o domínio a uma estratégia específica de automação do WhatsApp. O contrato deve aceitar adapters distintos: API oficial, BSP ou cliente auxiliar autorizado, com capability flags e health status por canal.

### Ainda a validar
- fluxo completo de Facebook/Instagram;
- reconexão e expiração de sessão;
- retry/backoff de envio;
- rate limits;
- handoff humano multicanal;
- sincronização de conversas e anexos.