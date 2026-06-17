---
name: golang-toolchain
description: Go toolchain commands — gofmt, go vet, go mod, go build, go test, go generate, ldflags, go:embed
license: MIT
compatibility: opencode
---

# Go Toolchain

Essential Go commands and toolchain usage.

## When to Activate

Activate when the task involves:
- Running Go commands (build, test, vet, fmt, mod)
- Setting up a new Go module
- Configuring build flags or ldflags
- Formatting code with gofmt/goimports
- Static analysis with go vet
- Code generation via go:generate
- Profiling with pprof/trace

Do not activate for: project structure layout, dependency management decisions, or writing Go code logic.

## Core Commands

| Command | Purpose |
|---------|---------|
| `go mod init <module>` | Initialize a new module |
| `go mod tidy` | Clean and add dependencies |
| `go mod vendor` | Create vendor directory |
| `go build ./...` | Build all packages |
| `go run .` | Run the application |
| `go test ./...` | Run all tests |
| `go test -cover ./...` | Tests with coverage |
| `go vet ./...` | Static analysis |
| `gofmt -w .` | Format all Go files |
| `goimports -w .` | Format + sort imports |
| `go generate ./...` | Run codegen directives |
| `go clean -cache` | Clear build cache |

## Formatting

- `gofmt` is mandatory — unformatted code is invalid
- `goimports` = gofmt + auto-sort/import resolution
- Use **tabs only**, tab width = 8 (Go standard)
- Group imports: stdlib → third-party → local (blank lines between)

## Build

```go
// Inject version via ldflags
go build -ldflags="-X main.Version=v1.0.0 -X main.Commit=$(git rev-parse HEAD)" ./cmd/app
```

## go:embed

```go
//go:embed templates/*
var templateFS embed.FS

//go:embed config.yaml
var configData []byte
```

## go vet

Run before every commit. Catches: unused variables, incorrect printf formats, dead code, lock issues.

## go mod

- Semantic versioning: `v1.2.3`
- Module path = repository + optional suffix (`/v2` for major versions)
- Always run `go mod tidy` after changing imports
- If using vendor: `go mod vendor` after tidy

## Profiling

- `go tool pprof` — CPU/memory profiling
- `go tool trace` — goroutine tracing

## Integration

- golang-dependencies — go.mod, replace, vendor details
- golang-codegen — go:generate, stringer, sqlc
- golang-project-structure — how to organize cmd/, internal/, pkg/