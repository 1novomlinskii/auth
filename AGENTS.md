# AI Agent Instructions

Этот репозиторий содержит проект на языке голнаг. Микросервис Авторизации и аутентификации по grpc и rest API
который будет использоваться разными микросервисами


##  Инструкция
- Оспаривай мои решения. Если подход плохой, скажи прямо.
- Не соглашайся просто потому, что я спросил.
- Проводи ревью задачи, что сделано, что нужно от меня и какой следующий шаг.
- После каждой задачи предлагай, что улучшить или автоматизировать.

- это пет проект чтоб прокачаться в программировании до сеньера


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
```

## Технологии

- **Язык:** Go 1.23
- **API:** gRPC + REST (через grpc-gateway)
- **БД:** PostgreSQL 16 + Redis 7
- **Миграции:** Goose
- **Конфиг:** Viper (YAML + env)
- **Логгер:** zerolog

## Структура проекта

```
├── cmd/
│   └── auth/
│       └── main.go          # Точка входа
├── internal/
│   ├── app/
│   │   └── app.go           # HTTP сервер, healthcheck, graceful shutdown
│   └── config/
│       └── config.go        # Конфиг (Viper + YAML)
├── migrations/              # SQL миграции Goose
│   ├── 001_users.sql
│   ├── 002_user_identities.sql
│   ├── 003_user_password.sql
│   ├── 004_user_sessions.sql
│   ├── 005_api_keys.sql
│   ├── 006_roles.sql
│   └── 007_user_roles.sql
├── config.yml               # Конфиг по умолчанию
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## План развития

Подробный план на 8 недель — [01-auth-service-plan.md](01-auth-service-plan.md).
Dсе измениея по базе данных нужно отражать в [database-schema.md](doc/info/database-schema.md).
Как будет работать  [jwt-auth-flow.md](doc/info/jwt-auth-flow.md).
