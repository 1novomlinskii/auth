---
name: dba-sql-optimization
description: SQL query optimization — rewriting inefficient queries, CTE vs subqueries vs window functions, JOIN strategies, anti-patterns (N+1, type casting, function calls in WHERE)
license: MIT
compatibility: opencode
---

# SQL Optimization

## When to Activate

Activate when the task involves:
- Rewriting slow SQL queries
- Choosing between CTE (WITH), subqueries, and window functions
- Fixing N+1 query problems
- Optimizing JOINs, aggregations, and GROUP BY
- Detecting anti-patterns: implicit type casting, function calls in WHERE, unnecessary SELECT DISTINCT
- Converting correlated subqueries to JOINs or window functions

Do not activate for: PostgreSQL config tuning, index design, PostGIS, or SQL formatting.

## CTE vs Subquery vs Window Function

| Construct | Use Case |
|-----------|----------|
| CTE (WITH) | Readability, recursive queries (`WITH RECURSIVE`), multiple references to same subquery |
| Subquery | Simple filtering, exists checks, single-use |
| Window function | Running totals, row_number, lag/lead — most efficient for per-group calculations |
| LATERAL | Reference previous FROM items, call set-returning functions per row |

### CTE Pitfall

CTEs are **optimization fences** in PostgreSQL — planner materializes them by default.
Use `WITH ... AS MATERIALIZED` or `WITH ... AS NOT MATERIALIZED` to control behavior.
Prefer subqueries or `LATERAL` for better plan optimization when CTE is used once.

## Anti-Patterns

### 1. N+1 Queries

```sql
-- BAD: 1 query for users + N queries for their orders
SELECT *
  FROM users;                                    -- 1 query
-- then for each user: SELECT * FROM orders WHERE user_id = ?  -- N queries

-- GOOD: single query (or two with JOIN)
SELECT u.*,
       o.*
  FROM users AS u
       LEFT JOIN orders AS o
       ON o.user_id = u.id;
```

### 2. Function Calls in WHERE

```sql
-- BAD: function call prevents index usage
SELECT *
  FROM orders
 WHERE DATE(created_at) = '2024-01-01';

-- GOOD: range query uses index
SELECT *
  FROM orders
 WHERE created_at >= '2024-01-01'
   AND created_at < '2024-01-02';
```

### 3. Implicit Type Casting

```sql
-- BAD: text column compared to integer → seq scan
SELECT *
  FROM users
 WHERE id = '42';  -- id is INTEGER

-- GOOD: match types
SELECT *
  FROM users
 WHERE id = 42;
```

### 4. Unnecessary DISTINCT

```sql
-- BAD: DISTINCT just to deduplicate bad JOIN
SELECT DISTINCT u.*
  FROM users AS u
       JOIN orders AS o
       ON o.user_id = u.id;

-- GOOD: no duplicates in well-formed JOIN, or use EXISTS
SELECT u.*
  FROM users AS u
 WHERE EXISTS (SELECT 1
                 FROM orders AS o
                WHERE o.user_id = u.id);
```

### 5. SELECT *

```sql
-- BAD: read + transfer unused columns, prevents index-only scans
SELECT *
  FROM orders;

-- GOOD: list only needed columns
SELECT id,
       status,
       total
  FROM orders;
```

## JOIN Optimization

- Filter before joining: apply WHERE conditions as early as possible
- Use `EXISTS` vs `IN`: `EXISTS` is typically faster with large subquery result sets
- Prefer `LEFT JOIN` only when NULL matches are needed — regular `JOIN` is faster
- `LATERAL` for complex per-row computations

## Aggregation Tips

- Filter with `WHERE` before aggregation, not `HAVING` after
- `COUNT(*)` vs `COUNT(col)`: `COUNT(*)` counts all rows, `COUNT(col)` skips NULLs
- `EXPLAIN ANALYZE` to check aggregation strategy (HashAggregate vs GroupAggregate)

## Integration

- dba-postgresql-performance — EXPLAIN analysis of rewritten queries
- dba-schema-design — schema-level decisions that impact query performance
- dba-sql-guide — formatting optimized SQL