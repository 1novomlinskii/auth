# Прогресс: Неделя 1 — Go-фундамент

## Сделано

### Неделя 1: Go-фундамент

#### Структура проекта
- [x] Инициализирован Go module (`github.com/user/auth`)
- [x] Создана структура папок: `cmd/auth/`, `internal/config/`, `internal/app/`, `migrations/`

#### Конфиг (Viper + YAML)
- [x] `internal/config/config.go` — структура Config с тегами mapstructure
- [x] Чтение `.yml` файла через Viper
- [x] Env-префикс `AUTH_` для переопределения параметров
- [x] Defaults: порт 8080, postgres/localhost, redis/localhost

#### Логгер (zerolog)
- [x] Выбор JSON/pretty format через config
- [x] Запросный middleware через zerolog/hlog (method, path, status, duration)
- [x] Уровень логирования из конфига

#### Приложение
- [x] `internal/app/app.go` — App с HTTP сервером, graceful shutdown. Принимает готовый pool (не создаёт)

#### Graceful shutdown
- [x] `signal.NotifyContext(SIGTERM, SIGINT)`
- [x] HTTP Shutdown с таймаутом 10s

#### Healthcheck
- [x] `GET /health` — пингует БД, возвращает `{"status":"ok"}` (200) или `{"status":"unavailable"}` (503)

#### Миграции (Goose, 7 таблиц)
- [x] `001_users.sql` — users (UUID PK, tenant_id, login, email и т.д.)
- [x] `002_user_identities.sql` — OAuth2 identities (Google/GitHub)
- [x] `003_user_password.sql` — bcrypt password hash
- [x] `004_user_sessions.sql` — refresh tokens с partial index
- [x] `005_api_keys.sql` — service-to-service ключи
- [x] `006_roles.sql` — RBAC роли с permissions массивом
- [x] `007_user_roles.sql` — M2M user ↔ role

#### Docker
- [x] `Dockerfile` — multistage (golang:1.23-alpine → alpine:3.19 + curl)
- [x] `.dockerignore`
- [x] `docker-compose.yml` — app + postgres:18 + redis:7 с healthcheck

#### Makefile
- [x] `make run` — локальный запуск
- [x] `make build` — сборка бинарника
- [x] `make gen` — генерация proto + qtc
- [x] `make lint` — go vet + golangci-lint
- [x] `make test` — go test
- [x] `make local-up` — docker compose up (с build)
- [x] `make local-down` — docker compose down
- [x] `make local-clear` — docker compose down -v

#### Точка входа
- [x] `cmd/auth/main.go` — config → logger → pool → migrations → app.New → app.Run → pool.Close
- [x] Миграции накатываются в `main` (composition root), а не в `internal/app`

#### Исправлено после code review (Неделя 1)
- [x] Порядок middleware zerolog — NewHandler → AccessHandler → RequestIDHandler
- [x] runMigrations переиспользует пул через `stdlib.OpenDBFromPool` (вместо нового подключения)
- [x] Удалён мёртвый закомментированный код + json-iterator
- [x] Добавлены ReadTimeout/WriteTimeout/IdleTimeout на http.Server
- [x] Добавлен healthcheck для app в docker-compose
- [x] Удалён дубликат индекса в миграции 002_user_identities
- [x] `db.Close()` проверяет ошибку (errcheck)
- [x] Слои: pool + миграции в `main`, HTTP + shutdown в `internal/app`

### Неделя 2-3: gRPC + REST + бизнес-логика + SQL

#### Proto-контракт
- [x] `proto/auth/v1/auth.proto` — 5 gRPC сервисов (Auth, User, Token, Role, APIKey), 30+ сообщений
- [x] `buf.yaml` / `buf.gen.yaml` — конфигурация buf
- [x] `buf generate` → `pkg/authv1/auth.pb.go` + `auth_grpc.pb.go`

#### Entity (сущности)
- [x] `internal/entity/user.go` — User, Session
- [x] `internal/entity/role.go` — Role
- [x] `internal/entity/apikey.go` — APIKey

#### Usecase (бизнес-логика)
- [x] `internal/usecase/auth.go` — AuthUsecase + интерфейсы UserRepository, SessionRepository, PasswordRepository
- [x] `internal/usecase/user.go` — UserUsecase + интерфейс UserRepository
- [x] `internal/usecase/role.go` — RoleUsecase + интерфейс RoleRepository
- [x] `internal/usecase/apikey.go` — APIKeyUsecase + интерфейс APIKeyRepository

#### Repository (SQL)
- [x] `internal/repository/pg/pg.go` — NewRepo, go:generate qtc
- [x] `internal/repository/pg/query/user.qtpl` — INSERT/SELECT/UPDATE/LIST/COUNT
- [x] `internal/repository/pg/query/session.qtpl` — INSERT/SELECT/REVOKE
- [x] `internal/repository/pg/query/password.qtpl` — INSERT/SELECT/UPDATE
- [x] `internal/repository/pg/query/role.qtpl` — CRUD + Assign/Revoke/CheckPermission
- [x] `internal/repository/pg/query/apikey.qtpl` — INSERT/SELECT/LIST/REVOKE
- [x] `internal/repository/pg/user.go` — UserRepo impl
- [x] `internal/repository/pg/session.go` — SessionRepo + PasswordRepo impl
- [x] `internal/repository/pg/role.go` — RoleRepo impl
- [x] `internal/repository/pg/apikey.go` — APIKeyRepo impl

#### Delivery: gRPC
- [x] `internal/delivery/grpc/auth_server.go` — AuthServiceServer
- [x] `internal/delivery/grpc/user_server.go` — UserServiceServer + TokenServiceServer
- [x] `internal/delivery/grpc/role_server.go` — RoleServiceServer
- [x] `internal/delivery/grpc/apikey_server.go` — APIKeyServiceServer
- [x] `internal/delivery/grpc/interceptors.go` — Recovery, Logging, Auth

#### Delivery: REST
- [x] `internal/delivery/http/auth_handler.go` — 6 endpoints (Register, Login, OAuth, Refresh, Logout)
- [x] `internal/delivery/http/user_handler.go` — 3 endpoints (Get, Update, List)
- [x] `internal/delivery/http/role_handler.go` — 9 endpoints (CRUD + Assign/Revoke/Check)
- [x] `internal/delivery/http/apikey_handler.go` — 4 endpoints (Create, List, Validate, Revoke)
- [x] `internal/delivery/http/middleware.go` — CORS + Auth middleware
- [x] `internal/delivery/http/errors.go` — Google-style JSON error format

#### Config + App
- [x] `internal/config/config.go` — добавлены GRPCConfig, JWTConfig
- [x] `config.yml` — секции grpc (port: 50051), jwt (access_ttl, refresh_ttl)
- [x] `internal/app/app.go` — запуск gRPC (50051) + HTTP (8080), graceful shutdown обоих
- [x] `cmd/auth/main.go` — полное wiring: pool → repo → usecase → handlers → app