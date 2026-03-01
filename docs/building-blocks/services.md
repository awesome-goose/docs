# Services

Services contain the business logic of your application, separating concerns from controllers and data access.

## Creating a Service

### Basic Service

```go
package app

type UserService struct {
    db *sql.Db `inject:""`
}

func (s *UserService) GetAllUsers() ([]User, error) {
    var users []User
    err := s.db.Find(&users).Error
    return users, err
}

func (s *UserService) GetUserByID(id string) (*User, error) {
    var user User
    err := s.db.First(&user, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    user := &User{
        ID:    uuid.New().String(),
        Name:  dto.Name,
        Email: dto.Email,
    }
    err := s.db.Create(user).Error
    return user, err
}
```

### Registering Services

Register in your module:

```go
func (m *UserModule) Declarations() []any {
    return []any{
        &UserController{},
        &UserService{},
    }
}

func (m *UserModule) Exports() []any {
    return []any{
        &UserService{},  // Available to other modules
    }
}
```

## Service Patterns

### Repository Pattern

Separate data access from business logic:

```go
// Repository (data access)
type UserRepository struct {
    db *sql.Db `inject:""`
}

func (r *UserRepository) FindByID(id string) (*User, error) {
    var user User
    err := r.db.First(&user, "id = ?", id).Error
    return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
    var user User
    err := r.db.First(&user, "email = ?", email).Error
    return &user, err
}

func (r *UserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

// Service (business logic)
type UserService struct {
    repo  *UserRepository `inject:""`
    cache *cache.Cache    `inject:""`
}

func (s *UserService) GetUser(id string) (*User, error) {
    // Check cache first
    cached, err := cache.GetAs[User](s.cache, "user:"+id)
    if err == nil {
        return &cached, nil
    }

    // Fetch from database
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // Cache for next time
    s.cache.Set("user:"+id, user, 10*time.Minute)

    return user, nil
}
```

### Service with Dependencies

```go
type OrderService struct {
    orderRepo   *OrderRepository   `inject:""`
    userService *UserService       `inject:""`
    productService *ProductService `inject:""`
    emailService *EmailService     `inject:""`
    log         types.Log          `inject:""`
}

func (s *OrderService) CreateOrder(dto CreateOrderDTO) (*Order, error) {
    // Validate user exists
    user, err := s.userService.GetUserByID(dto.UserID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // Validate products
    var items []OrderItem
    var total float64
    for _, item := range dto.Items {
        product, err := s.productService.GetByID(item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product not found: %s", item.ProductID)
        }

        items = append(items, OrderItem{
            ProductID: product.ID,
            Quantity:  item.Quantity,
            Price:     product.Price,
        })
        total += product.Price * float64(item.Quantity)
    }

    // Create order
    order := &Order{
        ID:     uuid.New().String(),
        UserID: user.ID,
        Items:  items,
        Total:  total,
        Status: "pending",
    }

    if err := s.orderRepo.Create(order); err != nil {
        return nil, err
    }

    // Send confirmation email
    go s.emailService.SendOrderConfirmation(user.Email, order)

    s.log.Info("Order created", "order_id", order.ID, "user_id", user.ID)

    return order, nil
}
```

### Transaction Support

Handle database transactions:

```go
func (s *TransferService) Transfer(fromID, toID string, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Debit from source account
        var from Account
        if err := tx.First(&from, "id = ?", fromID).Error; err != nil {
            return err
        }
        if from.Balance < amount {
            return fmt.Errorf("insufficient balance")
        }
        from.Balance -= amount
        if err := tx.Save(&from).Error; err != nil {
            return err
        }

        // Credit to destination account
        var to Account
        if err := tx.First(&to, "id = ?", toID).Error; err != nil {
            return err
        }
        to.Balance += amount
        if err := tx.Save(&to).Error; err != nil {
            return err
        }

        return nil
    })
}
```

## Service Lifecycle

### OnBoot Hook

Initialize service state on boot:

```go
type CacheService struct {
    db    *sql.Db      `inject:""`
    cache *cache.Cache `inject:""`
}

func (s *CacheService) OnBoot(kernel types.Kernel) error {
    // Warm up cache on startup
    var settings []Setting
    s.db.Find(&settings)

    for _, setting := range settings {
        s.cache.Set("setting:"+setting.Key, setting.Value, 0)
    }

    return nil
}
```

### OnShutdown Hook

Cleanup on shutdown:

```go
func (s *QueueService) OnShutdown() error {
    // Wait for pending jobs
    s.WaitForCompletion(30 * time.Second)
    return nil
}
```

## Best Practices

### 1. Keep Services Focused

```go
// ✅ Good: Single responsibility
type UserService struct { ... }      // User operations only
type OrderService struct { ... }     // Order operations only
type EmailService struct { ... }     // Email operations only

// ❌ Bad: Too many responsibilities
type AppService struct { ... }       // Users + Orders + Emails + ...
```

### 2. Use Interfaces for Dependencies

```go
// ✅ Good: Interface dependency
type UserServiceInterface interface {
    GetByID(id string) (*User, error)
}

type OrderService struct {
    users UserServiceInterface `inject:""`
}

// ❌ Bad: Concrete dependency (harder to test)
type OrderService struct {
    users *UserService `inject:""`
}
```

### 3. Handle Errors Properly

```go
// ✅ Good: Wrap errors with context
func (s *UserService) GetByID(id string) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}

// ❌ Bad: Swallow errors or lose context
func (s *UserService) GetByID(id string) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err  // Lost context
    }
    return user, nil
}
```

### 4. Log Important Operations

```go
func (s *PaymentService) ProcessPayment(orderID string, amount float64) error {
    s.log.Info("Processing payment",
        "order_id", orderID,
        "amount", amount)

    err := s.gateway.Charge(amount)
    if err != nil {
        s.log.Error("Payment failed",
            "order_id", orderID,
            "error", err)
        return err
    }

    s.log.Info("Payment successful", "order_id", orderID)
    return nil
}
```

### 5. Use DTOs at Service Boundaries

```go
// ✅ Good: DTOs at boundaries
func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    user := &User{
        Name:  dto.Name,
        Email: dto.Email,
    }
    return user, s.repo.Create(user)
}

// ❌ Bad: Entity at boundaries
func (s *UserService) CreateUser(user *User) error {
    return s.repo.Create(user)  // Caller creates entity
}
```

## Testing Services

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    mockCache := &MockCache{}

    service := &UserService{
        repo:  mockRepo,
        cache: mockCache,
    }

    // Test
    dto := CreateUserDTO{
        Name:  "John",
        Email: "john@example.com",
    }

    user, err := service.CreateUser(dto)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
    assert.True(t, mockRepo.CreateCalled)
}
```

## Next Steps

- [Controllers](controllers.md) - Request handling
- [Modules](../core-concepts/modules.md) - Module organization
- [Dependency Injection](../core-concepts/dependency-injection.md) - DI patterns
