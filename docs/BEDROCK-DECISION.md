# Decisão: Bedrock + Cross-play (CampusWorld)

## Contexto

A maioria do público universitário joga **Minecraft Bedrock** (celular/console). Parte joga **Java** (PC). Amigos que querem jogar **juntos** (PC + celular no mesmo mapa) precisam de **cross-play**.

CampusWorld já tem backend Go, plugin Paper (Java) e fluxo de convites/whitelist. Esta decisão define como suportar Bedrock sem abandonar Java.

## Arquitetura: 2 mundos, 3 entradas

| Mundo | Runtime | Público | Porta / URL |
|-------|---------|---------|-------------|
| **Mundo 1 — Bedrock nativo** | [Dragonfly](https://github.com/df-mc/dragonfly) (Go) | Maioria mobile | UDP `19132` (`bedrock.campus...`) |
| **Mundo 2 — Cross-play** | Paper 26 + **Geyser** + **Floodgate** | Java + Bedrock no mesmo mapa Java | TCP `25565` for Java`, UDP `19132` via Geyser (`play.campus...`) |

### Regras

1. **Geyser leva Bedrock → Java**, não o inverso. Dragonfly é servidor Bedrock nativo, separado do Paper.
2. Quem quer **jogar com amigos PC + celular no mesmo mundo** entra pelo **Mundo 2** (Geyser).
3. Quem joga **só no celular**, em servidor otimizado para Bedrock, entra pelo **Mundo 1** (Dragonfly).
4. **Conta única CampusWorld** no backend: mesmo `player_id`, identidades por plataforma (`java` UUID, `bedrock` XUID).
5. **Sem monetização** no escopo atual; cosméticos/VIP ficam na Fase 6 ([PHASE-6-VIP.md](./PHASE-6-VIP.md)).

## Por que Dragonfly (Mundo 1)

| Critério | Dragonfly | Alternativa (só Geyser) |
|----------|-----------|-------------------------|
| Stack | Go — alinhado ao backend | Só Java/Paper |
| Bedrock nativo | Sim, protocolo Bedrock | Proxy para Java |
| Controle / extensão | Código Go custom | Plugins Java |
| Público mobile | Experiência Bedrock pura | Traduz protocolo |

Dragonfly foi escolhido porque há tempo de dev disponível, a stack Go é familiar e permite evoluir features Bedrock (action bar, presença, bridge API) sem depender só do ecossistema Paper.

## Identidade de jogadores

```
players (conta CampusWorld)
  └── player_identities
        platform: java | bedrock
        external_id: UUID Java | XUID Bedrock
```

Bedrock não tem UUID Java. Usamos UUID determinístico para compatibilidade com tabelas existentes:

```
minecraft_uuid = SHA1(URL namespace, "bedrock:" + xuid)
```

O XUID continua sendo a chave canônica Bedrock nas rotas `/v1/internal/whitelist/bedrock/{xuid}`.

## API interna (servidores → backend)

| Rota | Cliente |
|------|---------|
| `GET /v1/internal/whitelist/{minecraftUuid}` | Plugin Paper (Java) |
| `GET /v1/internal/whitelist/bedrock/{xuid}` | Servidor Dragonfly |
| `POST /v1/internal/players/upsert` | Plugin Paper |
| `POST /v1/internal/players/bedrock/upsert` | Servidor Dragonfly |

Autenticação: header `X-Plugin-Key` (mesmo `PLUGIN_API_KEY`).

## Game servers (slug)

| Slug | Descrição |
|------|-----------|
| `bedrock` | Mundo 1 — Dragonfly |
| `crossplay` | Mundo 2 — Paper + Geyser |
| `vanilla` | Paper Java-only (dev/local) |

## Fora de escopo (fases futuras)

- Redis presença ao vivo / guildas online no frontend
- Link automático Java ↔ Bedrock (mesma pessoa, duas plataformas)
- Economia, VIP, cosméticos pagos
- Sincronização de mundos entre Dragonfly e Paper

## Referências

- [MULTIPLATFORM-ROADMAP.md](./MULTIPLATFORM-ROADMAP.md) — fases de implementação
- [SETUP-PAPER.md](./SETUP-PAPER.md) — Paper local
- [RAILWAY-DEPLOY.md](./RAILWAY-DEPLOY.md) — deploy produção
