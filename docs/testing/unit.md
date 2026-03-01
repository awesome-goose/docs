# Unit Testing

Test individual components in isolation.

## What is Unit Testing?

Unit tests verify that individual functions and methods work correctly in isolation from the rest of the system.

## Testing Services

### Basic Service Test

```go
// users.service_test.go
package users

import (
    "testing"
)

func TestUserService_ValidateEmail(t *testing.T) {
    service := &UserService{}

    tests := []struct {
        email    string
        valid    bool
    }{
        {"user@example.com", true},
        {"invalid", false},
        {"", false},
        {"user@", false},
    }

    for _, tc := range tests {
        t.Run(tc.email, func(t *testing.T) {
            result := service.ValidateEmail(tc.email)
            if result != tc.valid {
                t.Errorf("ValidateEmail(%q) = %v, expected %v", tc.email, result, tc.valid)
            }
        })
    }
}
```

### Service with Dependencies

```go
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
    mockDB := &MockDB{
        users: map[string]*User{
            "123": {ID: "123", Email: "test@example.com"},
        },
    }

    service := &UserService{db: mockDB}

    t.Run("existing user", func(t *testing.T) {
        user, err := service.GetByID("123")

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if user.Email != "test@example.com" {
            t.Errorf("expected email test@example.com, got %s", user.Email)
        }
    })

    t.Run("non-existing user", func(t *testing.T) {
        _, err := service.GetByID("999")

        if err == nil {
            t.Error("expected error for non-existing user")
        }
    })
}
```

## Testing Business Logic

### Pure Functions

```go
func TestCalculateDiscount(t *testing.T) {
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

            if result != tc.expected {
                t.Errorf("calculateDiscount(%.2f, %.2f) = %.2f, expected %.2f",
                    tc.price, tc.percent, result, tc.expected)
            }
        })
    }
}
```

### Complex Logic

```go
func TestOrderService_CalculateTotal(t *testing.T) {
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
    expected := 50.00

    if total != expected {
        t.Errorf("expected total %.2f, got %.2f", expected, total)
    }
}
```

## Testing Error Handling

```go
func TestUserService_Create_ValidationErrors(t *testing.T) {
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
                if err == nil {
                    t.Error("expected error but got nil")
                } else if !strings.Contains(err.Error(), tc.errorMsg) {
                    t.Errorf("expected error containing %q, got %q", tc.errorMsg, err.Error())
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

## Testing Edge Cases

```go
func TestUserService_EdgeCases(t *testing.T) {
    service := &UserService{}

    t.Run("nil input", func(t *testing.T) {
        _, err := service.Process(nil)
        if err == nil {
            t.Error("expected error for nil input")
        }
    })

    t.Run("empty slice", func(t *testing.T) {
        result := service.ProcessMany([]User{})
        if len(result) != 0 {
            t.Error("expected empty result for empty input")
        }
    })

    t.Run("boundary values", func(t *testing.T) {
        // Test max int
        result := service.Calculate(math.MaxInt64)
        if result < 0 {
            t.Error("overflow not handled correctly")
        }
    })
}
```

## Testing Structs

### Constructor Tests

```go
func TestNewUser(t *testing.T) {
    user := NewUser("test@example.com", "Test User")

    if user.ID == "" {
        t.Error("expected ID to be generated")
    }

    if user.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", user.Email)
    }

    if user.CreatedAt.IsZero() {
        t.Error("expected CreatedAt to be set")
    }
}
```

### Method Tests

```go
func TestUser_IsActive(t *testing.T) {
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
            user:     User{Active: true, DeletedAt: &time.Time{}},
            expected: false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := tc.user.IsActive()
            if result != tc.expected {
                t.Errorf("expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

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
