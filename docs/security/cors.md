# CORS Configuration

Configure Cross-Origin Resource Sharing (CORS) for your Goose API.

## Overview

CORS controls which origins can access your API from browsers. Without proper CORS configuration, browsers block cross-origin requests.

## Quick Setup

### Basic CORS Middleware

```go
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Set CORS headers
    ctx.SetHeader("Access-Control-Allow-Origin", "*")
    ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

    // Handle preflight
    if ctx.Method() == "OPTIONS" {
        return ctx.Status(204).Send("")
    }

    return next()
}
```

### Apply Globally

```go
func (c *AppController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/", Handler: c.Index, Middlewares: []any{&CORSMiddleware{}}},
        // Apply to all routes...
    }
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

func (m *CORSMiddleware) Handle(ctx types.Context, next types.Next) any {
    origin := ctx.Header("Origin")

    // Check if origin is allowed
    if !m.isOriginAllowed(origin) {
        return next()
    }

    // Set CORS headers
    ctx.SetHeader("Access-Control-Allow-Origin", origin)
    ctx.SetHeader("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
    ctx.SetHeader("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))

    if len(m.config.ExposedHeaders) > 0 {
        ctx.SetHeader("Access-Control-Expose-Headers", strings.Join(m.config.ExposedHeaders, ", "))
    }

    if m.config.AllowCredentials {
        ctx.SetHeader("Access-Control-Allow-Credentials", "true")
    }

    if m.config.MaxAge > 0 {
        ctx.SetHeader("Access-Control-Max-Age", strconv.Itoa(m.config.MaxAge))
    }

    // Handle preflight
    if ctx.Method() == "OPTIONS" {
        return ctx.Status(204).Send("")
    }

    return next()
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
func LoadCORSConfig() *CORSConfig {
    appEnv := env.String("APP_ENV", "development")

    if appEnv == "production" {
        origins := strings.Split(env.String("CORS_ORIGINS", ""), ",")
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

Browsers send OPTIONS requests before certain requests:

```go
func (m *CORSMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Set CORS headers first
    m.setCORSHeaders(ctx)

    // Handle preflight (OPTIONS) immediately
    if ctx.Method() == "OPTIONS" {
        return ctx.Status(204).Send("")
    }

    return next()
}
```

## Route-Specific CORS

Apply different CORS settings to specific routes:

```go
func (c *Controller) Routes() types.Routes {
    publicCORS := NewCORSMiddleware(&CORSConfig{
        AllowedOrigins: []string{"*"},
    })

    privateCORS := NewCORSMiddleware(&CORSConfig{
        AllowedOrigins:   []string{"https://admin.myapp.com"},
        AllowCredentials: true,
    })

    return types.Routes{
        // Public API - allow all origins
        {Method: "GET", Path: "/api/public", Handler: c.PublicData, Middlewares: []any{publicCORS}},

        // Admin API - restricted origins
        {Method: "GET", Path: "/api/admin", Handler: c.AdminData, Middlewares: []any{privateCORS, &AuthMiddleware{}}},
    }
}
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
func (m *CORSMiddleware) Handle(ctx types.Context, next types.Next) any {
    origin := ctx.Header("Origin")
    method := ctx.Method()

    log.Printf("CORS request: origin=%s method=%s path=%s", origin, method, ctx.Path())

    if !m.isOriginAllowed(origin) {
        log.Printf("CORS blocked: origin %s not in allowed list", origin)
    }

    // ... rest of middleware
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
