# Infrastructure — Salva Food

## Decisão

O Salva Food não usa Supabase nem outro BaaS como dependência central.

A base é self-hosted e composta por PostgreSQL, Redis e armazenamento S3-compatible.
A autenticação, autorização, sessão e regras de negócio continuam dentro da API Go.
O Next.js permanece apenas como camada web/BFF.

## Topologia local

- PostgreSQL 17: source of truth e transações.
- Redis 8: cache, rate limiting, presença efêmera, locks curtos e futuro realtime.
- MinIO: armazenamento S3-compatible para imagens, documentos e exports.
- Go API: domínio, HTTP, workers leves e integração com os serviços de infraestrutura.
- Next.js: backoffice, PDV, KDS, salão e storefront.

## Portas locais

A stack usa portas altas para não conflitar com serviços já existentes na máquina:

- PostgreSQL: 127.0.0.1:55432
- Redis: 127.0.0.1:56379
- MinIO API: 127.0.0.1:59000
- MinIO Console: 127.0.0.1:59001
## Persistência e migrações

As migrations versionadas ficam em apps/api/db/migrations.
O projeto usa Goose como runner de desenvolvimento local.
O source of truth do schema continua sendo SQL versionado no repositório.

Comandos principais:

make infra-up
make db-migrate
make db-test
make infra-down

## Realtime

O PostgreSQL continua sendo autoritativo.
Eventos de domínio serão persistidos via transactional outbox.
Redis será usado como camada de distribuição/coordenação para fan-out de atualizações em tempo real,
sem transformar Redis em banco primário.

## Object storage

MinIO implementa a API S3 localmente.
A aplicação deve depender de uma porta S3 abstrata, permitindo trocar MinIO por R2, AWS S3 ou outro
provedor compatível sem acoplar o domínio.

## Produção

A mesma separação pode ser implantada em uma VPS dedicada ou em hosts separados:
PostgreSQL gerenciado/self-hosted, Redis, S3-compatible e uma ou mais instâncias da API/web.
Nenhum desses componentes exige um BaaS específico.
