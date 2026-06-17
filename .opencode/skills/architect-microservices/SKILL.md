---
name: architect-microservices
description: Microservice architecture design — bounded contexts, service decomposition, sync/async communication, Saga, event-driven, deployment strategies
compatibility: opencode
---

# Microservices

## When to Activate

Activate when the task involves:
- Decomposing a monolith into microservices
- Defining service boundaries and ownership
- Inter-service communication (sync gRPC/HTTP vs async events)
- Saga pattern for distributed transactions
- Event-driven architecture design
- Anti-corruption Layer for legacy integration
- Service mesh topology
- Deployment strategies (blue/green, canary, rolling)

Do not activate for: clean architecture layers within a single service (`architect-clean-architecture`), general scaling (`architect-application-scaling`), or API design (`architect-api-design`). Microservices should be the *result* of architectural analysis, not the starting assumption.

## Core Concepts

### Decomposition Is the Hardest Problem

The single most important decision: **where to draw service boundaries**. Get this wrong and all messaging infrastructure, CI/CD, and team structure just accelerates a bad design.

**Good boundaries come from business capabilities and bounded contexts (DDD), not technical layers.**

### Monolith-First

Do NOT start with microservices. Start with a modular monolith with clean internal boundaries. Extract services only when:
- Team size requires independent deploys (>2-3 teams)
- Different scalability requirements for different modules
- Different tech stacks justified per module
- Deployment velocity blocked by monolithic release cadence

"If you can't build a monolith, you can't build microservices." — S. Newman

## Detailed Topics

### Bounded Context (DDD)

A bounded context is a logical boundary around a domain model. Each microservice owns exactly one bounded context.

| Context | Model | Example |
|---------|-------|---------|
| **Sales** | Order, Quote, Customer | Accepting orders |
| **Shipping** | Package, Route, Carrier | Fulfilling orders |
| **Billing** | Invoice, Payment, Ledger | Charging customers |

The same entity may appear in multiple contexts with different attributes (a Customer in Sales has a Cart; in Billing, has a BalanceDue).

### Communication Patterns

**Synchronous (request/response):**
```
Service A ──HTTP/gRPC──→ Service B
```
- Pros: Simple, familiar, real-time response
- Cons: Coupling (A depends on B's availability), cascading failures
- Best for: Queries, commands that need immediate confirmation

**Asynchronous (events/messaging):**
```
Service A ──publish──→ [Broker] ──consume──→ Service B
```
- Pros: Loose coupling, buffering, fan-out
- Cons: Eventual consistency, debugging complexity
- Best for: Side effects, cross-service notifications, audit

**Event-Driven decision:**
- Event type: Domain Event (something happened) vs Command (do something)
- Delivery guarantee: At-most-once, at-least-once, exactly-once
- Schema evolution: Schema registry (Avro, Protobuf) vs shared type library vs self-describing JSON

### Saga Pattern

A saga coordinates a distributed transaction across services using a sequence of local transactions with compensating actions on failure.

**Choreography (no coordinator):**
```
OrderSvc: create order → emit "OrderCreated"
PaymentSvc: listen "OrderCreated" → charge → emit "PaymentCharged" | emit "PaymentFailed"
ShippingSvc: listen "PaymentCharged" → ship → emit "OrderShipped"
On failure: each service emits compensating events
```
Pros: Simple, decentralised. Cons: Event flow is implicit, hard to trace.

**Orchestration (coordinator):**
```
Orchestrator → OrderSvc: create()
Orchestrator → PaymentSvc: charge()
  on fail → OrderSvc: compensate()
Orchestrator → ShippingSvc: ship()
```
Pros: Explicit flow, visible in code. Cons: Orchestrator is a single point of complexity.

### Anti-corruption Layer (ACL)

Translate between a legacy system's model and your modern model. Usually implemented as a separate adapter service.

```
[Modern Service (Customer)] ──ACL──→ [Legacy CRM (Kunde)]
   │                                          │
   │ Translates: name → vorname+nachname      │
   │            email → emailAddr             │
   │            id → customerNbr              │
```

## Practical Guidance

### Service Decomposition Checklist

- [ ] Each service maps to one bounded context
- [ ] Teams are aligned to service boundaries (Conway's Law)
- [ ] Services can be deployed independently
- [ ] Database per service (shared DB is an anti-pattern)
- [ ] Contracts are versioned and backward-compatible
- [ ] Synchronous calls have timeouts + circuit breakers
- [ ] Asynchronous events have schema registry + schema evolution

### When NOT to Extract

- Module has no independent scalability requirements
- Module is deeply entangled (shared DB, transactions spanning boundaries)
- Team can't operate the service independently
- API contract changes too frequently (consumers become tightly coupled)
- Data consistency requirements can't tolerate eventual consistency

## Guidelines

1. Database-per-service is non-negotiable. Shared DB between services is a distributed monolith.
2. Event contracts are harder to change than API contracts. Invest in schema evolution strategy early.
3. Service mesh (Istio, Linkerd) solves cross-cutting concerns (mTLS, observability, retries) — prefer it to per-service boilerplate.
4. Version APIs (URL path or header), not just events.
5. Service template (cookiecutter, project scaffold) reduces decision fatigue for new services.
6. Integration tests against real dependencies, not mocks — the network is the weakest link.

## Gotchas

1. **Distributed monolith**: Multiple services that must be deployed together, share a DB, or have circular dependencies. It's worse than a monolithic — costs of distribution without benefits.
2. **Chatty services**: Service A calls B, then calls B again, then calls C, with 10ms latency each → 30ms just in network overhead. Batch requests into one call when possible.
3. **Saga without compensation**: If a saga step fails but no compensation is implemented, the system is in an inconsistent state. Every saga step must have a compensating action.
4. **Shared type library**: Every service depends on a shared JAR/NPM package for models. Changes require coordinated deploys across services. Use schema registry instead.
5. **Too granular**: A service with one endpoint and 100 lines of code. The overhead of service infrastructure (CI/CD, monitoring, team) is never justified for this granularity. Start coarse, split only when needed.

## Integration

- `architect-clean-architecture` — each service follows clean architecture internally
- `architect-api-design` — service-to-service API contracts
- `architect-observability` — distributed tracing across services
- `architect-application-scaling` — independent scaling of services