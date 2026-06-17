---
name: architect-design-patterns
description: GoF patterns (creational, structural, behavioral) and enterprise patterns (Repository, Unit of Work, Saga, Outbox, Domain Event)
compatibility: opencode
---

# Design Patterns

## When to Activate

Activate when the task involves:
- Selecting a pattern for a recurring design problem
- Reviewing pattern usage in existing code
- Replacing ad-hoc solutions with well-known patterns
- Choosing between similar patterns (Strategy vs State, Factory vs Builder)

Do not activate for: defining project structure or API design. Adjacent skills: `architect-clean-architecture` for layer organisation, `architect-api-design` for API patterns.

## Core Concepts

A pattern is a **proven solution to a recurring problem in a context**. Patterns are not goals — apply them only when the problem actually exists.

### Decision Principles

- **Patterns are vocabulary, not prescription** — naming a pattern helps communication more than the implementation
- **YAGNI over patterns** — if you don't have the problem, don't apply the pattern
- **Prefer language-idiomatic solutions** — Go channels replace Observer, Python decorators replace Proxy
- **Wrong pattern is worse than no pattern** — forces wrong abstractions

## Detailed Topics

### Creational — Object Creation

| Pattern | Problem | When to apply |
|---------|---------|--------------|
| **Factory Method** | Client needs a product, defers instantiation to subclasses | Multiple implementations of the same interface, determined at runtime |
| **Abstract Factory** | Families of related products must be created together | UI themes, DB providers, cross-platform components |
| **Builder** | Complex object with many optional parts | SQL query builders, HTTP request builders, config objects |
| **Singleton** | Exactly one instance must exist | Logging, configuration, connection pools (but pass as dependency, don't use global) |
| **Prototype** | Objects are expensive to create, copy is cheaper | Caching, document cloning |

**Anti-pattern:** Random AbstractFactoryFactory. One product type with one implementation needs no factory. Add a factory when the second implementation appears.

### Structural — Composition

| Pattern | Problem | When to apply |
|---------|---------|--------------|
| **Adapter** | Incompatible interfaces need to work together | Wrapping a third-party library behind a project interface |
| **Facade** | Complex subsystem needs a simple entry point | Simplifying a convoluted API, legacy system wrapper |
| **Proxy** | Control access to an object | Lazy loading, access control, logging wrapper |
| **Decorator** | Add behaviour without modifying the object | Middleware chains, logging/metrics wrappers |
| **Composite** | Treat individual and composition uniformly | Tree structures, UI component trees |

**Anti-pattern:** Adapter proliferation — every external library wrapped in an interface. Wrap only at architectural boundaries (ports).

### Behavioral — Interaction

| Pattern | Problem | When to apply |
|---------|---------|--------------|
| **Strategy** | Different algorithms for the same task, interchangeable | Payment processors, sorting strategies, validation rules |
| **Observer (Pub/Sub)** | Objects need notification of state changes | Event systems, message queues, UI event handling |
| **Command** | Encapsulate a request as an object | Undo/redo, task queues, transaction scripting |
| **Template Method** | Skeleton of algorithm, steps vary | Batch processing frameworks, report generation |
| **State** | Object behaviour changes with internal state | Workflow engines, connection state machine |

**Anti-pattern:** Over-engineering with Strategy — a single `if/switch` with 3 branches does NOT need 3 strategy classes. Add Strategy when there are >5 branches or branches come from plugins.

### Enterprise Patterns

| Pattern | Problem | When to apply |
|---------|---------|--------------|
| **Repository** | Abstract data access behind collection-like interface | All data access should go through repositories (domain doesn't know DB) |
| **Unit of Work** | Track changes and commit atomically | Multiple repository operations in one transaction |
| **Domain Event** | Side effects on domain state changes | Loose coupling between aggregates, audit logs, notifications |
| **Saga** | Distributed transaction across services | Microservices coordination (see `architect-microservices`) |
| **Outbox** | Reliable message publishing with DB transaction | Exactly-once message delivery, event-driven systems |
| **CQRS** | Separate read and write models | Different read/write loads, complex queries (see `architect-application-scaling`) |

## Practical Guidance

### Pattern Selection Flow

1. Do I actually have this problem? If no → stop.
2. Is there a simpler language-idiomatic solution? If yes → use it.
3. Does the pattern map to a known concept in the team vocabulary? If yes → document the intent.
4. Can I defer the pattern until I see the second variation? If yes → defer.

### Golden Rules

- **Factory**: add when the second concrete type appears
- **Interface**: extract when the second implementation appears
- **Strategy**: add when the fifth branch appears or strategies come from plugins
- **Repository**: always — it's not a pattern, it's a boundary

## Guidelines

1. Name patterns in code only when the structure exactly matches the canonical form. Mismatched naming confuses readers.
2. Document the *intent* — why this pattern here, what problem it solves.
3. Prefer composition over inheritance (Go doesn't have inheritance — this is built-in).
4. Singletons must still be passed as dependencies, never accessed globally.
5. Patterns in tests are welcome — Test Builder, Object Mother, Fake.

## Gotchas

1. **Pattern-of-the-week architecture**: Applying every known pattern to a simple CRUD app. Each pattern adds complexity — justify its cost.
2. **Framework-enforced patterns**: Some frameworks force a pattern (e.g., Java EE forces Service Locator). Don't fight it, but isolate it behind ports.
3. **Pattern without the problem**: "We should use CQRS because it's modern." Measure the actual read/write ratio and query complexity first.
4. **Singleton as global state**: In Go, `sync.Once` in a package-level `var` makes testing painful. Use dependency injection even for "singletons."

## Integration

- `architect-clean-architecture` — patterns implement layer boundaries (Repository, Adapter)
- `architect-microservices` — Saga, Outbox, Domain Event in distributed context
- `architect-api-design` — REST patterns (Repository Controller, DTO Assembler)