---
name: dba-postgresql-performance
description: PostgreSQL performance diagnostics — EXPLAIN ANALYZE, index strategies (B-tree, GIN, GiST, BRIN, partial, covering), configuration tuning (work_mem, shared_buffers), pg_stat_statements, pg_stat_activity
license: MIT
compatibility: opencode
---

# PostgreSQL Performance

## When to Activate

Activate when the task involves:
- Analyzing slow queries with EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
- Recommending index types (B-tree, GIN, GiST, BRIN, partial, covering, composite)
- Tuning PostgreSQL configuration (work_mem, shared_buffers, effective_cache_size)
- Investigating locks, long-running queries, or high CPU via `pg_stat_activity`
- Identifying Seq Scan, Nested Loop, or Hash Join performance issues
- Reading wait events from `pg_stat_activity`

Do not activate for: SQL formatting (use dba-sql-guide), schema design, or PostGIS spatial queries.

## EXPLAIN Analysis

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT ...;
```

Key signals:
- **Seq Scan on large table** → missing index
- **Nested Loop with many rows** → consider Hash Join or Merge Join
- **Sort Method: external merge** → increase `work_mem`
- **Buffers: shared hit** → cached; **shared read** → disk I/O
- **Rows Removed by Filter** high → index needed

## Index Strategy

| Index Type | Use Case |
|-----------|----------|
| B-tree (default) | Equality and range queries (`=`, `<`, `>`, `BETWEEN`, `LIKE 'prefix%'`) |
| GIN | JSONB (`@>`, `?`), full-text search (`tsvector`), arrays |
| GiST | Full-text search, range types, spatial (PostGIS), fuzzy matching |
| BRIN | Large tables with naturally ordered data (time-series, logs) |
| Partial | `CREATE INDEX ... WHERE status = 'active'` — saves space |
| Covering | `CREATE INDEX ... INCLUDE (col1, col2)` — index-only scans |
| Composite | Column order matters: equality first, then range |
| Hash | Equality only, rarely needed (B-tree handles equality too) |

### Index Guidelines

- Index columns used in `WHERE`, `JOIN`, `ORDER BY`, `GROUP BY`
- B-tree default for most cases
- Column order in composite index: high selectivity → low selectivity
- Avoid over-indexing — each index slows `INSERT`/`UPDATE`/`DELETE`
- `CREATE INDEX CONCURRENTLY` — non-blocking creation
- `REINDEX INDEX` — rebuild bloated indexes
- `pg_stat_user_indexes` — unused indexes

## Configuration Tuning

| Parameter | Recommendation |
|-----------|---------------|
| `shared_buffers` | 25% of RAM (but not > 8GB without testing) |
| `work_mem` | 4-32MB per operation (too high → OOM, too low → disk sorts) |
| `effective_cache_size` | 50-75% of RAM (helps query planner) |
| `random_page_cost` | 1.1 for SSD, 4.0 for HDD |
| `max_connections` | Lower = better. Use connection pooler (PgBouncer) |

## Monitoring

```sql
-- Top queries by total time
  SELECT queryid, calls, total_exec_time, mean_exec_time, rows
    FROM pg_stat_statements
ORDER BY total_exec_time DESC
   LIMIT 10;

-- Active queries with locks
  SELECT pid, wait_event_type, wait_event, state, query
    FROM pg_stat_activity
   WHERE state != 'idle'
ORDER BY wait_event_type NULLS LAST;

-- Index usage
  SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read
    FROM pg_stat_user_indexes
ORDER BY idx_scan;
```

## Integration

- dba-sql-optimization — rewriting queries after EXPLAIN analysis
- dba-schema-design — designing indexes alongside table schemas
- dba-postgis — spatial index types (GiST, SP-GiST)
- dba-sql-guide — formatting SQL output