---
name: golang-project-structure
description: Go project layout — cmd/, internal/, pkg/ structure, layered architecture, single-package projects
license: MIT
compatibility: opencode
---

# Go Project Structure

## When to Activate

Activate when the task involves:
- Setting up a new Go project or module
- Organizing packages (cmd/, internal/, pkg/)
- Adding new layers (handler, service, repository)
- Deciding where to place code in an existing project
- Config structure and initialization

Do not activate for: writing individual functions, tests, or toolchain commands. Use golang-conventions for naming, golang-dependencies for go.mod setup.

## Standard Layout

```
myproject/
├── cmd/
│   └── nameapp/
│       └── main.go              # thin entry point only
├── internal/
│   ├── app/
│   │   ├── app.go               # app startup, dependency init
│   │   └── server.go            # HTTP/gRPC server setup
│   ├── config/
│   │   └── config.go            # config structs
│   ├── delivery/
│   │   ├── http/                # HTTP handlers, middleware, routes
│   │   └── grpc/                # gRPC services, protobuf contracts
│   ├── entity/                  # domain models/entities
│   ├── err/                     # custom application errors
│   ├── metrics/                 # Prometheus etc.
│   ├── services/                # external API clients
│   ├── usecase/
│   │   └── user.go              # business logic
│   └── repository/
│       └── pg/
│           ├── postgres.go      # DB connection, transactions
│           └── query/           # raw SQL queries
├── pkg/
│   └── model/
│       └── types.go             # public types for external consumers
├── migrations/
├── config/
├── local_deploy/
├── go.mod
├── Makefile
└── README.md
```

## Principles

### cmd/ — thin main only

```go
func main() {
    srv := server.New()
    if err := srv.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### internal/ — private code

Packages under `internal/` cannot be imported from outside the module. All internal implementation goes here.

### pkg/ — public code

Only code safe to import from other projects.

### Layer Separation

- **handler/transport** — HTTP/gRPC handlers, request/response parsing
- **service/usecase** — business logic
- **repository** — storage (DB, API, files)
- **entity/domain** — pure data structures (no dependencies)

## Single Package

For small projects/utilities — all files in module root:

```
myproject/
├── main.go
├── handler.go
├── store.go
├── go.mod
└── Makefile
```

Don't overcomplicate structure for small projects.

## Configuration

```go
type Config struct {
    Port    int           `env:"PORT" default:"8080"`
    DB      DBConfig
}

type DBConfig struct {
    DSN string `env:"DB_DSN" required:"true"`
}
```

- Use `os.Getenv` or libraries like `envconfig`/`viper`
- Don't hardcode config values (exception: tests)

## Integration

- golang-conventions — package naming, file organization
- golang-error-handling — custom error types in err/ package
- golang-dependencies — go module setup