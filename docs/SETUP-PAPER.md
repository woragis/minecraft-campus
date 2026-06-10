# Setup Paper — CampusWorld Fase 1

Guia para testar whitelist + convites in-game. **Não precisa do launcher Mojang** — só JDK 21 e o jar do Paper.

## Pré-requisitos

- JDK 21
- Docker (para backend + Postgres)
- Download do [Paper 1.21](https://papermc.io/downloads/paper)

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

```bash
mkdir campus-paper && cd campus-paper
# copie o paper-*.jar baixado
echo "eula=true" > eula.txt
java -jar paper-*.jar nogui
```

Na primeira execução o Paper gera `server.properties`. Reinicie após instalar o plugin.

## 4. Teste do fluxo

1. Fundador entra no servidor
2. `/campus status` → deve mostrar API online
3. `/invite NomeDoAmigo` → recebe código `CW-XXXXXX`
4. Amigo tenta entrar → whitelist aceita (probation)
5. Verificar no backend:

```bash
curl http://127.0.0.1:8080/v1/players/minecraft/<uuid-do-amigo>
```

## Troubleshooting

| Problema | Solução |
|----------|---------|
| Kick "precisa de convite" | Rodar `/invite` antes; username deve bater |
| Kick "indisponível" | API fora do ar ou `plugin-key` errado |
| Fundador não entra | Definir `BOOTSTRAP_MINECRAFT_UUID` e reiniciar API |
| Plugin não conecta | Checar `base-url` e firewall |

## Próximo passo

- Velocity (proxy multi-servidor)
- BlueMap
- Frontend web
