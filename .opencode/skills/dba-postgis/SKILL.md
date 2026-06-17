---
name: dba-postgis
description: PostGIS spatial queries and optimization — GiST/SP-GiST indexes, ST_DWithin, ST_Intersects, SRID, geometry vs geography, spatial query tuning
license: MIT
compatibility: opencode
---

# PostGIS / Spatial Data

## When to Activate

Activate when the task involves:
- Writing spatial SQL queries with PostGIS functions
- Creating spatial indexes (GiST, SP-GiST)
- Choosing between geometry and geography types
- Working with SRID and coordinate transformations (ST_Transform)
- Optimizing spatial queries (ST_DWithin, ST_Intersects, ST_Distance)
- Spatial JOINs between tables with geometry columns

Do not activate for: general PostgreSQL performance, schema design, or SQL formatting.

## Spatial Indexes

```sql
CREATE INDEX idx_locations_geom ON locations USING GIST (geom);
CREATE INDEX idx_locations_geom_spgist ON locations USING SPGIST (geom);
```

- **GiST** — general-purpose, most common for spatial (point, polygon, line)
- **SP-GiST** — better for point clouds with non-overlapping partitions
- Always index the geometry column used in spatial WHERE conditions
- `CLUSTER` on spatial index for better locality

## Geometry vs Geography

| Type | Use Case |
|------|----------|
| `GEOMETRY` | Cartesian plane — local projections, equal-area, small regions |
| `GEOGRAPHY` | Earth sphere — lat/lon, great-circle distances, global coverage |

- Use `GEOGRAPHY` for GPS coordinates (lat/lon) — automatic great-circle calculations
- Use `GEOMETRY` for projected coordinate systems (UTM, Mercator) — faster, planar math
- Transform between them: `geom::geography` or `ST_Transform(geom, target_srid)`

## SRID (Spatial Reference ID)

- 4326 — WGS84 (GPS, lat/lon) — most common
- 3857 — Web Mercator (Google Maps, OpenStreetMap tiles)
- Use `ST_SetSRID(geom, 4326)` to assign, `ST_Transform(geom, 3857)` to reproject
- Always specify SRID explicitly — don't rely on defaults

## Key Functions

### Proximity

```sql
-- Find points within 1km radius (geography = meters)
SELECT *
  FROM locations
 WHERE ST_DWithin(geom::geography, target::geography, 1000);

-- Distance between two points (geography returns meters)
SELECT ST_Distance(a.geom::geography, b.geom::geography) AS dist_meters
  FROM locations AS a,
       locations AS b;
```

### Intersection

```sql
-- Find all points inside a polygon
SELECT *
  FROM points,
       polygons
 WHERE ST_Intersects(points.geom, polygons.geom);

-- Geofencing: find if current location is inside zone
SELECT ST_Within(current_loc::geography, zone::geography);
```

### Transformation

```sql
-- Transform GPS to UTM for planar distance calculation
SELECT ST_Transform(geom, 32633)
  FROM locations;  -- UTM zone 33N
```

## Optimization Tips

1. Always use `ST_DWithin` (index-enabled) over `ST_Distance < threshold` (no index)
2. Add `&&` bounding box operator — it's automatically used with GiST but verify with EXPLAIN
3. Use `ST_Simplify` for rendering — reduces geometry complexity
4. `ST_Subdivide` — split large polygons for faster spatial JOINs
5. `VACUUM ANALYZE` after bulk spatial inserts to update stats
6. Filter with non-spatial conditions first before spatial to reduce index lookups

## Integration

- dba-postgresql-performance — EXPLAIN analysis of spatial queries
- dba-schema-design — spatial column types and SRID choices
- dba-sql-guide — formatting spatial SQL