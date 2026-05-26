# Mundo Invest - Client Management & Pipefy Integration

API em Go para gerenciar clientes, simular o mapeamento de cards no Pipefy e processar webhooks de atualização de card com idempotência.

## Requisitos

- Go
- Docker + Docker Compose

## Stack

- Go
- Gin
- PostgreSQL via Docker Compose
- `database/sql` com driver `lib/pq`
- `godotenv` para carregar variáveis locais

## Como executar localmente

Suba o banco:

```bash
sudo docker compose up -d
```

Crie um `.env` a partir do exemplo:

```bash
cp .env.example .env
```

Execute a API:

```bash
go run ./cmd/api
```

A API sobe em:

```text
http://localhost:8080
```

O schema e as seeds de status/prioridade são aplicados automaticamente quando a aplicação inicializa.

## Como rodar os testes

```bash
go test ./...
```

Caso seu ambiente bloqueie escrita no cache padrão do Go, use:

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```

## Endpoints

### Health check

```bash
curl http://localhost:8080/health
```

### Criar cliente

```bash
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "Atualização cadastral",
    "valor_patrimonio": 250000
  }'
```

O cliente é salvo com status inicial `Aguardando Análise`. A mutation GraphQL do Pipefy é montada internamente na camada de integração, sem ser exposta na resposta HTTP.

O campo `valor_patrimonio` é tratado como inteiro em reais no payload da API, alinhado ao exemplo do teste técnico (`250000` significa R$ 250.000). Internamente, a aplicação converte esse valor para centavos antes de persistir no banco para integridade dos dados.

### Simular webhook do Pipefy

```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "evt_123",
    "card_id": "card_456",
    "cliente_email": "joao.silva@example.com",
    "timestamp": "2026-05-18T12:00:00Z"
  }'
```

Regras aplicadas:

- `valor_patrimonio >= 200000`: prioridade `prioridade_alta`
- `valor_patrimonio < 200000`: prioridade `prioridade_normal`
- `event_id` duplicado retorna `409 Conflict` e não processa novamente

## Mutations GraphQL do Pipefy

As mutations estão em `internal/infra/pipefy/card.go` e `internal/infra/pipefy/field.go`.

Referências oficiais usadas:

- `createCard`: https://developers.pipefy.com/reference/create-a-card-with-the-required-fields-fulfilled
- `updateCardField`: https://api-docs.pipefy.com/reference/mutations/updateCardField
- `UpdateCardFieldInput`: https://api-docs.pipefy.com/reference/inputObjects/UpdateCardFieldInput

O `createCard` segue o formato documentado com `pipe_id` e `fields_attributes`, contendo `field_id` e `field_value`.

O update usa `updateCardField(input: { card_id, field_id, new_value })`, com aliases para atualizar status e prioridade no mesmo payload GraphQL simulado.
