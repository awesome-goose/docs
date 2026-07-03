# Frequently Asked Questions

Common questions about the Goose framework.

## General

### What is Goose?

Goose is a modular Go framework for building API, Web, and CLI applications. It provides dependency injection, routing, middleware, and built-in modules for common tasks.

### Why use Goose over other frameworks?

- **Multi-platform** - Build API, Web, and CLI from one codebase
- **Modular** - Use only what you need
- **Dependency injection** - Clean, testable architecture
- **Go conventions** - Follows Go idioms and best practices

### What Go version is required?

Goose requires Go 1.21 or later.

### Is Goose production-ready?

Yes, Goose is designed for production use with features like graceful shutdown, health checks, and proper error handling.

## Getting Started

### How do I install Goose?

```bash
# Install CLI
go install github.com/awesome-goose/cli@latest

# Create new project
goose app --name=myapp --template=api
```

### How do I create my first API?

See the [Quick Start Guide](../getting-started/quick-start.md) for a step-by-step tutorial.

### Can I use Goose without the CLI?

Yes, you can manually set up a project:

```go
package main

import (
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    platform := api.NewPlatform(api.WithPort(8080))
    module := &AppModule{}

    stop, _ := goose.Start(goose.API(platform, module, nil))
    defer stop()
}
```

## Architecture

### How does dependency injection work?

Use struct tags to inject dependencies:

```go
type UserController struct {
    service *UserService `inject:""`
}
```

The container automatically resolves and injects dependencies.

### What's the difference between Imports, Exports, and Declarations?

- **Imports** - External modules this module depends on
- **Exports** - Services this module provides to others
- **Declarations** - Internal controllers and services

### Can I have circular dependencies?

No, circular dependencies are not supported. Refactor using a shared module or interfaces.

## Routing

### How do I define routes?

Build a route module with `router.ForRoutes` and import it from your module:

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{Controller{}, "List"}),
    router.Post("/users", []any{Controller{}, "Create"}),
)
```

### How do I access route parameters?

Declare them on a DTO with the `param` tag — the kernel binds them for you:

```go
type ShowDto struct {
    ID string `param:"id"`
}

func (c *Controller) Show(dto *ShowDto) types.Output {
    return output.JSON(c.service.GetByID(dto.ID))
}
```

### Can I use regex in routes?

Not directly. Use path parameters and validate in your handler:

```go
router.Get("/users/:id", []any{Controller{}, "Show"})
```

## Database

### Which databases are supported?

Through GORM: PostgreSQL, MySQL, SQLite, SQL Server.

### How do I run migrations?

Use AutoMigrate in the `Boot` startup hook (implements `types.Bootable`):

```go
func (s *AppService) Boot(k types.Kernel) error {
    return s.db.AutoMigrate(&User{}, &Post{})
}
```

### How do I use transactions?

```go
tx := s.db.Begin()
// Operations...
tx.Commit() // or tx.Rollback()
```

## Testing

### How do I test my application?

Use standard Go testing:

```go
func TestUserService(t *testing.T) {
    service := &UserService{}
    result, err := service.Create(dto)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

### How do I mock dependencies?

Create interfaces and mock implementations:

```go
type MockDB struct {
    users map[string]*User
}
```

## Deployment

### How do I deploy to production?

1. Build: `go build -o myapp main.go`
2. Configure environment variables
3. Run behind a reverse proxy (nginx)

See [Deployment Guide](../deployment/overview.md).

### How do I use Docker?

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server main.go

FROM alpine:3.19
COPY --from=builder /app/server /server
CMD ["/server"]
```

## Performance

### How do I improve performance?

1. Use caching for frequent queries
2. Add database indexes
3. Use connection pooling
4. Enable response compression
5. Use background jobs for slow operations

### How do I enable caching?

```go
cacheModule := cache.NewModule(
    cache.WithDriver("redis"),
    cache.WithHost("localhost"),
)
```

## Troubleshooting

### "Module not found" error

Ensure the module is properly imported:

```go
func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        &users.UsersModule{},
    }
}
```

### "Dependency not resolved" error

Check that the dependency is declared or exported:

```go
func (m *UsersModule) Declarations() []any {
    return []any{
        &UserService{},
    }
}
```

### Application won't start

Check:

1. Port not already in use
2. Database connection string correct
3. Environment variables set
4. No circular dependencies

## Getting Help

- **Documentation** - Read the docs thoroughly
- **GitHub Issues** - Report bugs or request features
- **Discussions** - Ask questions
- **Examples** - Check example projects

## Contributing

See [Contributing Guide](contributing.md) for how to contribute to Goose.
