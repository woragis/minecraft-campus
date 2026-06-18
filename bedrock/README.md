# CampusWorld Bedrock (Dragonfly)

Servidor **Minecraft Bedrock nativo** para o Mundo 1 do CampusWorld. Consulta a API Go para whitelist/convites antes de permitir entrada.

## Pré-requisitos

- Backend API rodando ([`backend/`](../backend/))
- Conta XBOX Live (AuthEnabled no Dragonfly)
- Cliente Bedrock apontando para `localhost:19132`

## Variáveis

| Variável | Default | Descrição |
|----------|---------|-----------|
| `CAMPUS_API_URL` | `http://127.0.0.1:8080` | URL da API CampusWorld |
| `PLUGIN_API_KEY` | — | Mesma chave do Paper (`dev-plugin-key` em dev) |
| `SERVER_SLUG` | `bedrock` | Slug em `game_servers` |
| `SERVER_NAME` | `CampusWorld Bedrock` | MOTD |

Copie `.env.example` para `.env` ou exporte as variáveis.

## Dev local (Go)

```bash
cd bedrock
cp .env.example .env   # ajuste se necessário
export $(grep -v '^#' .env | xargs)
go run .
```

## Docker

```bash
cd bedrock
docker compose up -d --build
```

## Fluxo de join

1. Jogador conecta na porta UDP `19132`
2. Servidor lê XUID + gamertag
3. `GET /v1/internal/whitelist/bedrock/{xuid}?username=...`
4. Se não permitido → kick com mensagem amigável
5. Se permitido → `POST /v1/internal/players/bedrock/upsert` (presença)

## Documentação

- [BEDROCK-DECISION.md](../docs/BEDROCK-DECISION.md)
- [MULTIPLATFORM-ROADMAP.md](../docs/MULTIPLATFORM-ROADMAP.md)
