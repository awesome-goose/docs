# Requests

Learn how to handle incoming requests and extract data in Goose controllers.

## Request Context

Every handler receives a `types.Context` containing all request information:

```go
func (c *UserController) Handle(ctx types.Context) any {
    // Access request data through ctx.Request()
    request := ctx.Request()
}
```

## Context Interface

The `Context` interface provides:

```go
type Context interface {
    Request() Request      // Access request data
    Response() Response    // Access response methods
    SetValue(key, value any)  // Store context values
    GetValue(key any) any     // Retrieve context values
}
```

## Request Interface

The `Request` interface provides:

```go
type Request interface {
    Headers() map[string][]string
    Method() Method
    Paths() []string
    Queries() map[string]string
    Params() map[string]string
    Body() ([]byte, error)
}
```

## Path Parameters

Extract dynamic path segments:

```go
// Route: /users/:id/posts/:postId
func (c *PostController) Show(ctx types.Context) any {
    params := ctx.Request().Params()
    userId := params["id"]       // First parameter
    postId := params["postId"]   // Second parameter

    return c.service.GetUserPost(userId, postId)
}
```

## Query Parameters

Access URL query string values:

```go
// GET /users?page=2&limit=10&sort=name&order=asc
func (c *UserController) List(ctx types.Context) any {
    queries := ctx.Request().Queries()
    page := queries["page"]             // "2"
    limit := queries["limit"]           // "10"
    sort := queries["sort"]             // "name"
    order := queries["order"]           // "asc"

    return c.service.ListUsers(page, limit, sort, order)
}
```

### Helper for Query Defaults

```go
// Helper function for query defaults
func getQueryOr(queries map[string]string, key, defaultVal string) string {
    if val, ok := queries[key]; ok && val != "" {
        return val
    }
    return defaultVal
}

func (c *UserController) List(ctx types.Context) any {
    queries := ctx.Request().Queries()
    page := getQueryOr(queries, "page", "1")
    limit := getQueryOr(queries, "limit", "10")
    // ...
}
```

## Request Body

### JSON Binding

Parse JSON request body to a struct:

```go
import "encoding/json"

type CreateUserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": "Invalid JSON: " + err.Error()}
    }

    return c.service.CreateUser(dto)
}
```

### Form Data

Handle form submissions (parse manually or use a form parsing library):

```go
import (
    "encoding/json"
    "net/url"
)

type ContactFormDTO struct {
    Name    string
    Email   string
    Message string
}

func (c *WebController) SubmitContact(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    values, err := url.ParseQuery(string(body))
    if err != nil {
        return map[string]string{"error": "Invalid form data"}
    }

    dto := ContactFormDTO{
        Name:    values.Get("name"),
        Email:   values.Get("email"),
        Message: values.Get("message"),
    }

    c.service.SendContactEmail(dto)
    return map[string]string{"redirect": "/contact/thanks"}
}
```

### Raw Body

Access the raw request body:

```go
func (c *WebhookController) Handle(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    // Process raw webhook payload
    return c.processWebhook(body)
}
```

## Headers

Access request headers:

```go
func (c *ApiController) Handle(ctx types.Context) any {
    headers := ctx.Request().Headers()

    // Get header (returns []string for multi-value headers)
    auth := headers["Authorization"]
    contentType := headers["Content-Type"]
    userAgent := headers["User-Agent"]

    // Custom headers
    apiKey := headers["X-API-Key"]
    requestId := headers["X-Request-ID"]

    // Get first value if exists
    var authToken string
    if len(auth) > 0 {
        authToken = auth[0]
    }

    return nil
}
```

## Context Values

Store and retrieve values in the request context:

```go
// Set a value (e.g., in middleware)
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    user := m.authService.GetUser(ctx)
    ctx.SetValue("user", user)
    return nil
}

// Get a value (e.g., in controller)
func (c *UserController) Profile(ctx types.Context) any {
    user := ctx.GetValue("user").(*User)
    return user
}
```

## Request Metadata

Access request metadata through the Request interface:

```go
func (c *LoggingController) Handle(ctx types.Context) any {
    req := ctx.Request()

    // HTTP Method
    method := req.Method()  // types.Method (GET, POST, etc.)

    // Request paths (URL segments)
    paths := req.Paths()    // []string{"users", "123"}

    return nil
}
```

## File Uploads

Handle file uploads by parsing the request body:

```go
import (
    "mime/multipart"
    "net/http"
)

func (c *UploadController) Upload(ctx types.Context) any {
    // Access raw response to get multipart form
    raw := ctx.Response().Raw()

    // For HTTP platforms, cast to http.ResponseWriter
    // then parse multipart from the original request
    // This is platform-specific handling

    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    // Process the file data from body
    // File handling depends on your HTTP server implementation

    return map[string]string{
        "message": "File uploaded",
    }
}
```

## Request Validation

Combine binding with validation:

```go
import "encoding/json"

type CreateProductDTO struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Category    string  `json:"category" validate:"required,oneof=electronics clothing food"`
}

func (c *ProductController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    var dto CreateProductDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": "Invalid JSON"}
    }

    // Validate using your validation library
    if err := c.validate(&dto); err != nil {
        return map[string]any{"error": "validation_failed", "details": err.Error()}
    }

    return c.service.CreateProduct(dto)
}
```

## CLI Arguments

For CLI applications, access command arguments through the request:

```go
// Command: mycli greet --name=John --times=3
func (c *CliController) Greet(ctx types.Context) any {
    queries := ctx.Request().Queries()
    name := queries["name"]   // "John"
    times := queries["times"] // "3" (as string, convert as needed)

    // Convert to int
    count, _ := strconv.Atoi(times)
    for i := 0; i < count; i++ {
        fmt.Printf("Hello, %s!\n", name)
    }

    return nil
}
```

## Best Practices

### 1. Always Validate Input

```go
import "encoding/json"

func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]string{"error": "Failed to read body"}
    }

    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]string{"error": "Invalid request body"}
    }

    // Validate before processing
    if err := c.validate(&dto); err != nil {
        return map[string]any{"error": "validation_failed", "details": err.Error()}
    }

    return c.service.CreateUser(dto)
}
```

### 2. Use Strong Types

```go
// ✅ Good: Typed DTO
type CreateUserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ❌ Bad: Untyped map
func (c *UserController) Create(ctx types.Context) any {
    var data map[string]interface{}
    body, _ := ctx.Request().Body()
    json.Unmarshal(body, &data)
    // No type safety
}
```

### 3. Handle Missing Values

```go
func (c *UserController) Show(ctx types.Context) any {
    params := ctx.Request().Params()
    id := params["id"]
    if id == "" {
        return map[string]string{"error": "ID is required"}
    }

    queries := ctx.Request().Queries()
    page := getQueryOr(queries, "page", "1")
    limit := getQueryOr(queries, "limit", "10")

    return c.service.GetUser(id)
}
```

## Next Steps

- [Responses](responses.md) - Sending responses
- [Validation](validation.md) - Input validation
- [Controllers](controllers.md) - Controller patterns
