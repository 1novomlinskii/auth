---
name: dba-schema-design
description: PostgreSQL schema design — normalization, data types, primary/foreign keys, constraints, partitioning, migration strategies
license: MIT
compatibility: opencode
---

# Schema Design

## When to Activate

Activate when the task involves:
- Designing new database tables and relationships
- Choosing PostgreSQL data types (UUID, BIGSERIAL, JSONB, etc.)
- Defining primary keys, foreign keys, and constraints
- Table partitioning (range, list, hash)
- Normalization vs denormalization decisions
- Planning database migrations

Do not activate for: query optimization, index tuning, PostGIS, or SQL formatting.

## Data Types

| Type | When to Use |
|------|------------|
| `UUID` | Primary keys in distributed systems, microservices, public IDs |
| `BIGSERIAL` / `BIGINT` | Auto-increment PKs for single-server apps |
| `TEXT` | Variable-length strings (no length limit, no performance penalty) |
| `VARCHAR(n)` | Only when length limit is a business requirement |
| `TIMESTAMPTZ` | Always for timestamps (stores UTC, converts to session TZ) |
| `DATE` | Date-only values (no time component) |
| `NUMERIC(p,s)` | Exact decimal arithmetic (money, calculations) |
| `JSONB` | Flexible schemas, queryable via `@>`, `->`, `->>` |
| `BOOLEAN` | TRUE / FALSE / NULL (avoid integer flags) |
| `INTEGER` / `BIGINT` | Whole numbers |
| `BYTEA` | Binary data (rare — store files externally) |

## Primary Keys

- Use `UUID` or `BIGSERIAL` — avoid natural keys (email, SSN) as PK
- Prefer UUID for distributed systems (no sequential guesses)
- `BIGSERIAL` for high-write single-server apps (smaller index)
- Always define PK explicitly: `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`

## Foreign Keys

```sql
CREATE TABLE orders (
    PRIMARY KEY (id),
    id         UUID          NOT NULL DEFAULT gen_random_uuid(),
    user_id    UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);
```

- Always use `ON DELETE` (CASCADE, SET NULL, RESTRICT)
- FK constraints enforce referential integrity — don't skip them
- Add indexes on FK columns (needed for JOIN performance)

## Constraints

```sql
CREATE TABLE products (
    PRIMARY KEY (id),
    id     UUID           NOT NULL DEFAULT gen_random_uuid(),
    sku    TEXT           NOT NULL UNIQUE,
    price  NUMERIC(10,2)  NOT NULL CHECK (price > 0),
    status TEXT           NOT NULL DEFAULT 'active'
                          CHECK (status IN ('active', 'disabled', 'archived'))
);
```

- `NOT NULL` by default for required fields
- `CHECK` constraints for business rules (cheaper than app-layer validation)
- `UNIQUE` for alternate keys (sku, slug)
- `EXCLUDE` for complex uniqueness (overlapping date ranges)

## Partitioning

```sql
CREATE TABLE logs (
    id         BIGSERIAL    NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL,
    level      TEXT         NOT NULL,
    message    TEXT
) PARTITION BY RANGE (created_at);

-- Create partitions
CREATE TABLE logs_2024_q1 PARTITION OF logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
CREATE TABLE logs_2024_q2 PARTITION OF logs
    FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');
```

Use cases: time-series (logs, events), large tables that need periodic cleanup.
Partition types: `RANGE` (dates), `LIST` (regions), `HASH` (sharding).

## Normalization vs Denormalization

- Normalize by default (3NF) — reduces redundancy, maintains consistency
- Denormalize only when measured performance requires it
- Common denormalization: pre-computed totals, materialized views, JSONB for flexible fields

## Integration

- dba-postgresql-performance — index recommendations alongside schema
- dba-sql-optimization — query patterns that depend on schema design
- dba-sql-guide — formatting CREATE TABLE and DDL