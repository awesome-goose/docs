# Modules

Modules are the fundamental building blocks of Goose applications. They encapsulate related functionality and promote code organization.

## Module Interface

Every Goose module implements the `Module` interface:

```go
type Module interface {
    Imports() []Module      // Modules this module depends on
    Exports() []any         // Services available to other modules
    Declarations() []any    // Components provided by this module
}
```

## Creating a Module

### Basic Module

```go
package users

import "github.com/awesome-goose/goose/types"

type UsersModule struct{}

func (m *UsersModule) Imports() []types.Module {
    return []types.Module{}
}

func (m *UsersModule) Exports() []any {
    return []any{
        &UsersService{},
    }
}

func (m *UsersModule) Declarations() []any {
    return []any{
        &UsersController{},
        &UsersService{},
    }
}
```

### Using the CLI

Generate a module with the Goose CLI:

```bash
# Plain module
goose g module --name=users --type=plain

# Resource module with entity
goose g module --name=products --type=resource
```

## Module Components

### Declarations

Declarations are components that belong to this module:

```go
func (m *UsersModule) Declarations() []any {
    return []any{
        &UsersController{},
        &UsersService{},
        &UsersRepository{},
        &AuthMiddleware{},
    }
}
```

### Exports

Exports are services available to other modules:

```go
func (m *UsersModule) Exports() []any {
    return []any{
        &UsersService{}, // Other modules can inject this
    }
}
```

### Imports

Imports are dependencies on other modules:

```go
func (m *OrdersModule) Imports() []types.Module {
    return []types.Module{
        &users.UsersModule{},     // Need users service
        &products.ProductsModule{}, // Need products service
    }
}
```

## Module Dependencies

### Importing Services

When you import a module, its exported services become available:

```go
// orders/orders.service.go
type OrdersService struct {
    usersService    *users.UsersService    `inject:""`
    productsService *products.ProductsService `inject:""`
}

func (s *OrdersService) CreateOrder(userID string, productIDs []string) (*Order, error) {
    user, err := s.usersService.GetByID(userID)
    if err != nil {
        return nil, err
    }
    // Use the imported services...
}
```

### Circular Dependencies

Avoid circular dependencies between modules:

```
❌ Bad:
UsersModule → OrdersModule → UsersModule

✅ Good:
UsersModule → SharedModule
OrdersModule → SharedModule
```

Solution: Extract shared functionality into a separate module.

## Global Modules

Global modules are automatically available to all modules:

```go
type SharedModule struct{}

func (m *SharedModule) IsGlobal() bool {
    return true
}

func (m *SharedModule) Exports() []any {
    return []any{
        &LoggingService{},
        &CacheService{},
    }
}
```

## Built-in Modules

Goose provides several built-in modules:

### SQL Module

```go
import "github.com/awesome-goose/goose/modules/sql"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.NewModule(
            sql.WithDriver("sqlite"),
            sql.WithDatabase("./data/app.db"),
        ),
    }
}
```

### Cache Module

```go
import "github.com/awesome-goose/goose/modules/cache"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        cache.NewModule(
            cache.WithDefaultTTL(15 * time.Minute),
        ),
    }
}
```

### Queue Module

```go
import "github.com/awesome-goose/goose/modules/queues"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        queues.NewModule(
            queues.WithQueue("default"),
        ),
    }
}
```

### KV Module

```go
import "github.com/awesome-goose/goose/modules/kv"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        kv.NewModule(
            kv.WithGroup("app"),
        ),
    }
}
```

### Cron Module

```go
import "github.com/awesome-goose/goose/modules/cron"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        cron.NewModule(),
    }
}
```

## Module Best Practices

### 1. Single Responsibility

Each module should handle one feature:

```go
// ✅ Good: Focused module
type AuthModule struct{}     // Authentication only
type UsersModule struct{}    // User management only
type EmailModule struct{}    // Email sending only

// ❌ Bad: Too many responsibilities
type UtilsModule struct{}    // Authentication + Users + Email
```

### 2. Minimal Exports

Only export what other modules need:

```go
func (m *UsersModule) Exports() []any {
    return []any{
        &UsersService{},      // ✅ Public API
        // &UsersRepository{}, // ❌ Internal, don't export
    }
}
```

### 3. Clear Dependencies

Import only what you need:

```go
func (m *OrdersModule) Imports() []types.Module {
    return []types.Module{
        &users.UsersModule{},    // ✅ Actually needed
        // &email.EmailModule{}, // ❌ Not used here
    }
}
```

### 4. Testable Design

Design modules for easy testing:

```go
// Module with interface for testing
type UsersServiceInterface interface {
    GetByID(id string) (*User, error)
}

// Concrete implementation
type UsersService struct {
    db *sql.Db `inject:""`
}

func (s *UsersService) GetByID(id string) (*User, error) {
    // Implementation
}
```

## Module Structure

Recommended file organization:

```
users/
├── users.module.go        # Module definition
├── users.controller.go    # HTTP handlers
├── users.service.go       # Business logic
├── users.repository.go    # Data access (optional)
├── users.routes.go        # Route definitions
├── users.dtos.go          # Data transfer objects
├── users.entity.go        # Database entity
├── users.middleware.go    # Module middleware
└── users_test.go          # Tests
```

## Next Steps

- [Dependency Injection](dependency-injection.md) - How services are injected
- [Controllers](../building-blocks/controllers.md) - Building controllers
- [Services](../building-blocks/services.md) - Service layer patterns
