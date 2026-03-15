# FIAP X - Sistema de Processamento de Vídeos

Este repositório contém o código-fonte do backend da plataforma FIAP X, baseada na arquitetura sugerida (Event-Driven, Microserviços e Clean Architecture) com foco em alta escalabilidade e resiliência (Zero perda de requisições).

## Tecnologias e Atendimento aos Requisitos
- **Linguagem:** Go 1.21 (Garante concorrência severa para múltiplos vídeos simultâneos).
- **Banco de Dados:** PostgreSQL (Persistência garantida, Particionamento, JSONB para Metadados, Multi-tenant RLS).
- **Mensageria e Filas:** RabbitMQ (Atua como "buffer", evitando perda de requisições em horários de pico e lidando com DLQ/Retries).
- **Cache e Sessão:** Redis (Desempenho no controle de requisições).
- **Storage:** MinIO / S3 (Armazenamento dos frames e zips).
- **Monitoramento:** Prometheus e Grafana.
- **Autenticação:** JWT (Sistema protegido por usuário e senha).
- **CI/CD:** Github Actions (Lint, Testes automatizados e Build - `.github/workflows/ci.yml`).
- **Infraestrutura:** Docker e Docker Compose nativos.

## Funcionalidades Atendidas (MVP)
✅ **Múltiplos Vídeos ao mesmo tempo:** A arquitetura baseada em Workers escutando fila (`RabbitMQ`) atende N videos de forma assíncrona.
✅ **Resiliência em Picos:** A API Gateway apenas enfileira o Job; os picos não derrubam a aplicação.
✅ **Proteção por Usuário/Senha:** Handlers de Auth preparados com suporte futuro a JWT e BCrypt.
✅ **Listagem de Status dos Vídeos:** Rota configurada para o usuário acompanhar `PENDING -> PROCESSING -> COMPLETED`.
✅ **Notificações:** O Domínio contempla notificações `TipoErro` e `TipoSucesso` pós-processamento.
✅ **Testes de Qualidade:** Ampla bateria de Mocks e Unit Tests.

## Como rodar o ambiente de desenvolvimento
O projeto conta com um `docker-compose.yml` que sobe os serviços fundamentais.

1. Suba a infraestrutura, a API e o Monitoramento:
   ```bash
   make docker-up
   # ou
   docker-compose up --build -d
   ```

2. Execute o setup dos módulos Go (Se rodar via host nativo):
   ```bash
   go mod tidy
   ```

3. Rodar a API Principal localmente:
   ```bash
   make run
   ```

4. Rodar o Worker (em um terminal separado):
   ```bash
   make worker
   ```

5. Rodar os testes automatizados com cobertura:
   ```bash
   make coverage
   ```

## Portas de Serviços em Execução (Docker Compose)
- **API principal:** `:8080`
- **PostgreSQL:** `:5432`
- **RabbitMQ Dashboard:** `:15672` (Usuário: `guest`, Senha: `guest`)
- **MinIO S3 UI:** `:9001` (Usuário: `minioadmin`, Senha: `minioadmin`)
- **Prometheus:** `:9090`
- **Grafana:** `:3000` (Usuário: `admin`, Senha: `admin`)
