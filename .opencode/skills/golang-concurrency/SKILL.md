---
name: golang-concurrency
description: Go concurrency — goroutines, errgroup, channels, sync.Mutex/RWMutex, sync.Once, context, goleak, anti-patterns
license: MIT
compatibility: opencode
---

# Go Concurrency

## When to Activate

Activate when the task involves:
- Writing goroutines
- Parallel task execution with errgroup
- Channel-based producer/consumer patterns
- Mutex/RWMutex for shared state protection
- Lazy initialization with sync.Once
- Context propagation with timeouts/cancellation
- Goroutine leak detection in tests
- Worker pool patterns

Do not activate for: general error handling, basic project setup, or standard CRUD operations without concurrency.

## Goroutines

- Launch with `go fn()`
- Always know when a goroutine will complete
- Use `sync.WaitGroup` or `errgroup` to wait for goroutine groups

```go
var wg sync.WaitGroup
for _, task := range tasks {
    wg.Add(1)
    go func(t Task) {
        defer wg.Done()
        t.Process()
    }(task)
}
wg.Wait()
```

## errgroup

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())
for _, task := range tasks {
    t := task
    g.Go(func() error {
        return t.Process(ctx)
    })
}
if err := g.Wait(); err != nil {
    log.Printf("one or more tasks failed: %v", err)
}
```

- First error cancels context → remaining goroutines get stop signal
- Preferred for parallel processing

## Channels

- Use channels for communication, not synchronization
- Owner creates, closes, and sends to channel
- Receiver only reads

```go
ch := make(chan Result, 100)

// producer
go func() {
    defer close(ch)
    for _, item := range items {
        ch <- process(item)
    }
}()

// consumer
for result := range ch {
    results = append(results, result)
}
```

## sync.Mutex / RWMutex

- Mutex for shared state protection
- RWMutex when many readers, few writers
- Always use defer for unlock

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    return item, ok
}

func (c *Cache) Set(key string, item Item) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = item
}
```

## sync.Once

For lazy initialization:

```go
var (
    once   sync.Once
    config *Config
)

func GetConfig() *Config {
    once.Do(func() {
        config = loadConfig()
    })
    return config
}
```

## Context

- Pass `context.Context` as first argument
- Use for timeouts, cancellation, metadata

```go
func FetchData(ctx context.Context, id string) (*Data, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    return db.Query(ctx, id)
}
```

## goleak

Test for goroutine leaks:

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

## Anti-patterns

- **Don't** use `time.Sleep` for synchronization — use channels or WaitGroup
- **Don't** use `sync/atomic` for high-level synchronization
- **Don't** spawn goroutines in a loop without control — use worker pool or errgroup with limit
- **Don't** close a channel from multiple goroutines — causes panic

## Integration

- golang-error-handling — panic recovery in goroutines
- golang-testing — goleak for testing
- golang-patterns — worker pool pattern