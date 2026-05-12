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
    api.WithTimeout(30),                // Request timeout (seconds)
    api.WithName("My API"),             // API name
    api.WithVersion("0.0.0"),           // API version
    api.WithAuthor("Your Name"),        // Author
    api.WithDescription("API description"),
)
```

### Available Options

```go
type Config struct {
    Name        string  // API name
    Version     string  // API version
    Author      string  // Author
    Description string  // Description
    Host        string  // Listen address
    Port        int     // Port
    Timeout     int     // Request timeout
}
```

### Environment Configuration

Build the platform inside an initializer where the `types.Env` is available:

```go
initializers := []func(types.Container) error{
    func(c types.Container) error {
        e := env.NewEnv()
        platform := api.NewPlatform(
            api.WithHost(e.GetWithDefault("HOST", "localhost")),
            api.WithPort(e.GetInt("PORT")),
            api.WithTimeout(e.GetInt("TIMEOUT")),
        )
        return c.Register(func() *api.Platform { return platform }, "", true)
    },
}
```

## Response Handling

### JSON Responses

```go
func (c *Controller) Index(ctx types.Context) any {
    users := c.service.GetAll()
    return users  // Automatically serialized to JSON
}

func (c *Controller) Show(ctx types.Context) any {
    user := c.service.GetByID(ctx.Request().Params()["id"])
    if user == nil {
        return map[string]any{
            "error": "User not found",
            "_status": 404,
        }
    }
    return user
}
```

### Custom Status Codes

```go
func (c *Controller) Create(ctx types.Context) any {
    user, err := c.service.Create(dto)
    if err != nil {
        return map[string]any{
            "error": err.Error(),
            "_status": 400,
        }
    }
    return map[string]any{
        "data": user,
        "_status": 201,
    }
}
```

### Empty Response

```go
func (c *Controller) Delete(ctx types.Context) any {
    c.service.Delete(ctx.Request().Params()["id"])
    return map[string]any{
        "_status": 204,
    }
}
```

## Request Handling

### Request Body

```go
import "encoding/json"

type CreateUserDTO struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

func (c *Controller) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]any{
            "error": "Failed to read body",
            "_status": 400,
        }
    }
    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]any{
            "error": "Invalid request body",
            "_status": 400,
        }
    }
    // Use dto...
}
```

### Query Parameters

```go
func (c *Controller) Index(ctx types.Context) any {
    queries := ctx.Request().Queries()
    page := queries["page"]
    limit := queries["limit"]
    sort := queries["sort"]

    return c.service.GetPaginated(page, limit, sort)
}
```

### Path Parameters

```go
func (c *Controller) Show(ctx types.Context) any {
    id := ctx.Request().Params()["id"]
    return c.service.GetByID(id)
}
```

## Middleware

Middleware implements the `types.Middleware` interface:

```go
type Middleware interface {
    Handle(ctx Context) error
}
```

### Authentication

```go
type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    headers := ctx.Request().Headers()
    token := ""
    if auth := headers["Authorization"]; len(auth) > 0 {
        token = auth[0]
    }

    user, err := m.authService.ValidateToken(token)
    if err != nil {
        // Return error to stop the chain
        return err
    }

    ctx.SetValue("user", user)
    return nil  // Continue to next middleware/handler
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

func (m *VersionMiddleware) Handle(ctx types.Context) error {
    headers := ctx.Request().Headers()
    version := "v1"
    if v := headers["API-Version"]; len(v) > 0 {
        version = v[0]
    }
    ctx.SetValue("api_version", version)
    return nil
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
            return map[string]any{
                "error":   "Validation failed",
                "code":    "VALIDATION_ERROR",
                "details": validationErr.Fields,
                "_status": 400,
            }
        }
        return map[string]any{
            "error":   "Internal server error",
            "code":    "INTERNAL_ERROR",
            "_status": 500,
        }
    }
    return map[string]any{
        "data":    user,
        "_status": 201,
    }
}
```

## CORS

```go
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(ctx types.Context) error {
    resp := ctx.Response()
    resp.SetHeader("Access-Control-Allow-Origin", "*")
    resp.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    resp.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if ctx.Request().Method() == types.OPTIONS {
        // Handle preflight request
        resp.Write(types.SerialJSON, []byte(""), 204)
        return nil
    }

    return nil
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
