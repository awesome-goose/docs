# Requests

In Goose you don't pull data out of a request object inside your handler. Instead, each handler declares a **DTO** — a struct whose fields are annotated with binding tags — and the kernel populates it from the request before the handler runs.

```go
type ShowPostDto struct {
    UserID string `param:"id"`
    PostID string `param:"postId"`
}

func (c *PostController) Show(dto *ShowPostDto) types.Output {
    return output.JSON(c.service.GetUserPost(dto.UserID, dto.PostID))
}
```

## Binding Tags

A field is filled from the first source whose tag is present. Supported tags:

| Tag        | Source                                   | Example                          |
| ---------- | ---------------------------------------- | -------------------------------- |
| `param`    | Path parameter (`/users/:id`)            | `ID string \`param:"id"\``       |
| `query`    | Query string (`?page=2`)                 | `Page string \`query:"page"\``   |
| `header`   | Request header                           | `Auth string \`header:"Authorization"\`` |
| `json`     | JSON body field                          | `Name string \`json:"name"\``    |
| `form`     | URL-encoded / form body field            | `Email string \`form:"email"\``  |
| `flag`     | CLI flag (`--name=John`)                 | `Name string \`flag:"name"\``    |
| `context`  | Value set earlier via `ctx.SetValue`     | `User *User \`context:"user"\``  |

Fields are converted to the field's type automatically — `string`, `int`/`uint`, `float`, `bool`, and `[]string` (comma-separated) are all supported. Complex `json`/`form` fields are unmarshaled into structs, slices, and maps.

## Path Parameters

```go
// Route: /users/:id/posts/:postId
type ShowDto struct {
    UserID string `param:"id"`
    PostID string `param:"postId"`
}

func (c *PostController) Show(dto *ShowDto) types.Output {
    return output.JSON(c.service.GetUserPost(dto.UserID, dto.PostID))
}
```

## Query Parameters

```go
// GET /users?page=2&limit=10&sort=name&order=asc
type ListDto struct {
    Page  int    `query:"page"`
    Limit int    `query:"limit"`
    Sort  string `query:"sort"`
    Order string `query:"order"`
}

func (c *UserController) List(dto *ListDto) types.Output {
    return output.JSON(c.service.ListUsers(dto.Page, dto.Limit, dto.Sort, dto.Order))
}
```

If a query parameter is absent, the field keeps its zero value. Apply your own defaults in the handler:

```go
func (c *UserController) List(dto *ListDto) types.Output {
    if dto.Page == 0 {
        dto.Page = 1
    }
    if dto.Limit == 0 {
        dto.Limit = 10
    }
    return output.JSON(c.service.ListUsers(dto.Page, dto.Limit, dto.Sort, dto.Order))
}
```

## JSON Body

Annotate fields with `json` tags — the body is parsed once and bound for you:

```go
type CreateUserDto struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func (c *UserController) Create(dto *CreateUserDto) types.Output {
    return output.Created(c.service.CreateUser(dto))
}
```

To capture the entire JSON body into a single field (e.g. a nested entity), use the `,merge` suffix:

```go
type CreateOrderDto struct {
    Order OrderPayload `json:"body,merge"`
}
```

## Form Data

For URL-encoded form submissions, use the `form` tag:

```go
type ContactFormDto struct {
    Name    string `form:"name"`
    Email   string `form:"email"`
    Message string `form:"message"`
}

func (c *WebController) SubmitContact(dto *ContactFormDto) types.Output {
    c.service.SendContactEmail(dto)
    return output.Redirect("/contact/thanks")
}
```

## Headers

Bind a header directly to a field:

```go
type ApiDto struct {
    Authorization string `header:"Authorization"`
    APIKey        string `header:"X-API-Key"`
    RequestID     string `header:"X-Request-ID"`
}

func (c *ApiController) Handle(dto *ApiDto) types.Output {
    return output.JSON(c.service.Process(dto.APIKey))
}
```

Multi-value headers are bound to their first value. Declare `[]string` to receive comma-separated values as a slice.

## Context Values (from Middleware)

Middleware still works with `types.Context` and can stash values for later handlers:

```go
type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    user := m.authService.GetUser(ctx)
    ctx.SetValue("user", user)
    return nil
}
```

A handler receives that value by declaring a field with the `context` tag:

```go
type ProfileDto struct {
    User *User `context:"user"`
}

func (c *UserController) Profile(dto *ProfileDto) types.Output {
    return output.JSON(dto.User)
}
```

## CLI Arguments

CLI flags bind through the `flag` tag:

```go
// Command: mycli greet --name=John --times=3
type GreetDto struct {
    Name  string `flag:"name"`
    Times int    `flag:"times"`
}

func (c *CliController) Greet(dto *GreetDto) types.Output {
    var buf strings.Builder
    for i := 0; i < dto.Times; i++ {
        buf.WriteString(fmt.Sprintf("Hello, %s!\n", dto.Name))
    }
    return output.Console(buf.String())
}
```

## Request Validation

Add `validate` tags alongside your binding tags and validate the populated DTO — see [Validation](validation.md):

```go
type CreateProductDto struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Category    string  `json:"category" validate:"required,oneof=electronics clothing food"`
}

func (c *ProductController) Create(dto *CreateProductDto) types.Output {
    if err := c.validator.Validate(dto); err != nil {
        return output.UnprocessableEntity("Validation failed", err)
    }
    return output.Created(c.service.CreateProduct(dto))
}
```

## Raw Request Access

The underlying `types.Context` and `types.Request` are still available where you have a context (for example, in middleware). The `Request` interface exposes:

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

For advanced cases (webhooks with signature verification, multipart uploads), read the raw body in middleware and hand the parsed result to the handler via a context value.

## Best Practices

### 1. One DTO Per Handler

Declare exactly the fields a handler needs. It documents the endpoint's contract and keeps binding predictable.

### 2. Use Strong Types

```go
// ✅ Good: typed DTO with tags
type CreateUserDto struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ❌ Bad: reading and unmarshaling raw maps by hand
```

### 3. Validate Before Using

Always run validation on the populated DTO before touching your services — see the validation example above.

## Next Steps

- [Responses](responses.md) - Sending responses
- [Validation](validation.md) - Input validation
- [Controllers](controllers.md) - Controller patterns
