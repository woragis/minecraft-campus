# CampusWorld — Deploy econômico (profile `budget`)

Guia para rodar na nuvem **sem AWS** e com custo mínimo.

## Profile padrão

```env
CAMPUSWORLD_PROFILE=budget
WORKER_ENABLED=0
BACKUP_ENABLED=0
BACKUP_STORAGE=none
ALERTS_ENABLED=0
```

A API e o plugin funcionam normalmente. Rollback (Fase 4) continua disponível.

## Stack mínima

```bash
docker compose up -d postgres api
```

Não sobe o worker (`docker compose --profile worker` só quando necessário).

## Backup barato (só Postgres, disco local)

```env
WORKER_ENABLED=1
BACKUP_ENABLED=1
BACKUP_DATABASE_ENABLED=1
BACKUP_STORAGE=local
BACKUP_LOCAL_PATH=/var/backups/campusworld
BACKUP_DAILY_ENABLED=1
BACKUP_DATABASE_RETENTION_DAYS=7
```

```bash
docker compose --profile worker up -d worker
```

`pg_dump` roda dentro do container worker. Arquivos ficam no volume `backup_data` — **sem S3**.

## O que fica desligado em budget

| Feature | Var | Motivo |
|---------|-----|--------|
| Backup de mapa | `BACKUP_WORLD_ENABLED=0` | Arquivos `.mca` enormes |
| Cloud storage | `BACKUP_STORAGE=none` | Evita AWS/B2 |
| Alertas automáticos | `ALERTS_ENABLED=0` | Menos CPU |
| Refresh de métricas | `METRICS_REFRESH_ENABLED=0` | Queries on-demand via API |

Métricas públicas: `GET /v1/metrics/overview` (sem worker).

## Alternativas baratas a AWS

- **Disco da VPS** — `BACKUP_STORAGE=local`
- **Backblaze B2** — futuro; `BACKUP_STORAGE=b2` (não implementado ainda)
- **rsync para PC** — copiar volume `backup_data` manualmente

## Produção gradual

```env
CAMPUSWORLD_PROFILE=production
WORKER_ENABLED=1
BACKUP_ENABLED=1
BACKUP_DATABASE_ENABLED=1
BACKUP_STORAGE=local
ALERTS_ENABLED=1
```

Ainda sem mapa nem S3 até você ter orçamento.
