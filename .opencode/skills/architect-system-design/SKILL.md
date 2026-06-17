---
name: architect-system-design
description: CAP theorem, PACELC, consistency models, resilience patterns (circuit breaker, bulkhead, retry, timeout), SLA/SLO/SLI, trade-off analysis
compatibility: opencode
---

# System Design

## When to Activate

Activate when the task involves:
- CAP theorem analysis for distributed systems
- Consistency model choices (strong, eventual, causal)
- Resilience pattern design (circuit breaker, bulkhead, retry, timeout)
- Defining SLA / SLO / SLI for a system
- Trade-off analysis between design alternatives
- Rate limiting strategy
- Request tracing context (correlation IDs)

Do not activate for: specific database or application scaling (`architect-database-scaling`, `architect-application-scaling`), or individual API design (`architect-api-design`).

## Core Concepts

### CAP Theorem (Brewer)

A distributed data store can provide at most **two** of three guarantees:

| Property | Meaning |
|----------|---------|
| **Consistency** | Every read receives the most recent write |
| **Availability** | Every request receives a non-error response |
| **Partition Tolerance** | System continues operating despite network failures |

**Reality:** Network partitions are inevitable (P is always required). So choose **CP vs AP**:
- **CP** (PostgreSQL primary, Zookeeper): Sacrifice availability during partition to guarantee consistency
- **AP** (DynamoDB, Cassandra): Accept inconsistency during partition to stay available

### PACELC Extension

If partition (P), trade between A and C; if Else (E), trade between Latency (L) and Consistency (C):
- **PA/EL**: DynamoDB — available during partition, low latency / weak consistency normally
- **PC/EC**: Spanner — consistent during partition, consistent / higher latency normally

## Detailed Topics

### Consistency Models (strongest → weakest)

| Model | Guarantee | Use case |
|-------|-----------|----------|
| **Linearisability** | Operations appear atomic, instant, in real-time order | Leader election, distributed locks |
| **Sequential** | Operations appear in some order, consistent per-client | RDBMS primary |
| **Causal** | Causally related operations seen in order | Social feeds, collaborative editing |
| **Eventual** | All replicas converge eventually, no time bound | DNS, CDN, analytics |

### Resilience Patterns

**Circuit Breaker:**
```
[Client] → [Circuit Breaker] → [Upstream]
            OPEN: fail fast
            CLOSED: pass through
            HALF-OPEN: test one request
```
States: CLOSED (normal) → OPEN (after N failures) → HALF-OPEN (after timeout) → CLOSED or OPEN.
Key: detect upstream failure, fail fast instead of hanging.

**Bulkhead:**
Isolate resources into pools so failure in one pool doesn't drain others.
```
[Connection pool A] → [Service A]
[Connection pool B] → [Service B] (independent)
```
Key: OS processes, thread pools, connection pools are natural bulkheads.

**Retry with Backoff:**
```
retries = 3
for i in range(retries):
    try: call()
    except: sleep(2^i + jitter)  # exponential backoff
```
Key: Always add jitter to prevent thundering herd. Use exponential backoff, not linear.

**Timeout:**
Every external call must have a timeout. No timeout = resource leak under upstream failure.
```
http.Client{Timeout: 5 * time.Second}
db.ExecContext(ctx, "SELECT ...")  // uses ctx deadline
```

### Rate Limiting

| Algorithm | Characteristic | Use |
|-----------|---------------|-----|
| **Token Bucket** | Allows bursts, steady average | Most APIs |
| **Leaky Bucket** | Smooths bursts, fixed rate | Downstream protection |
| **Sliding Window** | Fair per time window | Per-user rate limiting |
| **Fixed Window** | Simple, boundary spikes | Coarse throttling |

### SLA / SLO / SLI

| Term | Meaning | Example |
|------|---------|---------|
| **SLI** (Indicator) | Actual measurement | p99 latency = 250ms, availability = 99.95% |
| **SLO** (Objective) | Target value | p99 latency ≤ 300ms, availability ≥ 99.9% |
| **SLA** (Agreement) | Contract with consequences | If SLO missed, refund 10% of monthly fee |

**Burn rate:** How fast you're exhausting error budget.
- SLO = 99.9% → 0.1% error budget/month
- Burn rate > 2× for 6 hours → urgent intervention

## Practical Guidance

### Trade-off Analysis Framework

When comparing two design options, document:

```
## Option A: ...
### Pro: ... Con: ...
## Option B: ...
### Pro: ... Con: ...
## Decision: ...
## Rationale: ...
## Trade-offs accepted: ...
```

Always identify what is **not** gained. A decision log without trade-offs is a marketing document.

### Resilience Stack (outer → inner)

1. **Timeout** (always, first line)
2. **Retry** (with backoff, for transient failures)
3. **Circuit Breaker** (for persistent failures)
4. **Bulkhead** (isolation, prevent cascade)
5. **Graceful degradation** (return partial results)

## Guidelines

1. Every external call must have a timeout. No exceptions.
2. Retry only idempotent operations. Non-idempotent + retry = duplicates.
3. Circuit breaker recovery: HALF-OPEN after timeout, not based on request count.
4. SLO must be < SLA (leave error budget for unexpected failures).
5. Rate limiting at two levels: global (infrastructure) and per-key (application).
6. Trade-off documentation is mandatory for every architectural decision.

## Gotchas

1. **Retry storm**: Upstream is slow → all clients retry → upstream dies completely. Solution: jitter + circuit breaker + limit retry count to 3.
2. **No timeout on DB queries**: A long-running query holds a connection, blocking the pool. Solution: `context.WithTimeout` on every DB call.
3. **CAP misunderstanding — picking AP for consistency-critical data**: Banking transaction ledger needs CP, not AP. Full ACID matters.
4. **Five-nines chasing**: 99.999% uptime (5.26 min/year downtime) costs exponentially more than 99.9%. Measure business impact, not vanity metrics.
5. **Error budget not used**: The team doesn't deploy because they fear SLO violations. Error budget is meant to be spent — it tells you when to deploy and when to hold.

## Integration

- `architect-database-scaling` — CAP theorem in DB context, consistency for sharding
- `architect-microservices` — resilience patterns for service-to-service calls
- `architect-observability` — SLI measurement for SLO tracking