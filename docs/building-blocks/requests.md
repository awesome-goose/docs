# Requests

Learn how to handle incoming requests and extract data in Goose controllers.

## Request Context

Every handler receives a `types.Context` containing all request information:

```go
func (c *UserController) Handle(ctx types.Context) any {
    // All request data is available through ctx
}
```

## Path Parameters

Extract dynamic path segments:

```go
// Route: /users/:id/posts/:postId
func (c *PostController) Show(ctx types.Context) any {
    userId := ctx.Param("id")       // First parameter
    postId := ctx.Param("postId")   // Second parameter

    return c.service.GetUserPost(userId, postId)
}
```

## Query Parameters

Access URL query string values:

```go
// GET /users?page=2&limit=10&sort=name&order=asc
func (c *UserController) List(ctx types.Context) any {
    page := ctx.Query("page")             // "2"
    limit := ctx.Query("limit")           // "10"
    sort := ctx.Query("sort")             // "name"
    order := ctx.Query("order")           // "asc"

    // With defaults
    pageNum := ctx.QueryDefault("page", "1")
    limitNum := ctx.QueryDefault("limit", "10")

    return c.service.ListUsers(page, limit, sort, order)
}
```

### Multiple Values

```go
// GET /search?tags=go&tags=api&tags=web
func (c *SearchController) Search(ctx types.Context) any {
    tags := ctx.QueryArray("tags")  // ["go", "api", "web"]
    return c.service.SearchByTags(tags)
}
```

## Request Body

### JSON Binding

Bind JSON request body to a struct:

```go
type CreateUserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func (c *UserController) Create(ctx types.Context) any {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Error(400, "Invalid JSON: " + err.Error())
    }

    return c.service.CreateUser(dto)
}
```

### Form Data

Handle form submissions:

```go
type ContactFormDTO struct {
    Name    string `form:"name"`
    Email   string `form:"email"`
    Message string `form:"message"`
}

func (c *WebController) SubmitContact(ctx types.Context) any {
    var dto ContactFormDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Error(400, "Invalid form data")
    }

    c.service.SendContactEmail(dto)
    return ctx.Redirect("/contact/thanks")
}
```

### Raw Body

Access the raw request body:

```go
func (c *WebhookController) Handle(ctx types.Context) any {
    body := ctx.RawBody()  // []byte

    // Process raw webhook payload
    return c.processWebhook(body)
}
```

## Headers

Access request headers:

```go
func (c *ApiController) Handle(ctx types.Context) any {
    // Get single header
    auth := ctx.Header("Authorization")
    contentType := ctx.Header("Content-Type")
    userAgent := ctx.Header("User-Agent")

    // Custom headers
    apiKey := ctx.Header("X-API-Key")
    requestId := ctx.Header("X-Request-ID")

    return nil
}
```

## Cookies

Read cookies from the request:

```go
func (c *WebController) Dashboard(ctx types.Context) any {
    // Get cookie value
    sessionToken := ctx.Cookie("session")
    preferences := ctx.Cookie("prefs")

    if sessionToken == "" {
        return ctx.Redirect("/login")
    }

    return ctx.Render("dashboard.html", nil)
}
```

## Request Metadata

Access request metadata:

```go
func (c *LoggingController) Handle(ctx types.Context) any {
    // HTTP Method
    method := ctx.Method()  // "GET", "POST", etc.

    // Request path
    path := ctx.Path()      // "/users/123"

    // Full URL
    url := ctx.URL()        // "http://localhost:8080/users/123?foo=bar"

    // Client IP address
    ip := ctx.ClientIP()    // "192.168.1.1"

    // Protocol
    protocol := ctx.Protocol()  // "HTTP/1.1"

    return nil
}
```

## Context Values

Store and retrieve values within the request lifecycle:

```go
// In middleware
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    user := m.authService.GetUser(token)
    ctx.Set("user", user)      // Store user
    ctx.Set("role", user.Role) // Store role
    return nil
}

// In controller
func (c *UserController) Profile(ctx types.Context) any {
    user := ctx.Get("user").(*User)   // Retrieve user
    role := ctx.Get("role").(string)  // Retrieve role

    return map[string]any{
        "user": user,
        "role": role,
    }
}
```

## File Uploads

Handle file uploads:

```go
func (c *UploadController) Upload(ctx types.Context) any {
    file, err := ctx.File("document")
    if err != nil {
        return ctx.Error(400, "No file uploaded")
    }

    // File metadata
    filename := file.Filename     // "report.pdf"
    size := file.Size             // 1024 bytes
    contentType := file.ContentType // "application/pdf"

    // Save file
    path := "/uploads/" + filename
    if err := file.Save(path); err != nil {
        return ctx.Error(500, "Failed to save file")
    }

    return map[string]string{
        "message": "File uploaded",
        "path": path,
    }
}
```

### Multiple Files

```go
func (c *UploadController) UploadMultiple(ctx types.Context) any {
    files, err := ctx.Files("images")
    if err != nil {
        return ctx.Error(400, "No files uploaded")
    }

    for _, file := range files {
        path := "/uploads/" + file.Filename
        file.Save(path)
    }

    return map[string]int{
        "uploaded": len(files),
    }
}
```

## Request Validation

Combine binding with validation:

```go
type CreateProductDTO struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Category    string  `json:"category" validate:"required,oneof=electronics clothing food"`
}

func (c *ProductController) Create(ctx types.Context) any {
    var dto CreateProductDTO

    // Bind JSON
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Error(400, "Invalid JSON")
    }

    // Validate
    if err := ctx.Validate(&dto); err != nil {
        return ctx.Error(422, err.Error())
    }

    return c.service.CreateProduct(dto)
}
```

## CLI Arguments

For CLI applications, access command arguments:

```go
// Command: mycli greet --name=John --times=3
func (c *CliController) Greet(ctx types.Context) any {
    name := ctx.Arg("name")      // "John"
    times := ctx.ArgInt("times") // 3

    for i := 0; i < times; i++ {
        fmt.Printf("Hello, %s!\n", name)
    }

    return nil
}
```

## Best Practices

### 1. Always Validate Input

```go
func (c *UserController) Create(ctx types.Context) any {
    var dto CreateUserDTO

    if err := ctx.Bind(&dto); err != nil {
        return ctx.Error(400, "Invalid request body")
    }

    // Validate before processing
    if err := c.validate(&dto); err != nil {
        return ctx.Error(422, err.Error())
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
    ctx.Bind(&data)
    // No type safety
}
```

### 3. Handle Missing Values

```go
func (c *UserController) Show(ctx types.Context) any {
    id := ctx.Param("id")
    if id == "" {
        return ctx.Error(400, "ID is required")
    }

    page := ctx.QueryDefault("page", "1")
    limit := ctx.QueryDefault("limit", "10")

    return c.service.GetUser(id)
}
```

## Next Steps

- [Responses](responses.md) - Sending responses
- [Validation](validation.md) - Input validation
- [Controllers](controllers.md) - Controller patterns
