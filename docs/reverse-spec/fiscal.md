# Fiscal / NFC-e

## FIS-001 — Relatório de notas fiscais
Status: `CONFIRMED_AUTH`
Rota: `/main/nf/reports/orders`

### Estrutura observada
- filtro por período;
- modos `Delivery e Balcão` e `Mesas e comandas`;
- estado vazio quando não há dados;
- alerta quando existem produtos sem categoria tributária;
- CTA para configurar tributação;
- orientação específica sobre reforma tributária.

## FIS-002 — Configuração NFC-e
Status: `CONFIRMED_AUTH`
Rota base: `/main/nf/nfc-settings`

Subseções observadas:
1. Configuração.
2. Categorias tributárias.
3. Permissões.

Na configuração existem abas `Dados da empresa` e `Dados fiscais`. A etapa empresarial solicita identificação fiscal e endereço do estabelecimento antes de avançar.

### Acceptance criteria iniciais
1. Dados fiscais devem ser versionados e auditáveis.
2. Produto sem classificação tributária precisa ser detectável antes da emissão.
3. Emissão fiscal deve ser adapter isolado do domínio de pedidos.
4. Erros SEFAZ precisam preservar código, mensagem, tentativa e correlação.
## FIS-003 — Inutilização
Status: `CONFIRMED_AUTH`
Rota: `/main/nf/unusability`

### Estrutura observada
- exibe estado do certificado digital;
- informa ausência de certificado quando aplicável;
- lista inutilizações existentes;
- estado vazio quando não há inutilizações;
- CTA `Solicitar inutilização`;
- explica que a ação comunica à SEFAZ numerações que não serão utilizadas;
- acesso à Central de ajuda.

### Ainda a validar
- upload/renovação do certificado;
- campos da solicitação de inutilização;
- emissão/reemissão/cancelamento de NFC-e;
- contingência e retries;
- autorização por colaborador;
- download/envio/impressão de XML/DANFE.