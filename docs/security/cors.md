# CORS Configuration

Configure Cross-Origin Resource Sharing (CORS) for your Goose API.

## Overview

CORS controls which origins can access your API from browsers. Without proper CORS configuration, browsers block cross-origin requests.

## Quick Setup

### Basic CORS Middleware

Middleware sets response headers through `ctx.Response()` and returns `nil` to continue. It cannot itself write a preflight response, so pair it with an explicit `OPTIONS` route (see [Preflight Requests](#preflight-requests)).

```go
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(ctx types.Context) error {
    resp := ctx.Response()
    resp.SetHeader("Access-Control-Allow-Origin", "*")
    resp.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    resp.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
    return nil
}
```

### Apply Globally

Attach the middleware to a root group so it runs on every route, and add an OPTIONS handler for preflight:

```go
var ROUTES = router.ForRoutes(
    types.Route{
        Path:        "/",
        Middlewares: types.Middlewares{&CORSMiddleware{}},
        Children: types.Routes{
            router.Get("/", []any{AppController{}, "Index"}),
            // Preflight — answers any OPTIONS request with 204
            {Method: "OPTIONS", Path: "/*", Handler: []any{AppController{}, "Preflight"}},
        },
    },
)

func (c *AppController) Preflight(dto *EmptyDto) types.Output {
    return output.NoContent()
}
```

## Configurable CORS

### Configuration Struct

```go
type CORSConfig struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    AllowCredentials bool
    MaxAge           int
}

func DefaultCORSConfig() *CORSConfig {
    return &CORSConfig{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
        ExposedHeaders:   []string{},
        AllowCredentials: false,
        MaxAge:           86400, // 24 hours
    }
}
```

### Configurable Middleware

```go
type CORSMiddleware struct {
    config *CORSConfig
}

func NewCORSMiddleware(config *CORSConfig) *CORSMiddleware {
    if config == nil {
        config = DefaultCORSConfig()
    }
    return &CORSMiddleware{config: config}
}

func (m *CORSMiddleware) Handle(ctx types.Context) error {
    origin := ""
    if h := ctx.Request().Headers()["Origin"]; len(h) > 0 {
        origin = h[0]
    }

    // Check if origin is allowed
    if !m.isOriginAllowed(origin) {
        return nil
    }

    // Set CORS headers on the response
    resp := ctx.Response()
    resp.SetHeader("Access-Control-Allow-Origin", origin)
    resp.SetHeader("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
    resp.SetHeader("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))

    if len(m.config.ExposedHeaders) > 0 {
        resp.SetHeader("Access-Control-Expose-Headers", strings.Join(m.config.ExposedHeaders, ", "))
    }

    if m.config.AllowCredentials {
        resp.SetHeader("Access-Control-Allow-Credentials", "true")
    }

    if m.config.MaxAge > 0 {
        resp.SetHeader("Access-Control-Max-Age", strconv.Itoa(m.config.MaxAge))
    }

    return nil
}

func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
    if len(m.config.AllowedOrigins) == 0 {
        return false
    }

    for _, allowed := range m.config.AllowedOrigins {
        if allowed == "*" || allowed == origin {
            return true
        }
    }

    return false
}
```

## Common Configurations

### Development (Allow All)

```go
cors := NewCORSMiddleware(&CORSConfig{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"*"},
})
```

### Production (Specific Origins)

```go
cors := NewCORSMiddleware(&CORSConfig{
    AllowedOrigins:   []string{
        "https://myapp.com",
        "https://www.myapp.com",
        "https://admin.myapp.com",
    },
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           86400,
})
```

### With Credentials

When using cookies or HTTP authentication:

```go
cors := NewCORSMiddleware(&CORSConfig{
    AllowedOrigins:   []string{"https://myapp.com"}, // Cannot use "*" with credentials
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
    AllowCredentials: true,
    MaxAge:           3600,
})
```

### Expose Custom Headers

Allow client to read custom response headers:

```go
cors := NewCORSMiddleware(&CORSConfig{
    AllowedOrigins: []string{"*"},
    ExposedHeaders: []string{
        "X-Request-ID",
        "X-Total-Count",
        "X-Page",
        "X-Per-Page",
    },
})
```

## Environment-Based Configuration

```go
// e is a *env.Env (e.g. env.NewEnv()) obtained from injection.
func LoadCORSConfig(e *env.Env) *CORSConfig {
    appEnv := e.GetWithDefault("APP_ENV", "development")

    if appEnv == "production" {
        origins := strings.Split(e.Get("CORS_ORIGINS"), ",")
        return &CORSConfig{
            AllowedOrigins:   origins,
            AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
            AllowedHeaders:   []string{"Content-Type", "Authorization"},
            AllowCredentials: true,
            MaxAge:           86400,
        }
    }

    // Development: allow all
    return &CORSConfig{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
    }
}
```

**.env for production:**

```env
CORS_ORIGINS=https://myapp.com,https://admin.myapp.com
```

## Preflight Requests

Browsers send an `OPTIONS` request before certain cross-origin requests. Because middleware can't emit its own response, register an `OPTIONS` route that returns `204` — the CORS middleware still runs first and attaches the headers:

```go
var ROUTES = router.ForRoutes(
    types.Route{
        Path:        "/api",
        Middlewares: types.Middlewares{NewCORSMiddleware(nil)},
        Children: types.Routes{
            router.Get("/users", []any{UserController{}, "List"}),
            router.Post("/users", []any{UserController{}, "Create"}),
            // Preflight handler
            {Method: "OPTIONS", Path: "/*", Handler: []any{UserController{}, "Preflight"}},
        },
    },
)

func (c *UserController) Preflight(dto *EmptyDto) types.Output {
    return output.NoContent()
}
```

## Route-Specific CORS

Apply different CORS settings to specific routes:

```go
var (
    publicCORS = NewCORSMiddleware(&CORSConfig{
        AllowedOrigins: []string{"*"},
    })

    privateCORS = NewCORSMiddleware(&CORSConfig{
        AllowedOrigins:   []string{"https://admin.myapp.com"},
        AllowCredentials: true,
    })

    ROUTES = router.ForRoutes(
        // Public API - allow all origins
        router.Get("/api/public", []any{Controller{}, "PublicData"}, publicCORS),

        // Admin API - restricted origins
        router.Get("/api/admin", []any{Controller{}, "AdminData"}, privateCORS, &AuthMiddleware{}),
    )
)
```

## Wildcard Subdomains

```go
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
    for _, pattern := range m.config.AllowedOrigins {
        if pattern == "*" {
            return true
        }

        // Support wildcard subdomains
        if strings.HasPrefix(pattern, "*.") {
            domain := strings.TrimPrefix(pattern, "*.")
            if strings.HasSuffix(origin, domain) || origin == "https://"+domain || origin == "http://"+domain {
                return true
            }
        }

        if pattern == origin {
            return true
        }
    }
    return false
}

// Usage
cors := NewCORSMiddleware(&CORSConfig{
    AllowedOrigins: []string{"*.myapp.com"}, // Allows app.myapp.com, admin.myapp.com, etc.
})
```

## Debugging CORS

### Log CORS Requests

```go
func (m *CORSMiddleware) Handle(ctx types.Context) error {
    req := ctx.Request()
    origin := ""
    if h := req.Headers()["Origin"]; len(h) > 0 {
        origin = h[0]
    }

    log.Printf("CORS request: origin=%s method=%s path=%v", origin, req.Method().String(), req.Paths())

    if !m.isOriginAllowed(origin) {
        log.Printf("CORS blocked: origin %s not in allowed list", origin)
    }

    // ... rest of middleware
    return nil
}
```

### Common Issues

1. **"No 'Access-Control-Allow-Origin' header"** - CORS middleware not applied or origin not allowed
2. **"Credentials not supported"** - Using `*` origin with credentials
3. **"Method not allowed"** - Method not in AllowedMethods list
4. **"Header not allowed"** - Custom header not in AllowedHeaders list

## Security Considerations

1. **Don't use `*` in production** with credentials
2. **Whitelist specific origins** for sensitive APIs
3. **Limit allowed methods** to what's needed
4. **Be careful with exposed headers** - don't expose sensitive information
5. **Set reasonable MaxAge** to reduce preflight requests

## Best Practices

1. Configure CORS based on environment
2. Use specific origins in production
3. Only allow necessary methods and headers
4. Enable credentials only when needed
5. Test CORS configuration thoroughly

## Next Steps

- [Authentication](authentication.md) - Secure your API
- [Rate Limiting](rate-limiting.md) - Prevent abuse
- [Security Overview](overview.md) - Security best practices
