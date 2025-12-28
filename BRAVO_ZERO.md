# Ent-Contrib in the Bravo Zero Ecosystem

> **DeepCreative Fork:** `github.com/DeepCreative/ent-contrib`
>
> This is DeepCreative's fork of the official [ent-contrib](https://github.com/ent/ent) repository, containing extensions and tooling for the [Ent ORM](https://entgo.io) framework.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Core Extensions](#core-extensions)
3. [Integration with Bravo Zero Services](#integration-with-bravo-zero-services)
4. [Repository Structure](#repository-structure)
5. [Onboarding Guide](#onboarding-guide)
6. [Development Workflow](#development-workflow)
7. [Troubleshooting](#troubleshooting)

---

## Overview

**ent-contrib** provides a collection of extensions for the [Ent](https://entgo.io) entity framework for Go. In the Bravo Zero ecosystem, it serves as the foundational layer for:

- **GraphQL API Generation** - Automatic GraphQL schema and resolver generation from Ent schemas
- **OpenAPI Specification** - RESTful API documentation generation
- **Protocol Buffers / gRPC** - Service definition generation for microservice communication
- **Schema AST Manipulation** - Programmatic schema modifications

### Why This Matters

Bravo Zero's Go services (`www/server`, `dreamscape-go`) use Ent as their ORM layer. This fork allows DeepCreative to:

1. Customize GraphQL generation behavior
2. Add Bravo Zero-specific annotations and features
3. Fix bugs faster than waiting for upstream releases
4. Maintain compatibility across all Bravo Zero services

---

## Core Extensions

### 1. `entgql` - GraphQL Integration

**Location:** `entgql/`

Generates GraphQL schemas and resolvers from Ent entity schemas. This is the **most heavily used** extension in Bravo Zero.

**Key Features:**
- Automatic `Node` interface implementation for Relay-compatible APIs
- Pagination with Relay cursor connections
- Where input filters for complex queries
- Transaction support for mutations
- Edge/relationship resolution

**Usage in Bravo Zero:**

```go
// In ent/schema/user.go
import "entgo.io/contrib/entgql"

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("chats", Chat.Type).
            Ref("members").
            Through("chat_memberships", ChatMembership.Type).
            Annotations(entgql.RelayConnection()), // ← Generates Relay connection
    }
}
```

**Generated GraphQL:**
```graphql
type User implements Node {
    id: ID!
    chats(first: Int, after: Cursor, last: Int, before: Cursor): ChatConnection!
}
```

### 2. `entoas` - OpenAPI/Swagger Generation

**Location:** `entoas/`

Generates OpenAPI 3.0 specification documents from Ent schemas, enabling:
- Swagger UI documentation
- Client SDK generation
- API contract validation

### 3. `entproto` - Protocol Buffers / gRPC

**Location:** `entproto/`

Generates `.proto` files and gRPC service definitions from Ent schemas for:
- Microservice communication
- Language-agnostic API contracts
- High-performance inter-service calls

### 4. `schemast` - Schema AST Manipulation

**Location:** `schemast/`

Programmatic manipulation of Ent schema files for:
- Code generation tools
- Schema migrations
- Automated refactoring

---

## Integration with Bravo Zero Services

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         BRAVO ZERO ECOSYSTEM                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐   │
│  │   www/server    │     │  dreamscape-go  │     │  Future Go Svc  │   │
│  │  (Main Backend) │     │   (Dreamscape)  │     │                 │   │
│  └────────┬────────┘     └────────┬────────┘     └────────┬────────┘   │
│           │                       │                       │            │
│           └───────────────────────┼───────────────────────┘            │
│                                   │                                    │
│                    ┌──────────────▼──────────────┐                     │
│                    │        ent-contrib          │                     │
│                    │   (DeepCreative Fork)       │                     │
│                    │                             │                     │
│                    │  ┌────────┐  ┌─────────┐   │                     │
│                    │  │ entgql │  │ entoas  │   │                     │
│                    │  └────────┘  └─────────┘   │                     │
│                    │  ┌─────────┐ ┌──────────┐  │                     │
│                    │  │entproto│ │ schemast │  │                     │
│                    │  └─────────┘ └──────────┘  │                     │
│                    └─────────────────────────────┘                     │
│                                   │                                    │
│                    ┌──────────────▼──────────────┐                     │
│                    │      entgo.io/ent           │                     │
│                    │    (Ent ORM Framework)      │                     │
│                    └─────────────────────────────┘                     │
│                                   │                                    │
│                    ┌──────────────▼──────────────┐                     │
│                    │        PostgreSQL           │                     │
│                    │        (Database)           │                     │
│                    └─────────────────────────────┘                     │
└─────────────────────────────────────────────────────────────────────────┘
```

### Services Using ent-contrib

| Service | Repository | Uses entgql | Uses entoas | Uses entproto |
|---------|------------|-------------|-------------|---------------|
| Main Backend | `www/server` | ✅ | ❌ | ✅ (planned) |
| Dreamscape | `dreamscape-go` | ✅ | ❌ | ❌ |

### Go Module Configuration

Services reference the DeepCreative fork via `go.mod` replace directives:

```go
// www/server/go.mod
replace entgo.io/contrib => github.com/DeepCreative/ent-contrib v1.2.2

// dreamscape-go/go.mod
replace entgo.io/contrib => github.com/DeepCreative/ent-contrib v0.0.0-20240508033148-e2589ff89e2e
```

---

## Repository Structure

```
ent-contrib/
├── entgql/                    # GraphQL extension (PRIMARY)
│   ├── extension.go           # Main extension entry point
│   ├── annotation.go          # Schema annotations (@entgql.*)
│   ├── schema.go              # GraphQL schema generation
│   ├── pagination.go          # Relay cursor pagination
│   ├── transaction.go         # Mutation transaction handling
│   ├── template/              # Go template files for codegen
│   │   ├── collection.tmpl
│   │   ├── edge.tmpl
│   │   ├── mutation_input.tmpl
│   │   ├── node.tmpl
│   │   ├── pagination.tmpl
│   │   └── where_input.tmpl
│   └── internal/              # Test fixtures and examples
│       ├── todo/              # Basic example
│       ├── todofed/           # Federation example
│       ├── todoglobalid/      # Global ID example
│       └── todopulid/         # PULID example
│
├── entoas/                    # OpenAPI extension
│   ├── extension.go
│   ├── annotation.go
│   └── generator.go
│
├── entproto/                  # Protocol Buffers extension
│   ├── extension.go
│   ├── message.go
│   ├── service.go
│   ├── field.go
│   └── cmd/
│       ├── entproto/          # CLI tool
│       ├── protoc-gen-ent/    # Ent schema from proto
│       └── protoc-gen-entgrpc/# gRPC service generator
│
├── schemast/                  # Schema AST manipulation
│   ├── mutate.go
│   ├── load.go
│   └── print.go
│
├── go.mod                     # Module dependencies
├── go.sum
└── README.md
```

---

## Onboarding Guide

### For New Engineers / AI Agents

#### Prerequisites

1. **Go 1.23+** installed
2. **PostgreSQL** for local development
3. Familiarity with:
   - Go language basics
   - GraphQL concepts
   - ORM patterns

#### Step 1: Understand the Ent Framework

Before diving into ent-contrib, understand Ent itself:

```bash
# Read the Ent documentation
open https://entgo.io/docs/getting-started
```

Key Ent concepts:
- **Schemas** define entity types (`ent/schema/*.go`)
- **Fields** define columns
- **Edges** define relationships
- **Mixins** provide reusable schema components
- **Hooks** intercept mutations

#### Step 2: Explore a Bravo Zero Service

```bash
# Look at www/server's Ent usage
cd /Users/ibdrew/Documents/www/server

# Explore the schema definitions
ls internal/ent/schema/

# See how entgql annotations are used
cat internal/ent/schema/user.go
```

#### Step 3: Understand Code Generation

Ent and entgql use code generation:

```bash
# In any Ent-using service
go generate ./...

# This runs ent/generate.go which invokes:
# 1. Ent core codegen (creates internal/ent/*.go)
# 2. entgql extension (creates ent.graphql, resolvers)
# 3. gqlgen (creates GraphQL server code)
```

#### Step 4: Make a Change

1. Modify an Ent schema (add a field, edge, or annotation)
2. Run `go generate ./...`
3. Observe the generated changes
4. Test with the GraphQL playground

### Common Tasks

#### Adding a New Entity

```go
// 1. Create ent/schema/my_entity.go
package schema

import (
    "entgo.io/contrib/entgql"
    "entgo.io/ent"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
)

type MyEntity struct {
    ent.Schema
}

func (MyEntity) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
    }
}

func (MyEntity) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),     // Expose in GraphQL queries
        entgql.Mutations(        // Generate mutation inputs
            entgql.MutationCreate(),
            entgql.MutationUpdate(),
        ),
    }
}
```

```bash
# 2. Regenerate code
go generate ./...

# 3. Run migrations
go run ./cmd/migrate
```

#### Adding GraphQL Relay Pagination

```go
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type).
            Annotations(entgql.RelayConnection()), // ← Add this
    }
}
```

#### Customizing Where Filters

```go
func (MyEntity) Fields() []ent.Field {
    return []ent.Field{
        field.String("email").
            Annotations(
                entgql.OrderField("EMAIL"),  // Allow ordering
                entgql.Skip(entgql.SkipWhereInput), // Skip filter generation
            ),
    }
}
```

---

## Development Workflow

### Making Changes to ent-contrib

When you need to modify ent-contrib itself:

```bash
# 1. Clone the fork (already in workspace)
cd /Users/ibdrew/Documents/ent-contrib

# 2. Create a branch
git checkout -b feature/my-change

# 3. Make changes

# 4. Run tests
go test ./...

# 5. Test in a consuming service
# Edit the service's go.mod to use local path:
# replace entgo.io/contrib => /Users/ibdrew/Documents/ent-contrib

# 6. Commit and push
git add .
git commit -m "feat: description of change"
git push origin feature/my-change

# 7. Create PR and get review

# 8. After merge, update consuming services to new version
```

### Updating Services After Fork Changes

```bash
# In www/server or dreamscape-go
go get github.com/DeepCreative/ent-contrib@latest

# Or pin to specific commit
go get github.com/DeepCreative/ent-contrib@abc123def

# Regenerate code
go generate ./...
```

---

## Troubleshooting

### Common Issues

#### "undefined: entgql.X"

**Cause:** Version mismatch between ent and ent-contrib

**Solution:**
```bash
go get entgo.io/ent@latest
go get github.com/DeepCreative/ent-contrib@latest
go mod tidy
go generate ./...
```

#### GraphQL schema not updating

**Cause:** Stale generated files

**Solution:**
```bash
rm -f internal/ent/*.go
rm -f internal/graph/generated.go
rm -f internal/graph/ent.graphql
go generate ./...
```

#### "cannot find module providing package"

**Cause:** Replace directive not working

**Solution:** Ensure the replace path is correct and the module exists:
```bash
# Check your go.mod
cat go.mod | grep replace

# Verify the fork is accessible
go list -m github.com/DeepCreative/ent-contrib@latest
```

---

## Related Resources

- [Ent Documentation](https://entgo.io)
- [gqlgen Documentation](https://gqlgen.com)
- [Upstream ent-contrib](https://github.com/ent/contrib)
- [GraphQL Relay Spec](https://relay.dev/graphql/connections.htm)

---

## Contact

For questions about ent-contrib usage in Bravo Zero:
- Check the `#backend` channel
- Review existing PRs in `www` and `dreamscape-go` repos
- Consult this documentation














