# Auth Service

Микросервис авторизации и аутентификации. Предоставляет gRPC и REST API для других сервисов.


## Конфигурация

Конфиг читается из YAML файла (по умолчанию `./config.yml`). Путь можно переопределить через переменную `CONFIG_PATH`.

Любой параметр можно переопределить переменной окружения с префиксом `AUTH_`:

| Параметр | Env | По умолчанию | Описание |
|----------|-----|-------------|----------|
| `http.port` | `AUTH_HTTP_PORT` | `8080` | Порт HTTP сервера |
| `db.url` | `AUTH_DB_URL` | `postgres://postgres:postgres@localhost:5432/auth?sslmode=disable` | URL подключения к PostgreSQL |
| `redis.url` | `AUTH_REDIS_URL` | `redis://localhost:6379/0` | URL подключения к Redis |
| `log.level` | `AUTH_LOG_LEVEL` | `info` | Уровень логирования |
| `log.pretty` | `AUTH_LOG_PRETTY` | `true` | Человекочитаемый вывод (иначе JSON) |

Пример `config.yml`:

```yaml
http:
  port: "8080"
db:
  url: "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
redis:
  url: "redis://localhost:6379/0"
log:
  level: "info"
  pretty: true
```

## Запуск

### Локально (нужен PostgreSQL + Redis)

```bash
make run
```

### Через Docker

```bash
# Первый запуск (сборка образа)
make-local-up 1

# Последующие запуски
make-local-up

# Остановка
make-local-down

# Полная очистка (удалить volumes)
make-local-clear
```

## Makefile цели

| Цель | Описание |
|------|----------|
| `make run` | Запуск локально |
| `make build` | Сборка бинарника в `bin/auth` |
| `make test` | Запуск тестов |
| `make lint` | go vet + golangci-lint |
| `make-local-up [1]` | docker compose up (с --build) |
| `make-local-down` | docker compose down |
| `make-local-clear` | docker compose down -v |

## Healthcheck

```bash
curl localhost:8080/health
```

Ответ: `{"status":"ok"}` (200) или `{"status":"unavailable"}` (503).

## База данных

7 таблиц:

- **users** — учётные записи
- **user_identities** — OAuth2 identity (Google, GitHub и т.д.)
- **user_local_passwords** — хеши паролей (bcrypt)
- **user_sessions** — refresh токены
- **api_keys** — ключи для service-to-service
- **roles** — роли + permissions (RBAC)
- **user_roles** — связь пользователей с ролями

