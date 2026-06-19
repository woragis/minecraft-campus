# Fase 7 — Stats, HUD in-game e perfis web

Plano para playtime/mob kills, action bar + guild tag, e páginas de perfil/guilda no frontend.

## Visão geral

```text
Paper / Dragonfly
  join  → session start + HUD fetch
  quit  → playtime ingest + presence offline
  kill  → mob kill ingest (Paper; Bedrock parcial)

Backend
  server_players.play_time_secs + mob_kills
  GET /v1/players/{id}/stats
  GET /v1/internal/players/{id}/hud

Frontend
  /players/[id]  /guilds/[slug]  home melhorada
```

## Fase 7.0 — Documentação

- [x] Este arquivo

**Commit:** parent `faculdade/`

---

## Fase 7.1 — Backend: stats + HUD

### Entregas

- Migration `000014_player_stats.sql` — `mob_kills` em `server_players`
- `stats` service: ingest + aggregate
- `POST /v1/internal/stats/ingest`
- `GET /v1/players/{id}/stats`
- `GET /v1/internal/players/{id}/hud` (guilda, status, online da guilda)

**Commit:** submodule `backend/`

---

## Fase 7.2 — Plugin Paper: stats + action bar

### Entregas

- Playtime: join timestamp → quit ingest
- `MobKillListener` (EntityDeathEvent)
- `ActionBarTask` + tab prefix com guild tag
- Config `features.stats`, `features.hud`

**Commit:** submodule `plugin/`

---

## Fase 7.3 — Bedrock: stats + HUD tip

### Entregas

- Playtime no quit
- Mob kill via handler de ataque (living non-player)
- `SendTip` periódico com guilda/status/online

**Commit:** parent `bedrock/`

---

## Fase 7.4 — Frontend: perfis e guildas

### Entregas

- `/players/[id]` — perfil, guilda, stats
- `/guilds/[slug]` — membros, online, trust
- Home com links para guildas

**Commit:** submodule `frontend/`

---

## Fase 7.5 — Parent integration

- README + roadmap update
- Submodule pointers

**Commit:** parent `faculdade/`

---

## Modelo de dados

```sql
server_players (
  play_time_secs BIGINT,  -- já existia
  mob_kills BIGINT        -- novo
)
```

Ingest body:

```json
{
  "playerId": "uuid",
  "serverSlug": "vanilla",
  "sessionSeconds": 360,
  "mobKills": 1
}
```

HUD response:

```json
{
  "username": "woragis",
  "status": "active",
  "guildName": "Campus Crew",
  "guildSlug": "campus-crew",
  "guildOnlineCount": 2
}
```
