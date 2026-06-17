---
name: golang-error-handling
description: Go error handling — wrapping with %w, sentinel errors, custom error types, panic/recover, structured logging
license: MIT
compatibility: opencode
---

# Go Error Handling

## When to Activate

Activate when the task involves:
- Writing error handling code (if err != nil)
- Defining sentinel errors (ErrNotFound, ErrConflict)
- Creating custom error types
- Wrapping errors with context
- Using panic/recover in HTTP middleware or goroutines
- Logging errors at system boundaries

Do not activate for: general code style, testing patterns, or concurrency patterns.

## Basic Principles

- **Always handle errors.** Never ignore a returned error.
- Pattern: `if err != nil { return err }`
- Error is always the last return value.

## Wrapping Errors

Use `fmt.Errorf("context: %w", err)` to create error chains:

```go
if err != nil {
    return fmt.Errorf("get user %q: %w", id, err)
}
```

`%w` is the only way to correctly propagate an error with context.

## Sentinel Errors

Define sentinel errors for expected edge cases. Use `errors.Is()` to check:

```go
var (
    ErrNotFound = errors.New("user not found")
    ErrConflict = errors.New("user already exists")
)

if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

## Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Check via errors.As
var valErr *ValidationError
if errors.As(err, &valErr) {
    log.Printf("validation failed on %s: %s", valErr.Field, valErr.Message)
}
```

## Panic / Recover

- Panic only for truly exceptional situations (not for normal errors)
- HTTP middleware with recover:

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("panic recovered: %v", rec)
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

- In goroutines, always handle panic at top level:

```go
go func() {
    defer func() {
        if rec := recover(); rec != nil {
            log.Printf("goroutine panic: %v", rec)
        }
    }()
    doWork()
}()
```

## Logging Errors

- Log errors at system boundary (handler/server), not at depth (service/repo)
- At depth — wrap with context and propagate up
- Use structured logging (slog, logrus, zap):

```go
slog.Error("failed to process request",
    "error", err,
    "user_id", userID,
)
```

## Integration

- golang-conventions — error naming
- golang-concurrency — panic handling in goroutines
- golang-testing — testing error cases