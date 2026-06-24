# Fase 8 — Afiliação acadêmica, HUD e visitantes

Afiliação (estudante, servidor, egresso, visitante), catálogo UFPB, HUD in-game com curso/guilda, convites guest e restrições para visitantes.

## Visão geral

```text
Backend
  migration 000015_affiliation + catalog tables
  player.affiliation_type + university/faculty/course slugs
  GET /v1/catalog/{universities,faculties,courses}
  PATCH /v1/me/affiliation
  invites com affiliationType (student | guest)
  HUD + /me com labels do catálogo
  guests: sem guild lead, city, claim; afiliação bloqueada

Paper plugin
  /invite [guest] <nick>
  action bar, tab e chat prefix ([Visit], curso, guilda)
  join title por afiliação

Frontend
  /conta — picker de afiliação + convite visitante
  perfis e guildas com badge Visitante

Bedrock (parent)
  join title + tip com afiliação
```

## Tipos de afiliação

| Tipo    | Catálogo obrigatório | Convite        | Regras extras                          |
|---------|----------------------|----------------|----------------------------------------|
| student | universidade, centro, curso | padrão   | probation configurável               |
| staff   | opcional             | —              | —                                      |
| alumni  | como student         | —              | —                                      |
| guest   | nenhum               | `guest` no invite | probation mais longa; não altera afiliação; não cria guilda/cidade/claim |

## API (resumo)

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/v1/catalog/universities` | público | Lista universidades |
| GET | `/v1/catalog/faculties?universitySlug=` | público | Centros da universidade |
| GET | `/v1/catalog/courses?facultySlug=` | público | Cursos do centro |
| PATCH | `/v1/me/affiliation` | sessão web | Atualiza afiliação (não guest) |
| POST | `/v1/me/invites` | sessão web | `{ targetUsername, affiliationType? }` |
| POST | `/v1/internal/invites` | plugin key | Idem para Paper |
| GET | `/v1/internal/players/{id}/hud` | plugin key | HUD com affiliation + cores/nomes |

Variáveis: `GUEST_PROBATION_DAYS` (default 14), `PROBATION_DAYS` (default 7).

## Commits

- `backend/` — Phase 8 affiliation backend
- `plugin/` — Phase 8 HUD and guest invite
- `frontend/` — Phase 8 affiliation UI
- `faculdade/` — docs + bedrock + submodule bumps
