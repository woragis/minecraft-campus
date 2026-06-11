# Setup Paper — CampusWorld Fase 1

Guia para testar whitelist + convites in-game.

## Pré-requisitos

- **JDK 25** (Paper 26.1+; `./gradlew build` no plugin baixa via toolchain)
- Docker (para backend + Postgres)
- Download do [Paper 26.1](https://papermc.io/downloads/paper) (ex.: `paper-26.1.2-69.jar`)
- Cliente **Minecraft Java 26.1** no launcher

## 1. Backend

```bash
cd backend
cp .env.example .env
```

Edite `.env` e defina o fundador:

```env
PLUGIN_API_KEY=dev-plugin-key
BOOTSTRAP_MINECRAFT_UUID=<seu-uuid-minecraft>
BOOTSTRAP_USERNAME=SeuNick
```

Suba a stack:

```bash
./scripts/dev-up.sh
# ou: docker compose up -d --build
curl http://127.0.0.1:8080/health
```

## 2. Plugin

```bash
cd plugin
./gradlew build
cp build/libs/CampusWorld-0.1.0.jar /caminho/do/paper/plugins/
```

Configure `plugins/CampusWorld/config.yml`:

```yaml
api:
  base-url: "http://127.0.0.1:8080"
  plugin-key: "dev-plugin-key"
server:
  slug: "vanilla"
```

> Se Paper rodar em Docker e API no host, use `http://host.docker.internal:8080`.

## 3. Paper

Pasta de dev local: `plugin/paper-server/` (já tem `start.bat`).

```bash
cd plugin
./gradlew build
cp build/libs/CampusWorld-0.1.0.jar paper-server/plugins/
cp paper-26.1.2-69.jar paper-server/
cd paper-server
echo "eula=true" > eula.txt
# Windows: start.bat
# ou Java 25 explícito:
java -Xms2G -Xmx4G -jar paper-26.1.2-69.jar --nogui
```

**Seed:** edite `server.properties` → `level-seed=` **antes** da primeira geração. Apague `world/` se já existir.

Na primeira execução o Paper gera `server.properties`. Reinicie após instalar o plugin.

## 4. Teste do fluxo

1. Fundador entra no servidor
2. `/campus status` → deve mostrar API online
3. `/invite NomeDoAmigo` → recebe código `CW-XXXXXX`
4. Amigo tenta entrar → whitelist aceita (probation)
5. Verificar no backend:

```bash
curl http://127.0.0.1:8080/v1/lookup/players/minecraft/<uuid-do-amigo>
```

## Troubleshooting

| Problema | Solução |
|----------|---------|
| Kick "precisa de convite" | Rodar `/invite` antes; username deve bater |
| Kick "indisponível" | API fora do ar ou `plugin-key` errado (deve bater com `PLUGIN_API_KEY` no `.env`) |
| Fundador não entra | Definir `BOOTSTRAP_MINECRAFT_UUID` no `backend/.env`, `docker compose up -d --build api` |
| API reiniciando no Docker | Ver `docker compose logs api`; schema é só via SQL migrations (sem GORM AutoMigrate) |
| Plugin não conecta | Checar `base-url` e firewall |
| `requires Java 25` | Usar Temurin 25 ou `paper-server/start.bat` |
| Versão incompatível no cliente | Launcher deve estar em **26.1**, não 1.21.x |

## Próximo passo

- Velocity (proxy multi-servidor)
- BlueMap
- Frontend web
