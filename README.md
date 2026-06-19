# CampusWorld

Monorepo raiz do projeto **CampusWorld** — plataforma social persistente construída sobre Minecraft.

## Repositórios

Este projeto usa **git submodules**:

| Submodule | Repositório | Descrição |
|-----------|-------------|-----------|
| [`backend/`](backend/) | [minecraft-campus-backend](https://github.com/woragis/minecraft-campus-backend) | API Spring Boot, trust engine, analytics |
| [`frontend/`](frontend/) | [minecraft-campus-frontend](https://github.com/woragis/minecraft-campus-frontend) | Site Next.js — perfis, guildas, dashboard |
| [`plugin/`](plugin/) | [minecraft-campus-plugin](https://github.com/woragis/minecraft-campus-plugin) | Plugin Paper — whitelist, convites, claims |
| [`bedrock/`](bedrock/) | *(neste repo)* | Servidor Dragonfly — Bedrock nativo (Mundo 1) |

## Documentação

- Especificação e roadmap: [`CAMPUSWORLD.md`](CAMPUSWORLD.md)
- Arquitetura: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
- Setup Paper: [`docs/SETUP-PAPER.md`](docs/SETUP-PAPER.md)
- Bedrock (Dragonfly): [`bedrock/README.md`](bedrock/README.md)
- Cross-play (Geyser): [`docs/SETUP-GEYSER.md`](docs/SETUP-GEYSER.md)
- Presença ao vivo: [`docs/PHASE-5-PRESENCE.md`](docs/PHASE-5-PRESENCE.md)
- Stats, HUD e perfis web: [`docs/PHASE-7-STATS-HUD-FRONTEND.md`](docs/PHASE-7-STATS-HUD-FRONTEND.md)
- Decisão multiplataforma: [`docs/BEDROCK-DECISION.md`](docs/BEDROCK-DECISION.md)
- Roadmap multiplataforma: [`docs/MULTIPLATFORM-ROADMAP.md`](docs/MULTIPLATFORM-ROADMAP.md)
- Deploy Railway: [`docs/RAILWAY-DEPLOY.md`](docs/RAILWAY-DEPLOY.md)

## Clone

```bash
git clone --recurse-submodules git@github.com:woragis/minecraft-campus.git
cd minecraft-campus
```

Se já clonou sem submodules:

```bash
git submodule update --init --recursive
```

## Estrutura

```text
faculdade/
├── CAMPUSWORLD.md
├── README.md
├── backend/      # submodule — API Go
├── bedrock/      # Dragonfly Bedrock (Mundo 1)
├── frontend/     # submodule
└── plugin/       # submodule — Paper + Geyser cross-play
```
