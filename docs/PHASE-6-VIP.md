# Fase 6 — VIP, cosméticos e apoio à comunidade

> Planejamento — **não implementado**. Documento de referência para quando o projeto retomar monetização/cosméticos.

**CampusWorld** — plataforma social sobre Minecraft. VIP aqui é **patrocínio da comunidade**, não pay-to-win.

---

## Princípio central

> **Apoiar o servidor ≠ comprar vantagem competitiva.**

A filosofia do projeto prioriza reputação, não economia abusiva. VIP reforça **status social** e **visibilidade**, sem quebrar trust, convites ou território.

Alinhado com a [EULA da Mojang](https://www.minecraft.net/en-us/usage-guidelines): cosméticos e acesso ao servidor são permitidos; vantagem injusta no survival não.

---

## O que incluir (e o que evitar)

### ✅ Bom para o CampusWorld

| Categoria | Exemplos | Por quê |
|-----------|----------|---------|
| **Cosméticos in-game** | Partículas, aura, trail ao andar, pet cosmético (sem dano), som ao entrar | Percepção alta, custo baixo, zero impacto em PvP/survival |
| **Identidade social** | Tag no chat (`[Apoiador]`, `[CC]`), cor de nome no tab, prefixo custom (moderado) | Combina com trust/guilda — status visível |
| **Perfil web** | Moldura no perfil, badge “Apoiador desde…”, destaque na lista de membros da guilda | Frontend vira vitrine; reforça a plataforma social |
| **Mapa (BlueMap)** | Ícone custom no mapa, pin dourado na cidade, tour “Patrimônio” | Liga com patrimônio histórico sem vender proteção extra |
| **QoL leve** | +1 home, fila prioritária no Velocity, `/nick` cosmético | Aceitável se for **marginal** |
| **Eventos exclusivos** | Cosmético de temporada (Halloween, formatura), não permanente | Receita recorrente sem inflar perks permanentes |

### ❌ Evitar

- Mais área de claim por dinheiro
- Bypass de probation / trust
- Convites extras pagos
- Proteção anti-rollback / anti-audit
- Itens de survival (diamante, kits, fly em survival)
- Qualquer perk que faça jogador free com trust alto se sentir inferior

**Regra de ouro:** se um jogador free com trust alto se sentir “inferior” no jogo, o perk está errado.

---

## Modelo de tiers

**2 tiers + cosméticos avulsos** no início.

### Tier 1 — Apoiador (~R$ 9–15/mês)

- Tag + cor no chat/tab
- 1 cosmético ativo (partícula ou pet)
- Badge no perfil web
- Fila prioritária (se houver fila no Velocity)

### Tier 2 — Patrono (~R$ 25–40/mês)

- Tudo do Apoiador
- Até 3 cosméticos simultâneos
- Moldura de perfil web + destaque no mapa
- +1 `/home`
- Acesso antecipado a eventos / cosméticos de temporada

### Avulso (one-shot)

- Cosmético específico (R$ 5–15)
- “Destaque patrimônio” por 30 dias no mapa (não é proteção permanente)
- Pacote “formatura” / guilda (cosmético compartilhado da guilda no mapa)

### Caminho gratuito (obrigatório)

Cosméticos também desbloqueáveis por **reputação**:

| Conquista | Recompensa |
|-----------|------------|
| Trust ≥ 80 por 30 dias | Tag “Confiável” |
| Sponsor score alto | Tag “Mentor” |
| Fundador de cidade | Ícone no mapa |
| Construção virou patrimônio | Badge permanente “Legado” |

VIP paga por **exclusividade e conveniência**, não por status que só dinheiro compra.

---

## Arquitetura técnica

```text
Frontend / Loja  →  Stripe ou Mercado Pago
                        ↓ webhook
                   Backend Go (entitlements)
                        ↓ heartbeat / sync
                   Plugin Paper (CosmeticEngine)
```

### Novas entidades (backend)

```text
products          # apoiador_mensal, particulas_aurora, etc.
entitlements      # player_id, product_id, expires_at, source
cosmetic_loadout  # quais cosméticos estão equipados
purchases         # auditoria financeira (webhook id, valor, status)
```

### Fluxo

1. Jogador compra na web (ou ganha por conquista/admin).
2. Webhook confirma pagamento → backend cria `entitlement`.
3. UUID já vinculado (pareamento in-game ou login web).
4. Plugin, no join/heartbeat, busca entitlements e aplica cosméticos.
5. Expiração: worker diário revoga; plugin limpa no próximo sync.

### Plugin — módulo `cosmetics`

- `ParticleTrailListener`
- `TabListDecorator`
- `ChatFormatHook` (limite de caracteres no prefix)
- Comando `/cosmetic` para equipar/desequipar
- Toggles em `config.yml` (padrão Fase 5)

### Pagamentos (deploy barato)

| Opção | Quando usar |
|-------|-------------|
| **Mercado Pago** | Público BR, Pix |
| **Stripe** | Cartão internacional |
| **Ko-fi / Buy Me a Coffee** | MVP rápido |
| **Manual (Pix + admin)** | Primeiros apoiadores, zero integração |

**Recomendação:** Pix manual + painel admin na Fase 6a; integração automática na 6b.

---

## Roadmap (Fase 6 dividida)

### Fase 6a — Cosméticos e apoiadores (MVP)

- [ ] Tabelas `products`, `entitlements`, `cosmetic_loadout`
- [ ] Admin: conceder/revogar tier manualmente
- [ ] Plugin: 3–5 cosméticos (partícula, tag, tab color)
- [ ] Perfil web: badge “Apoiador”
- **Critério:** 1 apoiador real usando cosmético in-game

### Fase 6b — Loja web + pagamento

- [ ] Página `/support` ou `/loja`
- [ ] Webhook Mercado Pago ou Stripe
- [ ] Renovação/expiração automática

### Fase 6c — Patrimônio + cosméticos sociais

- [ ] Destaque no BlueMap (pago ou conquistado)
- [ ] Cosméticos de guilda/cidade
- [ ] Integração com trust milestones (gratuito)

### Fase 6d — Economia leve (opcional)

- [ ] Moeda virtual **não comprável com dinheiro real**
- [ ] Ganha jogando/eventos; gasta em cosméticos in-game
- [ ] Dinheiro real → só tiers e cosméticos diretos na web

Relacionado ao item original do roadmap: *Patrimônio histórico* e *economia leve* ficam em **6c** e **6d**.

---

## Receita realista

Servidor universitário, 30–80 jogadores ativos:

| Cenário | Apoiadores | Receita/mês |
|---------|------------|-------------|
| Conservador | 5 × R$ 15 | ~R$ 75 |
| Saudável | 15 × R$ 20 | ~R$ 300 |
| Forte | 30 × R$ 25 | ~R$ 750 |

Objetivo: **pagar VPS + domínio**, não virar produto SaaS.

---

## Recomendações

1. Começar só com cosméticos + 1 tier “Apoiador”.
2. Sempre ter caminho gratuito via trust/guilda/patrimônio.
3. Pagamento manual primeiro; automação depois.
4. Não vender território nem convites.
5. Frontend (loja + perfil com badges) é pré-requisito prático da 6b.

---

## Decisões em aberto

1. **Público:** só universitários ou abrir com VIP para sustentar infra?
2. **Preço:** R$ 10–15 (acessível) ou R$ 25+ (menos gente, mais receita)?
3. **Pagamento:** Pix manual no início ou Mercado Pago desde o dia 1?
4. **Cosméticos iniciais:** partículas + tags ou pet/aura desde o dia 1?
5. **Fila prioritária:** Velocity terá fila ou perk irrelevante?

---

## Referências

- [CAMPUSWORLD.md](../CAMPUSWORLD.md) — visão e roadmap geral
- [BUDGET-DEPLOY.md](./BUDGET-DEPLOY.md) — deploy barato na nuvem
- [SETUP-PAPER.md](./SETUP-PAPER.md) — testar plugin in-game
