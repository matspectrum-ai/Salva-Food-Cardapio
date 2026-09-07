# Reverse Spec — Anota AI -> Salva Food

## Objetivo
Documentar somente comportamento observável do produto-alvo em sessão autenticada e converter cada observação em contratos verificáveis para o Salva Food.

## Evidência
- Data da primeira sessão autenticada: 2026-09-06.
- Painel observado: `https://admin.anota.ai`.
- Evidências visuais são mantidas temporariamente fora do Git quando podem conter dados da conta.
- Nenhum cookie, token, senha ou segredo do alvo deve ser persistido no repositório.

## Estados de confiança
- `CONFIRMED_AUTH`: comportamento ou elemento observado diretamente após autenticação.
- `CONFIRMED_PUBLIC`: confirmado em documentação/material público.
- `INFERRED`: hipótese derivada de UI ou naming, ainda sem execução do fluxo.
- `UNKNOWN`: ainda não inspecionado.

## Regra de implementação
Uma tela semelhante não significa paridade. Cada fluxo deve registrar rota, estados, ações, validações, transições, erros, efeitos colaterais e acceptance criteria antes de ser marcado como `PARITY`.

## Artefatos
- `navigation.md`: taxonomia, rotas e shell global.
- `orders.md`: pedidos em tempo real e agendados.
- `catalog.md`: gestor de cardápio e seus submódulos.
- `hall.md`: mesas, comandas e App do Garçom.
- `pdv.md`: PDV e atalhos operacionais.
- `delivery.md`: entregas e regiões.
- `kds.md`: mapa/telas de cozinha.
- `payments.md`: onboarding e métodos de pagamento online.
- `robot.md`: intents, personalidades e configurações do robô.
- `growth.md`: recuperador, cashback, cupons e promoções.
- `integrations.md`: catálogo de integrações e canais.
- `settings.md`: cardápio digital e estabelecimento.