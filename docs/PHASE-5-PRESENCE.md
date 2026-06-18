# Fase 5 — Presença ao vivo (Redis)

## Objetivo

Mostrar **quem está online** por servidor e por guilda — no frontend e via API pública.

## Arquitetura

```text
Paper / Dragonfly
  join  → upsert + MarkOnline (Redis)
  quit  → MarkOffline
  tick  → Heartbeat (TTL refresh)

Backend API
  Redis ZSET/HASH por servidor e guilda
  GET /v1/presence/overview
  GET /v1/presence/guilds/{id}

Frontend Next.js
  Landing com total online + guildas ativas
```

## Redis (backend)

| Variável | Default | Descrição |
|----------|---------|-----------|
| `REDIS_ENABLED` | `0` | Liga presença ao vivo |
| `REDIS_URL` | `redis://redis:6379/0` | URL Redis |
| `PRESENCE_TTL_SECONDS` | `120` | Expira jogador sem heartbeat |

### Docker local com presença

```bash
cd backend
# .env: REDIS_ENABLED=1
docker compose up -d redis api
```

## Rotas internas (plugin / bedrock)

| Método | Rota | Quando |
|--------|------|--------|
| POST | `/v1/internal/presence/offline` | Player quit |
| POST | `/v1/internal/presence/heartbeat` | A cada ~60s |

O upsert (`/players/upsert` e `/players/bedrock/upsert`) já marca online automaticamente quando Redis está ativo.

## Plugin Paper

`features.presence` em `config.yml`:

```yaml
features:
  presence:
    enabled: true
    heartbeat-interval-ticks: 1200  # 60s
```

## Bedrock (Dragonfly)

Heartbeat em goroutine (60s) + `HandleQuit` para offline.

## Frontend

```bash
cd frontend && npm install && npm run dev
```

## Checklist

- [ ] `REDIS_ENABLED=1` + Redis no compose
- [ ] Jogador entra no Paper → aparece em `/v1/presence/overview`
- [ ] Jogador sai → some em até TTL (imediato com quit hook)
- [ ] Guilda com membro online → `/v1/presence/guilds/{id}` retorna lista
- [ ] Frontend em `:3000` mostra contadores

## Próximo (Fase 6+)

- Playtime agregado no frontend
- Mob kills (stats ingest)
- Action bar in-game com guild tag
- Ver [PHASE-6-VIP.md](./PHASE-6-VIP.md) (adiado)
