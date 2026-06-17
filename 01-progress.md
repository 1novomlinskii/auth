# Прогресс: Неделя 1 — Go-фундамент

## Сделано

### Структура проекта
- [x] Инициализирован Go module (`github.com/user/auth`)
- [x] Создана структура папок: `cmd/auth/`, `internal/config/`, `internal/app/`, `migrations/`

### Конфиг (Viper + YAML)
- [x] `internal/config/config.go` — структура Config с тегами mapstructure
- [x] Чтение `.yml` файла через Viper
- [x] Env-префикс `AUTH_` для переопределения параметров
- [x] Defaults: порт 8080, postgres/localhost, redis/localhost

### Логгер (zerolog)
- [x] Выбор JSON/pretty format через config
- [x] Запросный middleware через zerolog/hlog (method, path, status, duration)
- [x] Уровень логирования из конфига

### Приложение
- [x] `internal/app/app.go` — App с HTTP сервером, graceful shutdown. Принимает готовый pool (не создаёт)

### Graceful shutdown
- [x] `signal.NotifyContext(SIGTERM, SIGINT)`
- [x] HTTP Shutdown с таймаутом 10s

### Healthcheck
- [x] `GET /health` — пингует БД, возвращает `{"status":"ok"}` (200) или `{"status":"unavailable"}` (503)

### Миграции (Goose, 7 таблиц)
- [x] `001_users.sql` — users (UUID PK, tenant_id, global_handle, email и т.д.)
- [x] `002_user_identities.sql` — OAuth2 identities (Google/GitHub)
- [x] `003_user_local_passwords.sql` — bcrypt password hash
- [x] `004_user_sessions.sql` — refresh tokens с partial index
- [x] `005_api_keys.sql` — service-to-service ключи
- [x] `006_roles.sql` — RBAC роли с permissions массивом
- [x] `007_user_roles.sql` — M2M user ↔ role

### Docker
- [x] `Dockerfile` — multistage (golang:1.23-alpine → alpine:3.19 + curl)
- [x] `.dockerignore`
- [x] `docker-compose.yml` — app + postgres:16 + redis:7 с healthcheck

### Makefile
- [x] `make run` — локальный запуск
- [x] `make build` — сборка бинарника
- [x] `make lint` — go vet + golangci-lint
- [x] `make test` — go test
- [x] `make-local-up [1]` — docker compose up (с build)
- [x] `make-local-down` — docker compose down
- [x] `make-local-clear` — docker compose down -v

### Точка входа
- [x] `cmd/auth/main.go` — config → logger → pool → migrations → app.New → app.Run → pool.Close
- [x] Миграции накатываются в `main` (composition root), а не в `internal/app`

## Исправлено после code review
- [x] Порядок middleware zerolog — NewHandler → AccessHandler → RequestIDHandler
- [x] runMigrations переиспользует пул через `stdlib.OpenDBFromPool` (вместо нового подключения)
- [x] Удалён мёртвый закомментированный код + json-iterator
- [x] Добавлены ReadTimeout/WriteTimeout/IdleTimeout на http.Server
- [x] Добавлен healthcheck для app в docker-compose
- [x] Удалён дубликат индекса в миграции 002_user_identities
- [x] `db.Close()` проверяет ошибку (errcheck)
- [x] Слои: pool + миграции в `main`, HTTP + shutdown в `internal/app`