# Testing Overview

Write tests for your Goose applications to ensure reliability and maintainability.

## Testing Approach

Goose applications can be tested at multiple levels:

| Type                                | Description               | Speed  | Coverage    |
| ----------------------------------- | ------------------------- | ------ | ----------- |
| [Unit Tests](unit.md)               | Test individual functions | Fast   | Narrow      |
| [Integration Tests](integration.md) | Test module interactions  | Medium | Broad       |
| [HTTP Tests](http.md)               | Test API endpoints        | Medium | API surface |
| [E2E Tests](e2e.md)                 | Test full user flows      | Slow   | Full stack  |

## Quick Start

### Set Up Test File

```go
// app/users/users_test.go
package users

import (
    "testing"
)

func TestUserService_Create(t *testing.T) {
    // Arrange
    service := NewUserService()

    // Act
    user, err := service.Create(CreateUserDTO{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Assert
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if user.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", user.Email)
    }
}
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./app/users/...

# Run specific test
go test -run TestUserService_Create ./app/users/...
```

## Test Structure

### Arrange-Act-Assert Pattern

```go
func TestCalculateTotal(t *testing.T) {
    // Arrange - Set up test data
    items := []Item{
        {Price: 10.00, Quantity: 2},
        {Price: 5.00, Quantity: 3},
    }

    // Act - Execute the code under test
    total := calculateTotal(items)

    // Assert - Verify the result
    expected := 35.00
    if total != expected {
        t.Errorf("expected %.2f, got %.2f", expected, total)
    }
}
```

### Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "user@example.com", true},
        {"no at sign", "userexample.com", false},
        {"no domain", "user@", false},
        {"empty", "", false},
        {"with subdomain", "user@mail.example.com", true},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := validateEmail(tc.email)
            if result != tc.expected {
                t.Errorf("validateEmail(%q) = %v, expected %v",
                    tc.email, result, tc.expected)
            }
        })
    }
}
```

## Testing with Dependencies

### Using Interfaces

```go
// Define interface
type EmailSender interface {
    Send(to, subject, body string) error
}

// Service uses interface
type UserService struct {
    emailSender EmailSender
}

func (s *UserService) Register(dto RegisterDTO) (*User, error) {
    user := &User{Email: dto.Email, Name: dto.Name}
    // Save user...

    s.emailSender.Send(dto.Email, "Welcome!", "Welcome to our app")
    return user, nil
}
```

### Mock Implementation

```go
// Mock for testing
type MockEmailSender struct {
    SentEmails []SentEmail
}

type SentEmail struct {
    To      string
    Subject string
    Body    string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.SentEmails = append(m.SentEmails, SentEmail{to, subject, body})
    return nil
}

// Test
func TestUserService_Register_SendsWelcomeEmail(t *testing.T) {
    mockSender := &MockEmailSender{}
    service := &UserService{emailSender: mockSender}

    service.Register(RegisterDTO{Email: "test@example.com", Name: "Test"})

    if len(mockSender.SentEmails) != 1 {
        t.Errorf("expected 1 email sent, got %d", len(mockSender.SentEmails))
    }

    if mockSender.SentEmails[0].To != "test@example.com" {
        t.Errorf("expected email to test@example.com")
    }
}
```

## Test Helpers

### Setup and Teardown

```go
func TestMain(m *testing.M) {
    // Setup before all tests
    setupTestDatabase()

    // Run tests
    code := m.Run()

    // Teardown after all tests
    teardownTestDatabase()

    os.Exit(code)
}

func setupTestDatabase() {
    // Initialize test database
}

func teardownTestDatabase() {
    // Cleanup test database
}
```

### Per-Test Setup

```go
func setupTest(t *testing.T) (*gorm.DB, func()) {
    db := createTestDB()

    // Return cleanup function
    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }

    return db, cleanup
}

func TestUserService(t *testing.T) {
    db, cleanup := setupTest(t)
    defer cleanup()

    service := &UserService{db: db}
    // Test...
}
```

## Assertions

### Standard Library

```go
func TestExample(t *testing.T) {
    result := 42

    if result != 42 {
        t.Errorf("expected 42, got %d", result)
    }

    if result == 0 {
        t.Fatal("result should not be zero")
    }
}
```

### Custom Assertion Helpers

```go
func assertEqual(t *testing.T, expected, actual interface{}) {
    t.Helper()
    if expected != actual {
        t.Errorf("expected %v, got %v", expected, actual)
    }
}

func assertNil(t *testing.T, value interface{}) {
    t.Helper()
    if value != nil {
        t.Errorf("expected nil, got %v", value)
    }
}

func assertNotNil(t *testing.T, value interface{}) {
    t.Helper()
    if value == nil {
        t.Error("expected non-nil value")
    }
}

// Usage
func TestWithHelpers(t *testing.T) {
    result, err := someFunction()

    assertNil(t, err)
    assertEqual(t, "expected", result)
}
```

## Test Coverage

```bash
# Generate coverage report
go test -cover ./...
```

## Best Practices

1. **Test behavior, not implementation** - Focus on what, not how
2. **Use descriptive test names** - Clearly state what's being tested
3. **Keep tests independent** - Tests shouldn't affect each other
4. **Use table-driven tests** - For multiple input variations
5. **Mock external dependencies** - Don't call real APIs in tests
6. **Test edge cases** - Empty inputs, nulls, boundaries
7. **Run tests in CI** - Automate test execution
8. **Maintain test coverage** - Aim for meaningful coverage

## Test Organization

```
app/
├── users/
│   ├── users.service.go
│   ├── users.service_test.go      # Unit tests
│   ├── users.controller.go
│   └── users.controller_test.go   # Controller tests
├── orders/
│   ├── orders.service.go
│   └── orders.service_test.go
└── integration/
    └── api_test.go                # Integration tests
```

## Next Steps

- [Unit Testing](unit.md) - Test individual components
- [Integration Testing](integration.md) - Test interactions
- [HTTP Testing](http.md) - Test API endpoints
- [Mocking](mocking.md) - Create test doubles
