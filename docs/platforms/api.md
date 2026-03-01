# API Platform

Build JSON REST APIs with the Goose API platform.

## Overview

The API platform serves JSON responses for building RESTful APIs.

## Quick Start

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    platform := api.NewPlatform(
        api.WithHost("localhost"),
        api.WithPort(8080),
    )

    module := &app.AppModule{}

    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Configuration

### Platform Options

```go
platform := api.NewPlatform(
    api.WithHost("0.0.0.0"),           // Listen address
    api.WithPort(8080),                 // Port
    api.WithPrefix("/api/v1"),          // URL prefix
    api.WithReadTimeout(30*time.Second),
    api.WithWriteTimeout(30*time.Second),
)
```

### Environment Configuration

```go
platform := api.NewPlatform(
    api.WithHost(env.String("HOST", "localhost")),
    api.WithPort(env.Int("PORT", 8080)),
    api.WithPrefix(env.String("API_PREFIX", "/api")),
)
```

## Response Handling

### JSON Responses

```go
func (c *Controller) Index(ctx types.Context) any {
    users := c.service.GetAll()
    return users  // Automatically serialized to JSON
}

func (c *Controller) Show(ctx types.Context) any {
    user := c.service.GetByID(ctx.Param("id"))
    if user == nil {
        return ctx.Status(404).JSON(map[string]string{
            "error": "User not found",
        })
    }
    return user
}
```

### Custom Status Codes

```go
func (c *Controller) Create(ctx types.Context) any {
    user, err := c.service.Create(dto)
    if err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": err.Error(),
        })
    }
    return ctx.Status(201).JSON(user)
}
```

### Empty Response

```go
func (c *Controller) Delete(ctx types.Context) any {
    c.service.Delete(ctx.Param("id"))
    return ctx.Status(204).Send("")
}
```

## Request Handling

### Request Body

```go
type CreateUserDTO struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

func (c *Controller) Create(ctx types.Context) any {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": "Invalid request body",
        })
    }
    // Use dto...
}
```

### Query Parameters

```go
func (c *Controller) Index(ctx types.Context) any {
    page := ctx.QueryInt("page", 1)
    limit := ctx.QueryInt("limit", 10)
    sort := ctx.Query("sort", "created_at")

    return c.service.GetPaginated(page, limit, sort)
}
```

### Path Parameters

```go
func (c *Controller) Show(ctx types.Context) any {
    id := ctx.Param("id")
    return c.service.GetByID(id)
}
```

## Middleware

### Authentication

```go
type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context, next types.Next) any {
    token := ctx.Header("Authorization")

    user, err := m.authService.ValidateToken(token)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Unauthorized",
        })
    }

    ctx.Set("user", user)
    return next()
}
```

### Apply Middleware

```go
func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.Index},
        {Method: "GET", Path: "/users/:id", Handler: c.Show},
        {Method: "POST", Path: "/users", Handler: c.Create,
            Middlewares: []any{&AuthMiddleware{}}},
    }
}
```

## API Versioning

### URL Path Versioning

```go
// v1 module
type V1Module struct{}

func (m *V1Module) Declarations() []any {
    return []any{
        &v1.UsersController{},
    }
}

// v2 module
type V2Module struct{}

func (m *V2Module) Declarations() []any {
    return []any{
        &v2.UsersController{},
    }
}

// Register with prefixes
platform := api.NewPlatform(api.WithPort(8080))

goose.Start(
    goose.API(platform, &V1Module{}, nil, api.WithPrefix("/api/v1")),
    goose.API(platform, &V2Module{}, nil, api.WithPrefix("/api/v2")),
)
```

### Header Versioning

```go
type VersionMiddleware struct{}

func (m *VersionMiddleware) Handle(ctx types.Context, next types.Next) any {
    version := ctx.Header("API-Version")
    if version == "" {
        version = "v1"
    }
    ctx.Set("api_version", version)
    return next()
}
```

## Error Handling

### Standard Error Response

```go
type ErrorResponse struct {
    Error   string            `json:"error"`
    Code    string            `json:"code,omitempty"`
    Details map[string]string `json:"details,omitempty"`
}

func (c *Controller) Create(ctx types.Context) any {
    user, err := c.service.Create(dto)
    if err != nil {
        if validationErr, ok := err.(*ValidationError); ok {
            return ctx.Status(400).JSON(ErrorResponse{
                Error:   "Validation failed",
                Code:    "VALIDATION_ERROR",
                Details: validationErr.Fields,
            })
        }
        return ctx.Status(500).JSON(ErrorResponse{
            Error: "Internal server error",
            Code:  "INTERNAL_ERROR",
        })
    }
    return ctx.Status(201).JSON(user)
}
```

## CORS

```go
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(ctx types.Context, next types.Next) any {
    ctx.SetHeader("Access-Control-Allow-Origin", "*")
    ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if ctx.Method() == "OPTIONS" {
        return ctx.Status(204).Send("")
    }

    return next()
}
```

## Best Practices

1. **Use consistent response format** across all endpoints
2. **Version your API** from the start
3. **Validate input** thoroughly
4. **Handle errors** with appropriate status codes
5. **Document your API** with OpenAPI/Swagger
6. **Use authentication** for sensitive endpoints
7. **Implement rate limiting** for public APIs

## Next Steps

- [Routing](../building-blocks/routing.md) - Route definitions
- [Controllers](../building-blocks/controllers.md) - Request handling
- [Authentication](../security/authentication.md) - Secure your API
