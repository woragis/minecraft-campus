# Setup cross-play (Geyser + Floodgate)

Mundo 2 do CampusWorld: **Java e Bedrock no mesmo mapa Paper**, via [Geyser](https://geysermc.org/) + [Floodgate](https://geysermc.org/wiki/floodgate/).

## Quando usar

| Cenário | Entrada |
|---------|---------|
| Só celular, servidor Bedrock nativo | Dragonfly (`bedrock/`) — porta `19132` |
| PC + celular **juntos** no mesmo mundo | Cross-play (este guia) — Java `:25565`, Bedrock `:19132` via Geyser |
| Só Java (dev) | Paper vanilla — [`SETUP-PAPER.md`](./SETUP-PAPER.md) |

## Pré-requisitos

- Backend API rodando
- Docker Desktop (Windows/Mac) ou Docker + `host.docker.internal`
- Conta Microsoft/XBOX no cliente Bedrock

## Subir cross-play

```bash
# Terminal 1 — API
cd backend && docker compose up -d

# Terminal 2 — Paper + Geyser + Floodgate
cd plugin
docker compose -f docker-compose.crossplay.yml up -d --build
```

## Conectar

| Cliente | Endereço | Porta |
|---------|----------|-------|
| Minecraft Java 26.1 | `localhost` | `25565` |
| Minecraft Bedrock | `localhost` | `19132` |

O slug enviado à API é `crossplay` (`SERVER_SLUG`).

## Como funciona

1. **Paper** roda o plugin CampusWorld (whitelist, guildas, claims).
2. **Floodgate** permite contas Bedrock autenticadas sem conta Java premium.
3. **Geyser** traduz protocolo Bedrock → Java no mesmo processo Paper.
4. Bedrock players aparecem no mundo Java; whitelist usa UUID Floodgate (rota Java existente).

## Primeira execução

Na primeira subida, Geyser e Floodgate geram configs em `/data/plugins/`. Se Bedrock não conectar:

1. Confirme UDP `19132` aberto no firewall
2. Verifique logs: `docker compose -f docker-compose.crossplay.yml logs -f`
3. Em `plugins/Geyser-Spigot/config.yml`, `auth-type` deve ser `floodgate` (automático com Floodgate instalado)

## Produção

- URLs separadas na landing: `java.…`, `play.…` (Geyser), `bedrock.…` (Dragonfly)
- Não distribua `key.pem` do Floodgate
- Ver [`RAILWAY-DEPLOY.md`](./RAILWAY-DEPLOY.md) para API; Paper cross-play segue o mesmo padrão de volume `/data`

## Referências

- [BEDROCK-DECISION.md](./BEDROCK-DECISION.md)
- [MULTIPLATFORM-ROADMAP.md](./MULTIPLATFORM-ROADMAP.md)
