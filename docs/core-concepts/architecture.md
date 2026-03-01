# Architecture Overview

Goose follows a modular, layered architecture that promotes separation of concerns and maintainability.

## High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                        Application                            │
├──────────────────────────────────────────────────────────────┤
│  Platform Layer                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   API       │  │   Web       │  │   CLI       │          │
│  │  Platform   │  │  Platform   │  │  Platform   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
├──────────────────────────────────────────────────────────────┤
│  Module Layer                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Users     │  │  Products   │  │   Orders    │          │
│  │   Module    │  │   Module    │  │   Module    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
├──────────────────────────────────────────────────────────────┤
│  Core Layer                                                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │ Kernel   │ │ Router   │ │Container │ │Traverser │        │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
├──────────────────────────────────────────────────────────────┤
│  Infrastructure Layer                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │   SQL    │ │  Cache   │ │  Queue   │ │   Cron   │        │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
└──────────────────────────────────────────────────────────────┘
```

## Core Components

### Kernel

The kernel is the heart of Goose. It orchestrates:

- Application bootstrapping
- Module traversal
- Route registration
- Request handling

```go
// The kernel starts your application
stop, err := goose.Start(goose.API(platform, module, initializers))
```

### Container (IoC)

The dependency injection container:

- Registers services
- Resolves dependencies
- Manages singletons

```go
// Services are automatically injected
type UserController struct {
    service  *UserService `inject:""`
    log      types.Log    `inject:""`
    cache    *cache.Cache `inject:""`
}
```

### Router

The router matches incoming requests to handlers:

- Method-based routing (GET, POST, etc.)
- Parameter extraction (`:id`)
- Middleware aggregation
- Route caching for performance

### Traverser

The traverser walks the module tree:

- Discovers modules
- Registers declarations
- Handles imports/exports
- Executes lifecycle hooks

## Request Flow

```
   Request
      │
      ▼
┌─────────────────┐
│    Platform     │  ← Receives HTTP/CLI request
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Router      │  ← Finds matching route
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Middleware    │  ← Executes middleware chain
│     Stack       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Controller    │  ← Handles the request
│    Handler      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Service      │  ← Business logic
│     Layer       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Response      │  ← Sends response
│   Serializer    │
└─────────────────┘
```

## Module System

Modules are the building blocks of Goose applications:

```go
type Module interface {
    Imports() []Module      // Dependencies
    Exports() []any         // Public services
    Declarations() []any    // Local components
}
```

### Module Lifecycle

1. **Import Resolution** - Resolve module dependencies
2. **Declaration Registration** - Register controllers, services
3. **Export Publishing** - Make services available
4. **Hook Execution** - Run lifecycle hooks

## Platform Abstraction

Goose supports multiple platforms through a unified interface:

```go
type Platform interface {
    Type() PlatformType
    Name() string
    Boot(container Container) (App, error)
}
```

### API Platform

Handles HTTP requests, returns JSON responses.

### Web Platform

Handles HTTP requests, renders HTML templates.

### CLI Platform

Handles command-line arguments, outputs to console.

## Design Patterns

### Dependency Injection

All dependencies are injected automatically:

```go
type OrderService struct {
    userService    *UserService    `inject:""`
    productService *ProductService `inject:""`
    db             *sql.Db         `inject:""`
}
```

### Repository Pattern

Data access through repositories:

```go
type UserRepository struct {
    db *sql.Db `inject:""`
}

func (r *UserRepository) FindByID(id string) (*User, error) {
    // Database query
}
```

### Service Layer

Business logic in services:

```go
type OrderService struct {
    repo *OrderRepository `inject:""`
}

func (s *OrderService) CreateOrder(data CreateOrderDTO) (*Order, error) {
    // Business logic
}
```

### Controller Layer

Request handling in controllers:

```go
type OrderController struct {
    service *OrderService `inject:""`
}

func (c *OrderController) Create(ctx types.Context) any {
    var dto CreateOrderDTO
    ctx.Bind(&dto)
    return c.service.CreateOrder(dto)
}
```

## Configuration Architecture

```
Environment Variables (.env)
         │
         ▼
┌─────────────────┐
│   Env Package   │  ← Loads env vars
└────────┬────────┘
         │
         ▼
YAML Files (config/*.yaml)
         │
         ▼
┌─────────────────┐
│ Config Package  │  ← Parses YAML
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Container    │  ← Injects config
└─────────────────┘
```

## Next Steps

- [Modules](modules.md) - Deep dive into modules
- [Dependency Injection](dependency-injection.md) - IoC container
- [Lifecycle](lifecycle.md) - Application lifecycle
