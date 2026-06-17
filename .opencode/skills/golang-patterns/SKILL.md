---
name: golang-patterns
description: Go idiomatic patterns — Functional Options, Repository, Middleware chain, Factory, Builder, Worker Pool
license: MIT
compatibility: opencode
---

# Go Patterns

## When to Activate

Activate when the task involves:
- Creating configurable constructors (Options pattern)
- Implementing data access layer (Repository pattern)
- HTTP middleware chaining
- Object creation with multiple variants (Factory)
- Building complex objects step by step (Builder)
- Concurrent task processing (Worker Pool)

Do not activate for: simple constructors, basic CRUD, or writing tests. When in doubt, prefer simple code over patterns.

## Options Pattern (Functional Options)

```go
type Server struct {
    port    int
    timeout time.Duration
    logger  *slog.Logger
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(t time.Duration) Option {
    return func(s *Server) { s.timeout = t }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        port:    8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(WithPort(9000), WithTimeout(10*time.Second))
```

## Repository Pattern

```go
type UserRepository interface {
    Get(ctx context.Context, id string) (*User, error)
    List(ctx context.Context, filter UserFilter) ([]*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}
```

## Middleware Pattern

```go
type Middleware func(http.Handler) http.Handler

func WithLogging(logger *slog.Logger) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Info("request", "method", r.Method, "path", r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
```

## Factory Pattern

```go
type Parser interface {
    Parse(data []byte) (*Document, error)
}

type JSONParser struct{}
type XMLParser struct{}

func NewParser(format string) (Parser, error) {
    switch format {
    case "json":
        return &JSONParser{}, nil
    case "xml":
        return &XMLParser{}, nil
    default:
        return nil, fmt.Errorf("unsupported format: %s", format)
    }
}
```

## Builder Pattern

```go
type QueryBuilder struct {
    q Query
}

func Select(table string) *QueryBuilder {
    return &QueryBuilder{q: Query{table: table}}
}

func (b *QueryBuilder) Where(field string, value any) *QueryBuilder { ... }
func (b *QueryBuilder) OrderBy(field string) *QueryBuilder { ... }
func (b *QueryBuilder) Limit(n int) *QueryBuilder { ... }
func (b *QueryBuilder) Build() Query { return b.q }

// Usage
q := Select("users").Where("age", 18).OrderBy("name").Limit(10).Build()
```

## Worker Pool

```go
func WorkerPool(ctx context.Context, jobs <-chan Job, n int) <-chan Result {
    results := make(chan Result, n)
    var wg sync.WaitGroup

    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case results <- job.Process():
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

## When NOT to Use Patterns

- **Don't** use interfaces for everything — if one implementation won't change, skip the interface
- **Don't** factory without need — `NewX()` is enough without variants
- **Go is not Java/TypeScript.** Visitor, Abstract Factory, Observer are almost always over-engineering

## Integration

- golang-conventions — interface naming, composition
- golang-concurrency — worker pool patterns
- golang-error-handling — error wrapping in patterns