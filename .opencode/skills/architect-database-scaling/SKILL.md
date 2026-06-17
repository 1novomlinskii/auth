---
name: architect-database-scaling
description: Database scaling strategies — sharding, replication, partitioning, read replicas, NoSQL vs SQL, migration patterns
compatibility: opencode
---

# Database Scaling

## When to Activate

Activate when the task involves:
- Scaling a database under growth
- Read replicas for read-heavy workloads
- Sharding for write scalability
- Table partitioning for data management
- NoSQL vs SQL trade-off analysis
- Online schema migration strategies
- Connection pooling configuration

Do not activate for: general application scaling (`architect-application-scaling`), individual query optimisation, or index tuning. Adjacent: `architect-system-design`, `architect-application-scaling`.

## Core Concepts

Database scaling follows a hierarchy — apply in order:

1. **Optimise queries + indexes** (cheapest, most effective)
2. **Connection pooling** (reduce connection overhead)
3. **Read replicas** (offload reads from primary)
4. **Caching** (see `architect-application-scaling`)
5. **Partitioning** (manage large tables)
6. **Sharding** (most expensive, last resort)

### The CAP Trade-off

When scaling databases, you trade consistency, availability, or partition tolerance:
- **CP**: Consistent + Partition-tolerant (traditional RDBMS primary)
- **AP**: Available + Partition-tolerant (DynamoDB, Cassandra)
- **CA**: Consistent + Available (doesn't survive partitions)

For most OLTP systems: prefer consistency (CP). Accept eventual consistency only when the business process allows stale reads.

## Detailed Topics

### Read Replicas

```
[Primary: writes] ──async/sync──→ [Replica 1: reads]
                                    [Replica 2: reads]
```

- **Synchronous**: COMMIT waits for replica ack. Data-safe, higher latency.
- **Asynchronous**: Primary doesn't wait. Lower latency, replica lag possible.

**Use for:** Reporting, analytics, read-heavy workloads, geographic distribution.
**Risk:** Stale reads (replica lag). Route critical reads to primary if consistency required.

### Partitioning (Table-Level)

Split a large table into smaller physical chunks using a partition key.

```sql
CREATE TABLE orders (
    id BIGINT, created_at DATE, ...
) PARTITION BY RANGE (created_at);

CREATE TABLE orders_2024_q1 PARTITION OF orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
```

**When:** Tables >100GB, time-series data, easy data retirement.
**Types:** Range, List, Hash.
**Benefit:** Query pruning, easier maintenance (drop partition vs DELETE), parallel scans.

### Sharding (Database-Level)

Distribute data across multiple independent database instances.

```
Shard 1: users [0-999]    Shard 2: users [1000-1999]
Shard 3: users [2000-2999] Shard 4: users [3000-3999]
```

**Shard key selection is the critical decision:**
- Good: `user_id`, `customer_id`, `tenant_id` (natural, even distribution)
- Bad: `created_at` (hot spot on latest shard), `status` (skewed distribution)

**Strategies:**
- **Hash-based**: Consistent hashing (rebalancing friendly)
- **Range-based**: Simple, but hot spots
- **Directory-based**: Lookup table maps key → shard (flexible, extra hop)

**Cost:** Cross-shard queries impossible, JOINs across shards require application orchestration, rebalancing is painful.

### NoSQL vs SQL Decision

| Factor | Choose SQL | Choose NoSQL |
|--------|-----------|-------------|
| Data shape | Structured, relational, ACID | Flexible schema, nested documents |
| Operations | Complex JOINs, aggregations | Simple key lookups, denormalised |
| Consistency | Strong consistency required | Eventual consistency acceptable |
| Scale | Read replicas + partitioning sufficient | Global scale, multi-region writes |
| Examples | PostgreSQL, MySQL, SQLite | DynamoDB, MongoDB, Cassandra |

**Common trap:** Choosing NoSQL "for scale" when the data is inherently relational. You will reimplement JOINs in application code.

### Migration Strategies

| Strategy | Downtime | Risk | When |
|----------|----------|------|------|
| **Expand-Contract** (dual writes) | Zero | Medium | Any production migration |
| **Online schema change** (pt-online-schema-change, gh-ost) | Near-zero | Low | DDL on large tables |
| **Blue/Green DB** | Cutover window | Low | Re-platforming to new DB |
| **Backup-restore** | Full downtime | Low | Small databases, maintenance windows |

**Expand-Contract pattern:**
1. Add new column/table (allow both old and new)
2. Dual-write to both, backfill data
3. Migrate reads to new structure
4. Remove old structure

## Practical Guidance

### Scaling Ladder

1. Index audit + query optimisation → 2. Connection pooling → 3. Read replicas → 4. Application cache → 5. Partitioning → 6. Vertical scale → 7. Sharding (last)

### Sharding Decision Checklist

- [ ] Write throughput exceeds single node capacity after vertical scaling?
- [ ] Data size exceeds practical single-node storage?
- [ ] Can you choose a good shard key (evenly distributed, natural to queries)?
- [ ] Can you avoid cross-shard queries at the application level?
- [ ] Team has operational capacity for sharding complexity?

If any answer is "no," explore partitioning or read replicas first.

## Guidelines

1. Premature sharding is the #1 database scaling mistake. Exhaust simpler options first.
2. Read replicas are low-hanging fruit — add them as soon as read load impacts writes.
3. Choose partition key based on query pattern, not just data distribution.
4. Connection pooling (PgBouncer / pgbouncer-sidecar) prevents connection exhaustion.
5. Document shard key and routing logic explicitly — it's invisible in the schema.
6. Test replica lag under load before relying on read replicas for consistency-critical paths.

## Gotchas

1. **Cross-shard transactions**: ACID across shards is impossible. Use Saga pattern (`architect-microservices`) or accept that some operations are best-effort.
2. **Shard rebalancing**: Adding or removing shards requires data movement. Consistent hashing minimises but doesn't eliminate this. Plan for rebalancing windows.
3. **Sequence ID on sharded DB**: Auto-increment breaks across shards. Use UUID v7 (time-ordered), Snowflake IDs, or a central sequence service.
4. **Replica lag under load**: The replica may lag seconds behind during bursts. Build read staleness tolerance into the application or route critical reads to primary.
5. **Connection pool saturation**: Each application instance opens N connections. 100 instances × 10 connections = 1000 connections. Tune pool size: pool ≈ (connections per instance × instances) + spare.

## Integration

- `architect-application-scaling` — caching as a cheaper alternative to DB scaling
- `architect-microservices` — database-per-service, saga for distributed transactions
- `architect-system-design` — CAP theorem, consistency models