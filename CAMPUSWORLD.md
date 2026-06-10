# CampusWorld

> Plataforma social persistente construída sobre Minecraft.  
> O Minecraft é uma interface — o produto é a comunidade, a reputação e o grafo social.

**Público inicial:** universitários (Estatística, Ciência da Computação, Física e áreas relacionadas).  
**Visão:** arquitetura expansível para qualquer comunidade.

---

## Índice

1. [Visão e Objetivos](#visão-e-objetivos)
2. [Filosofia](#filosofia)
3. [Arquitetura Geral](#arquitetura-geral)
4. [Rede de Servidores](#rede-de-servidores)
5. [Sistemas de Domínio](#sistemas-de-domínio)
6. [Frontend e Mapa](#frontend-e-mapa)
7. [Dados, Grafos e Estatísticas](#dados-grafos-e-estatísticas)
8. [Modelo de Dados](#modelo-de-dados)
9. [Estrutura de Repositórios](#estrutura-de-repositórios)
10. [Roadmap](#roadmap)
11. [Planejamento — Análise e Decisões](#planejamento--análise-e-decisões)

---

## Visão e Objetivos

### Objetivo principal

Criar uma comunidade online duradoura onde:

- jogadores possam construir cidades e projetos
- exista proteção eficiente contra griefers
- a reputação tenha valor
- a confiança seja rastreável
- a história do mundo seja preservada

### Objetivos secundários

- incentivar colaboração e planejamento urbano
- reduzir vandalismo
- criar métricas sociais interessantes
- produzir dados para análises estatísticas
- criar um ambiente visualmente organizado

### Diferencial

Não é "mais um servidor de Minecraft". É a combinação de:

- comunidade
- reputação
- grafos sociais
- análise social
- visualização de dados
- persistência histórica

---

## Filosofia

O servidor não deve ser apenas um survival comum. A ideia é transformar o mundo em uma **sociedade digital**.

| Não é foco | É foco |
|------------|--------|
| PvP | comunidade |
| economia abusiva | reputação |
| grind excessivo | cidades, construções, história |

---

## Arquitetura Geral

```text
Internet
    │
    ▼
Velocity Proxy
    │
 ┌──┼──────────────┐
 │  │              │
 ▼  ▼              ▼

Vanilla       Pixelmon      Modpack A
Paper         Forge         Forge

 │
 ▼

Plugin CampusWorld

 │
 ▼

Backend Go (server/)

 │
 ├─ PostgreSQL
 ├─ Analytics Engine
 ├─ Trust Engine
 ├─ Rollback Engine
 └─ API REST /v1

 │
 ▼

Frontend Next.js

 ├─ Mapa (BlueMap embed / overlay)
 ├─ Perfil
 ├─ Guildas
 ├─ Cidades
 ├─ Convites
 ├─ Dashboard
 └─ Administração
```

### Princípios arquiteturais

1. **Identidade única** — um jogador, uma conta, múltiplos servidores.
2. **Backend como fonte da verdade** — plugin sincroniza estado; não decide regras de negócio sozinho.
3. **Eventos imutáveis** — auditoria e rollback partem de um log append-only.
4. **Separação de leitura e escrita** — API para mutações; views/materialized views para dashboards.
5. **Expansível por servidor** — cada instância Minecraft registra `server_id` em todas as entidades.

---

## Rede de Servidores

O sistema suporta múltiplos servidores com identidade compartilhada.

| Servidor (exemplo) | Stack |
|--------------------|-------|
| Vanilla Survival   | Paper |
| Pixelmon           | Forge |
| Modpack Create     | Forge |
| Modpack Técnico    | Forge |
| Hardcore           | Paper |
| Eventos temporários| Paper / Forge |

O usuário mantém **mesma conta, reputação, guilda e histórico** independente do servidor.

### Tabelas de ligação servidor ↔ jogador

- `server_players` — presença, último login, tempo jogado por servidor
- `server_events` — eventos específicos de instância (manutenção, wipe, temporada)

---

## Sistemas de Domínio

### Sistema de Confiança

#### Trust Score

Confiabilidade do jogador com base em:

- comportamento
- denúncias
- punições
- histórico

#### Sponsor Score

Qualidade dos convites enviados.

```text
A convida B → B convida C → C causa problemas
→ Sponsor Score de B é reduzido
```

#### Guild Trust

Reputação agregada da guilda.

Exemplo: Estatística 95 · Computação 91 · Física 88

---

### Sistema de Convites

Todo convite gera uma relação no grafo social.

```text
João
 └── Maria
       └── Pedro
```

Permite rastrear:

- origem de jogadores
- clusters de risco
- grupos confiáveis
- propagação de griefing

---

### Probation

Jogadores novos entram em período de avaliação.

Restrições possíveis:

- área limitada
- poucas permissões
- sem construções públicas
- sem criação de guilda

Após comportamento adequado → acesso completo.

---

### Guildas

Comunidades organizadas.

Exemplos: Estatística, Computação, Física, Matemática, Engenharia, Externos.

Cada guilda possui:

- líder
- reputação
- membros
- cidades
- alianças

---

### Sistema de Cidades

Incentiva urbanismo e projetos coletivos.

Exemplo — Cidade de Estatística:

- praça
- biblioteca
- mercado
- universidade
- observatório

Cada cidade possui:

- fundador
- população
- guilda
- data de criação
- área

---

### Sistema de Claims

Claims protegem terrenos.

| Campo | Descrição |
|-------|-----------|
| dono | jogador |
| coordenadas | região no mundo |
| tamanho | área em blocos² |
| cidade | opcional |
| guilda | opcional |

Crescimento de área pode depender de:

- tempo jogado
- confiança (trust score)
- atividade

#### Progressão territorial (exemplo)

| Tempo de conta | Área máxima |
|----------------|-------------|
| Nova           | 1 000 blocos² |
| 1 semana       | 1 500 blocos² |
| 1 mês          | 3 000 blocos² |
| 6 meses        | 10 000 blocos² |

---

### Regras de Construção

Objetivos: evitar caos visual, preservar estética, incentivar cidades.

Possíveis regras:

- construções gigantes exigem aprovação
- farms industriais em áreas específicas
- castelos e monumentos com categoria especial
- áreas históricas protegidas

#### Zonas

| Zona | Uso |
|------|-----|
| Urbana | casas, lojas, apartamentos |
| Rural | fazendas, celeiros |
| Industrial | redstone, farms |
| Histórica | castelos, monumentos, museus |

---

### Alianças

Guildas podem criar alianças.

Exemplo: Computação + Física.

Benefícios possíveis:

- cidades compartilhadas
- projetos conjuntos
- eventos

---

### Patrimônio Histórico

Construções especiais recebem:

- proteção extra
- destaque no mapa
- histórico permanente

Exemplos: castelos, monumentos, bibliotecas, museus.

---

### Sistema de Auditoria

Registrar eventos importantes (append-only):

- Block Place / Block Break
- Chest Open
- Item Transfer
- Player Join / Quit

Objetivos: investigações, rollback, estatísticas.

---

### Sistema de Rollback

Desfazer ações específicas sem restaurar o mundo inteiro.

Exemplo:

```text
Pedro destruiu 3 000 blocos
→ rollback(Pedro, janela_temporal)
→ sistema reconstrói apenas os blocos afetados
```

---

### Sistema de Backups

| Frequência | Retenção |
|------------|----------|
| Diário     | 7 dias |
| Semanal    | 8 semanas |
| Mensal     | 12 meses |

Backups enviados para armazenamento externo (S3, Backblaze, etc.).

---

## Frontend e Mapa

### Mapa online — BlueMap

- casas públicas
- claims
- cidades
- monumentos
- guildas (camadas / filtros)

### Páginas

| Rota | Conteúdo |
|------|----------|
| `/` | jogadores online, servidores, eventos |
| `/players/[id]` | reputação, guilda, cidade, convites, histórico |
| `/guilds/[id]` | membros, reputação, cidades |
| `/cities/[id]` | mapa, população, fundador |
| `/dashboard` | estatísticas, gráficos, crescimento |
| `/admin` | moderação, rollback, whitelist |

---

## Dados, Grafos e Estatísticas

### Métricas de produto

- jogadores ativos (DAU / WAU / MAU)
- retenção por coorte de convite
- crescimento de cidades
- distribuição territorial
- histórico de guildas

### Grafo social

Armazenar arestas:

- convites (`invited_by`)
- guildas (`member_of`)
- alianças (`allied_with`)
- comércio (fase futura)
- interações (chat, co-presença — fase futura)

Cálculos possíveis:

- centralidade
- PageRank
- detecção de comunidades (Louvain / Leiden)
- influência e propagação de risco

---

## Modelo de Dados

### Tabelas principais

```text
players
guilds
cities
claims
invites
trust_scores
sponsor_scores
audit_events
rollbacks
world_snapshots
alliances
server_events
server_players
```

### Entidades sugeridas (campos mínimos)

<details>
<summary><code>players</code></summary>

- `id`, `minecraft_uuid`, `username`, `email` (opcional, web)
- `trust_score`, `sponsor_score`
- `probation_until`, `status` (active, banned, probation)
- `invited_by_player_id`, `created_at`

</details>

<details>
<summary><code>guilds</code></summary>

- `id`, `slug`, `name`, `leader_player_id`
- `trust_score`, `created_at`

</details>

<details>
<summary><code>cities</code></summary>

- `id`, `name`, `founder_player_id`, `guild_id`
- `server_id`, `world`, `center_x`, `center_z`, `area_blocks`
- `created_at`

</details>

<details>
<summary><code>claims</code></summary>

- `id`, `owner_player_id`, `city_id`, `guild_id`
- `server_id`, `world`, `min_x`, `min_z`, `max_x`, `max_z`
- `zone_type`, `created_at`

</details>

<details>
<summary><code>invites</code></summary>

- `id`, `sponsor_player_id`, `invited_player_id`
- `code`, `status`, `created_at`, `accepted_at`

</details>

<details>
<summary><code>audit_events</code></summary>

- `id`, `server_id`, `world`, `player_id`
- `event_type`, `block_x`, `block_y`, `block_z`
- `block_type`, `metadata` (JSONB), `occurred_at`

</details>

---

## Estrutura de Repositórios

Monorepo recomendado em `minecraft/faculdade/`:

```text
faculdade/
├── CAMPUSWORLD.md          # este arquivo
├── docs/
│   ├── architecture.md
│   ├── api-contract.md
│   └── trust-algorithm.md
├── plugin/                 # Paper / Spigot
│   ├── src/main/java/...
│   │   ├── MainPlugin.java
│   │   ├── commands/
│   │   ├── listeners/
│   │   ├── services/
│   │   ├── managers/
│   │   ├── models/
│   │   └── api/
│   └── resources/
│       ├── plugin.yml
│       └── config.yml
├── backend/                # Go (padrão Lingo)
│   ├── migrations/
│   ├── server/
│   │   ├── cmd/server/main.go
│   │   └── internal/
│   │       ├── httpserver/
│   │       ├── player/
│   │       ├── invite/
│   │       ├── trust/
│   │       └── ...
│   └── docker-compose.yml  # Postgres local
├── frontend/               # Next.js
│   └── app/
│       ├── map/
│       ├── players/
│       ├── guilds/
│       ├── cities/
│       ├── dashboard/
│       └── admin/
└── infra/
    ├── velocity/
    ├── paper/
    └── bluemap/
```

### Stack proposta (MVP)

| Camada | Tecnologia | Motivo |
|--------|------------|--------|
| Proxy | Velocity | moderno, forwarding, multi-backend |
| Servidor | Paper 1.21+ | API estável, performance |
| Plugin | Java 21 + Paper API | ecossistema maduro |
| Backend | Go 1.22 + GORM | rápido, leve em RAM, padrão handler/service/repository |
| Banco | PostgreSQL 16 | JSONB, grafos via extensões, analytics |
| Fila (fase 4+) | Redis ou RabbitMQ | auditoria assíncrona |
| Frontend | Next.js 15 + Tailwind | SSR, dashboard, mapa embed |
| Auth web | Supabase Auth ou Keycloak | login web separado do Minecraft |
| Mapa | BlueMap | open source, integrável |
| Deploy | VPS (Hetzner/OCI) + Docker | custo baixo para universitários |

---

## Roadmap

### Fase 1 — Fundação (MVP jogável)

- [ ] Velocity + 1 servidor Paper funcional
- [ ] Plugin: heartbeat, whitelist, comando `/invite`
- [ ] Backend: CRUD players, invites, auth plugin ↔ API
- [ ] Frontend: landing, perfil básico, lista de convites
- [ ] BlueMap básico online
- [ ] PostgreSQL + migrations

**Critério de sucesso:** 10–20 jogadores convidados jogando com whitelist e mapa público.

---

### Fase 2 — Confiança e comunidade

- [ ] Trust Score v1 (regras simples, transparentes)
- [ ] Sponsor Score v1
- [ ] Probation (restrições no plugin)
- [ ] Guildas (criação, membros, líder)
- [ ] Páginas web de guilda e perfil enriquecido

**Critério de sucesso:** moderação consegue identificar origem de jogador problemático via grafo de convites.

---

### Fase 3 — Território

- [ ] Claims (plugin + API)
- [ ] Cidades (vinculação claim ↔ cidade ↔ guilda)
- [ ] Progressão territorial
- [ ] Zonas e regras de construção (mínimo viável)
- [ ] Alianças entre guildas
- [ ] Overlay de claims/cidades no mapa

---

### Fase 4 — Auditoria e rollback

- [ ] Pipeline de audit events (async)
- [ ] Rollback por jogador + janela temporal
- [ ] Painel admin de investigação
- [ ] Retenção e particionamento de `audit_events`

---

### Fase 5 — Operações e visibilidade

- [ ] Backups automatizados (diário / semanal / mensal)
- [ ] Dashboard de métricas (retenção, cidades, território)
- [ ] Alertas (picos de griefing, clusters de risco)

---

### Fase 6 — Patrimônio e economia (opcional)

- [ ] Patrimônio histórico (proteção + destaque)
- [ ] Economia leve (se fizer sentido para a comunidade)

---

### Fase 7 — Analytics avançado

- [ ] Motor de grafos (PageRank, comunidades)
- [ ] Visualizações interativas (grafo de convites, influência)
- [ ] Export para análise estatística (CSV / API)

---

## Planejamento — Análise e Decisões

### O que está muito bom na especificação

1. **Visão clara de produto** — não é "servidor de Minecraft", é plataforma social.
2. **Grafo de convites** — mecanismo orgânico de growth + anti-griefing; raro em servidores MC.
3. **Separação plugin / backend / frontend** — escala bem para multi-servidor.
4. **Roadmap em fases** — permite lançar cedo e validar com a faculdade.
5. **Alinhamento com Estatística/CC** — dados reais para projetos acadêmicos.

### Riscos e como mitigar

| Risco | Mitigação |
|-------|-----------|
| Escopo grande demais | Fase 1 enxuta: 1 servidor, convites, whitelist, mapa |
| Auditoria gera volume absurdo | Amostragem + batch + só blocos em claims / probation |
| Rollback complexo | Começar com rollback de blocos em claims protegidos apenas |
| Multi-servidor (Forge + Paper) | Adiar modpacks; só Paper na Fase 1 |
| Reputação injusta | Publicar fórmula v1; painel de apelação; score explicável |
| Poucos jogadores no início | Convites fechados; guildas por curso; eventos semanais |

### Decisões em aberto (resolver antes de codar)

#### 1. Autenticação

| Opção | Prós | Contras |
|-------|------|---------|
| **A) Só Minecraft (UUID + whitelist)** | Simples | Web desconectada do jogo |
| **B) Microsoft/Mojang OAuth na web** | Conta única | Mais complexo |
| **C) Email (Supabase) + vincular UUID no primeiro login** | Bom para dashboard | Dois passos |

**Recomendação Fase 1:** opção A no plugin + opção C na web (vincular UUID manualmente ou via código de pareamento in-game).

#### 2. Onde roda o quê

- **VPS única (4–8 GB RAM)** para Fase 1: Velocity + Paper + Postgres + backend + BlueMap.
- Separar banco e backend em VPS diferente quando passar de ~30 jogadores simultâneos.

#### 3. Claims — plugin próprio ou existente?

| Opção | Nota |
|-------|------|
| GriefPrevention / Lands | Rápido, menos controle |
| Plugin custom integrado ao backend | Mais trabalho, modelo de dados unificado |

**Recomendação:** custom desde o início se claims dependem de trust score e cidades — senão a integração vira dívida técnica.

#### 4. Fórmula Trust Score v1 (proposta inicial)

```text
trust = 100
  - (denúncias_confirmadas × 15)
  - (rollbacks_aplicados × 5)
  + (dias_jogados_em_probation_sem_incidente × 0.5)
  clamp(0, 100)
```

Sponsor Score: média ponderada do trust dos convidados diretos, com penalidade se convidado for banido em < 30 dias.

> Ajustar com dados reais após Fase 2.

#### 5. Nome e branding

- **CampusWorld** — bom, genérico o suficiente para expandir além da faculdade.
- Domínio, Discord, e identidade visual: definir antes do lançamento Fase 1.

### MVP mínimo absoluto (se quiser ir ainda mais enxuto)

```text
Semana 1–2: Paper + whitelist manual + Discord
Semana 3–4: Plugin convites + API + Postgres
Semana 5–6: Site estático com mapa BlueMap + perfis
```

Só depois disso: trust, guildas, claims.

### Próximos passos concretos

1. **Validar com o grupo** — quantos jogadores reais na Fase 1? Qual curso pilota?
2. **Escolher auth** — pareamento UUID ↔ conta web.
3. **Subir repositório** — monorepo com `plugin`, `backend`, `frontend`.
4. **Contrato API v0** — endpoints mínimos: `POST /invites`, `GET /players/{uuid}`, `PATCH /players/{uuid}/probation`.
5. **Infra Fase 1** — Docker Compose local; depois VPS com script de deploy.
6. **Regras do servidor** — documento curto de conduta (complementar este arquivo).

---

## Glossário

| Termo | Significado |
|-------|-------------|
| Trust Score | Confiabilidade individual do jogador |
| Sponsor Score | Qualidade dos convites que o jogador emitiu |
| Guild Trust | Reputação agregada da guilda |
| Probation | Período de avaliação para jogadores novos |
| Claim | Região protegida de construção |
| Patrimônio | Construção com proteção histórica permanente |

---

*Última atualização: junho/2026*
