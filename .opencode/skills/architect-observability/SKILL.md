---
name: architect-observability
description: Monitoring, logging, metrics, distributed tracing, OpenTelemetry, alerting, SLO burn rate, dashboards
compatibility: opencode
---

# Observability

## When to Activate

Activate when the task involves:
- Designing monitoring and observability strategy
- Metrics definition (RED/USE, application metrics)
- Structured logging
- Distributed tracing with OpenTelemetry
- Sampling strategies (head-based, tail-based)
- Alerting rules and SLO burn rate
- Dashboard design
- Observability infrastructure (Prometheus, Grafana, Loki, Tempo, OpenTelemetry Collector)

Do not activate for: general resilience patterns (`architect-system-design`), or deployment monitoring. Adjacent: `architect-system-design` for SLO definition.

## Core Concepts

### The Three Pillars

| Pillar | What | Tool examples |
|--------|------|--------------|
| **Metrics** | Numerical measurements over time | Prometheus, Grafana, Datadog |
| **Logs** | Discrete events with timestamps | Loki, ELK, CloudWatch |
| **Traces** | Request flow across distributed services | Tempo, Jaeger, Zipkin |

**Observability** is the ability to understand the system's internal state from its external outputs. Having all three pillars is necessary but not sufficient — the *questions you can answer* matter more.

## Detailed Topics

### Metrics

**Service-level (RED):**
- **Rate** — requests per second
- **Errors** — error rate (5xx, 4xx, exceptions)
- **Duration** — latency distribution (p50, p90, p99)

**Resource-level (USE):**
- **Utilisation** — % of resource busy (CPU, memory, disk, connections)
- **Saturation** — queue length or pressure
- **Errors** — error count for the resource

**Application metrics** (custom to the business):
- Orders placed per second
- Queue depth per partition
- Cache hit ratio
- Database connection pool utilisation

### Logging

**Structured logging** — machine-parseable format (JSON), not free text.

```
# Bad
"User 42 logged in successfully"

# Good
{"level": "info", "ts": "2024-03-15T10:00:00Z", "msg": "login_success", "user_id": 42, "source_ip": "10.0.0.1"}
```

**Log levels:** DEBUG (development), INFO (normal ops), WARN (anomaly, not error), ERROR (failure, needs attention).

**What to log:**
- Request start and end (with duration and status)
- External calls (DB queries, API calls, message publish)
- Business decisions (approval, rejection, state change)
- Errors with stack traces (at the boundary, not in every layer)

**What NOT to log:**
- Passwords, tokens, secrets, PII

### Distributed Tracing

A trace is a tree of spans. Each span represents one unit of work (HTTP request, DB query, message processing).

```
Trace: POST /orders/checkout
├─ Span: /orders/checkout (service: gateway)
│  ├─ Span: validate_stock (service: inventory)  → 15ms
│  ├─ Span: process_payment (service: payment)   → 120ms
│  │  ├─ Span: db.query (service: payment)       → 5ms
│  │  └─ Span: http.post (service: payment-gw)   → 110ms
│  └─ Span: create_order (service: orders)       → 20ms
```

**Context propagation:** Trace ID and Span ID must be passed through all calls (HTTP headers, message metadata, gRPC metadata).

**W3C Trace-Context header:**
```
traceparent: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
              |  |              trace-id             |   span-id   |flag
```

**Sampling strategies:**
- **Head-based**: Decision at the root span. Simple, consistent. May miss rare errors.
- **Tail-based**: Decide to store after the trace completes. Better for error capture, expensive.
- **Probability sampling**: Fixed % of traces (1%, 10%). Good for volume analysis.
- **Rate limiting**: Max N traces per second. Protects storage costs.

### Alerting

**Alert on symptoms, not causes.**
- Good: "p99 latency > 300ms" (symptom — users are affected)
- Bad: "CPU > 80%" (cause — may or may not affect users)

**SLO burn rate:**
- Alert when error budget is burning faster than expected
- Example: 2% error budget/month. If 1 hour uses 5% of budget → urgent alert

**Severity:**
- **P0/Critical**: Users are actively impacted. Response: immediate.
- **P1/High**: Degradation, SLO close to violation. Response: within hours.
- **P2/Medium**: No immediate user impact. Response: next working day.

### OpenTelemetry

The industry standard for observability telemetry collection. Collect → Process → Export.

```
[App] → [OTel SDK] → [OTel Collector] → [Backend (Prometheus, Tempo, Loki)]
```

- **One SDK** for traces, metrics, logs
- **Vendor-neutral** — pluggable exporters
- **Auto-instrumentation** for popular frameworks (HTTP, gRPC, DB, message queues)

## Practical Guidance

### Observability Maturity Levels

1. **Level 1 — Black box**: Alerts on "is it up?" (health checks). No idea what's happening inside.
2. **Level 2 — Basic**: Essential RED metrics, structured logging, searchable logs.
3. **Level 3 — Observable**: Distributed tracing, custom business metrics, SLO-based alerting.
4. **Level 4 — Proactive**: Burn rate alerts, dashboards per team/service, root cause analysis from traces.

### Dashboard Design Principles

- **One screen, one question** — don't cram everything into one dashboard
- **Left to right, top to bottom flow** — start with the most important
- **Red/Yellow/Green thresholds** — visual urgency
- **Correlation** — latency + error rate + throughput on the same graph

## Guidelines

1. Every service exports RED metrics (Rate, Errors, Duration) by default
2. Structured JSON logging — no free text log messages
3. Trace context propagated through ALL inter-service calls
4. Sampling rate should satisfy: `traces_stored × avg_spans × cost_per_span ≤ storage_budget`
5. Alert on SLO burn rate, not static thresholds
6. One dashboard per service, one overview dashboard per domain
7. Log levels: ERROR needs immediate action, WARN needs investigation later

## Gotchas

1. **Logging in hot path**: `fmt.Sprintf` in every request handler adds latency. Use structured loggers with lazy evaluation.
2. **Metrics cardinality explosion**: Tagging every trace ID or user email as a metric label creates millions of time series and crashes Prometheus. Use a limited set of label values.
3. **No trace context**: Without w3c traceparent, traces are disconnected spans. Every service must propagate the trace context.
4. **Dashboard overload**: A dashboard with 50 graphs where no one knows what to look at. Start with 5 key graphs per service.
5. **Everything is P0**: Alert fatigue. Every alert that doesn't get a response should be either downgraded or deleted.
6. **Auto-instrumentation magic**: Auto-instrumentation captures standard spans but misses business logic. Add manual spans for critical business operations (DB writes, order validation).