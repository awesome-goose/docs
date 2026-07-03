# Responses

Every handler returns a `types.Output`. You build these with the helpers in the `io/output` package — never by hand and never as a bare `map` or struct. The kernel reads the output's status code, headers, and content type, then serializes its body for the active platform.

```go
import "github.com/awesome-goose/goose/io/output"
```

## JSON Responses (API)

`output.JSON` wraps your data in a standard envelope: `{ "success": bool, "data": any, "message"?: string, "meta"?: any }`.

```go
// Simple object
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    return output.JSON(&User{ID: "1", Name: "John"})
}
// Body: {"success": true, "data": {"id": "1", "name": "John"}}

// Array
func (c *UserController) List(dto *EmptyDto) types.Output {
    return output.JSON([]User{
        {ID: "1", Name: "John"},
        {ID: "2", Name: "Jane"},
    })
}
// Body: {"success": true, "data": [{"id": "1", ...}, {"id": "2", ...}]}

// Map
func (c *UserController) Stats(dto *EmptyDto) types.Output {
    return output.JSON(map[string]any{
        "total":    100,
        "active":   85,
        "inactive": 15,
    })
}
```

If you need a **bare** payload without the `{success, data}` envelope, use `output.Raw`:

```go
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    return output.Raw(&User{ID: "1", Name: "John"})
}
// Body: {"id": "1", "name": "John"}
```

## HTML Responses (Web)

Render a template with `output.View`. The template path is resolved under your app's templates directory and wrapped in the default layout.

```go
func (c *WebController) Home(dto *EmptyDto) types.Output {
    return output.View("pages/home.html", map[string]any{
        "title": "Welcome",
        "user":  currentUser,
    })
}
```

### Template Data

```go
type UserProfileDto struct {
    ID string `param:"id"`
}

func (c *WebController) UserProfile(dto *UserProfileDto) types.Output {
    user := c.userService.GetByID(dto.ID)
    posts := c.postService.GetByUser(user.ID)

    return output.View("pages/users/profile.html", map[string]any{
        "user":  user,
        "posts": posts,
        "stats": map[string]int{
            "followers": 1000,
            "following": 500,
        },
    })
}
```

Options let you override the layout, status, or add partials:

```go
return output.View("pages/users/edit.html", data,
    output.WithLayout("base/admin.html"),
    output.WithHTMLCode(422),
)
```

For a raw HTML string (no template), use `output.HTML("<h1>Hi</h1>")`.

## Redirects

Use `output.Redirect` and friends. They set the `Location` header and the correct status code for you.

```go
// Simple redirect (302)
func (c *WebController) Logout(dto *EmptyDto) types.Output {
    c.authService.Logout()
    return output.Redirect("/")
}

// Permanent redirect (301)
func (c *WebController) OldPage(dto *EmptyDto) types.Output {
    return output.RedirectPermanent("/new-page")
}

// Redirect after create, with flash data
func (c *WebController) AfterCreate(dto *CreateUserDto) types.Output {
    user := c.service.CreateUser(dto)
    return output.Redirect("/users/" + user.ID).
        WithSuccess("User created")
}
```

Other helpers: `output.RedirectTemporary` (307), `output.RedirectPermanentPreserve` (308), `output.Back(fallback)` (redirect to the referrer), and `output.Away(url)` (external). Flash helpers `.With(...)`, `.WithError(...)`, `.WithErrors(...)`, and `.WithOldInput()` are chainable.

## Error Responses

Use the semantic helpers — they set the status code and a consistent error envelope (`{"success": false, "message": ...}`).

```go
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetByID(dto.ID)
    if err != nil {
        return output.NotFound("User not found")
    }
    return output.JSON(user)
}

func (c *UserController) Create(dto *CreateUserDto) types.Output {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        if errors.Is(err, ErrDuplicate) {
            return output.Conflict("User already exists")
        }
        return output.InternalServerError("Failed to create user")
    }
    return output.Created(user)
}
```

To attach structured details to an error, use `output.ErrorWithData` or `output.UnprocessableEntity`:

```go
return output.UnprocessableEntity("Validation failed", map[string][]string{
    "email": {"is required", "must be a valid email"},
})
```

## Status Codes

Prefer the semantic helpers over passing raw codes:

| Helper                          | Code | Meaning               |
| ------------------------------- | ---- | --------------------- |
| `output.OK(data)`               | 200  | Successful GET        |
| `output.Created(data)`          | 201  | Successful POST       |
| `output.Accepted(data)`         | 202  | Queued for processing |
| `output.NoContent()`            | 204  | Successful DELETE     |
| `output.BadRequest(msg)`        | 400  | Invalid input         |
| `output.Unauthorized(msg)`      | 401  | Not authenticated     |
| `output.Forbidden(msg)`         | 403  | Not authorized        |
| `output.NotFound(msg)`          | 404  | Resource not found    |
| `output.Conflict(msg)`          | 409  | Conflict              |
| `output.UnprocessableEntity(…)` | 422  | Validation error      |
| `output.InternalServerError(…)` | 500  | Server error          |

For an arbitrary code, use `output.JSONWithCode(data, code)` or `output.SuccessWithCode(msg, data, code)`.

## Headers

Attach headers via output options:

```go
func (c *ApiController) Handle(dto *EmptyDto) types.Output {
    return output.JSON(data,
        output.WithHeader("X-Custom-Header", "value"),
        output.WithHeaders(map[string]string{
            "Cache-Control": "max-age=3600",
            "ETag":          "abc123",
        }),
    )
}
```

The concrete output types also expose a chainable `SetHeader`:

```go
resp := output.JSON(data)
resp.SetHeader("X-Request-Id", requestID)
return resp
```

## Cookies

Cookies are just a `Set-Cookie` header. Build one with `net/http` and attach it:

```go
import "net/http"

func (c *AuthController) Login(dto *LoginDto) types.Output {
    token := c.authService.CreateToken(dto)

    cookie := &http.Cookie{
        Name:     "session",
        Value:    token,
        MaxAge:   3600,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
    }

    return output.JSON(map[string]string{"status": "logged in"},
        output.WithHeader("Set-Cookie", cookie.String()),
    )
}
```

## File Downloads

Serve files with `output.Download` (attachment) or `output.File` (inline). Content type is detected from the extension.

```go
type FileDto struct {
    ID string `param:"id"`
}

func (c *FileController) Download(dto *FileDto) types.Output {
    file := c.fileService.GetFile(dto.ID)
    return output.Download(file.Path)
    // or output.DownloadWithName(file.Path, "report.pdf")
}
```

To serve bytes you already have in memory, use `output.DownloadFromContent(bytes, "name.csv")` or `output.FileFromContent(...)`.

## Streaming Responses

For large or generated files, stream with a callback via `output.StreamDownload`:

```go
func (c *ExportController) ExportCSV(dto *EmptyDto) types.Output {
    return output.StreamDownload("export.csv", func(write func([]byte) error) error {
        for _, user := range c.userService.GetAll() {
            line := fmt.Sprintf("%s,%s\n", user.ID, user.Name)
            if err := write([]byte(line)); err != nil {
                return err
            }
        }
        return nil
    })
}
```

## Response Envelope & Metadata

`output.JSON` already provides the `{success, data}` envelope, so you don't hand-roll it. Add a message or metadata with the `Success*` helpers or options:

```go
func (c *UserController) List(dto *EmptyDto) types.Output {
    users := c.service.GetAll()
    return output.SuccessWithMeta("Users loaded", users, map[string]any{
        "count": len(users),
    })
}
// Body: {"success": true, "data": [...], "message": "Users loaded", "meta": {"count": ...}}
```

### Pagination

Use `output.Paginated`, which attaches pagination metadata automatically:

```go
type ListUsersDto struct {
    Page  string `query:"page"`
    Limit string `query:"limit"`
}

func (c *UserController) List(dto *ListUsersDto) types.Output {
    page, _ := strconv.Atoi(dto.Page)
    if page == 0 {
        page = 1
    }
    limit, _ := strconv.Atoi(dto.Limit)
    if limit == 0 {
        limit = 10
    }

    users, total := c.service.Paginate(page, limit)
    return output.Paginated(users, page, limit, total)
}
```

## CLI Output

For CLI applications, use the `output` console helpers instead of returning strings:

```go
func (c *CliController) ListUsers(dto *EmptyDto) types.Output {
    users := c.service.GetAll()
    return output.ConsoleSuccess(fmt.Sprintf("Found %d users", len(users)))
}

// Formatted table
func (c *CliController) ListUsers(dto *EmptyDto) types.Output {
    users := c.service.GetAll()

    rows := make([][]string, 0, len(users))
    for _, u := range users {
        rows = append(rows, []string{u.ID, u.Name, u.Email})
    }

    return output.Table([]string{"ID", "Name", "Email"}, rows)
}
```

Other console helpers include `output.Line`, `output.Info`, `output.Warning`, `output.ConsoleError`, `output.List`, `output.Box`, and `output.ProgressBar`. See [Platforms → CLI](../platforms/cli.md).

## Best Practices

### 1. Lean on the Envelope

Let `output.JSON` / `output.Created` / `output.Error` provide a consistent shape rather than assembling maps yourself.

```go
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetByID(dto.ID)
    if err != nil {
        return output.NotFound(err.Error())
    }
    return output.JSON(user)
}
```

### 2. Use Semantic Status Helpers

`output.Created(user)` reads better and stays consistent versus `output.JSONWithCode(user, 201)`.

### 3. Don't Leak Sensitive Data

```go
// ✅ Good: return a response DTO
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // Password omitted
}

func (c *UserController) Show(dto *ShowUserDto) types.Output {
    u := c.service.GetByID(dto.ID)
    return output.JSON(UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
}

// ❌ Bad: returning the entity directly may include a password hash
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    return output.JSON(c.service.GetByID(dto.ID))
}
```

## Next Steps

- [Requests](requests.md) - Handling requests
- [Error Handling](error-handling.md) - Error responses
- [Validation](validation.md) - Input validation
