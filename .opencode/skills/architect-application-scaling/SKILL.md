---
name: architect-application-scaling
description: Horizontal and vertical scaling, CQRS, Event Sourcing, caching strategies, CDN, stateless design
compatibility: opencode
---

# Application Scaling

## When to Activate

Activate when the task involves:
- Scaling an application to handle increased load
- Designing or reviewing caching strategy
- Separating read and write paths (CQRS)
- Event Sourcing for audit or rebuildability
- CDN integration for static/dynamic content
- Stateless vs stateful design decisions

Do not activate for: database scaling alone (`architect-database-scaling`), microservices decomposition (`architect-microservices`), or resilience patterns (`architect-system-design`).

## Core Concepts

### Scale Directions

| Direction | Mechanism | When |
|-----------|-----------|------|
| **Vertical (scale up)** | Bigger machine (CPU, RAM, disk) | Simple, no code changes, hits ceiling |
| **Horizontal (scale out)** | More instances behind a load balancer | Virtually unlimited, requires stateless design |

**Rule of thumb:** Vertical first ±30% capacity. Beyond that, horizontal is cheaper and more resilient.

### Stateless Design

The application instance holds no session state, no in-memory cache that cannot be rebuilt. State moves to:
- Database (PostgreSQL, Redis)
- Distributed cache (Redis Cluster, Memcached)
- Object storage (S3, MinIO)
- Client side (JWT, cookies)

Stateless instances can be scaled to zero, rolled out gradually, and replaced at any time.

## Detailed Topics

### Caching

| Layer | Technology | Cache scope | Eviction |
|-------|-----------|------------|----------|
| **Browser** | `Cache-Control`, `ETag` | Per-user | TTL |
| **CDN** | CloudFront, Cloudflare, Fastly | Global | TTL, invalidation |
| **Reverse proxy** | Nginx, Varnish, Envoy | Per-instance | LRU, TTL |
| **Application** | Redis, Memcached | Shared/distributed | TTL, LRU, LFU |
| **Database** | Buffer pool, query cache | Per-instance | Internal |

**Cache-Aside Pattern:**
```
1. Read from cache. Hit → return.
2. Miss → read from DB.
3. Store in cache with TTL.
4. Return.
```

**Write strategies:**
- **Write-through**: Write to cache + DB synchronously. Consistent, slower writes.
- **Write-behind**: Write to cache, async flush to DB. Fast writes, risk of loss.
- **Cache invalidation**: On update, delete cache key. Next read refreshes.

### CQRS (Command Query Responsibility Segregation)

Separate the write model from the read model when:
- Reads and writes have different shapes
- Read load >> write load
- Complex queries impact write performance

```
[Command] → Write Model (normalized) → [DB]
[Query]   → Read Model (denormalized) → [Cache/ReadDB]
```

**Entry-level CQRS:** Separate repository interfaces for read and write. **Full CQRS:** Separate databases, eventual consistency.

### Event Sourcing

Store state as a sequence of events. Current state = fold over events.

```
events = [UserCreated, EmailChanged, PasswordReset, ...]
current = fold(events, initial)
```

**When:** Full audit trail, temporal queries, complex state rebuild needed.
**Cost:** Event store, projection management, schema evolution.

### Message Brokers

| Broker | Strength | Avoid if |
|--------|----------|----------|
| **Kafka** | High throughput, replay, partitioning | Low latency <5ms, simple pub/sub |
| **RabbitMQ** | Complex routing, latency <1ms | Throughput >100K msg/s |
| **NATS** | Simplicity, low latency, at-most-once | Exactly-once, replay |

## Practical Guidance

### Scaling Decision Framework

1. Profile: what's the bottleneck? (CPU, memory, I/O, network)
2. Is the app stateless? If not → make it stateless first.
3. Add cache at the bottleneck layer (least specific cache first).
4. Add read replicas if DB is the bottleneck.
5. Consider CQRS only if read/write separation solves a measured problem.
6. Event Sourcing — only if you need full audit or temporal queries.

### Anti-patterns: Premature Distribution

A common mistake is adopting CQRS, Event Sourcing, or microservices "just in case." The cost of distribution (consistency, debugging, operations) is always higher than expected. Start monolithic with clear module boundaries.

## Guidelines

1. Profile before optimising. Cache without measurement is guessing.
2. Every cache must have a TTL (even if long). Infinite TTL = memory leak.
3. Stateless design enables horizontal scaling — invest here first.
4. Shared mutable state is the #1 scaling blocker. Pass state through DB/cache, not memory.
5. CDN for static assets is default unless there's a reason not to.
6. CQRS is a spectrum — start with separate repository interfaces, split DB only when needed.

## Gotchas

1. **Cache stampede**: When cache expires, thousands of requests hit DB simultaneously. Solution: request coalescing, early recomputation, or a global lock.
2. **Sticky sessions**: ELB session affinity ties a user to one instance. Enables stateful apps but breaks rolling deployments and auto-scaling. Remove when possible.
3. **CQRS without event bus**: Commands and queries share the same DB connection pool, negating the benefit. Separate connection pools minimum.
4. **Event Sourcing schema evolution**: Events are immutable once written. Changing an event schema needs upcasters (versioned event readers).
5. **Distributed caching as single point of failure**: Redis outage takes down your app. Have a fallback (degraded read from DB) or use Redis Cluster.

## Integration

- `architect-database-scaling` — read replicas as read model backend
- `architect-microservices` — event-driven communication, CQRS in services
- `architect-system-design` — resilience patterns for distributed caches