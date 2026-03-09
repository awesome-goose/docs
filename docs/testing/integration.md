# Integration Testing

Test how multiple modules and services work together.

## Import

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "github.com/awesome-goose/goose/types"
)
```

## What is Integration Testing?

Integration tests verify that different parts of your application work together correctly. Unlike unit tests that isolate individual components, integration tests exercise the interaction between modules, services, and external systems.

## When to Write Integration Tests

- Testing module dependencies and imports
- Testing service interactions
- Testing database operations with real connections
- Testing external API integrations
- Testing the full request/response cycle

## Skipping Integration Tests

Use test tags to run integration tests only when needed:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestIntegration_UserFlow(t *testing.T) {
    // Skip unless TEST_LEVEL=integration
    goosetest.SkipUnlessIntegration(t)

    // Your integration test code...
}
```

Run integration tests with:

```bash
TEST_LEVEL=integration go test ./...
```

## Testing Module Integration

### Testing Module Imports

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app"
    "myapp/app/users"
    "myapp/app/orders"
)

func TestAppModule_IntegrationWithSubModules(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    // Test that AppModule correctly imports all submodules
    mt := gt.Module(&app.AppModule{})

    imports := mt.TestImports()
    gt.Expect(len(imports)).ToBeGreaterThan(0)

    // Verify specific modules are imported
    hasUsersModule := false
    hasOrdersModule := false
    for _, mod := range imports {
        switch mod.(type) {
        case *users.UsersModule:
            hasUsersModule = true
        case *orders.OrdersModule:
            hasOrdersModule = true
        }
    }

    gt.Expect(hasUsersModule).ToBeTrue()
    gt.Expect(hasOrdersModule).ToBeTrue()
}
```

### Testing Module Exports

```go
func TestUsersModule_ExportsUserService(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    mt := gt.Module(&users.UsersModule{})
    exports := mt.TestExports()

    // Verify UserService is exported for other modules
    hasUserService := false
    for _, exp := range exports {
        if _, ok := exp.(*users.UserService); ok {
            hasUserService = true
            break
        }
    }

    gt.Expect(hasUserService).ToBeTrue()
}
```

## Testing Service Integration

### Testing Service Dependencies

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app/orders"
    "myapp/app/users"
    "myapp/app/inventory"
)

func TestOrderService_WithRealDependencies(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    // Set up real (or test) dependencies
    gt.WithContainer(func(c *goosetest.TestContainer) {
        c.Register(&users.UserService{})
        c.Register(&inventory.InventoryService{})
    })

    orderService := &orders.OrderService{}
    st := gt.Service(orderService)

    // Test with real service interactions
    results := st.Run("CreateOrder", orders.CreateOrderDTO{
        UserID: "user-123",
        Items: []orders.OrderItem{
            {ProductID: "prod-1", Quantity: 2},
        },
    })

    gt.Expect(results[1]).ToBeNil() // no error
    order := results[0].(*orders.Order)
    gt.Expect(order.UserID).ToEqual("user-123")
}
```

### Testing Cross-Service Communication

```go
func TestUserRegistration_TriggersNotifications(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    // Set up services with real notification service
    notificationService := &notifications.NotificationService{}
    gt.WithContainer(func(c *goosetest.TestContainer) {
        c.Register(notificationService)
    })

    userService := &users.UserService{}
    st := gt.Service(userService)

    // Register a new user
    results := st.Run("Register", users.RegisterDTO{
        Email: "new@example.com",
        Name:  "New User",
    })

    gt.Expect(results[1]).ToBeNil()

    // Verify notification was sent
    gt.Expect(notificationService.SentCount()).ToEqual(1)
}
```

## Testing with Test Database

### Database Setup

```go
import (
    "os"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app/users"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    // Setup test database
    testDB = setupTestDatabase()

    code := m.Run()

    // Cleanup
    testDB.Close()
    os.Exit(code)
}

func setupTestDatabase() *sql.DB {
    db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        panic(err)
    }

    // Run migrations
    runMigrations(db)

    return db
}

func TestUserRepository_Integration(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    // Clean up before test
    testDB.Exec("TRUNCATE users CASCADE")

    repo := &users.UserRepository{DB: testDB}

    // Create user
    user, err := repo.Create(&users.User{
        Email: "test@example.com",
        Name:  "Test User",
    })
    gt.Expect(err).ToBeNil()
    gt.Expect(user.ID).Not().ToEqual("")

    // Find user
    found, err := repo.FindByID(user.ID)
    gt.Expect(err).ToBeNil()
    gt.Expect(found.Email).ToEqual("test@example.com")

    // Update user
    user.Name = "Updated Name"
    err = repo.Update(user)
    gt.Expect(err).ToBeNil()

    // Delete user
    err = repo.Delete(user.ID)
    gt.Expect(err).ToBeNil()

    // Verify deletion
    _, err = repo.FindByID(user.ID)
    gt.Expect(err).Not().ToBeNil()
}
```

### Transaction Testing

```go
func TestOrderService_TransactionRollback(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    testDB.Exec("TRUNCATE orders, inventory CASCADE")

    orderService := &orders.OrderService{DB: testDB}

    // Try to create order with insufficient inventory
    _, err := orderService.CreateOrder(orders.CreateOrderDTO{
        UserID: "user-123",
        Items: []orders.OrderItem{
            {ProductID: "prod-1", Quantity: 1000}, // More than available
        },
    })

    // Should fail
    gt.Expect(err).Not().ToBeNil()

    // Verify no partial data was saved (transaction rolled back)
    var count int
    testDB.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = $1", "user-123").Scan(&count)
    gt.Expect(count).ToEqual(0)
}
```

## Testing HTTP Integration

### Full Stack HTTP Test

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app"
)

func TestFullStackAPI(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)

    // Create handler with real modules
    handler := createHandler(&app.AppModule{})
    ht := goosetest.NewHTTPTest(t, handler).WithServer()
    defer ht.Close()

    // Test full user creation flow
    var user map[string]any
    ht.POST("/api/users").
        WithJSON(map[string]string{
            "email": "integration@test.com",
            "name":  "Integration Test",
        }).
        Do().
        ExpectCreated().
        JSON(&user)

    userID := user["id"].(string)
    ht.T.Expect(userID).Not().ToEqual("")

    // Verify user exists in database
    ht.GET("/api/users/" + userID).
        Do().
        ExpectOK()

    // Create order for user
    ht.POST("/api/orders").
        WithJSON(map[string]any{
            "user_id": userID,
            "items": []map[string]any{
                {"product_id": "prod-1", "quantity": 2},
            },
        }).
        Do().
        ExpectCreated()

    // Get user's orders
    var orders []map[string]any
    ht.GET("/api/users/" + userID + "/orders").
        Do().
        ExpectOK().
        JSON(&orders)

    ht.T.Expect(len(orders)).ToEqual(1)
}
```

## Testing External Services

### Mocking External APIs

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestPaymentService_Integration(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    gt := goosetest.NewGooseTest(t)

    // Create mock payment gateway
    mockGateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"transaction_id": "txn-123", "status": "success"}`))
    }))
    defer mockGateway.Close()

    // Create payment service with mock gateway URL
    paymentService := &payments.PaymentService{
        GatewayURL: mockGateway.URL,
    }

    // Test payment flow
    result, err := paymentService.ProcessPayment(payments.PaymentRequest{
        Amount:   100.00,
        Currency: "USD",
        CardID:   "card-123",
    })

    gt.Expect(err).ToBeNil()
    gt.Expect(result.TransactionID).ToEqual("txn-123")
    gt.Expect(result.Status).ToEqual("success")
}
```

## Best Practices

1. **Use test tags** - Skip integration tests by default with `SkipUnlessIntegration`
2. **Clean test data** - Reset database state before each test
3. **Use transactions** - Wrap tests in transactions for easy rollback
4. **Mock external services** - Use mock servers for external APIs
5. **Test realistic scenarios** - Test complete user flows, not just happy paths
6. **Isolate tests** - Each test should be independent

## Running Integration Tests

```bash
# Run only integration tests
TEST_LEVEL=integration go test ./...

# Run with verbose output
TEST_LEVEL=integration go test -v ./...

# Run specific integration test
TEST_LEVEL=integration go test -run TestUserFlow ./tests/...

# Run with timeout (integration tests may be slower)
TEST_LEVEL=integration go test -timeout 5m ./...
```

## Next Steps

- [Unit Testing](unit.md) - Test individual components
- [HTTP Testing](http.md) - Test API endpoints
- [E2E Testing](e2e.md) - Test full user flows
