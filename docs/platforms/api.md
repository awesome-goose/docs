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

Handlers return a `types.Output`. Use `output.JSON` and the semantic status helpers:

```go
type ShowDto struct {
    ID string `param:"id"`
}

func (c *Controller) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}

func (c *Controller) Show(dto *ShowDto) types.Output {
    user := c.service.GetByID(dto.ID)
    if user == nil {
        return output.NotFound("User not found")
    }
    return output.JSON(user)
}
```

### Custom Status Codes

```go
func (c *Controller) Create(dto *CreateUserDTO) types.Output {
    user, err := c.service.Create(dto)
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.Created(user)
}
```

### Empty Response

```go
func (c *Controller) Delete(dto *ShowDto) types.Output {
    c.service.Delete(dto.ID)
    return output.NoContent()
}
```

## Request Handling

### Request Body

The body is bound to your DTO automatically via `json` tags — no manual reading or unmarshaling:

```go
type CreateUserDTO struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

func (c *Controller) Create(dto *CreateUserDTO) types.Output {
    // dto is already populated. Validate, then use it.
    return output.Created(c.service.Create(dto))
}
```

### Query Parameters

```go
type IndexDto struct {
    Page  int    `query:"page"`
    Limit int    `query:"limit"`
    Sort  string `query:"sort"`
}

func (c *Controller) Index(dto *IndexDto) types.Output {
    return output.JSON(c.service.GetPaginated(dto.Page, dto.Limit, dto.Sort))
}
```

### Path Parameters

```go
type ShowDto struct {
    ID string `param:"id"`
}

func (c *Controller) Show(dto *ShowDto) types.Output {
    return output.JSON(c.service.GetByID(dto.ID))
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

Pass middleware as trailing arguments to the route helpers:

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "Index"}),
    router.Get("/users/:id", []any{UserController{}, "Show"}),
    router.Post("/users", []any{UserController{}, "Create"}, &AuthMiddleware{}),
)
```

## API Versioning

### URL Path Versioning

Nest each version's routes under a prefix with `router.Mount`, then start a single instance:

```go
// v1 module — imports its own routes
type V1Module struct{}

func (m *V1Module) Imports() []types.Module {
    return []types.Module{
        router.ForRoutes(
            router.Get("/users", []any{v1.UsersController{}, "Index"}),
        ),
    }
}
func (m *V1Module) Exports() []any      { return []any{} }
func (m *V1Module) Declarations() []any { return []any{} }

// Root module mounts each version under a prefix
type RootModule struct{}

func (m *RootModule) Imports() []types.Module {
    return []types.Module{
        router.Mount("/api/v1", &V1Module{}),
        router.Mount("/api/v2", &V2Module{}),
    }
}
func (m *RootModule) Exports() []any      { return []any{} }
func (m *RootModule) Declarations() []any { return []any{} }

// Start once
platform := api.NewPlatform(api.WithPort(8080))
goose.Start(goose.API(platform, &RootModule{}, nil))
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
func (c *Controller) Create(dto *CreateUserDTO) types.Output {
    user, err := c.service.Create(dto)
    if err != nil {
        if validationErr, ok := err.(*ValidationError); ok {
            return output.UnprocessableEntity("Validation failed", validationErr.Fields)
        }
        return output.InternalServerError("Internal server error")
    }
    return output.Created(user)
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
    return nil
}
```

Middleware sets headers on the shared response and returns `nil` to continue the chain (or an `error` to stop it). For a ready-made implementation, see [Security → CORS](../security/cors.md).

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
