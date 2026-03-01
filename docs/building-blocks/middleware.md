# Middleware

Middleware intercepts requests before they reach your handlers, allowing you to add cross-cutting functionality like authentication, logging, or rate limiting.

## Creating Middleware

### Basic Middleware

```go
package middleware

import "github.com/awesome-goose/goose/types"

type LoggingMiddleware struct {
    log types.Log `inject:""`
}

func (m *LoggingMiddleware) Handle(ctx types.Context) error {
    // Before request
    m.log.Info("Request started", ctx.Method(), ctx.Path())

    // Continue to next middleware/handler
    // Returning nil allows the request to proceed
    return nil
}
```

### Authentication Middleware

```go
type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    token := ctx.Header("Authorization")
    if token == "" {
        ctx.SetStatus(401)
        return fmt.Errorf("unauthorized: no token provided")
    }

    user, err := m.authService.ValidateToken(token)
    if err != nil {
        ctx.SetStatus(401)
        return fmt.Errorf("unauthorized: invalid token")
    }

    // Store user in context for handlers
    ctx.Set("user", user)

    return nil  // Allow request to proceed
}
```

### Rate Limiting Middleware

```go
type RateLimitMiddleware struct {
    cache  *cache.Cache `inject:""`
    limit  int
    window time.Duration
}

func (m *RateLimitMiddleware) Handle(ctx types.Context) error {
    ip := ctx.ClientIP()
    key := "ratelimit:" + ip

    count, _ := cache.GetAs[int](m.cache, key)
    if count >= m.limit {
        ctx.SetStatus(429)
        return fmt.Errorf("too many requests")
    }

    m.cache.Set(key, count+1, m.window)
    return nil
}
```

## Applying Middleware

### Global Middleware

Apply to all routes in the application:

```go
func (c *AppController) Routes() types.Routes {
    return types.Routes{
        {
            Path:        "/",
            Middlewares: types.Middlewares{&LoggingMiddleware{}},
            Children:    allRoutes,
        },
    }
}
```

### Route Group Middleware

Apply to specific route groups:

```go
func (c *AppController) Routes() types.Routes {
    return types.Routes{
        // Public routes
        {Method: "GET", Path: "/public", Handler: c.Public},

        // Protected routes
        {
            Path:        "/api",
            Middlewares: types.Middlewares{&AuthMiddleware{}},
            Children: types.Routes{
                {Method: "GET", Path: "/users", Handler: c.ListUsers},
                {Method: "POST", Path: "/users", Handler: c.CreateUser},
            },
        },

        // Admin routes
        {
            Path:        "/admin",
            Middlewares: types.Middlewares{
                &AuthMiddleware{},
                &AdminMiddleware{},
            },
            Children: adminRoutes,
        },
    }
}
```

### Single Route Middleware

Apply to specific routes:

```go
func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.List},  // No middleware
        {
            Method:      "POST",
            Path:        "/users",
            Handler:     c.Create,
            Middlewares: types.Middlewares{&AuthMiddleware{}},
        },
    }
}
```

## Middleware Chain

Middleware executes in order:

```go
Middlewares: types.Middlewares{
    &LoggingMiddleware{},   // 1st
    &AuthMiddleware{},       // 2nd
    &RateLimitMiddleware{},  // 3rd
}

// Execution order:
// Request → Logging → Auth → RateLimit → Handler
```

## Middleware Patterns

### Request Validation

```go
type ValidationMiddleware struct{}

func (m *ValidationMiddleware) Handle(ctx types.Context) error {
    contentType := ctx.Header("Content-Type")
    if ctx.Method() == "POST" && contentType != "application/json" {
        ctx.SetStatus(415)
        return fmt.Errorf("unsupported media type")
    }
    return nil
}
```

### CORS Middleware

```go
type CORSMiddleware struct {
    allowedOrigins []string
}

func (m *CORSMiddleware) Handle(ctx types.Context) error {
    origin := ctx.Header("Origin")

    for _, allowed := range m.allowedOrigins {
        if origin == allowed || allowed == "*" {
            ctx.SetHeader("Access-Control-Allow-Origin", origin)
            ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
            break
        }
    }

    // Handle preflight
    if ctx.Method() == "OPTIONS" {
        ctx.SetStatus(204)
        return nil
    }

    return nil
}
```

### Request ID Middleware

```go
type RequestIDMiddleware struct{}

func (m *RequestIDMiddleware) Handle(ctx types.Context) error {
    requestID := ctx.Header("X-Request-ID")
    if requestID == "" {
        requestID = uuid.New().String()
    }

    ctx.Set("requestID", requestID)
    ctx.SetHeader("X-Request-ID", requestID)

    return nil
}
```

### Timeout Middleware

```go
type TimeoutMiddleware struct {
    timeout time.Duration
}

func (m *TimeoutMiddleware) Handle(ctx types.Context) error {
    // Set deadline
    ctx.SetDeadline(time.Now().Add(m.timeout))
    return nil
}
```

## Middleware with Dependency Injection

Middleware can receive injected dependencies:

```go
type MetricsMiddleware struct {
    metrics *MetricsService `inject:""`
    log     types.Log       `inject:""`
}

func (m *MetricsMiddleware) Handle(ctx types.Context) error {
    start := time.Now()

    // After handler completes (using defer-like pattern)
    defer func() {
        duration := time.Since(start)
        m.metrics.RecordRequest(ctx.Method(), ctx.Path(), duration)
        m.log.Info("Request completed",
            "method", ctx.Method(),
            "path", ctx.Path(),
            "duration", duration)
    }()

    return nil
}
```

## Error Handling in Middleware

When middleware returns an error, the request is terminated:

```go
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    token := ctx.Header("Authorization")
    if token == "" {
        ctx.SetStatus(401)
        return fmt.Errorf("unauthorized")  // Request stops here
    }

    return nil  // Request continues
}
```

## Best Practices

### 1. Keep Middleware Focused

Each middleware should do one thing:

```go
// ✅ Good: Single responsibility
type AuthMiddleware struct{}      // Only authentication
type LoggingMiddleware struct{}   // Only logging
type RateLimitMiddleware struct{} // Only rate limiting

// ❌ Bad: Too many responsibilities
type AllInOneMiddleware struct{}  // Auth + Logging + RateLimit
```

### 2. Order Matters

Place middleware in logical order:

```go
Middlewares: types.Middlewares{
    &RequestIDMiddleware{},  // 1. Assign ID first
    &LoggingMiddleware{},    // 2. Log with ID
    &RateLimitMiddleware{},  // 3. Rate limit before auth
    &AuthMiddleware{},       // 4. Authenticate
    &ValidationMiddleware{}, // 5. Validate input
}
```

### 3. Don't Block

Avoid heavy operations in middleware:

```go
// ❌ Bad: Heavy operation blocks all requests
func (m *HeavyMiddleware) Handle(ctx types.Context) error {
    result := someHeavyComputation()  // Blocks
    return nil
}

// ✅ Good: Defer heavy work or use async
func (m *AsyncMiddleware) Handle(ctx types.Context) error {
    go someHeavyComputation()  // Non-blocking
    return nil
}
```

### 4. Handle Errors Gracefully

Set appropriate status codes:

```go
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    if !authenticated {
        ctx.SetStatus(401)  // Set status before error
        return fmt.Errorf("unauthorized")
    }
    return nil
}
```

## Next Steps

- [Routing](routing.md) - Route configuration
- [Controllers](controllers.md) - Request handlers
- [Error Handling](error-handling.md) - Error management
