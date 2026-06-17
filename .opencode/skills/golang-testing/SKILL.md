---
name: golang-testing
description: Go testing — table-driven tests, subtests, test helpers, mocks, golden files, integration tests, coverage
license: MIT
compatibility: opencode
---

# Go Testing

## When to Activate

Activate when the task involves:
- Writing unit tests (table-driven, subtests)
- Creating test helpers and fixtures
- Mocking interfaces (mockgen or manual mocks)
- Golden file tests for complex output
- Integration tests with build tags
- Measuring test coverage
- Testing error cases and edge cases

Do not activate for: writing production code, toolchain commands, or dependency management.

## Table-Driven Tests

Primary pattern for multiple test cases:

```go
func TestParsePhone(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {name: "valid mobile", input: "+79001234567", want: "79001234567"},
        {name: "invalid short", input: "123", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParsePhone(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Subtests

- `t.Run()` for grouping and isolation
- Each subtest runs independently
- Run single subtest: `go test -run TestParsePhone/valid_mobile`

## Test Helpers

```go
func SetupDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("postgres", testDSN)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    t.Cleanup(func() { db.Close() })
    return db
}
```

## Mocks

```go
//go:generate mockgen -source=service.go -destination=mocks/service.go -package=mocks

type UserRepository interface {
    Get(ctx context.Context, id string) (*User, error)
}
```

Or manual mock structs:

```go
type mockRepo struct {
    getUserFn func(ctx context.Context, id string) (*User, error)
}

func (m *mockRepo) Get(ctx context.Context, id string) (*User, error) {
    return m.getUserFn(ctx, id)
}
```

## Golden Files

For complex output tests (JSON, HTML):

```go
func TestHandler(t *testing.T) {
    golden := filepath.Join("testdata", t.Name()+".golden")
    got := doRequest(t, "/api/users")
    if *update {
        os.WriteFile(golden, got, 0644)
    }
    expected, _ := os.ReadFile(golden)
    assert.Equal(t, string(expected), string(got))
}
```

## Integration Tests

- Mark with build tag: `//go:build integration`
- Run separately: `go test -tags=integration ./...`

```go
//go:build integration

func TestPostgresIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
}
```

## Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

For CI: `go test -coverprofile=coverage.out -covermode=atomic ./...`

## Integration

- golang-concurrency — goleak for goroutine leak checks
- golang-codegen — mockgen for mock generation
- golang-conventions — test helper naming