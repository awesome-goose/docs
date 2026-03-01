# Dependency Injection

Goose includes a powerful IoC (Inversion of Control) container for dependency injection, making it easy to manage service dependencies.

## How It Works

Dependencies are injected automatically using struct tags:

```go
type UserController struct {
    service *UserService `inject:""`  // Injected automatically
    log     types.Log    `inject:""`  // Framework service
    cache   *cache.Cache `inject:""`  // Module service
}
```

## Basic Injection

### Service Injection

```go
// Define a service
type UserService struct {
    db *sql.Db `inject:""`  // Database injected
}

// Use the service
func (s *UserService) GetUser(id string) (*User, error) {
    var user User
    err := s.db.First(&user, "id = ?", id).Error
    return &user, err
}
```

### Controller Injection

```go
type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) Show(ctx types.Context) any {
    id := ctx.Param("id")
    return c.service.GetUser(id)
}
```

## Injection Tags

### Default Injection

The empty `inject:""` tag resolves by type:

```go
type MyService struct {
    db *sql.Db `inject:""` // Resolves *sql.Db type
}
```

### Named Injection

Use `inject:"name"` for named services:

```go
type MyService struct {
    primaryDB   *sql.Db `inject:"primary"`
    secondaryDB *sql.Db `inject:"secondary"`
}
```

### Type Injection

Use `inject:"type"` to explicitly resolve by type:

```go
type MyService struct {
    logger types.Log `inject:"type"`
}
```

## Registering Services

### In Module Declarations

Most services are registered through module declarations:

```go
func (m *AppModule) Declarations() []any {
    return []any{
        &UserService{},      // Registered automatically
        &UserController{},
        &OrderService{},
    }
}
```

### Using Initializers

For custom registration, use initializers:

```go
func main() {
    initializers := []func(types.Container) error{
        func(c types.Container) error {
            // Register with factory function
            c.Register(func() *MyCustomService {
                return &MyCustomService{
                    apiKey: os.Getenv("API_KEY"),
                }
            }, "", true)
            return nil
        },
    }

    stop, err := goose.Start(goose.API(platform, module, initializers))
}
```

### Factory Functions

Register services using factory functions:

```go
// Register a singleton
container.Register(func(db *sql.Db) *UserRepository {
    return &UserRepository{db: db}
}, "", true)  // true = singleton

// Register a transient (new instance each time)
container.Register(func() *RequestContext {
    return &RequestContext{}
}, "", false)  // false = transient
```

## Resolving Dependencies

### Automatic Resolution

Dependencies are resolved automatically when needed:

```go
type OrderService struct {
    userService    *UserService    `inject:""`  // Auto-resolved
    productService *ProductService `inject:""`  // Auto-resolved
    db             *sql.Db         `inject:""`  // Auto-resolved
}

// When OrderService is created, all dependencies are injected
```

### Manual Resolution

You can also resolve manually:

```go
func customInitializer(container types.Container) error {
    // Resolve by type
    userService := container.Resolve(&UserService{})

    // Use the service
    return nil
}
```

## Service Lifetime

### Singleton

One instance shared across the application:

```go
// Registered as singleton (default)
container.Register(func() *DatabaseConnection {
    return &DatabaseConnection{}
}, "", true)  // true = singleton
```

### Transient

New instance created each time:

```go
container.Register(func() *RequestScope {
    return &RequestScope{}
}, "", false)  // false = transient
```

## Dependency Graph

Goose automatically resolves the dependency graph:

```go
// These dependencies are resolved in correct order
type OrderController struct {
    service *OrderService `inject:""`
}

type OrderService struct {
    repo    *OrderRepository `inject:""`
    events  *EventDispatcher `inject:""`
}

type OrderRepository struct {
    db *sql.Db `inject:""`
}

// Resolution order: sql.Db → OrderRepository → EventDispatcher → OrderService → OrderController
```

## Interface Binding

Bind interfaces to concrete implementations:

```go
// Define interface
type UserRepositoryInterface interface {
    FindByID(id string) (*User, error)
}

// Implement interface
type SqlUserRepository struct {
    db *sql.Db `inject:""`
}

func (r *SqlUserRepository) FindByID(id string) (*User, error) {
    // Implementation
}

// Register binding
container.Register(func(db *sql.Db) UserRepositoryInterface {
    return &SqlUserRepository{db: db}
}, "", true)
```

## Circular Dependencies

Goose detects and prevents circular dependencies:

```go
// ❌ This will cause an error
type ServiceA struct {
    serviceB *ServiceB `inject:""`
}

type ServiceB struct {
    serviceA *ServiceA `inject:""`
}

// ✅ Solution: Use interfaces or refactor
type ServiceA struct {
    sharedService *SharedService `inject:""`
}

type ServiceB struct {
    sharedService *SharedService `inject:""`
}
```

## Testing with DI

The DI system makes testing easy:

```go
func TestUserController(t *testing.T) {
    // Create mock service
    mockService := &MockUserService{
        users: map[string]*User{
            "1": {ID: "1", Name: "Test User"},
        },
    }

    // Create controller with mock
    controller := &UserController{
        service: mockService,
    }

    // Test the controller
    result := controller.Show(mockContext("1"))
    // Assert...
}
```

## Best Practices

### 1. Constructor Functions

Use constructor functions for complex initialization:

```go
func NewUserService(db *sql.Db, cache *cache.Cache) *UserService {
    return &UserService{
        db:    db,
        cache: cache,
    }
}

// Register with constructor
container.Register(NewUserService, "", true)
```

### 2. Interface Segregation

Depend on interfaces, not concrete types:

```go
// ✅ Good: Depends on interface
type OrderService struct {
    repo OrderRepositoryInterface `inject:""`
}

// ❌ Bad: Depends on concrete type
type OrderService struct {
    repo *SqlOrderRepository `inject:""`
}
```

### 3. Minimal Dependencies

Keep dependencies minimal:

```go
// ✅ Good: Only needed dependencies
type UserService struct {
    db *sql.Db `inject:""`
}

// ❌ Bad: Too many dependencies
type UserService struct {
    db      *sql.Db      `inject:""`
    cache   *cache.Cache `inject:""`
    queue   *Queue       `inject:""`
    mailer  *Mailer      `inject:""`
    // ... more
}
```

## Next Steps

- [Lifecycle](lifecycle.md) - Application lifecycle hooks
- [Modules](modules.md) - Module system
- [Services](../building-blocks/services.md) - Service patterns
