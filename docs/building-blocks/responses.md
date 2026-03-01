# Responses

Learn how to send responses from your Goose handlers.

## Basic Responses

Controllers return values that are automatically serialized:

```go
func (c *UserController) Show(ctx types.Context) any {
    user := &User{ID: "1", Name: "John"}
    return user  // Automatically serialized
}
```

## JSON Responses (API)

For API platforms, return values are serialized to JSON:

```go
// Simple object
func (c *UserController) Show(ctx types.Context) any {
    return &User{ID: "1", Name: "John"}
}
// Output: {"id": "1", "name": "John"}

// Array
func (c *UserController) List(ctx types.Context) any {
    return []User{
        {ID: "1", Name: "John"},
        {ID: "2", Name: "Jane"},
    }
}
// Output: [{"id": "1", "name": "John"}, {"id": "2", "name": "Jane"}]

// Map
func (c *UserController) Stats(ctx types.Context) any {
    return map[string]any{
        "total": 100,
        "active": 85,
        "inactive": 15,
    }
}
// Output: {"total": 100, "active": 85, "inactive": 15}
```

## HTML Responses (Web)

For web platforms, render HTML templates:

```go
func (c *WebController) Home(ctx types.Context) any {
    return ctx.Render("pages/home.html", map[string]any{
        "title": "Welcome",
        "user": currentUser,
    })
}
```

### Template Data

Pass data to templates:

```go
func (c *WebController) UserProfile(ctx types.Context) any {
    user := c.userService.GetByID(ctx.Param("id"))
    posts := c.postService.GetByUser(user.ID)

    return ctx.Render("users/profile.html", map[string]any{
        "user":  user,
        "posts": posts,
        "stats": map[string]int{
            "followers": 1000,
            "following": 500,
        },
    })
}
```

## Redirects

Redirect to another URL:

```go
// Simple redirect
func (c *WebController) Logout(ctx types.Context) any {
    c.authService.Logout(ctx)
    return ctx.Redirect("/")
}

// Redirect with status code
func (c *WebController) OldPage(ctx types.Context) any {
    return ctx.RedirectPermanent("/new-page")  // 301
}

// Redirect to named route
func (c *WebController) AfterCreate(ctx types.Context) any {
    return ctx.Redirect("/users/" + user.ID)
}
```

## Error Responses

Return error responses:

```go
// Simple error
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetByID(ctx.Param("id"))
    if err != nil {
        return ctx.Error(404, "User not found")
    }
    return user
}

// Structured error
func (c *UserController) Create(ctx types.Context) any {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Error(400, map[string]any{
            "error": "validation_failed",
            "message": "Invalid request body",
            "details": err.Error(),
        })
    }
    return c.service.CreateUser(dto)
}
```

## Status Codes

### Setting Status Codes

```go
func (c *UserController) Create(ctx types.Context) any {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return ctx.Error(500, "Failed to create user")
    }

    ctx.SetStatus(201)  // Created
    return user
}
```

### Common Status Codes

| Code | Meaning               | Usage              |
| ---- | --------------------- | ------------------ |
| 200  | OK                    | Successful GET     |
| 201  | Created               | Successful POST    |
| 204  | No Content            | Successful DELETE  |
| 400  | Bad Request           | Invalid input      |
| 401  | Unauthorized          | Not authenticated  |
| 403  | Forbidden             | Not authorized     |
| 404  | Not Found             | Resource not found |
| 422  | Unprocessable Entity  | Validation error   |
| 500  | Internal Server Error | Server error       |

## Headers

Set response headers:

```go
func (c *ApiController) Handle(ctx types.Context) any {
    // Single header
    ctx.SetHeader("X-Custom-Header", "value")

    // Cache headers
    ctx.SetHeader("Cache-Control", "max-age=3600")
    ctx.SetHeader("ETag", "abc123")

    // CORS headers
    ctx.SetHeader("Access-Control-Allow-Origin", "*")

    return data
}
```

## Cookies

Set cookies in the response:

```go
func (c *AuthController) Login(ctx types.Context) any {
    token := c.authService.CreateToken(user)

    // Simple cookie
    ctx.SetCookie("session", token)

    // Cookie with options
    ctx.SetCookie("session", token, &CookieOptions{
        MaxAge:   3600,        // 1 hour
        Path:     "/",
        Domain:   "example.com",
        Secure:   true,        // HTTPS only
        HttpOnly: true,        // No JavaScript access
        SameSite: "Strict",
    })

    return map[string]string{"status": "logged in"}
}
```

## File Downloads

Send files as downloads:

```go
func (c *FileController) Download(ctx types.Context) any {
    fileID := ctx.Param("id")
    file := c.fileService.GetFile(fileID)

    return ctx.Download(file.Path, file.Name)
}
```

## Streaming Responses

For large data, use streaming:

```go
func (c *ExportController) ExportCSV(ctx types.Context) any {
    ctx.SetHeader("Content-Type", "text/csv")
    ctx.SetHeader("Content-Disposition", "attachment; filename=export.csv")

    return ctx.Stream(func(w io.Writer) error {
        users := c.userService.GetAll()
        for _, user := range users {
            fmt.Fprintf(w, "%s,%s\n", user.ID, user.Name)
        }
        return nil
    })
}
```

## Response Helpers

### Success Response Pattern

```go
func Success(data any) map[string]any {
    return map[string]any{
        "success": true,
        "data":    data,
    }
}

func (c *UserController) List(ctx types.Context) any {
    users := c.service.GetAll()
    return Success(users)
}
// Output: {"success": true, "data": [...]}
```

### Pagination Response

```go
func Paginated(data any, page, limit, total int) map[string]any {
    return map[string]any{
        "data": data,
        "meta": map[string]int{
            "page":       page,
            "limit":      limit,
            "total":      total,
            "totalPages": (total + limit - 1) / limit,
        },
    }
}

func (c *UserController) List(ctx types.Context) any {
    page := ctx.QueryInt("page", 1)
    limit := ctx.QueryInt("limit", 10)

    users, total := c.service.Paginate(page, limit)
    return Paginated(users, page, limit, total)
}
```

## CLI Output

For CLI applications:

```go
func (c *CliController) ListUsers(ctx types.Context) any {
    users := c.service.GetAll()

    // Simple string output
    return fmt.Sprintf("Found %d users", len(users))
}

// Formatted table
func (c *CliController) ListUsers(ctx types.Context) any {
    users := c.service.GetAll()

    var output strings.Builder
    output.WriteString("ID\tName\tEmail\n")
    output.WriteString("--\t----\t-----\n")

    for _, user := range users {
        output.WriteString(fmt.Sprintf("%s\t%s\t%s\n",
            user.ID, user.Name, user.Email))
    }

    return output.String()
}
```

## Best Practices

### 1. Consistent Response Structure

```go
// Define standard response format
type APIResponse struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
}

func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetByID(ctx.Param("id"))
    if err != nil {
        return APIResponse{Success: false, Error: err.Error()}
    }
    return APIResponse{Success: true, Data: user}
}
```

### 2. Use Appropriate Status Codes

```go
func (c *UserController) Create(ctx types.Context) any {
    user, err := c.service.Create(dto)
    if err != nil {
        if errors.Is(err, ErrDuplicate) {
            return ctx.Error(409, "User already exists")  // Conflict
        }
        return ctx.Error(500, "Failed to create user")
    }

    ctx.SetStatus(201)  // Created
    return user
}
```

### 3. Don't Leak Sensitive Data

```go
// ✅ Good: Use DTOs for responses
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // Password omitted
}

// ❌ Bad: Returning entity directly
func (c *UserController) Show(ctx types.Context) any {
    return c.service.GetByID(id)  // May include password hash
}
```

## Next Steps

- [Requests](requests.md) - Handling requests
- [Error Handling](error-handling.md) - Error responses
- [Validation](validation.md) - Input validation
