# Middleware

Middleware intercepts requests **before** they reach your handler, allowing you to add cross-cutting functionality like authentication, logging, or rate limiting.

## The Middleware Contract

```go
type Middleware interface {
    Handle(ctx types.Context) error
}
```

- Return `nil` to let the request continue to the next middleware / the handler.
- Return an `error` to **abort** the request — the platform renders an error response and the handler never runs.
- Middleware runs before the handler, so it can read the request (`ctx.Request()`), set response headers (`ctx.Response().SetHeader(...)`), and stash values for the handler (`ctx.SetValue(...)`).

## Creating Middleware

### Basic Middleware

```go
package middleware

import "github.com/awesome-goose/goose/types"

type LoggingMiddleware struct {
    log types.Log `inject:""`
}

func (m *LoggingMiddleware) Handle(ctx types.Context) error {
    req := ctx.Request()
    m.log.Info("Request started", "method", req.Method().String(), "path", req.Paths())
    return nil // allow the request to proceed
}
```

### Authentication Middleware

Read the header off `ctx.Request()`, stash the user with `ctx.SetValue`, and return an error to reject:

```go
import "errors"

type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    token := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        token = h[0]
    }
    if token == "" {
        return errors.New("unauthorized: no token provided")
    }

    user, err := m.authService.ValidateToken(token)
    if err != nil {
        return errors.New("unauthorized: invalid token")
    }

    // Store the user for the handler to read via a `context` DTO tag
    ctx.SetValue("user", user)

    return nil
}
```

The handler reads that value with a `context`-tagged DTO field:

```go
type ProfileDto struct {
    User *User `context:"user"`
}

func (c *UserController) Profile(dto *ProfileDto) types.Output {
    return output.JSON(dto.User)
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
    ip := clientIP(ctx) // e.g. from the X-Forwarded-For header
    key := "ratelimit:" + ip

    count, _ := cache.GetAs[int](m.cache, key)
    if count >= m.limit {
        return errors.New("too many requests")
    }

    m.cache.Set(key, count+1, m.window)
    return nil
}
```

## Applying Middleware

Middleware is attached to routes in your `router` module.

### Single Route Middleware

The verb helpers accept middleware as trailing arguments:

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "List"}),                     // no middleware
    router.Post("/users", []any{UserController{}, "Create"}, &AuthMiddleware{}), // protected
)
```

### Route Group Middleware

Attach middleware to a parent route; it applies to every child:

```go
var ROUTES = router.ForRoutes(
    // Public
    router.Get("/public", []any{AppController{}, "Public"}),

    // Protected group
    types.Route{
        Path:        "/api",
        Middlewares: types.Middlewares{&AuthMiddleware{}},
        Children: types.Routes{
            router.Get("/users", []any{AppController{}, "ListUsers"}),
            router.Post("/users", []any{AppController{}, "CreateUser"}),
        },
    },

    // Admin group with two middlewares
    types.Route{
        Path:        "/admin",
        Middlewares: types.Middlewares{&AuthMiddleware{}, &AdminMiddleware{}},
        Children:    adminRoutes,
    },
)
```

### Global Middleware

Wrap all routes under a root path:

```go
var ROUTES = router.ForRoutes(
    types.Route{
        Path:        "/",
        Middlewares: types.Middlewares{&LoggingMiddleware{}},
        Children:    allRoutes,
    },
)
```

## Middleware Chain

Middleware executes in order — a parent's middleware runs before its children's:

```go
Middlewares: types.Middlewares{
    &LoggingMiddleware{},   // 1st
    &AuthMiddleware{},      // 2nd
    &RateLimitMiddleware{}, // 3rd
}

// Execution order:
// Request → Logging → Auth → RateLimit → Handler
```

## Middleware Patterns

### Request Validation

```go
type ValidationMiddleware struct{}

func (m *ValidationMiddleware) Handle(ctx types.Context) error {
    req := ctx.Request()
    contentType := ""
    if h := req.Headers()["Content-Type"]; len(h) > 0 {
        contentType = h[0]
    }
    if req.Method().Is("POST") && contentType != "application/json" {
        return errors.New("unsupported media type")
    }
    return nil
}
```

### CORS Middleware

Set headers through `ctx.Response()`:

```go
type CORSMiddleware struct {
    allowedOrigins []string
}

func (m *CORSMiddleware) Handle(ctx types.Context) error {
    origin := ""
    if h := ctx.Request().Headers()["Origin"]; len(h) > 0 {
        origin = h[0]
    }

    resp := ctx.Response()
    for _, allowed := range m.allowedOrigins {
        if origin == allowed || allowed == "*" {
            resp.SetHeader("Access-Control-Allow-Origin", origin)
            resp.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            resp.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
            break
        }
    }
    return nil
}
```

### Request ID Middleware

```go
type RequestIDMiddleware struct{}

func (m *RequestIDMiddleware) Handle(ctx types.Context) error {
    requestID := ""
    if h := ctx.Request().Headers()["X-Request-ID"]; len(h) > 0 {
        requestID = h[0]
    }
    if requestID == "" {
        requestID = uuid.New().String()
    }

    ctx.SetValue("requestID", requestID)
    ctx.Response().SetHeader("X-Request-ID", requestID)

    return nil
}
```

## Middleware with Dependency Injection

Middleware can receive injected dependencies just like controllers:

```go
type MetricsMiddleware struct {
    metrics *MetricsService `inject:""`
    log     types.Log       `inject:""`
}

func (m *MetricsMiddleware) Handle(ctx types.Context) error {
    req := ctx.Request()
    m.metrics.CountRequest(req.Method().String(), req.Paths())
    m.log.Info("Request received", "method", req.Method().String(), "path", req.Paths())
    return nil
}
```

> Middleware runs **before** the handler and returns immediately, so a `defer` inside `Handle` fires before the handler executes — it can't measure handler duration. To time the full request, use the platform's request timeout or an outer HTTP layer.

## Error Handling in Middleware

Returning an error terminates the request:

```go
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    token := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        token = h[0]
    }
    if token == "" {
        return errors.New("unauthorized") // request stops here
    }
    return nil // request continues
}
```

## Best Practices

### 1. Keep Middleware Focused

Each middleware should do one thing:

```go
// ✅ Good: single responsibility
type AuthMiddleware struct{}      // Only authentication
type LoggingMiddleware struct{}   // Only logging
type RateLimitMiddleware struct{} // Only rate limiting

// ❌ Bad: too many responsibilities
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

Avoid heavy synchronous work in middleware — it runs on every matching request.

### 4. Pass Data via Context Values

Use `ctx.SetValue(key, value)` in middleware and read it in the handler through a `context`-tagged DTO field, rather than re-computing it.

## Next Steps

- [Routing](routing.md) - Route configuration
- [Controllers](controllers.md) - Request handlers
- [Error Handling](error-handling.md) - Error management
