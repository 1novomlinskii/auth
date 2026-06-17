---
name: golang-libraries
description: Recommended Go libraries for HTTP, logging, database, migrations, Redis, config, testing, validation — and what to avoid
license: MIT
compatibility: opencode
---

# Go Libraries

## When to Activate

Activate when the task involves:
- Choosing a library for HTTP routing, logging, database access, etc.
- Deciding between stdlib vs third-party library
- Checking which Go libraries are recommended vs avoided in this project

Do not activate for: writing code logic, toolchain commands, or project structure. Activate golang-dependencies for adding/updating dependencies.

## Recommended Libraries

| Purpose | Library |
|---------|---------|
| HTTP router | `net/http` (stdlib, Go 1.22+) / `chi` / `gorilla/mux` |
| Logging | `log/slog` (stdlib, Go 1.21+) / `zap` / `rs/zerolog` |
| PostgreSQL | `pgx` (preferred) / `lib/pq` |
| SQLite | `modernc.org/sqlite` (CGo-free) |
| Migrations | `golang-migrate/migrate` / `pressly/goose` |
| Redis | `redis/go-redis` |
| Config | `envconfig` / `kelseyhightower/envconfig` / `spf13/viper` |
| Testing | `stdlib testing` + `stretchr/testify` |
| Validation | `go-playground/validator` / `go-ozzo/ozzo-validation` |
| CLI flags | `spf13/pflag` |
| Protobuf validation | `envoyproxy/protoc-gen-validate` |
| Mock | `golang/mock` / `vektra/mockery` |
| Templates | `valyala/quicktemplate` |

## Libraries to Avoid

- **ORM** (GORM, go-pg) — hides SQL, causes N+1 problems. Use sqlc/pgx instead.
- **Heavy frameworks** — Go doesn't need Spring. `net/http` + chi solves 99% of cases.
- **Fragile libraries** — check stars, maintenance status, Go version compatibility.

## Integration

- golang-dependencies — how to add/update these libraries
- golang-toolchain — go get command usage
- golang-codegen — sqlc for type-safe database code