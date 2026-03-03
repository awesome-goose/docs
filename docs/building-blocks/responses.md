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

For web platforms, return data that your serializer processes:

```go
func (c *WebController) Home(ctx types.Context) any {
    return map[string]any{
        "_template": "pages/home.html",
        "title": "Welcome",
        "user": currentUser,
    }
}
```

### Template Data

Pass data to templates:

```go
func (c *WebController) UserProfile(ctx types.Context) any {
    user := c.userService.GetByID(ctx.Request().Params()["id"])
    posts := c.postService.GetByUser(user.ID)

    return map[string]any{
        "_template": "users/profile.html",
        "user":  user,
        "posts": posts,
        "stats": map[string]int{
            "followers": 1000,
            "following": 500,
        },
    }
}
```

## Redirects

Return redirect information that your platform/serializer handles:

```go
// Simple redirect
func (c *WebController) Logout(ctx types.Context) any {
    c.authService.Logout(ctx)
    return map[string]any{
        "_redirect": "/",
    }
}

// Redirect with status code (301 permanent)
func (c *WebController) OldPage(ctx types.Context) any {
    return map[string]any{
        "_redirect": "/new-page",
        "_status": 301,
    }
}

// Redirect after create
func (c *WebController) AfterCreate(ctx types.Context) any {
    return map[string]any{
        "_redirect": "/users/" + user.ID,
    }
}
```

## Error Responses

Return error responses as structured data:

```go
// Simple error
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetByID(ctx.Request().Params()["id"])
    if err != nil {
        return map[string]any{
            "error": "User not found",
            "_status": 404,
        }
    }
    return user
}

// Structured error
func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]any{
            "error": "validation_failed",
            "message": "Invalid request body",
            "details": err.Error(),
            "_status": 400,
        }
    }
    var dto CreateUserDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        return map[string]any{
            "error": "validation_failed",
            "message": "Invalid JSON",
            "_status": 400,
        }
    }
    return c.service.CreateUser(dto)
}
```

## Status Codes

### Setting Status Codes

Include status in your response for your serializer to handle:

```go
func (c *UserController) Create(ctx types.Context) any {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return map[string]any{
            "error": "Failed to create user",
            "_status": 500,
        }
    }

    return map[string]any{
        "data": user,
        "_status": 201,  // Created
    }
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

Set response headers using the Response interface:

```go
func (c *ApiController) Handle(ctx types.Context) any {
    resp := ctx.Response()

    // Single header
    resp.SetHeader("X-Custom-Header", "value")

    // Cache headers
    resp.SetHeader("Cache-Control", "max-age=3600")
    resp.SetHeader("ETag", "abc123")

    // CORS headers
    resp.SetHeader("Access-Control-Allow-Origin", "*")

    // Multiple headers at once
    resp.SetHeaders(map[string]string{
        "X-Custom-1": "value1",
        "X-Custom-2": "value2",
    })

    return data
}
```

## Cookies

Set cookies using the raw response (platform-specific):

```go
import "net/http"

func (c *AuthController) Login(ctx types.Context) any {
    token := c.authService.CreateToken(user)

    // For HTTP platforms, access raw response writer
    raw := ctx.Response().Raw()
    if w, ok := raw.(http.ResponseWriter); ok {
        http.SetCookie(w, &http.Cookie{
            Name:     "session",
            Value:    token,
            MaxAge:   3600,
            Path:     "/",
            Secure:   true,
            HttpOnly: true,
        })
    }

    return map[string]string{"status": "logged in"}
}
```

## File Downloads

Send files as downloads using raw response:

```go
import (
    "net/http"
    "os"
)

func (c *FileController) Download(ctx types.Context) any {
    fileID := ctx.Request().Params()["id"]
    file := c.fileService.GetFile(fileID)

    // Set download headers
    resp := ctx.Response()
    resp.SetHeader("Content-Disposition", "attachment; filename="+file.Name)
    resp.SetHeader("Content-Type", "application/octet-stream")

    // Read and return file content
    content, err := os.ReadFile(file.Path)
    if err != nil {
        return map[string]string{"error": "File not found"}
    }

    return content
}
```

## Streaming Responses

For large data, use raw response for streaming:

```go
import (
    "fmt"
    "net/http"
)

func (c *ExportController) ExportCSV(ctx types.Context) any {
    resp := ctx.Response()
    resp.SetHeader("Content-Type", "text/csv")
    resp.SetHeader("Content-Disposition", "attachment; filename=export.csv")

    // For HTTP platforms, stream directly to response writer
    raw := resp.Raw()
    if w, ok := raw.(http.ResponseWriter); ok {
        users := c.userService.GetAll()
        for _, user := range users {
            fmt.Fprintf(w, "%s,%s\n", user.ID, user.Name)
        }
    }

    return nil
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
    queries := ctx.Request().Queries()
    page, _ := strconv.Atoi(queries["page"])
    if page == 0 {
        page = 1
    }
    limit, _ := strconv.Atoi(queries["limit"])
    if limit == 0 {
        limit = 10
    }

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
    user, err := c.service.GetByID(ctx.Request().Params()["id"])
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
            return map[string]any{
                "error": "User already exists",
                "_status": 409,  // Conflict
            }
        }
        return map[string]any{
            "error": "Failed to create user",
            "_status": 500,
        }
    }

    return map[string]any{
        "data": user,
        "_status": 201,  // Created
    }
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
