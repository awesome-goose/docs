# Routing

Routing maps incoming requests to controller handlers. Routes are declared with the `router` module and registered by importing the resulting module into your module tree.

```go
import "github.com/awesome-goose/goose/modules/router"
```

## Defining Routes

Build a route module with `router.ForRoutes` and the HTTP verb helpers. Each handler is a `[]any{Controller{}, "MethodName"}` tuple:

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "List"}),
    router.Get("/users/:id", []any{UserController{}, "Show"}),
    router.Post("/users", []any{UserController{}, "Create"}),
    router.Put("/users/:id", []any{UserController{}, "Update"}),
    router.Delete("/users/:id", []any{UserController{}, "Delete"}),
)
```

Import the module so the kernel picks up its routes:

```go
func (m *AppModule) Imports() []types.Module {
    return []types.Module{ROUTES}
}
```

## HTTP Methods

Goose provides a helper per verb:

| Helper          | Method | Example                  |
| --------------- | ------ | ------------------------ |
| `router.Get`    | GET    | `GET /users`             |
| `router.Post`   | POST   | `POST /users`            |
| `router.Put`    | PUT    | `PUT /users/:id`         |
| `router.Patch`  | PATCH  | `PATCH /users/:id`       |
| `router.Delete` | DELETE | `DELETE /users/:id`      |
| `router.Cli`    | (CLI)  | `mycli users/list`       |

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "List"}),
    router.Post("/users", []any{UserController{}, "Create"}),
    router.Put("/users/:id", []any{UserController{}, "Update"}),
    router.Patch("/users/:id", []any{UserController{}, "Patch"}),
    router.Delete("/users/:id", []any{UserController{}, "Delete"}),
)
```

## Route Parameters

### Path Parameters

Use `:name` for dynamic segments. The value is bound to a DTO field via the `param` tag — handlers never read parameters off a context:

```go
// Route: /users/:id
type ShowDto struct {
    ID string `param:"id"`
}

func (c *UserController) Show(dto *ShowDto) types.Output {
    return output.JSON(c.service.GetUser(dto.ID))
}
```

### Multiple Parameters

```go
// Route: /users/:userId/posts/:postId
type ShowPostDto struct {
    UserID string `param:"userId"`
    PostID string `param:"postId"`
}

func (c *PostController) Show(dto *ShowPostDto) types.Output {
    return output.JSON(c.service.GetUserPost(dto.UserID, dto.PostID))
}
```

## Query Parameters

Query string values bind through the `query` tag:

```go
// GET /users?page=1&limit=10&search=john
type ListDto struct {
    Page   int    `query:"page"`
    Limit  int    `query:"limit"`
    Search string `query:"search"`
}

func (c *UserController) List(dto *ListDto) types.Output {
    if dto.Page == 0 {
        dto.Page = 1
    }
    return output.JSON(c.service.ListUsers(dto.Page, dto.Limit, dto.Search))
}
```

See [Requests](requests.md) for all binding tags.

## Route Groups

Nest routes under a shared prefix with `router.Group`:

```go
var ROUTES = router.ForRoutes(
    router.Group("/api",
        router.Group("/v1",
            router.Get("/users", []any{UserController{}, "List"}),
            router.Get("/users/:id", []any{UserController{}, "Show"}),
        ),
    ),
)
// Results in: /api/v1/users, /api/v1/users/:id
```

## Route Middleware

The verb helpers accept middleware as trailing arguments:

```go
var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UserController{}, "List"}),                    // Public
    router.Post("/users", []any{UserController{}, "Create"}, &AuthMiddleware{}), // Protected

    router.Group("/admin",
        router.Get("/users", []any{AdminController{}, "List"}),
        router.Delete("/users/:id", []any{AdminController{}, "Delete"}),
    ),
)
```

Middleware attached to a group route applies to all of its children. Middleware still receives `types.Context` and returns an `error`:

```go
func (m *AuthMiddleware) Handle(ctx types.Context) error {
    if !m.authService.IsAuthenticated(ctx) {
        return errors.New("unauthorized")
    }
    return nil
}
```

## Resource Routes

For standard CRUD, register all five RESTful routes at once with `router.Resource`:

```go
var ROUTES = router.ForRoutes(
    // POST /products, GET /products, GET/PATCH/DELETE /products/:id
    router.Resource("/products", ProductController{}).All(),
)
```

The controller must implement `Create`, `List`, `Get`, `Update`, and `Delete`. Narrow the set with `.Only(...)` or `.Except(...)`:

```go
router.Resource("/products", ProductController{}).Only("List", "Get")
router.Resource("/products", ProductController{}).Except("Delete")
```

## Mounting Modules Under a Prefix

`router.Mount` nests every route registered by a module (and its imports) under a prefix — the equivalent of NestJS's `RouterModule.register`:

```go
func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        router.Mount("/api/v1", &UsersModule{}),
        router.Mount("/api/v2", &UsersV2Module{}),
    }
}
```

## Route Struct

The verb helpers construct `types.Route` values; you can also build them literally:

```go
types.Route{
    Method:      types.GET,      // HTTP method
    Path:        "/users/:id",   // URL path pattern
    Handler:     []any{UserController{}, "Show"},
    Middlewares: types.Middlewares{&AuthMiddleware{}},
    Children:    types.Routes{}, // Nested routes
}
```

## CLI Routes (Commands)

For CLI apps, `router.Cli` maps a path to a command:

```go
var ROUTES = router.ForRoutes(
    router.Cli("/", []any{AppController{}, "Help"}),          // default command
    router.Cli("/greet", []any{GreetController{}, "Greet"}),  // mycli greet
    router.Group("/db",
        router.Cli("/migrate", []any{DbController{}, "Migrate"}), // mycli db/migrate
        router.Cli("/seed", []any{DbController{}, "Seed"}),       // mycli db/seed
    ),
)
```

## Route Resolution

Routes resolve in this order:

1. Exact path match
2. Pattern match with parameters
3. Wildcard match
4. 404 Not Found

```
Routes:
- GET /users           (exact)
- GET /users/:id       (parameter)

Request: GET /users/123
Match:   GET /users/:id  (id = "123")
```

## Best Practices

### 1. RESTful Conventions

```go
var ROUTES = router.ForRoutes(
    router.Get("/posts", []any{PostController{}, "Index"}),      // List
    router.Get("/posts/:id", []any{PostController{}, "Show"}),   // Show one
    router.Post("/posts", []any{PostController{}, "Create"}),    // Create
    router.Put("/posts/:id", []any{PostController{}, "Update"}), // Update
    router.Delete("/posts/:id", []any{PostController{}, "Delete"}), // Delete
)
```

### 2. Version Your API

```go
router.Group("/api",
    router.Group("/v1", v1Routes...),
    router.Group("/v2", v2Routes...),
)
```

### 3. Keep Routes in a `*.routes.go` File

Mirror the framework's convention: define a `ROUTES` module per feature and import it from the feature module.

## Next Steps

- [Controllers](controllers.md) - Building controllers
- [Middleware](middleware.md) - Request interceptors
- [Requests](requests.md) - Handling request data
