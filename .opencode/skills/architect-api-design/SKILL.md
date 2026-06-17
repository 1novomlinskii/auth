---
name: architect-api-design
description: API design for REST, GraphQL, gRPC — resource modelling, versioning, pagination, error handling, protocol selection
compatibility: opencode
---

# API Design

## When to Activate

Activate when the task involves:
- Designing a new API (REST, GraphQL, or gRPC)
- Choosing between API protocols
- API versioning strategy
- Error model design
- Pagination, filtering, sorting patterns
- API documentation (OpenAPI, Protobuf)
- Reviewing an existing API for consistency

Do not activate for: internal function design or database queries. Adjacent: `architect-microservices` for inter-service API contracts, `architect-observability` for API monitoring.

## Core Concepts

### Protocol Selection Framework

| Criterion | REST | GraphQL | gRPC |
|-----------|------|---------|------|
| Client diversity | Many (web, mobile, third-party) | Mobile/web (reduced payload) | Server-to-server |
| Data fetching | Fixed response shape | Client-defined | Fixed per RPC |
| Latency sensitivity | Moderate | Moderate | Low (binary, HTTP/2) |
| Streaming | SSE only | Subscriptions | Bidirectional streams |
| Schema strictness | Loose (OpenAPI docs) | Strong (Schema) | Strong (Protobuf) |
| Tooling ecosystem | Mature (curl, Postman) | Growing (Apollo, Relay) | Mature (grpcurl, protoc) |
| Caching | Natural (HTTP cache) | Challenging | Not built-in |

**Default choice:** REST. Move to GraphQL when clients need flexible queries. Move to gRPC for high-performance internal service-to-service communication.

## Detailed Topics

### RESTful API Design

**Resource naming:**
```
GET    /users              — list users
POST   /users              — create user
GET    /users/{id}         — get user
PATCH  /users/{id}         — partial update
DELETE /users/{id}         — delete user
GET    /users/{id}/orders  — sub-resource
```

**Plural nouns, not verbs.** `/getUser` is an RPC, not REST.

**HTTP methods by semantics:**
- `GET` — read, safe, idempotent, cacheable
- `POST` — create, non-idempotent
- `PUT` — full replacement, idempotent
- `PATCH` — partial update, may be non-idempotent
- `DELETE` — remove, idempotent

**Pagination:**
```
GET /users?page=2&limit=20
Response: { "data": [...], "meta": { "page": 2, "total": 142, "has_more": true } }
```

Cursor-based (preferred for large datasets):
```
GET /users?cursor=eyJpZCI6IDQyfQ==&limit=20
Response: { "data": [...], "meta": { "next_cursor": "eyJpZCI6IDYyfQ==" } }
```

**Error handling:**
```json
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "email must be a valid email address",
        "details": [
            { "field": "email", "reason": "invalid_format", "value": "not-an-email" }
        ]
    }
}
```

**Status codes:**
- 200 — success
- 201 — created
- 400 — bad request (client error)
- 401 — unauthenticated
- 403 — forbidden
- 404 — not found
- 409 — conflict (duplicate, stale version)
- 422 — validation error
- 429 — rate limited
- 5xx — server error (opaque to client)

**Versioning:**
- **URL prefix:** `/v1/users` — simple, explicit. Most common.
- **Header:** `Accept: application/vnd.api+json;version=2` — clean URL, more complex routing.
- **Never:** version via query param (`?v=1`). Pollution and poor caching.

### GraphQL

**Schema-first design:**
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  orders: [Order!]!
}

type Query {
  user(id: ID!): User
  users(page: Int, limit: Int): UserConnection!
}
```

**N+1 problem:** A query like `users { orders { items } }` causes N+1 DB queries. Solution: DataLoader (batches and caches per-request).

**Mutations:**
```graphql
type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
}
```

- Mutations return the created/updated object
- Use input types, not inline arguments
- One mutation = one semantic action

**GraphQL vs REST considerations:**
- GraphQL shifts complexity to the server (resolvers, batching, caching)
- Auth is harder (every field needs separate permission check)
- Rate limiting per query complexity, not endpoint

### gRPC

**Protobuf definition:**
```protobuf
service UserService {
  rpc GetUser (GetUserRequest) returns (User);
  rpc ListUsers (ListUsersRequest) returns (ListUsersResponse);
  rpc CreateUser (CreateUserRequest) returns (User);
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

**Streaming types:**
- **Unary** — request → response (standard RPC)
- **Server streaming** — request → stream of responses
- **Client streaming** — stream of requests → response
- **Bidirectional** — stream ↔ stream

**Error handling:**
```go
// gRPC errors use status codes + detailed messages
status.Error(codes.InvalidArgument, "email must be valid")
status.Errorf(codes.NotFound, "user %s not found", id)
```

**gRPC vs REST for internal APIs:**
- Stronger contract (Protobuf is machine-readable)
- Better performance (binary, HTTP/2 multiplexing)
- Code generation for clients
- Harder to debug (need grpcurl, binary payload)

## Practical Guidance

### API Design Checklist

- [ ] Resources named with plural nouns
- [ ] HTTP methods used semantically (GET=read, POST=create, etc.)
- [ ] Consistent error format across all endpoints
- [ ] Pagination on all list endpoints (cursor-based for large sets)
- [ ] Versioning strategy decided and documented
- [ ] Rate limiting with clear headers (`X-RateLimit-Remaining`, `Retry-After`)
- [ ] OpenAPI/Swagger spec for REST, Protobuf schema for gRPC

### REST vs gRPC Decision

Use REST if:
- API is public/third-party
- Clients are web browsers or mobile apps
- Need easy debugging with curl

Use gRPC if:
- Internal service-to-service
- High throughput, low latency critical
- Streaming (especially bidirectional)
- Strong typing across service boundaries is valuable

## Guidelines

1. REST: plural nouns, HTTP methods by semantics, consistent error format.
2. GraphQL: schema-first, DataLoader for N+1, query complexity limits.
3. gRPC: Protobuf-first, status codes for errors, streaming for bulk operations.
4. Version APIs from day one (even before the second version exists).
5. API documentation is part of the API — OpenAPI/Protobuf specs must be source of truth.
6. Rate limiting headers enable clients to handle throttling gracefully.

## Gotchas

1. **REST action endpoint**: `/users/{id}/activate` is an RPC in disguise. Use `PATCH /users/{id} { status: "active" }` or a proper command endpoint.
2. **GraphQL unlimited depth**: A deeply nested query (user → friends → friends → friends) can DoS the server. Set max query depth (e.g., 7 levels).
3. **gRPC breaking change**: Renaming a Protobuf field changes the wire format. Use field numbering carefully — never reuse numbers.
4. **Pagination without total count**: Clients building pagination UI need `total` count. Provide it for offset pagination, or don't for cursor-based (use `has_more`).
5. **Generic 500 for everything**: Errors like "not found" and "invalid input" returned as 500. Clients can't distinguish. Use appropriate status codes and structured error bodies.

## Integration

- `architect-microservices` — inter-service API contracts, service boundaries
- `architect-system-design` — rate limiting, circuit breaker for API gateways
- `architect-observability` — API metrics, tracing API calls