# Deploy no Railway — CampusWorld

Guia para rodar **API + Postgres + Paper** no [Railway](https://railway.app), com **volume persistente** para o mundo Minecraft.

## Visão geral

```text
Railway Project
├── Postgres (plugin gerenciado)
├── API (backend/server/Dockerfile)     → HTTPS público
└── Paper (plugin/Dockerfile.paper)     → TCP proxy + volume /data
         │
         └── plugin CampusWorld → CAMPUS_API_URL
```

| Serviço | Imagem | Persistência |
|---------|--------|--------------|
| Postgres | Plugin Railway | Automática (Railway) |
| API Go | `backend/server/Dockerfile` | Stateless |
| Paper 26.1 | `plugin/Dockerfile.paper` | **Volume em `/data`** |

> **Custo:** Paper 24/7 com 4 GB RAM no Railway costuma mais que uma VPS (Hetzner). Use Railway se priorizar simplicidade; para o mais barato, veja [BUDGET-DEPLOY.md](./BUDGET-DEPLOY.md) (API no Railway + Paper na VPS).

---

## Pré-requisitos

- Conta Railway
- Repositório GitHub conectado (`minecraft-campus` ou submodules)
- Cliente Minecraft **26.1** para jogadores

---

## 1. Postgres

1. Novo projeto Railway → **Add PostgreSQL**
2. Anote a variável `DATABASE_URL` (ou use referência `${{Postgres.DATABASE_URL}}` nos outros serviços)

Não use container Postgres com volume — o plugin gerenciado é mais simples.

---

## 2. API (backend)

### Criar serviço

1. **New Service** → Deploy from GitHub repo
2. **Root directory:** `backend` (se monorepo) ou raiz do repo `minecraft-campus-backend`
3. **Dockerfile path:** `server/Dockerfile`

### Variáveis de ambiente

```env
DATABASE_URL=${{Postgres.DATABASE_URL}}
PLUGIN_API_KEY=<gere-uma-chave-forte>
BOOTSTRAP_MINECRAFT_UUID=<uuid-oficial-mojang>
BOOTSTRAP_USERNAME=SeuNick
CAMPUSWORLD_PROFILE=budget
WORKER_ENABLED=0
HTTP_ADDR=:8080
```

### Rede

1. **Settings → Networking → Generate Domain**
2. Anote a URL pública, ex.: `https://campusworld-api-production.up.railway.app`
3. Teste: `curl https://SUA-URL/health` → `{"status":"ok"}`

### Referência privada (opcional)

Na mesma rede Railway, o Paper pode usar URL interna HTTP se configurado — para o plugin, a **URL pública HTTPS** da API é a opção mais simples.

---

## 3. Paper (Minecraft)

### Criar serviço

1. **New Service** → mesmo repositório
2. **Root directory:** `plugin`
3. **Dockerfile path:** `Dockerfile.paper`

### Volume do mundo

1. Serviço Paper → **Volumes** → **Add volume**
2. **Mount path:** `/data`
3. Tamanho sugerido: **5–20 GB** (mundo cresce com o tempo)

O que fica no volume:

```text
/data/
├── world/              ← overworld
├── world_nether/
├── world_the_end/
├── server.properties   ← seed, motd, etc.
├── plugins/
│   └── CampusWorld/
│       └── config.yml  ← gerado no boot
└── eula.txt
```

O que fica na **imagem** (atualiza a cada deploy):

- Java 25 + Paper 26.1.2
- `CampusWorld.jar`

### Variáveis de ambiente

```env
CAMPUS_API_URL=https://SUA-URL-API.up.railway.app
PLUGIN_API_KEY=<mesma-chave-da-api>
SERVER_SLUG=vanilla
JAVA_MEMORY_MIN=2G
JAVA_MEMORY_MAX=4G
LEVEL_SEED=                    # opcional — só na 1ª geração do mundo
```

`PLUGIN_API_KEY` **deve ser idêntica** à da API.

### TCP Proxy (jogadores)

1. **Settings → Networking → TCP Proxy**
2. **Application port:** `25565`
3. Railway gera algo como: `gondola.proxy.rlwy.net:25862`

Jogadores conectam em:

```text
gondola.proxy.rlwy.net:25862
```

(não é `localhost` nem porta 25565 na internet)

### RAM

1. **Settings → Resources**
2. Mínimo recomendado: **4 GB** para Paper 26.1 + alguns jogadores

---

## 4. Ordem de deploy

1. Postgres
2. API (aguardar healthy)
3. Paper (volume montado em `/data`)

Na primeira subida do Paper:

- Aceita EULA automaticamente
- Gera `world/` no volume
- Se `LEVEL_SEED` estiver definido **e** `world/` não existir, aplica a seed

---

## 5. Seed do mundo

A seed só vale **antes** de existir `world/`:

1. Pare o serviço Paper
2. Defina `LEVEL_SEED=123456789` nas variáveis
3. **Apague** o conteúdo do volume (ou crie volume novo — **apaga tudo**)
4. Suba de novo

Depois que `world/` existe, mudar `LEVEL_SEED` não regera o mapa.

---

## 6. Testar localmente (Docker)

Com a API rodando no host (`docker compose up` em `backend/`):

```bash
cd plugin
docker compose -f docker-compose.paper.yml up -d --build
```

- Mundo persiste no volume Docker `paper_data`
- API local: `CAMPUS_API_URL=http://host.docker.internal:8080`
- Minecraft: `localhost:25565`

---

## 7. Atualizar plugin ou Paper

1. Push no GitHub
2. Railway redeploya o serviço Paper
3. Novo `CampusWorld.jar` é copiado para `/data/plugins/` no boot
4. **Mundo no volume não é apagado**

Para atualizar versão do Paper, altere `PAPER_VERSION` / `PAPER_BUILD` em `Dockerfile.paper` e redeploy.

---

## 8. Troubleshooting

| Problema | Solução |
|----------|---------|
| Kick "CampusWorld indisponível" | `CAMPUS_API_URL` errada; API fora do ar; `PLUGIN_API_KEY` diferente |
| Fundador não entra | `BOOTSTRAP_MINECRAFT_UUID` na API com UUID **oficial** Mojang |
| Mundo sumiu após deploy | Volume não montado em `/data` ou mount path errado |
| Seed não aplicou | `world/` já existia — apague volume e redeploy |
| Não conecta no MC | Usar host:porta do **TCP Proxy**, cliente **26.1** |
| Out of memory | Aumentar `JAVA_MEMORY_MAX` e RAM do serviço Railway |

### Logs

- Railway dashboard → serviço → **Deployments → View logs**
- Paper: procure `CampusWorld ativo. API: ...`

---

## 9. Worker (opcional)

Backups e alertas (Fase 5) — serviço separado:

- **Root:** `backend`
- **Dockerfile:** `server/Dockerfile.worker`
- Variáveis: ver [BUDGET-DEPLOY.md](./BUDGET-DEPLOY.md)

Não é necessário para whitelist/convites funcionarem.

---

## 10. Checklist rápido

- [ ] Postgres criado
- [ ] API com domínio público e `/health` ok
- [ ] `PLUGIN_API_KEY` igual na API e no Paper
- [ ] `BOOTSTRAP_MINECRAFT_UUID` configurado
- [ ] Volume Paper em `/data`
- [ ] TCP Proxy na porta 25565
- [ ] RAM ≥ 4 GB no Paper
- [ ] Cliente Minecraft 26.1

---

## Referências

- [SETUP-PAPER.md](./SETUP-PAPER.md) — dev local sem Docker
- [BUDGET-DEPLOY.md](./BUDGET-DEPLOY.md) — deploy barato
- [Railway TCP Proxy](https://docs.railway.com/networking/tcp-proxy)
- [Railway Volumes](https://docs.railway.com/reference/volumes)
