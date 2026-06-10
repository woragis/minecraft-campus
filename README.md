# CampusWorld

Monorepo raiz do projeto **CampusWorld** — plataforma social persistente construída sobre Minecraft.

## Repositórios

Este projeto usa **git submodules**:

| Submodule | Repositório | Descrição |
|-----------|-------------|-----------|
| [`backend/`](backend/) | [minecraft-campus-backend](https://github.com/woragis/minecraft-campus-backend) | API Spring Boot, trust engine, analytics |
| [`frontend/`](frontend/) | [minecraft-campus-frontend](https://github.com/woragis/minecraft-campus-frontend) | Site Next.js — perfis, guildas, dashboard |
| [`plugin/`](plugin/) | [minecraft-campus-plugin](https://github.com/woragis/minecraft-campus-plugin) | Plugin Paper — whitelist, convites, claims |

## Documentação

Especificação completa e roadmap: [`CAMPUSWORLD.md`](CAMPUSWORLD.md)

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
├── backend/      # submodule
├── frontend/     # submodule
└── plugin/       # submodule
```
