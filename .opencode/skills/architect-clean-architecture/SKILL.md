---
name: architect-clean-architecture
description: Clean Architecture, Hexagonal (Ports & Adapters), Onion, DIP, слои и границы
compatibility: opencode
---

# Clean Architecture

## When to Activate

Activate when the task involves:
- Defining or reviewing architectural layers (domain, use cases, infrastructure)
- Enforcing Dependency Inversion Principle
- Implementing Hexagonal Architecture (Ports & Adapters)
- Onion Architecture layering
- Bounded Context mapping
- Detecting leaky abstractions, dependency cycles, or anemic domain models

Do not activate for: database schema design alone, individual class/function design, or API endpoint naming. Adjacent skills: `architect-design-patterns`, `architect-api-design`.

## Core Concepts

Clean Architecture organises code into concentric layers where dependencies point **inward** — toward the domain. Outer layers depend on inner layers, never the reverse.

The three dominant formulations are equivalent in their core rule:

| Form | Key mechanism |
|------|--------------|
| **Clean Architecture** (R. Martin) | Entities → Use Cases → Interface Adapters → Frameworks |
| **Hexagonal** (A. Cockburn) | Ports (interfaces) at the boundary, Adapters implement them |
| **Onion** (J. Palermo) | Domain Model → Domain Services → Application → Infrastructure |

All three enforce the same principle: **domain code has zero knowledge of frameworks, databases, or UI**.

### Dependency Rule

Source code dependencies can only point inward. Nothing in an inner circle can know about something in an outer circle.

```
Frameworks & Drivers → Interface Adapters → Use Cases → Entities
        (outer)                                      (inner)
```

When outer needs to communicate inward, use **Dependency Inversion** — define an interface in the inner layer, implement it in the outer layer.

## Detailed Topics

### Layers

**Entities (Domain):**
- Core business objects and rules
- Pure logic, no I/O, no framework annotations
- Most stable layer — changes rarely

**Use Cases (Application):**
- Orchestrate business workflows
- Depend on Entities, define interfaces (Ports) for external I/O
- One use case = one business operation

**Interface Adapters:**
- Controllers, presenters, gateways
- Convert between use case formats and external formats
- Implements Ports defined by Use Cases

**Frameworks & Drivers:**
- DB drivers, web framework, UI, message brokers
- Most volatile layer — changes often
- Kept at the edges

### Ports & Adapters (Hexagonal)

- **Port** — interface defined inside the application boundary (e.g., `UserRepository`)
- **Adapter** — implementation outside the boundary (e.g., `PostgresUserRepository`)
- Application core never imports infrastructure packages directly

```
[Adapter] → implements → [Port(interface)] ← [Application Core]
```

### Dependency Injection

The mechanism that wires layers together. Adapters are injected into the core at composition root (application entry point).

```
db := postgres.NewConnection(cfg)
userRepo := postgres.NewUserRepository(db)
svc := service.NewUserService(userRepo)
```

## Practical Guidance

### Decision flow for layer placement

1. Does this code express a business rule? → Entity
2. Does it orchestrate a business operation? → Use Case
3. Does it convert between formats or protocols? → Interface Adapter
4. Does it talk to a specific technology? → Framework/Driver

### When NOT to use Clean Architecture

- **CRUD-only app** (simple forms over data) — unnecessary overhead. Use layered architecture without full hexagonal.
- **Prototype / MVP** — start simple, refactor into clean architecture only when warranted.
- **Tiny service** (<500 lines total) — folder conventions are sufficient.

## Guidelines

1. Domain layer imports zero external packages (stdlib only)
2. Use Cases define interfaces for repositories, not implementations
3. No database annotations (ORM tags) in domain entities
4. Composition root creates and wires all dependencies
5. Tests for Use Cases mock Ports, not Adapters
6. Business logic is never in controllers, handlers, or DB repositories

## Gotchas

1. **Leaky abstraction**: Using `context.Context` from a web framework in domain code — domain should only know about its own abstractions. Pass only the data the domain needs.
2. **Anemic domain model**: All logic in use cases/services, entities are plain structs with getters/setters. The domain should encapsulate business rules, not just data.
3. **Dependency cycle**: Package A imports B, B imports A. In Go: move the interface to a third package both depend on, or place it in the consumer's package.
4. **Framework lock-in**: ORM-annotated entities, framework-specific base classes in domain code. The domain should be framework-agnostic by definition.

## Integration

- `architect-design-patterns` — Repository, Factory, Adapter patterns used within clean architecture
- `architect-api-design` — controllers/adapters for web layer
- `architect-microservices` — bounded contexts align with clean architecture boundaries