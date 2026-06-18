# План: Единый сервис авторизации (Auth Service)

**Период:** 8 недель (40 дней)
**Цель:** Создать standalone gRPC + REST сервис авторизации с multi-tenant, Google OAuth2, RBAC, API Keys и SDK для подключения другими сервисами (Tinder, URL Shortener, Telegram Chat).

---

## Архитектура

```
                        ┌──────────────────────────┐
                        │      External Client     │
                        │  (REST via grpc-gateway) │
                        └───────────┬──────────────┘
                                    │ REST (JSON)
                        ┌───────────▼──────────────┐
                        │    Auth Service (port     │
                        │    8080 REST, 50051 gRPC)│
                        │                         │
                        │  ┌─── Use Cases ───────┐ │
                        │  │ Register / Login     │ │
                        │  │ Google OAuth2        │ │
                        │  │ Token Management     │ │
                        │  │ RBAC (Roles/Perms)   │ │
                        │  │ API Keys             │ │
                        │  │ Multi-tenant         │ │
                        │  └──────────────────────┘ │
                        └────────┬──────────┬───────┘
                                 │          │
                           ┌────▼──┐  ┌────▼────┐
                           │  PG   │  │ Redis   │
                           └───────┘  └─────────┘

▲ gRPC (shared auth-client-go SDK)
│
┌──────────────┴──────────────────────────────┐
│        Future: Tinder / URL Shortener       │
│        Подключают auth-client-go SDK        │
└─────────────────────────────────────────────┘
```

**Микросервисов:** 1 (Auth Service)
**БД:** PostgreSQL 16 (users, identities, sessions, roles, api_keys) + Redis 7 (refresh tokens, rate limiter, blacklist)
**API:** gRPC (internal) + REST (external через grpc-gateway)
**Конфигурация:** .env + caarlos0/env
**Протокол auth:** JWT RS256 (access 15min + refresh 7d)
**SDK:** auth-client-go (Go-пакет для подключения auth в другие сервисы)

---

## 📅 Неделя 1: Go-фундамент + Clean Architecture

**Цель:** Скелет проекта, который компилируется, поднимается через docker-compose, имеет конфиг и логгер.

| День | Задача | Время |
|------|--------|-------|
| 1 | Инициализация Go module, структура проекта: `cmd/auth/`, `internal/domain/`, `internal/usecase/`, `internal/repository/`, `internal/handler/`, `internal/config/`, `pkg/` | 1.5ч |
| 2 | Конфиг из `.yml` через `Viper` + `zerolog` логгер (JSON, уровни, source) | 1.5ч |
| 3 | Graceful shutdown (`SIGTERM`, `SIGINT`) + base app runner | 1.5ч |
| 4 | Dockerfile (multistage build) + `.dockerignore` + `.gitignore` | 1.5ч |
| 5 | `docker-compose.yml`: app + postgres:16 + redis:7 | 1ч |
| 6 | Makefile: `lint`, `test`, `build`, `run`, `docker-up`, `docker-down`, `migrate-up`, `migrate-down` | 1.5ч |

**Изучить:**
- Go project layout: https://github.com/golang-standards/project-layout
- Clean Architecture в Go: https://github.com/bxcodec/go-clean-arch
- caarlos0/env: https://github.com/caarlos0/env
- slog: https://pkg.go.dev/log/slog
- logging guide: https://betterstack.com/community/guides/logging/logging-in-go/
- Graceful shutdown: https://grafana.com/blog/2024/02/28/how-we-handle-graceful-shutdowns-in-go/
- Multistage Docker Go: https://docs.docker.com/language/golang/build-images/
- Docker Compose: https://docs.docker.com/compose/gettingstarted/
- Makefile для Go: https://clavinjune.dev/en/blogs/makefile-for-golang-project/

**Результат:** Go-проект с чистой архитектурой, docker-compose (app + postgres + redis), Makefile.

---

## 📅 Неделя 2: PostgreSQL + миграции + User repository

**Цель:** Полноценная схема БД (7 таблиц), миграции, pgx pool.

| День | Задача | Время |
|------|--------|-------|
| 1 | **Миграция v1:** `users`, `user_identities`, `user_password_hash` | 2ч |
| 2 | **Миграция v2:** `user_sessions`, `api_keys`, `roles`, `user_roles` | 2ч |
| 3 | pgx pool: подключение, health check, конфигурация пула | 1.5ч |
| 4 | UserRepository: `Create`, `GetByID`, `GetByEmail`, `GetByIdentity` | 2ч |
| 5 | IdentityRepository: `Create`, `FindByProvider`, `UpdateLastLogin` | 1.5ч |
| 6 | SessionRepository + RoleRepository + APIKeyRepository (интерфейсы + реализация) | 2ч |

**Изучить:**
- Goose миграции: https://github.com/pressly/goose
- pgx pool: https://github.com/jackc/pgx
- pgx wiki: https://github.com/jackc/pgx/wiki
- PostgreSQL UUID: https://www.postgresql.org/docs/current/uuid-ossp.html
- PostgreSQL array: https://www.postgresql.org/docs/current/arrays.html
- INET type: https://www.postgresql.org/docs/current/datatype-net-types.html
- Repository pattern в Go: https://threedots.tech/post/repository-pattern-in-go/
- pgx batch queries: https://github.com/jackc/pgx?tab=readme-ov-file#batch-queries

**Результат:** 7 таблиц в PostgreSQL, полный набор repository-слоя.

---

## 📅 Неделя 3: gRPC + buf + Register/Login/JWT

**Цель:** Рабочий gRPC сервер с Register, Login, JWT генерацией.

| День | Задача | Время |
|------|--------|-------|
| 1 | Установка buf, создание `buf.gen.yaml`, `buf.yaml`, proto-директории | 1.5ч |
| 2 | Proto: `auth.proto` — AuthService (Register, Login, RefreshToken, Logout, ValidateToken) | 2ч |
| 3 | Генерация Go-кода через buf, подключение gRPC сервера в `cmd/auth/main.go` | 1.5ч |
| 4 | UseCase: Register (bcrypt cost=12, валидация email/username) | 2ч |
| 5 | UseCase: Login (verify password → JWT access + refresh RS256) | 2ч |
| 6 | gRPC handler: Register, Login, RefreshToken, Logout | 1.5ч |

**Изучить:**
- Buf CLI: https://buf.build/docs/installation/
- Buf + Go tutorial: https://buf.build/docs/tutorials/getting-started-with-buf-and-go/
- Protobuf style guide: https://protobuf.dev/programming-guides/style/
- gRPC concepts: https://grpc.io/docs/what-is-grpc/introduction/
- grpc-go quickstart: https://grpc.io/docs/languages/go/quickstart/
- bcrypt в Go: https://pkg.go.dev/golang.org/x/crypto/bcrypt
- go-playground/validator: https://github.com/go-playground/validator
- JWT RS256 в Go: https://github.com/golang-jwt/jwt
- JWT best practices: https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/
- gRPC error handling: https://grpc.io/docs/guides/error/

**Результат:** gRPC сервер с Register/Login, JWT (RS256 access + refresh).

---

## 📅 Неделя 4: Redis, refresh tokens, rate limiting, interceptorы

**Цель:** Redis-слой, полноценный token management, защита.

| День | Задача | Время |
|------|--------|-------|
| 1 | Redis подключение: go-redis client, health check | 1.5ч |
| 2 | Refresh token store в Redis (TTL 7d) + blacklist access tokens | 2ч |
| 3 | Token Bucket rate limiter на Redis (10 req/s, burst 20) | 2ч |
| 4 | gRPC unary interceptors: logging (slog), recovery (panic → error) | 2ч |
| 5 | Rate limiter interceptor + token validation interceptor | 1.5ч |
| 6 | UseCase: RefreshToken + Logout (blacklist access, delete refresh) | 2ч |

**Изучить:**
- go-redis: https://github.com/redis/go-redis
- Redis best practices: https://redis.io/docs/management/optimization/
- Redis TTL: https://redis.io/commands/expire/
- Token Bucket: https://en.wikipedia.org/wiki/Token_bucket
- Rate limiting with Redis: https://redis.io/glossary/rate-limiting/
- go-grpc-middleware: https://github.com/grpc-ecosystem/go-grpc-middleware
- gRPC interceptors: https://grpc.io/docs/guides/interceptors/
- gRPC metadata: https://grpc.io/docs/guides/metadata/
- Bearer token: https://datatracker.ietf.org/doc/html/rfc6750
- Refresh token rotation: https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation

**Результат:** Redis-слой, rate limiting, gRPC interceptors, полный цикл авторизации.

---

## 📅 Неделя 5: Google OAuth2 + multi-tenant

**Цель:** Вход через Google, tenant isolation.

| День | Задача | Время |
|------|--------|-------|
| 1 | Google OAuth2: создание проекта в Google Cloud Console, настройка redirect URI | 1ч |
| 2 | LoginOAuth2: state + code → token exchange → userinfo → create/find identity | 2.5ч |
| 3 | Sign-up via OAuth2: если identity не найден → создать user + identity | 1.5ч |
| 4 | Multi-tenant: tenant_id в контекст, фильтрация репозиториев по tenant_id | 2ч |
| 5 | Proto: LoginOAuth2 → gRPC handler → ручной тест через grpcurl | 2ч |
| 6 | Интеграционные тесты: OAuth callback mocking (httptest.Server) | 2ч |

**Изучить:**
- Google OAuth2 setup: https://developers.google.com/identity/protocols/oauth2
- OAuth2 в Go: https://pkg.go.dev/golang.org/x/oauth2
- Google OAuth2 Go example: https://developers.google.com/identity/protocols/oauth2/web-server#go
- Google sub claim: https://developers.google.com/identity/openid-connect/openid-connect#an-id-tokens-payload
- Multi-tenant patterns: https://learn.microsoft.com/en-us/azure/architecture/patterns/multitenant-saas
- Row-level security PostgreSQL: https://www.postgresql.org/docs/current/ddl-rowsecurity.html
- grpcurl: https://github.com/fullstorydev/grpcurl
- httptest.Server: https://pkg.go.dev/net/http/httptest#Server

**Результат:** Вход через Google, multi-tenant изоляция, тесты.

---

## 📅 Неделя 6: RBAC + API Keys

**Цель:** Полноценная авторизация: роли, permission'ы, сервис-сервис аутентификация.

| День | Задача | Время |
|------|--------|-------|
| 1 | Proto: RoleService (CreateRole, ListRoles, UpdateRole, DeleteRole, AssignRole, CheckPermission) | 1.5ч |
| 2 | RoleRepository + PermissionRepository (CRUD, связь user ↔ role) | 2ч |
| 3 | UseCase: AssignRole, RevokeRole, CheckPermission, GetUserRoles | 2ч |
| 4 | gRPC handler: CheckPermission (с локальным кэшированием ролей) | 2ч |
| 5 | API Keys: CreateAPIKey (генерация + SHA-256 хеш), ValidateAPIKey, RevokeAPIKey | 2.5ч |
| 6 | API Key interceptor: извлечение из gRPC metadata → валидация → user context | 1.5ч |

**Изучить:**
- RBAC model: https://en.wikipedia.org/wiki/Role-based_access_control
- PostgreSQL many-to-many: https://www.postgresql.org/docs/current/tutorial-join.html
- ABAC vs RBAC: https://www.okta.com/identity-101/role-based-access-control-vs-attribute-based-access-control/
- Caching в Go: https://towardsdev.com/caching-patterns-in-go-238ce578b3f4
- API key best practices: https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-api-key-registration
- crypto/rand: https://pkg.go.dev/crypto/rand

**Результат:** RBAC с кастомными ролями, API Keys для service-to-service auth.

---

## 📅 Неделя 7: REST (grpc-gateway) + auth-client-go SDK

**Цель:** REST API через grpc-gateway + Go-пакет для других сервисов.

| День | Задача | Время |
|------|--------|-------|
| 1 | grpc-gateway: настройка `buf.gen.yaml`, proto annotations (google.api.http) | 2ч |
| 2 | Google API HTTP annotations в proto: `POST /v1/auth/login`, `POST /v1/auth/register`, `POST /v1/auth/refresh` | 2ч |
| 3 | CORS middleware + запуск HTTP + gRPC в одном процессе | 1.5ч |
| 4 | **SDK:** `pkg/authclient/` — gRPC client interceptor (вставляет JWT в metadata) | 2ч |
| 5 | **SDK:** Server interceptor — `AuthUnaryServerInterceptor` (извлекает JWT, валидирует локально) | 2.5ч |
| 6 | **SDK:** Permission middleware — `RequirePermission(perm string)` interceptor + README | 2ч |

**Изучить:**
- grpc-gateway: https://github.com/grpc-ecosystem/grpc-gateway
- grpc-gateway tutorial: https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/
- HTTP annotation proto: https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
- CORS в Go: https://github.com/rs/cors
- gRPC client interceptor: https://github.com/grpc/grpc-go/blob/master/examples/features/interceptor/README.md
- JWT RS256 validation: https://golang-jwt.github.io/jwt/usage/signing_methods/#rsa
- gRPC middleware chain: https://itnext.io/chaining-grpc-interceptors-in-go-5b817b4803b3

**Результат:** REST API + gRPC, SDK для подключения auth в другие сервисы.

---

## 📅 Неделя 8: Деплой + CI/CD + документация

**Цель:** Pet project в production-like окружении.

| День | Задача | Время |
|------|--------|-------|
| 1 | OpenTelemetry: OTEL SDK, Prometheus метрики | 2.5ч |
| 2 | Prometheus: `grpc_requests_total`, `grpc_request_duration_seconds`, `db_queries_total` | 2ч |
| 3 | Jaeger: spans в gRPC handler + SQL запросы (pgx tracing) | 2ч |
| 4 | GitHub Actions: `.github/workflows/ci.yml` — lint + test + build + docker push | 2ч |
| 5 | Helm chart: `deploy/charts/auth-service/` — Deployment, Service, ConfigMap, Secrets | 2.5ч |
| 6 | README: архитектура, proto-спецификация, примеры запросов, как запустить | 1.5ч |

**Изучить:**
- OpenTelemetry Go: https://opentelemetry.io/docs/languages/go/getting-started/
- Prometheus Go client: https://github.com/prometheus/client_golang
- gRPC мониторинг: https://grafana.com/blog/2024/03/05/how-to-monitor-grpc-with-prometheus/
- OpenTelemetry tracing: https://opentelemetry.io/docs/languages/go/instrumentation/
- GitHub Actions Go: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go
- Docker GH Actions: https://docs.docker.com/build/ci/github-actions/
- Helm: https://helm.sh/docs/chart_template_guide/getting_started/
- Helm best practices: https://helm.sh/docs/chart_best_practices/

**Результат:** Helm chart, CI/CD, мониторинг, документация.

---

## 📦 Чеклист готовности Auth Service

- [ ] Clean Architecture структура (cmd/internal/domain/usecase/repository/handler/pkg)
- [ ] gRPC + REST (grpc-gateway) API
- [ ] Register / Login / RefreshToken / Logout
- [ ] Google OAuth2 (LoginOAuth2)
- [ ] JWT RS256 (access 15min + refresh 7d)
- [ ] bcrypt (cost=12)
- [ ] Redis: refresh token store + blacklist + rate limiter
- [ ] Multi-tenant изоляция (tenant_id)
- [ ] RBAC: кастомные роли + permission (CRUD через API)
- [ ] API Keys для service-to-service auth
- [ ] gRPC interceptors (logging, recovery, rate limit, auth)
- [ ] Graceful shutdown
- [ ] Health checks (Liveness + Readiness)
- [ ] auth-client-go SDK (gRPC interceptors, permission middleware)
- [ ] OpenTelemetry: Prometheus метрики + Jaeger traces
- [ ] GitHub Actions (lint + test + build + push)
- [ ] Docker + docker-compose
- [ ] Helm chart
- [ ] README с архитектурой и примерами
- [ ] Unit + Integration тесты, 80%+ coverage

---

## Схема БД

```sql
-- 1. Users (ядро)
CREATE TABLE users (
    user_id       UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id     UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001',
    login         VARCHAR(50) UNIQUE NOT NULL,
    display_name  VARCHAR(255),
    email_primary VARCHAR(255) UNIQUE NOT NULL,
    locale        VARCHAR(10),
    timezone      VARCHAR(50),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ,
    is_active     BOOLEAN DEFAULT TRUE
);

-- 2. User identities (OAuth2/SSO — Google, GitHub, Apple)
CREATE TABLE user_identities (
    identity_id      UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id          UUID NOT NULL REFERENCES users(user_id),
    provider         VARCHAR(20) NOT NULL,      -- 'google', 'github', 'apple'
    provider_user_id VARCHAR(255) NOT NULL,      -- sub от провайдера
    provider_email   VARCHAR(255),
    provider_name    VARCHAR(255),
    avatar_url       TEXT,
    last_login_at    TIMESTAMPTZ,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_user_id)
);
CREATE INDEX idx_identities_lookup ON user_identities (provider, provider_user_id);

-- 3. Password hash (только для email+password регистрации)
CREATE TABLE user_password_hash (
    user_id             UUID PRIMARY KEY REFERENCES users(user_id),
    password_hash       TEXT NOT NULL,                  -- bcrypt cost=12
    reset_token_hash    TEXT,                           -- SHA-256 хеш токена сброса
    reset_token_expires_at TIMESTAMPTZ,                 -- срок действия токена сброса
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ
);

-- 4. Sessions (refresh tokens)
CREATE TABLE user_sessions (
    session_id    UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id       UUID NOT NULL REFERENCES users(user_id),
    identity_id   UUID REFERENCES user_identities(identity_id),
    refresh_token TEXT UNIQUE NOT NULL,           -- SHA-256 хеш
    device_info   TEXT,
    ip_address    INET,
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    revoked_at    TIMESTAMPTZ                     -- NULL = активна
);
CREATE INDEX idx_sessions_user ON user_sessions (user_id) WHERE revoked_at IS NULL;

-- 5. API keys (service-to-service)
CREATE TABLE api_keys (
    api_key_id    UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id     UUID NOT NULL,
    service_name  TEXT NOT NULL,
    key_prefix    VARCHAR(8) NOT NULL,            -- первые 8 символов ключа (для идентификации)
    key_hash      TEXT NOT NULL,                  -- SHA-256 полного ключа
    permissions   TEXT[] NOT NULL DEFAULT '{}',   -- ['users:read', 'users:write']
    created_by    UUID REFERENCES users(user_id),
    expires_at    TIMESTAMPTZ,
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (key_hash)
);

-- 6. Roles & Permissions (RBAC)
CREATE TABLE roles (
    role_id       UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id     UUID NOT NULL,
    name          TEXT NOT NULL,                  -- 'admin', 'moderator', 'user'
    permissions   TEXT[] NOT NULL DEFAULT '{}',   -- ['users:delete', 'posts:moderate']
    UNIQUE (tenant_id, name)
);

CREATE TABLE user_roles (
    user_id       UUID NOT NULL REFERENCES users(user_id),
    role_id       UUID NOT NULL REFERENCES roles(role_id),
    assigned_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
```

---

## gRPC proto контракт

```protobuf
service AuthService {
  // Authentication
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc LoginOAuth2(OAuth2Request) returns (AuthResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
  rpc Logout(LogoutRequest) returns (Empty);

  // Token validation (для middleware других сервисов)
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc IntrospectToken(IntrospectTokenRequest) returns (IntrospectTokenResponse);

  // Authorization
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
  rpc GetUserRoles(GetUserRolesRequest) returns (GetUserRolesResponse);

  // API Keys (service-to-service)
  rpc CreateAPIKey(CreateAPIKeyRequest) returns (CreateAPIKeyResponse);
  rpc ValidateAPIKey(ValidateAPIKeyRequest) returns (ValidateAPIKeyResponse);
  rpc RevokeAPIKey(RevokeAPIKeyRequest) returns (Empty);

  // Admin / User management
  rpc GetUser(GetUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc AssignRole(AssignRoleRequest) returns (Empty);
  rpc RevokeRole(RevokeRoleRequest) returns (Empty);

  // Role management (CRUD)
  rpc CreateRole(CreateRoleRequest) returns (Role);
  rpc ListRoles(ListRolesRequest) returns (ListRolesResponse);
  rpc UpdateRole(UpdateRoleRequest) returns (Role);
  rpc DeleteRole(DeleteRoleRequest) returns (Empty);
}