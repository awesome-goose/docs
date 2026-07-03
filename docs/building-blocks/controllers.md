# Controllers

Controllers handle incoming requests and return responses. They act as the bridge between your routes and business logic.

## Handler Contract

Every controller method (handler) follows the same contract:

```go
func (c *Controller) HandlerName(dto *SomeDto) types.Output
```

- **dto** - A pointer to a struct whose fields are auto-populated from the request (path params, query string, headers, body, etc.) via struct tags. See [Requests](requests.md).
- **return** - A value implementing the `types.Output` interface. You never build these by hand — use the helpers in the `io/output` package (`output.JSON`, `output.Created`, `output.View`, `output.Redirect`, ...).

> Handlers do **not** receive `types.Context` and do **not** return a bare `any`. The kernel populates the DTO for you and serializes whatever `types.Output` you return.

## Creating a Controller

### Basic Controller

```go
package app

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type UserController struct {
    service *UserService `inject:""`
}

// DTOs describe the request shape for each handler.
type ShowUserDto struct {
    ID string `param:"id"`
}

type CreateUserDto struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (c *UserController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAllUsers())
}

func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetUser(dto.ID)
    if err != nil {
        return output.NotFound("User not found")
    }
    return output.JSON(user)
}

func (c *UserController) Create(dto *CreateUserDto) types.Output {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.Created(user)
}
```

`EmptyDto` is just an empty struct (`type EmptyDto struct{}`) for handlers that take no input.

### Registering Routes

Handlers are wired to paths with the `router` module, using the `[]any{Controller{}, "MethodName"}` handler tuple:

```go
import "github.com/awesome-goose/goose/modules/router"

var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "Index"}),
    router.Get("/users/:id", []any{UserController{}, "Show"}),
    router.Post("/users", []any{UserController{}, "Create"}),
)
```

Then import the route module and declare any services in your module:

```go
func (m *AppModule) Imports() []types.Module {
    return []types.Module{ROUTES}
}

func (m *AppModule) Declarations() []any {
    return []any{
        &UserService{},
    }
}
```

Controllers referenced in a route tuple are resolved and dependency-injected automatically — they don't need to appear in `Declarations()`.

## Common Patterns

```go
type ProductController struct {
    service *ProductService `inject:""`
}

type ProductIDDto struct {
    ID string `param:"id"`
}

type CreateProductDto struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type UpdateProductDto struct {
    ID    string  `param:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// List resources
func (c *ProductController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}

// Show single resource
func (c *ProductController) Show(dto *ProductIDDto) types.Output {
    product, err := c.service.GetByID(dto.ID)
    if err != nil {
        return output.NotFound("Product not found")
    }
    return output.JSON(product)
}

// Create resource
func (c *ProductController) Create(dto *CreateProductDto) types.Output {
    product, err := c.service.Create(dto)
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.Created(product)
}

// Update resource
func (c *ProductController) Update(dto *UpdateProductDto) types.Output {
    product, err := c.service.Update(dto.ID, dto)
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.JSON(product)
}

// Delete resource
func (c *ProductController) Delete(dto *ProductIDDto) types.Output {
    if err := c.service.Delete(dto.ID); err != nil {
        return output.BadRequest(err.Error())
    }
    return output.NoContent()
}
```

## Dependency Injection

Controllers receive dependencies through injection:

```go
type OrderController struct {
    orderService   *OrderService   `inject:""`
    userService    *UserService    `inject:""`
    productService *ProductService `inject:""`
    log            types.Log       `inject:""`
    cache          *cache.Cache    `inject:""`
}
```

## Accessing Request Data

Request data is bound to DTO fields via struct tags — you declare what you need and the kernel populates it before your handler runs:

```go
type HandleDto struct {
    ID            string `param:"id"`            // Path parameter  /users/:id
    Page          string `query:"page"`          // Query string    ?page=1
    Authorization string `header:"Authorization"` // Request header
    Name          string `json:"name"`           // JSON body field
    User          any    `context:"user"`         // Value set by middleware via ctx.SetValue
}

func (c *UserController) Handle(dto *HandleDto) types.Output {
    // dto.ID, dto.Page, dto.Authorization, dto.Name, dto.User are ready to use
    result := c.service.Do(dto.ID, dto.Name)

    // Set response headers via output options
    return output.JSON(result, output.WithHeader("X-Custom", "value"))
}
```

See [Requests](requests.md) for the full list of binding tags.

## Response Types

All responses are created with the `io/output` package. See [Responses](responses.md) for the complete catalog.

### JSON Response (API)

```go
func (c *ApiController) GetUser(dto *ShowUserDto) types.Output {
    return output.JSON(&User{ID: "1", Name: "John"})
}
// Body: {"success": true, "data": {"id": "1", "name": "John"}}
```

Use `output.OK`, `output.Created`, `output.Accepted`, etc. for specific status codes.

### HTML Response (Web)

```go
func (c *WebController) GetUser(dto *ShowUserDto) types.Output {
    user := c.service.GetUser(dto.ID)
    return output.View("pages/users/show.html", map[string]any{
        "user": user,
    })
}
```

### Redirect Response

```go
func (c *WebController) CreateUser(dto *CreateUserDto) types.Output {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return output.Back().WithError("Could not create user")
    }
    return output.Redirect("/users/" + user.ID).WithSuccess("User created")
}
```

### Error Response

```go
func (c *Controller) Handle(dto *ShowUserDto) types.Output {
    return output.NotFound("Resource not found")
    // or output.Error("..."), output.BadRequest("..."), output.InternalServerError("...")
}
```

### Custom Status Code

```go
func (c *Controller) Create(dto *CreateProductDto) types.Output {
    newResource := c.service.Create(dto)
    return output.JSONWithCode(newResource, 201)
    // or the semantic helper: output.Created(newResource)
}
```

## Controller Organization

### By Feature (Recommended)

```
app/
├── users/
│   ├── users.controller.go
│   ├── users.dtos.go
│   └── users.routes.go
├── products/
│   ├── products.controller.go
│   ├── products.dtos.go
│   └── products.routes.go
└── orders/
    ├── orders.controller.go
    ├── orders.dtos.go
    └── orders.routes.go
```

Keeping routes in a dedicated `*.routes.go` file (via `router.ForRoutes`) mirrors the framework's own conventions.

## Resource Controllers

For CRUD, implement the standard action set and register them all at once with `router.Resource`, which maps them to RESTful routes:

```go
type ProductController struct {
    service *ProductService `inject:""`
}

func (c *ProductController) List(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}

func (c *ProductController) Get(dto *ProductIDDto) types.Output {
    return output.JSON(c.service.GetByID(dto.ID))
}

func (c *ProductController) Create(dto *CreateProductDto) types.Output {
    return output.Created(c.service.Create(dto))
}

func (c *ProductController) Update(dto *UpdateProductDto) types.Output {
    return output.JSON(c.service.Update(dto.ID, dto))
}

func (c *ProductController) Delete(dto *ProductIDDto) types.Output {
    return output.NoContent()
}
```

Register the whole resource in one call:

```go
var ROUTES = router.ForRoutes(
    // GET/POST /products, GET/PATCH/DELETE /products/:id
    router.Resource("/products", ProductController{}).All(),
)
```

`router.Resource(...)` returns a builder — use `.All()`, `.Only("List", "Get")`, or `.Except("Delete")` to control which actions are registered.

## Best Practices

### 1. Keep Controllers Thin

Controllers should only handle HTTP concerns — bind input and shape output. Push logic into services:

```go
// ✅ Good: Thin controller
func (c *UserController) Create(dto *CreateUserDto) types.Output {
    user, err := c.service.CreateUser(dto)  // Delegate to service
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.Created(user)
}

// ❌ Bad: Fat controller with business logic
func (c *UserController) Create(dto *CreateUserDto) types.Output {
    if dto.Email == "" {
        return output.BadRequest("Email required")
    }
    user := &User{Name: dto.Name, Email: dto.Email}
    c.db.Create(user)                     // Direct DB access
    c.mailer.Send(user.Email, "Welcome!") // Email logic
    return output.Created(user)
}
```

### 2. Use DTOs for Input

Declare a DTO per handler and let the kernel bind and validate it:

```go
type CreateUserDto struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

func (c *UserController) Create(dto *CreateUserDto) types.Output {
    return output.Created(c.service.CreateUser(dto))
}
```

### 3. Handle Errors Consistently

Prefer the semantic `output.*` helpers so status codes stay consistent across the app:

```go
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetByID(dto.ID)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return output.NotFound("User not found")
        }
        return output.InternalServerError("Internal error")
    }
    return output.JSON(user)
}
```

## Next Steps

- [Routing](routing.md) - Defining routes
- [Requests](requests.md) - Handling request data
- [Responses](responses.md) - Sending responses
- [Services](services.md) - Business logic layer
