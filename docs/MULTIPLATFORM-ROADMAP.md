# Roadmap multiplataforma — CampusWorld

Plano de implementação Bedrock (Dragonfly) + cross-play (Geyser), dividido em fases com commits separados.

## Visão geral

```text
                    ┌─────────────────┐
                    │  Backend Go API │
                    │  Postgres       │
                    └────────┬────────┘
           ┌─────────────────┼─────────────────┐
           │                 │                 │
    ┌──────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐
    │  Dragonfly  │   │ Paper+Plugin│   │   Frontend  │
    │  Mundo 1    │   │ + Geyser    │   │  (futuro)   │
    │  bedrock    │   │ crossplay   │   │             │
    └─────────────┘   └─────────────┘   └─────────────┘
```

## Fase 0 — Documentação ✅

**Escopo:** decisão arquitetural e roadmap.

- [x] `docs/BEDROCK-DECISION.md`
- [x] `docs/MULTIPLATFORM-ROADMAP.md`

**Commit:** parent repo (`faculdade/`)

---

## Fase 1 — Backend: identidades Bedrock

**Escopo:** schema + API para whitelist/upsert Bedrock.

### Entregas

- Migration `000013_player_identities.sql`
  - Tabela `player_identities`
  - Seed `game_servers`: `bedrock`, `crossplay`
- Model `PlayerIdentity`
- Repositório: lookup por `(platform, external_id)`, upsert identity
- Service: `CheckBedrockWhitelist`, `UpsertBedrockFromServer`, `BedrockMinecraftUUID`
- Rotas:
  - `GET /v1/internal/whitelist/bedrock/{xuid}?username=`
  - `POST /v1/internal/players/bedrock/upsert`
- Testes unitários + integração (se `TEST_DATABASE_URL`)

**Commit:** submodule `backend/`

---

## Fase 2 — Servidor Dragonfly + bridge API

**Escopo:** servidor Bedrock nativo em Go.

### Entregas

- Pasta `bedrock/` no monorepo
  - Dragonfly baseado no [df-mc/template](https://github.com/df-mc/template)
  - Cliente HTTP para CampusWorld API
  - Join hook: whitelist → kick ou upsert presença
  - `config.toml`, `Dockerfile`, `.env.example`, `README.md`

### Variáveis

| Variável | Descrição |
|----------|-----------|
| `CAMPUS_API_URL` | URL da API (ex. `http://localhost:8080`) |
| `PLUGIN_API_KEY` | Chave interna |
| `SERVER_SLUG` | Default `bedrock` |

**Commit:** parent repo (`faculdade/bedrock/`)

---

## Fase 3 — Cross-play: Geyser + Floodgate

**Escopo:** Mundo 2 — Bedrock entra no Paper via Geyser.

### Entregas

- `plugin/docker-compose.crossplay.yml`
- `plugin/docker/geyser-config.yml` (template)
- `docs/SETUP-GEYSER.md`
- `SERVER_SLUG=crossplay` no compose

**Commit:** submodule `plugin/`

---

## Fase 4 — Parent: integração monorepo

**Escopo:** README, links, submodule pointers.

- Atualizar `README.md` com `bedrock/` e docs
- Commit parent com submodule refs

**Commit:** parent repo

---

## Fase 5 — Presença e stats ✅

- [x] Redis: jogadores online por servidor/guilda
- [x] Dragonfly/Paper: heartbeat + eventos join/quit
- [x] Frontend: guildas online, overview de presença

Ver [PHASE-5-PRESENCE.md](./PHASE-5-PRESENCE.md).

**Commits:** `backend/`, `plugin/`, `bedrock/`, `frontend/`, parent docs

---

## Fase 6 — Cosméticos / VIP (adiado)

Ver [PHASE-6-VIP.md](./PHASE-6-VIP.md).

---

## Ordem de deploy local

```bash
# 1. API
cd backend && docker compose up -d

# 2. Bedrock (Dragonfly)
cd bedrock && go run .   # ou docker compose up

# 3. Cross-play (opcional)
cd plugin && docker compose -f docker-compose.crossplay.yml up -d --build
```

## Checklist por fase

| Fase | Teste manual |
|------|--------------|
| 1 | `curl` whitelist bedrock com XUID + username convidado |
| 2 | Cliente Bedrock conecta `:19132`, convidado entra, não-convidado kick |
| 3 | Cliente Bedrock conecta Geyser `:19132`, aparece no Paper |
| 4 | README e clone com submodules OK |
