---
name: golang-dependencies
description: Go dependency management — go.mod, go.sum, vendoring, replace directives, govulncheck
license: MIT
compatibility: opencode
---

# Go Dependency Management

## When to Activate

Activate when the task involves:
- Adding or updating dependencies
- Setting up go.mod
- Vendoring dependencies
- Using replace directives for local development or forks
- Checking dependency security (govulncheck)
- Deciding whether to add a dependency vs using stdlib

Do not activate for: writing Go code logic, project structure, or toolchain commands. Use golang-libraries for specific library recommendations.

## Core Principles

- **Minimize dependencies.** Each dependency is a risk (security, breaking changes, binary size).
- Prefer the standard library — Go stdlib is very rich.
- If 20 lines of code can replace a library, write the code.

## go.mod

```go
module github.com/user/project

go 1.22

require (
    github.com/lib/pq v1.10.9
    golang.org/x/sync v0.7.0
)
```

## Managing Dependencies

```bash
go get github.com/lib/pq@v1.10.9   # add/update
go get -u ./...                     # update all
go mod tidy                          # clean up
go mod verify                        # verify integrity
go mod download                      # download to cache
```

## Vendor

```bash
go mod vendor
```

Use vendor when:
- CI has no internet access
- Full control over dependency versions required
- Reproducible build required

## Indirect Dependencies

Don't add `// indirect` manually — `go mod tidy` manages this automatically.

## Replace

Only for local development or forks:

```
replace github.com/orig/lib => ../local/lib
replace github.com/orig/lib => github.com/our-fork/lib v1.2.3-fix
```

**Don't commit replace without need.**

## Security

- Regularly update: `go get -u ./... && go mod tidy`
- Use `govulncheck`: `go run golang.org/x/vuln/cmd/govulncheck ./...`
- Don't use libraries with known unpatched CVEs

## Integration

- golang-libraries — specific library recommendations
- golang-toolchain — go mod commands
- golang-project-structure — module path conventions