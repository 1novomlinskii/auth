---
name: golang-conventions
description: Go naming, composition over inheritance, interfaces, zero values, formatting, comments
license: MIT
compatibility: opencode
---

# Go Conventions & Code Style

## When to Activate

Activate when the task involves:
- Naming variables, functions, types, interfaces, packages
- Structuring code with embedding vs inheritance
- Defining interfaces (consumer-side, small interfaces)
- Using zero values
- Writing doc comments
- Any Go code style decision

Do not activate for: toolchain commands, testing, error handling specifics.

## Naming

| Kind | Style | Example |
|------|-------|---------|
| Local variables | camelCase | `userCount`, `req`, `err` |
| Function params | camelCase | `ctx context.Context` |
| Unexported funcs | camelCase | `parseConfig()` |
| Exported funcs | PascalCase | `NewServer()` |
| Types | PascalCase | `UserService` |
| Interfaces | PascalCase + -er | `Reader`, `Writer` |
| Constants | camelCase / PascalCase | `maxRetries` / `MaxRetries` |

Rules:
- Shorter scope = shorter name (`i` not `index`, `ctx` not `context`)
- Single-letter only for loop indices (i, j, k) and receivers (t *testing.T)
- Receiver: 1-2 letters: `s *Server`, `u *UserService`
- No Hungarian notation, no snake_case, no getters/setters

## Composition over Inheritance

```go
// CORRECT: composition via embedding
type Logger struct{ *log.Logger }
type Server struct {
    *Logger
    db    *sql.DB
    cache *redis.Client
}
```

## Interfaces

- **Define interfaces where they're used** (consumer-side), not where implemented
- Keep interfaces small — fewer methods is better
- Accept interfaces, return concrete types

```go
// CORRECT: interface at consumer
type UserRepository interface {
    Get(ctx context.Context, id string) (*User, error)
}
```

## Zero Values

Use zero values instead of explicit initialization:

```go
var buf bytes.Buffer       // ready to use
var mu sync.Mutex          // ready to use
var slice []string         // nil, but append works
```

## Comments

- All exported identifiers **must** have doc comments
- Doc comments start with the identifier name:

```go
// Package server provides HTTP server utilities.
package server

// UserService handles user-related business logic.
type UserService struct{}
```

- Inline comments explain WHY, not WHAT

## Integration

- golang-error-handling — error wrapping conventions
- golang-project-structure — package naming, internal/ conventions
- golang-testing — test helpers and mock conventions