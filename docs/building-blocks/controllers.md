# Controllers

Controllers handle incoming requests and return responses. They act as the bridge between your routes and business logic.

## Creating a Controller

### Basic Controller

```go
package app

import "github.com/awesome-goose/goose/types"

type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) Index(ctx types.Context) any {
    return c.service.GetAllUsers()
}

func (c *UserController) Show(ctx types.Context) any {
    id := ctx.Request().Params()["id"]
    return c.service.GetUser(id)
}

func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": err.Error()}
    }
    return c.service.CreateUser(dto)
}
```

### Registering the Controller

Register in your module's declarations:

```go
func (m *AppModule) Declarations() []any {
    return []any{
        &UserController{},
        &UserService{},
    }
}
```

## Controller Methods

### Handler Signature

All handlers follow the same pattern:

```go
func (c *Controller) HandlerName(ctx types.Context) any
```

- **ctx** - The request context with all request data
- **return** - Any value (auto-serialized based on platform)

### Common Patterns

```go
// List resources
func (c *ProductController) Index(ctx types.Context) any {
    products := c.service.GetAll()
    return products
}

// Show single resource
func (c *ProductController) Show(ctx types.Context) any {
    id := ctx.Request().Params()["id"]
    product, err := c.service.GetByID(id)
    if err != nil {
        return map[string]string{"error": "Product not found"}
    }
    return product
}

// Create resource
func (c *ProductController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    var dto CreateProductDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": err.Error()}
    }

    product, err := c.service.Create(dto)
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    return product
}

// Update resource
func (c *ProductController) Update(ctx types.Context) any {
    id := ctx.Request().Params()["id"]
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    var dto UpdateProductDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": err.Error()}
    }

    product, err := c.service.Update(id, dto)
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    return product
}

// Delete resource
func (c *ProductController) Delete(ctx types.Context) any {
    id := ctx.Request().Params()["id"]
    if err := c.service.Delete(id); err != nil {
        return map[string]string{"error": err.Error()}
    }
    return map[string]bool{"deleted": true}
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

## Request Context

The context provides access to all request data:

```go
func (c *UserController) Handle(ctx types.Context) any {
    req := ctx.Request()
    resp := ctx.Response()

    // Path parameters
    id := req.Params()["id"]

    // Query parameters
    page := req.Queries()["page"]

    // Request body
    body, err := req.Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }
    var dto MyDTO
    json.Unmarshal(body, &dto)

    // Headers
    auth := req.Headers()["Authorization"]

    // Set response headers
    resp.SetHeader("X-Custom", "value")

    // Store context values
    ctx.SetValue("user", user)

    return result
}
```

## Response Types

### JSON Response (API)

```go
// Return any value - automatically serialized to JSON
func (c *ApiController) GetUser(ctx types.Context) any {
    return &User{ID: "1", Name: "John"}
}
// Output: {"id": "1", "name": "John"}
```

### HTML Response (Web)

```go
// For web platforms, return data and let the serializer handle rendering
func (c *WebController) GetUser(ctx types.Context) any {
    user := c.service.GetUser(ctx.Request().Params()["id"])
    return map[string]any{
        "user": user,
    }
}
```

### Redirect Response

```go
func (c *WebController) CreateUser(ctx types.Context) any {
    // Process...
    // For redirects, return data that your serializer handles
    return map[string]string{"redirect": "/users"}
}
```

### Error Response

```go
func (c *Controller) Handle(ctx types.Context) any {
    // Return error as a map or struct
    return map[string]any{
        "error": "Resource not found",
        "status": 404,
    }
}
```

### Custom Status Code

```go
func (c *Controller) Create(ctx types.Context) any {
    // Return with status information for serialization
    return map[string]any{
        "data": newResource,
        "_status": 201,  // Handle in your serializer
    }
}
```

## Controller Organization

### By Feature (Recommended)

```
app/
├── users/
│   ├── users.controller.go
│   └── users.routes.go
├── products/
│   ├── products.controller.go
│   └── products.routes.go
└── orders/
    ├── orders.controller.go
    └── orders.routes.go
```

### Single Controller Pattern

For simple apps, one controller may suffice:

```go
type AppController struct {
    userService    *UserService    `inject:""`
    productService *ProductService `inject:""`
}

func (c *AppController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.ListUsers},
        {Method: "GET", Path: "/products", Handler: c.ListProducts},
    }
}
```

## Resource Controllers

For CRUD operations, create resource controllers:

```go
type ResourceController[T any] interface {
    Index(ctx types.Context) any   // GET /resources
    Show(ctx types.Context) any    // GET /resources/:id
    Create(ctx types.Context) any  // POST /resources
    Update(ctx types.Context) any  // PUT /resources/:id
    Delete(ctx types.Context) any  // DELETE /resources/:id
}
```

Example implementation:

```go
type ProductController struct {
    service *ProductService `inject:""`
}

func (c *ProductController) Index(ctx types.Context) any {
    return c.service.GetAll()
}

func (c *ProductController) Show(ctx types.Context) any {
    return c.service.GetByID(ctx.Request().Params()["id"])
}

func (c *ProductController) Create(ctx types.Context) any {
    body, _ := ctx.Request().Body()
    var dto CreateProductDTO
    json.Unmarshal(body, &dto)
    return c.service.Create(dto)
}

func (c *ProductController) Update(ctx types.Context) any {
    body, _ := ctx.Request().Body()
    var dto UpdateProductDTO
    json.Unmarshal(body, &dto)
    return c.service.Update(ctx.Request().Params()["id"], dto)
}

func (c *ProductController) Delete(ctx types.Context) any {
    return c.service.Delete(ctx.Request().Params()["id"])
}

func (c *ProductController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/products", Handler: c.Index},
        {Method: "GET", Path: "/products/:id", Handler: c.Show},
        {Method: "POST", Path: "/products", Handler: c.Create},
        {Method: "PUT", Path: "/products/:id", Handler: c.Update},
        {Method: "DELETE", Path: "/products/:id", Handler: c.Delete},
    }
}
```

## Best Practices

### 1. Keep Controllers Thin

Controllers should only handle HTTP concerns:

```go
// ✅ Good: Thin controller
func (c *UserController) Create(ctx types.Context) any {
    body, _ := ctx.Request().Body()
    var dto CreateUserDTO
    json.Unmarshal(body, &dto)
    return c.service.CreateUser(dto)  // Delegate to service
}

// ❌ Bad: Fat controller with business logic
func (c *UserController) Create(ctx types.Context) any {
    body, _ := ctx.Request().Body()
    var dto CreateUserDTO
    json.Unmarshal(body, &dto)

    // Don't do business logic here
    if dto.Email == "" {
        return map[string]string{"error": "Email required"}
    }
    user := &User{Name: dto.Name, Email: dto.Email}
    c.db.Create(user)  // Direct DB access
    c.mailer.Send(user.Email, "Welcome!")  // Email logic
    return user
}
```

### 2. Use DTOs for Input

Always use DTOs for request binding:

```go
type CreateUserDTO struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": err.Error()}
    }
    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": err.Error()}
    }
    return c.service.CreateUser(dto)
}
```

### 3. Handle Errors Consistently

```go
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetByID(ctx.Request().Params()["id"])
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return map[string]any{"error": "User not found", "status": 404}
        }
        return map[string]any{"error": "Internal error", "status": 500}
    }
    return user
}
```

## Next Steps

- [Routing](routing.md) - Defining routes
- [Requests](requests.md) - Handling request data
- [Responses](responses.md) - Sending responses
- [Services](services.md) - Business logic layer
