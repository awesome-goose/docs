# Unit Testing

Test individual components in isolation.

## Import

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)
```

## What is Unit Testing?

Unit tests verify that individual functions and methods work correctly in isolation from the rest of the system.

## Testing Services

### Basic Service Test

```go
// users.service_test.go
package users

import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserService_ValidateEmail(t *testing.T) {
    test := goosetest.New(t)
    service := &UserService{}

    tests := []struct {
        email string
        valid bool
    }{
        {"user@example.com", true},
        {"invalid", false},
        {"", false},
        {"user@", false},
    }

    for _, tc := range tests {
        t.Run(tc.email, func(t *testing.T) {
            result := service.ValidateEmail(tc.email)
            test.Expect(result).ToEqual(tc.valid)
        })
    }
}
```

### Service with Dependencies Using TestContainer

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

// Mock database interface
type MockDB struct {
    users map[string]*User
}

func (m *MockDB) FindByID(id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func (m *MockDB) Create(user *User) error {
    m.users[user.ID] = user
    return nil
}

func TestUserService_GetByID(t *testing.T) {
    test := goosetest.New(t)
    container := goosetest.NewTestContainer()

    mockDB := &MockDB{
        users: map[string]*User{
            "123": {ID: "123", Email: "test@example.com"},
        },
    }
    container.Register(mockDB)

    service := &UserService{}
    container.Inject(service)

    t.Run("existing user", func(t *testing.T) {
        user, err := service.GetByID("123")
        test.Expect(err).ToBeNil()
        test.Expect(user.Email).ToEqual("test@example.com")
    })

    t.Run("non-existing user", func(t *testing.T) {
        _, err := service.GetByID("999")
        test.Expect(err).Not().ToBeNil()
    })
}
```

### Using ServiceTest Helper

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserServiceWithHelper(t *testing.T) {
    gt := goosetest.NewGooseTest(t)

    gt.WithContainer(func(c *goosetest.TestContainer) {
        c.Register(&MockDB{users: make(map[string]*User)})
    })

    service := &UserService{}
    st := gt.Service(service)

    // Call method and get results
    results := st.Run("Create", CreateUserDTO{Email: "test@example.com"})

    gt.Expect(results[0]).Not().ToBeNil() // user
    gt.Expect(results[1]).ToBeNil()       // error
}
```

## Testing Business Logic

### Pure Functions

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestCalculateDiscount(t *testing.T) {
    test := goosetest.New(t)

    tests := []struct {
        name     string
        price    float64
        percent  float64
        expected float64
    }{
        {"10% off $100", 100.00, 10, 90.00},
        {"25% off $80", 80.00, 25, 60.00},
        {"no discount", 50.00, 0, 50.00},
        {"100% off", 100.00, 100, 0.00},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := calculateDiscount(tc.price, tc.percent)
            test.Expect(result).ToEqual(tc.expected)
        })
    }
}
```

### Complex Logic

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestOrderService_CalculateTotal(t *testing.T) {
    test := goosetest.New(t)
    service := &OrderService{}

    order := &Order{
        Items: []OrderItem{
            {ProductID: "1", Price: 10.00, Quantity: 2},
            {ProductID: "2", Price: 25.00, Quantity: 1},
        },
        Discount:     5.00,
        ShippingCost: 10.00,
    }

    total := service.CalculateTotal(order)

    // Items: (10*2) + (25*1) = 45
    // Discount: -5
    // Shipping: +10
    // Total: 50
    test.Expect(total).ToEqual(50.00)
}
```

## Testing Error Handling

```go
import (
    "strings"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserService_Create_ValidationErrors(t *testing.T) {
    test := goosetest.New(t)
    service := &UserService{}

    tests := []struct {
        name        string
        dto         CreateUserDTO
        expectError bool
        errorMsg    string
    }{
        {
            name:        "empty email",
            dto:         CreateUserDTO{Email: "", Name: "Test"},
            expectError: true,
            errorMsg:    "email is required",
        },
        {
            name:        "invalid email",
            dto:         CreateUserDTO{Email: "invalid", Name: "Test"},
            expectError: true,
            errorMsg:    "invalid email format",
        },
        {
            name:        "empty name",
            dto:         CreateUserDTO{Email: "test@example.com", Name: ""},
            expectError: true,
            errorMsg:    "name is required",
        },
        {
            name:        "valid input",
            dto:         CreateUserDTO{Email: "test@example.com", Name: "Test"},
            expectError: false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            _, err := service.Create(tc.dto)

            if tc.expectError {
                test.Expect(err).Not().ToBeNil()
                test.Expect(strings.Contains(err.Error(), tc.errorMsg)).ToBeTrue()
            } else {
                test.Expect(err).ToBeNil()
            }
        })
    }
}
```

## Testing Edge Cases

```go
import (
    "math"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserService_EdgeCases(t *testing.T) {
    test := goosetest.New(t)
    service := &UserService{}

    t.Run("nil input", func(t *testing.T) {
        _, err := service.Process(nil)
        test.Expect(err).Not().ToBeNil()
    })

    t.Run("empty slice", func(t *testing.T) {
        result := service.ProcessMany([]User{})
        test.Expect(len(result)).ToEqual(0)
    })

    t.Run("boundary values", func(t *testing.T) {
        // Test max int
        result := service.Calculate(math.MaxInt64)
        test.Expect(result).ToBeGreaterThan(0)
    })
}
```

## Testing Controllers

Use `ControllerTest` to test controllers with dependency injection:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserController(t *testing.T) {
    gt := goosetest.NewGooseTest(t)

    // Set up mock dependencies
    gt.WithContainer(func(c *goosetest.TestContainer) {
        c.Register(&MockUserService{})
    })

    controller := &UserController{}
    ct := gt.Controller(controller)

    // Create mock context
    ctx := gt.NewContext()
    ctx.MockRequest().WithMethod("GET").WithPaths("users", "123")

    // Run controller method
    ct.Run("GetUser", ctx)

    // Assert response
    response := ctx.MockResponse()
    gt.Expect(response.StatusCode()).ToEqual(200)
}
```

## Testing Modules

Use `ModuleTest` to verify module structure:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app/users"
)

func TestUsersModule(t *testing.T) {
    gt := goosetest.NewGooseTest(t)
    mt := gt.Module(&users.UsersModule{})

    t.Run("should have valid imports", func(t *testing.T) {
        imports := mt.TestImports()
        gt.Expect(len(imports)).ToBeGreaterThan(0)
    })

    t.Run("should have valid exports", func(t *testing.T) {
        exports := mt.TestExports()
        gt.Expect(exports).Not().ToBeNil()
    })

    t.Run("should have valid declarations", func(t *testing.T) {
        declarations := mt.TestDeclarations()
        gt.Expect(len(declarations)).ToBeGreaterThan(0)
    })
}
```

## Testing Routes

Use `RouteTest` to verify route configuration:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app/users"
)

func TestUsersRoutes(t *testing.T) {
    gt := goosetest.NewGooseTest(t)
    rt := gt.Route(users.ROUTES)

    t.Run("should have GET /users route", func(t *testing.T) {
        rt.TestRouteExists("GET", []string{"users"})
    })

    t.Run("should have GET /users/:id route", func(t *testing.T) {
        rt.TestRouteExists("GET", []string{"users", "123"})
    })

    t.Run("should have POST /users route", func(t *testing.T) {
        rt.TestRouteExists("POST", []string{"users"})
    })

    t.Run("should extract route params", func(t *testing.T) {
        rt.TestRouteParams("GET", []string{"users", "456"}, map[string]string{"id": "456"})
    })

    t.Run("should return not found for invalid route", func(t *testing.T) {
        rt.TestRouteNotFound("DELETE", []string{"invalid"})
    })
}
```

## Testing Structs

### Constructor Tests

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestNewUser(t *testing.T) {
    test := goosetest.New(t)
    user := NewUser("test@example.com", "Test User")

    test.Expect(user.ID).Not().ToEqual("")
    test.Expect(user.Email).ToEqual("test@example.com")
    test.Expect(user.CreatedAt.IsZero()).ToBeFalse()
}
```

### Method Tests

```go
import (
    "testing"
    "time"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUser_IsActive(t *testing.T) {
    test := goosetest.New(t)
    now := time.Now()

    tests := []struct {
        name     string
        user     User
        expected bool
    }{
        {
            name:     "active user",
            user:     User{Active: true, DeletedAt: nil},
            expected: true,
        },
        {
            name:     "inactive user",
            user:     User{Active: false, DeletedAt: nil},
            expected: false,
        },
        {
            name:     "deleted user",
            user:     User{Active: true, DeletedAt: &now},
            expected: false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := tc.user.IsActive()
            test.Expect(result).ToEqual(tc.expected)
        })
    }
}
```

## Using Test Fixtures

Generate random test data with `Fixture`:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestWithRandomData(t *testing.T) {
    test := goosetest.New(t)
    fixture := goosetest.NewFixture()

    // Generate random user data
    user := &User{
        Email:  fixture.Email(),
        Name:   fixture.String(10),
        Age:    fixture.Int(18, 65),
        Active: fixture.Bool(),
    }

    // Test with random data
    err := service.Create(user)
    test.Expect(err).ToBeNil()
}
```

## Next Steps

- [Integration Testing](integration.md) - Test module interactions
- [HTTP Testing](http.md) - Test API endpoints
- [E2E Testing](e2e.md) - Test full user flows

## Testing Concurrency

```go
func TestCounter_ThreadSafe(t *testing.T) {
    counter := NewCounter()

    var wg sync.WaitGroup
    iterations := 1000

    // Increment concurrently
    for i := 0; i < iterations; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()

    if counter.Value() != iterations {
        t.Errorf("expected %d, got %d", iterations, counter.Value())
    }
}
```

## Benchmarks

```go
func BenchmarkUserService_ValidateEmail(b *testing.B) {
    service := &UserService{}
    email := "test@example.com"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.ValidateEmail(email)
    }
}

func BenchmarkCalculateDiscount(b *testing.B) {
    for i := 0; i < b.N; i++ {
        calculateDiscount(100.00, 15)
    }
}

// Run benchmarks: go test -bench=.
```

## Best Practices

1. **Test one thing at a time** - Each test should verify one behavior
2. **Use descriptive names** - `TestUserService_Create_WithInvalidEmail_ReturnsError`
3. **Keep tests fast** - Unit tests should run in milliseconds
4. **Don't test private methods** - Test through public API
5. **Use test helpers** - Reduce duplication with helper functions
6. **Test error paths** - Not just happy paths
7. **Avoid test interdependence** - Each test should be independent

## Next Steps

- [Integration Testing](integration.md) - Test module interactions
- [Mocking](mocking.md) - Create test doubles
- [HTTP Testing](http.md) - Test API endpoints
