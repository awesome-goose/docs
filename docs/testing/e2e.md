# End-to-End Testing

Test complete user flows from start to finish.

## Import

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)
```

## What is E2E Testing?

End-to-end (E2E) tests verify complete user workflows through your entire application stack. They simulate real user behavior and test the full system including:

- Frontend interactions (if applicable)
- API requests and responses
- Database operations
- External service integrations
- Authentication and authorization flows

## When to Write E2E Tests

- Critical user journeys (signup, checkout, etc.)
- Complex multi-step workflows
- Scenarios that span multiple services
- Regression testing for production bugs
- Smoke tests for deployments

## Skipping E2E Tests

E2E tests are typically slower and require more setup. Use test tags to run them only when needed:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestE2E_UserSignupFlow(t *testing.T) {
    // Skip unless TEST_LEVEL=e2e
    goosetest.SkipUnlessE2E(t)

    // Your E2E test code...
}
```

Run E2E tests with:

```bash
TEST_LEVEL=e2e go test ./...
```

## Testing Complete User Flows

### User Registration and Login

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
    "myapp/app"
)

func TestE2E_UserRegistrationFlow(t *testing.T) {
    goosetest.SkipUnlessE2E(t)

    // Start test server with full application
    handler := createFullHandler(&app.AppModule{})
    ht := goosetest.NewHTTPTest(t, handler).WithServer()
    defer ht.Close()

    // Step 1: Register new user
    var registerResp map[string]any
    ht.POST("/api/auth/register").
        WithJSON(map[string]string{
            "email":    "newuser@example.com",
            "password": "SecurePassword123!",
            "name":     "New User",
        }).
        Do().
        ExpectCreated().
        JSON(&registerResp)

    userID := registerResp["user"].(map[string]any)["id"].(string)
    ht.T.Expect(userID).Not().ToEqual("")

    // Step 2: User should receive verification email (check side effect)
    // In real tests, you might check a test email service

    // Step 3: Verify email
    verificationToken := registerResp["verification_token"].(string)
    ht.POST("/api/auth/verify-email").
        WithJSON(map[string]string{
            "token": verificationToken,
        }).
        Do().
        ExpectOK()

    // Step 4: Login with credentials
    var loginResp map[string]any
    ht.POST("/api/auth/login").
        WithJSON(map[string]string{
            "email":    "newuser@example.com",
            "password": "SecurePassword123!",
        }).
        Do().
        ExpectOK().
        JSON(&loginResp)

    accessToken := loginResp["access_token"].(string)
    ht.T.Expect(accessToken).Not().ToEqual("")

    // Step 5: Access protected resource with token
    var profile map[string]any
    ht.GET("/api/users/me").
        WithBearer(accessToken).
        Do().
        ExpectOK().
        JSON(&profile)

    ht.T.Expect(profile["email"]).ToEqual("newuser@example.com")
}
```

### E-Commerce Checkout Flow

```go
func TestE2E_CheckoutFlow(t *testing.T) {
    goosetest.SkipUnlessE2E(t)

    handler := createFullHandler(&app.AppModule{})
    ht := goosetest.NewHTTPTest(t, handler).WithServer()
    defer ht.Close()

    // Setup: Login as existing user
    var loginResp map[string]any
    ht.POST("/api/auth/login").
        WithJSON(map[string]string{
            "email":    "buyer@example.com",
            "password": "password123",
        }).
        Do().
        ExpectOK().
        JSON(&loginResp)

    token := loginResp["access_token"].(string)

    // Step 1: Browse products
    var products []map[string]any
    ht.GET("/api/products").
        WithQuery("category", "electronics").
        Do().
        ExpectOK().
        JSON(&products)

    ht.T.Expect(len(products)).ToBeGreaterThan(0)
    productID := products[0]["id"].(string)

    // Step 2: Add to cart
    ht.POST("/api/cart/items").
        WithBearer(token).
        WithJSON(map[string]any{
            "product_id": productID,
            "quantity":   2,
        }).
        Do().
        ExpectCreated()

    // Step 3: View cart
    var cart map[string]any
    ht.GET("/api/cart").
        WithBearer(token).
        Do().
        ExpectOK().
        JSON(&cart)

    ht.T.Expect(len(cart["items"].([]any))).ToEqual(1)

    // Step 4: Add shipping address
    ht.POST("/api/checkout/shipping").
        WithBearer(token).
        WithJSON(map[string]string{
            "street":  "123 Main St",
            "city":    "New York",
            "state":   "NY",
            "zip":     "10001",
            "country": "USA",
        }).
        Do().
        ExpectOK()

    // Step 5: Select payment method
    ht.POST("/api/checkout/payment").
        WithBearer(token).
        WithJSON(map[string]string{
            "card_id": "card-123",
        }).
        Do().
        ExpectOK()

    // Step 6: Place order
    var order map[string]any
    ht.POST("/api/checkout/complete").
        WithBearer(token).
        Do().
        ExpectCreated().
        JSON(&order)

    orderID := order["id"].(string)
    ht.T.Expect(orderID).Not().ToEqual("")
    ht.T.Expect(order["status"]).ToEqual("pending")

    // Step 7: Verify order in history
    var orders []map[string]any
    ht.GET("/api/users/me/orders").
        WithBearer(token).
        Do().
        ExpectOK().
        JSON(&orders)

    ht.T.Expect(len(orders)).ToBeGreaterThan(0)
    ht.T.Expect(orders[0]["id"]).ToEqual(orderID)

    // Step 8: Verify cart is cleared
    ht.GET("/api/cart").
        WithBearer(token).
        Do().
        ExpectOK().
        JSON(&cart)

    ht.T.Expect(len(cart["items"].([]any))).ToEqual(0)
}
```

### Multi-User Interaction Flow

```go
func TestE2E_CommentingFlow(t *testing.T) {
    goosetest.SkipUnlessE2E(t)

    handler := createFullHandler(&app.AppModule{})
    ht := goosetest.NewHTTPTest(t, handler).WithServer()
    defer ht.Close()

    // User 1: Create a post
    var user1Login map[string]any
    ht.POST("/api/auth/login").
        WithJSON(map[string]string{"email": "user1@example.com", "password": "pass1"}).
        Do().
        ExpectOK().
        JSON(&user1Login)
    token1 := user1Login["access_token"].(string)

    var post map[string]any
    ht.POST("/api/posts").
        WithBearer(token1).
        WithJSON(map[string]string{
            "title":   "My First Post",
            "content": "Hello world!",
        }).
        Do().
        ExpectCreated().
        JSON(&post)

    postID := post["id"].(string)

    // User 2: Comment on the post
    var user2Login map[string]any
    ht.POST("/api/auth/login").
        WithJSON(map[string]string{"email": "user2@example.com", "password": "pass2"}).
        Do().
        ExpectOK().
        JSON(&user2Login)
    token2 := user2Login["access_token"].(string)

    ht.POST("/api/posts/" + postID + "/comments").
        WithBearer(token2).
        WithJSON(map[string]string{
            "content": "Great post!",
        }).
        Do().
        ExpectCreated()

    // User 1: Should see the comment
    var postWithComments map[string]any
    ht.GET("/api/posts/" + postID).
        WithBearer(token1).
        Do().
        ExpectOK().
        JSON(&postWithComments)

    comments := postWithComments["comments"].([]any)
    ht.T.Expect(len(comments)).ToEqual(1)

    // User 1: Should receive notification (check notification endpoint)
    var notifications []map[string]any
    ht.GET("/api/notifications").
        WithBearer(token1).
        Do().
        ExpectOK().
        JSON(&notifications)

    ht.T.Expect(len(notifications)).ToBeGreaterThan(0)
    ht.T.Expect(notifications[0]["type"]).ToEqual("comment")
}
```

## Testing CLI Applications

### E2E CLI Test

```go
import (
    "bytes"
    "os/exec"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestE2E_CLIWorkflow(t *testing.T) {
    goosetest.SkipUnlessE2E(t)
    test := goosetest.New(t)

    // Step 1: Initialize project
    cmd := exec.Command("./myapp", "init", "--name", "test-project")
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    test.Expect(err).ToBeNil()
    test.Expect(stdout.String()).ToContain("Project initialized")

    // Step 2: Add a module
    cmd = exec.Command("./myapp", "generate", "module", "users")
    cmd.Dir = "test-project"
    stdout.Reset()
    cmd.Stdout = &stdout

    err = cmd.Run()
    test.Expect(err).ToBeNil()
    test.Expect(stdout.String()).ToContain("Module created")

    // Step 3: Build the project
    cmd = exec.Command("./myapp", "build")
    cmd.Dir = "test-project"
    err = cmd.Run()
    test.Expect(err).ToBeNil()

    // Cleanup
    exec.Command("rm", "-rf", "test-project").Run()
}
```

## Test Data Management

### Seeding Test Data

```go
func seedTestData(db *sql.DB) {
    // Create test users
    db.Exec(`
        INSERT INTO users (id, email, name, password_hash) VALUES
        ('user-1', 'buyer@example.com', 'Buyer', '$hash'),
        ('user-2', 'seller@example.com', 'Seller', '$hash')
    `)

    // Create test products
    db.Exec(`
        INSERT INTO products (id, name, price, stock) VALUES
        ('prod-1', 'Laptop', 999.99, 10),
        ('prod-2', 'Phone', 599.99, 20)
    `)
}

func cleanupTestData(db *sql.DB) {
    db.Exec("TRUNCATE users, products, orders, cart_items CASCADE")
}

func TestE2E_WithSeededData(t *testing.T) {
    goosetest.SkipUnlessE2E(t)

    db := getTestDB()
    seedTestData(db)
    t.Cleanup(func() {
        cleanupTestData(db)
    })

    // Run E2E tests with seeded data...
}
```

## Parallel E2E Tests

### Isolating Parallel Tests

```go
func TestE2E_ParallelFlows(t *testing.T) {
    goosetest.SkipUnlessE2E(t)

    t.Run("UserA_Checkout", func(t *testing.T) {
        t.Parallel()

        // Use unique test data
        email := fmt.Sprintf("usera_%d@test.com", time.Now().UnixNano())
        // ... test checkout flow for User A
    })

    t.Run("UserB_Checkout", func(t *testing.T) {
        t.Parallel()

        // Use unique test data
        email := fmt.Sprintf("userb_%d@test.com", time.Now().UnixNano())
        // ... test checkout flow for User B
    })
}
```

## Test Environment Setup

### Docker Compose for E2E

```yaml
# docker-compose.test.yml
version: "3.8"
services:
  app:
    build: .
    environment:
      - DATABASE_URL=postgres://test:test@db:5432/testdb
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    environment:
      - POSTGRES_USER=test
      - POSTGRES_PASSWORD=test
      - POSTGRES_DB=testdb

  redis:
    image: redis:7
```

```go
func TestMain(m *testing.M) {
    // Start test environment
    cmd := exec.Command("docker-compose", "-f", "docker-compose.test.yml", "up", "-d")
    cmd.Run()

    // Wait for services to be ready
    time.Sleep(5 * time.Second)

    code := m.Run()

    // Cleanup
    exec.Command("docker-compose", "-f", "docker-compose.test.yml", "down", "-v").Run()

    os.Exit(code)
}
```

## Best Practices

1. **Use unique test data** - Generate unique emails, IDs to avoid conflicts
2. **Clean up after tests** - Reset database state after E2E tests
3. **Test complete flows** - Don't skip steps, test the full user journey
4. **Handle async operations** - Add appropriate waits for background jobs
5. **Use realistic data** - Test with data similar to production
6. **Monitor test duration** - E2E tests should have timeout limits
7. **Run in isolation** - Use separate databases for E2E tests

## Running E2E Tests

```bash
# Run only E2E tests
TEST_LEVEL=e2e go test ./...

# Run with longer timeout
TEST_LEVEL=e2e go test -timeout 10m ./...

# Run specific E2E test
TEST_LEVEL=e2e go test -run TestE2E_CheckoutFlow ./tests/...

# Run with verbose output
TEST_LEVEL=e2e go test -v ./...

# Skip CI environment (for local-only tests)
TEST_LEVEL=e2e go test ./... # CI env var is checked by SkipInCI
```

## CI/CD Integration

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  e2e:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Run E2E Tests
        env:
          TEST_LEVEL: e2e
          DATABASE_URL: postgres://postgres:test@localhost:5432/postgres
        run: go test -timeout 10m ./...
```

## Next Steps

- [Unit Testing](unit.md) - Test individual components
- [Integration Testing](integration.md) - Test module interactions
- [HTTP Testing](http.md) - Test API endpoints
