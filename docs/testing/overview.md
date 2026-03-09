# Testing Overview

Write tests for your Goose applications to ensure reliability and maintainability.

## Import

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)
```

## Testing Approach

Goose applications can be tested at multiple levels:

| Type                                | Description               | Speed  | Coverage    |
| ----------------------------------- | ------------------------- | ------ | ----------- |
| [Unit Tests](unit.md)               | Test individual functions | Fast   | Narrow      |
| [Integration Tests](integration.md) | Test module interactions  | Medium | Broad       |
| [HTTP Tests](http.md)               | Test API endpoints        | Medium | API surface |
| [E2E Tests](e2e.md)                 | Test full user flows      | Slow   | Full stack  |

## Quick Start

### Set Up Test File with Goose Testing

```go
// app/users/users_test.go
package users

import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUserService_Create(t *testing.T) {
    // Create test helper with fluent assertions
    test := goosetest.New(t)

    // Arrange
    service := NewUserService()

    // Act
    user, err := service.Create(CreateUserDTO{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Assert with fluent API
    test.Expect(err).ToBeNil()
    test.Expect(user.Email).ToEqual("test@example.com")
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
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestCalculateTotal(t *testing.T) {
    test := goosetest.New(t)

    // Arrange - Set up test data
    items := []Item{
        {Price: 10.00, Quantity: 2},
        {Price: 5.00, Quantity: 3},
    }

    // Act - Execute the code under test
    total := calculateTotal(items)

    // Assert - Verify the result with fluent assertions
    test.Expect(total).ToEqual(35.00)
}
```

### Table-Driven Tests

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestValidateEmail(t *testing.T) {
    test := goosetest.New(t)

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
            test.Expect(result).ToEqual(tc.expected)
        })
    }
}
```

## Test Suite Management

Use `SuiteRunner` for organized test suites with setup/teardown:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

type UserServiceSuite struct {
    goosetest.Suite
    service *UserService
    db      *MockDB
}

func (s *UserServiceSuite) SetupSuite() {
    // Runs once before all tests
    s.db = NewMockDB()
}

func (s *UserServiceSuite) SetupTest() {
    // Runs before each test
    s.service = &UserService{db: s.db}
    s.db.Clear()
}

func (s *UserServiceSuite) TeardownTest() {
    // Runs after each test
    s.db.Clear()
}

func (s *UserServiceSuite) TestCreate() {
    user, err := s.service.Create(CreateUserDTO{Email: "test@example.com"})
    s.T.Expect(err).ToBeNil()
    s.T.Expect(user.Email).ToEqual("test@example.com")
}

func (s *UserServiceSuite) TestGetByID() {
    s.db.Insert(&User{ID: "123", Email: "test@example.com"})
    user, err := s.service.GetByID("123")
    s.T.Expect(err).ToBeNil()
    s.T.Expect(user.ID).ToEqual("123")
}

func TestUserServiceSuite(t *testing.T) {
    runner := goosetest.NewSuiteRunner(t, &UserServiceSuite{})
    runner.Run()
}
```

## Testing with Dependencies

### Using the Test Container

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

// Define interface
type EmailSender interface {
    Send(to, subject, body string) error
}

// Service uses interface
type UserService struct {
    EmailSender EmailSender
}

func (s *UserService) Register(dto RegisterDTO) (*User, error) {
    user := &User{Email: dto.Email, Name: dto.Name}
    s.EmailSender.Send(dto.Email, "Welcome!", "Welcome to our app")
    return user, nil
}

// Test with TestContainer for dependency injection
func TestUserService_Register(t *testing.T) {
    test := goosetest.New(t)
    container := goosetest.NewTestContainer()

    // Register mock dependency
    mockSender := &MockEmailSender{}
    container.Register(mockSender)

    // Create service and inject dependencies
    service := &UserService{}
    container.Inject(service)

    // Test
    user, err := service.Register(RegisterDTO{Email: "test@example.com", Name: "Test"})
    test.Expect(err).ToBeNil()
    test.Expect(user.Email).ToEqual("test@example.com")
}
```

### Using the Mock Helper

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestWithMock(t *testing.T) {
    mock := goosetest.NewMock(t)

    // Set up return values
    mock.On("GetUser", &User{ID: "123", Email: "test@example.com"}, nil)

    // Set up expectations
    mock.ExpectCall("GetUser", 1, "123")

    // Execute code that calls mock
    result := mock.Called("GetUser", "123")

    // Verify expectations were met
    mock.Verify()
}
```

## Test Fixtures

Generate random test data with `Fixture`:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestWithFixtures(t *testing.T) {
    fixture := goosetest.NewFixture()

    // Generate random data
    email := fixture.Email()           // e.g., "aBcDeFgHiJ@kLmNo.com"
    name := fixture.String(10)         // Random 10-char string
    age := fixture.Int(18, 65)         // Random int between 18-65
    active := fixture.Bool()           // Random boolean

    user := &User{
        Email:  email,
        Name:   name,
        Age:    age,
        Active: active,
    }
    // Use user in tests...
}
```

## Test Tagging

Control which tests run based on environment:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestUnitOnly(t *testing.T) {
    goosetest.SkipUnlessUnit(t)
    // This test only runs for unit tests (default)
}

func TestIntegrationOnly(t *testing.T) {
    goosetest.SkipUnlessIntegration(t)
    // Run with: TEST_LEVEL=integration go test ./...
}

func TestE2EOnly(t *testing.T) {
    goosetest.SkipUnlessE2E(t)
    // Run with: TEST_LEVEL=e2e go test ./...
}

func TestSkipInCI(t *testing.T) {
    goosetest.SkipInCI(t)
    // Skipped when CI=true environment variable is set
}
```

## Test Helpers

### Setup and Teardown

```go
import (
    "os"
    "testing"
)

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

## Fluent Assertions Reference

The `goosetest.T` wrapper provides fluent assertions:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestAssertions(t *testing.T) {
    test := goosetest.New(t)

    // Equality
    test.Expect(42).ToEqual(42)
    test.Expect("hello").ToBe("hello")
    test.Expect(obj1).ToDeepEqual(obj2)

    // Nil checks
    test.Expect(nil).ToBeNil()
    test.Expect(value).Not().ToBeNil()

    // Negation
    test.Expect(10).Not().ToEqual(20)

    // Boolean
    test.Expect(true).ToBeTrue()
    test.Expect(false).ToBeFalse()

    // Numeric comparisons
    test.Expect(10).ToBeGreaterThan(5)
    test.Expect(5).ToBeLessThan(10)

    // String assertions
    test.Expect("hello world").ToContain("world")
    test.Expect("hello").ToStartWith("he")
    test.Expect("hello").ToEndWith("lo")

    // Collections
    test.Expect(len(slice)).ToBeGreaterThan(0)
    test.Expect(slice).ToHaveLength(3)
}
```

## Next Steps

- [Unit Testing](unit.md) - Test individual components
- [Integration Testing](integration.md) - Test module interactions
- [HTTP Testing](http.md) - Test API endpoints
- [E2E Testing](e2e.md) - Test full user flows

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
