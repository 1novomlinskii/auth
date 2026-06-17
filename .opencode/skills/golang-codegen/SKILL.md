---
name: golang-codegen
description: Go code generation — go:generate, stringer, sqlc, mockgen/mockery, go:embed
license: MIT
compatibility: opencode
---

# Go Code Generation

## When to Activate

Activate when the task involves:
- Adding //go:generate directives
- Generating String() methods for enum types (stringer)
- Generating type-safe database code from SQL (sqlc)
- Generating mocks for interfaces (mockgen, mockery)
- Embedding static files with go:embed
- Deciding when to use (or not use) code generation

Do not activate for: writing regular Go code, dependency management, or testing patterns.

## //go:generate

Central mechanism for running generators:

```go
//go:generate go run golang.org/x/tools/cmd/stringer -type=Status
//go:generate go run github.com/vektra/mockery/v2@v2 --name=UserRepository
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
```

Run: `go generate ./...`

## stringer

Generate `String()` for enum types:

```go
type Status int

const (
    StatusUnknown Status = iota
    StatusActive
    StatusInactive
)

//go:generate stringer -type=Status -linecomment
```

## sqlc

Generate type-safe Go code from SQL:

```sql
-- query.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
```

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "migrations/"
    gen:
      go:
        package: "repository"
        out: "internal/repository"
```

```go
// Generated:
func (q *Queries) GetUser(ctx context.Context, id string) (User, error)
```

## mockgen / mockery

```go
//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks
//go:generate mockery --name=UserRepository --output=mocks --outpkg=mocks
```

## go:embed

```go
//go:embed templates/*
var templateFS embed.FS

//go:embed config.yaml
var configData []byte

//go:embed static/css/*.css
var cssFiles embed.FS
```

## When to Use

- **stringer** — always for enum types
- **sqlc** — for PostgreSQL projects (type safety, no ORM)
- **mockgen/mockery** — for interfaces with many methods
- **embed** — for templates, migrations, static files

## When NOT to Use

- **Don't** generate code that can be written in 5 lines by hand
- **Don't** generate boilerplate where a simple library suffices
- **Don't** use bash scripts for codegen — use `//go:generate` instead

## Integration

- golang-testing — mockgen for test mocks
- golang-dependencies — recommended codegen tools
- golang-project-structure — where to place generated code