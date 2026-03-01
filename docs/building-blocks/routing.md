# Routing

Routing in Goose maps incoming requests to handler functions. Routes are defined declaratively in controllers.

## Defining Routes

Routes are defined in a `Routes()` method on your controller:

```go
type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.List},
        {Method: "GET", Path: "/users/:id", Handler: c.Show},
        {Method: "POST", Path: "/users", Handler: c.Create},
        {Method: "PUT", Path: "/users/:id", Handler: c.Update},
        {Method: "DELETE", Path: "/users/:id", Handler: c.Delete},
    }
}
```

## HTTP Methods

Goose supports all standard HTTP methods:

| Method    | Usage              | Example           |
| --------- | ------------------ | ----------------- |
| `GET`     | Retrieve resources | `GET /users`      |
| `POST`    | Create resources   | `POST /users`     |
| `PUT`     | Update resources   | `PUT /users/1`    |
| `PATCH`   | Partial update     | `PATCH /users/1`  |
| `DELETE`  | Delete resources   | `DELETE /users/1` |
| `OPTIONS` | CORS preflight     | `OPTIONS /users`  |
| `HEAD`    | Headers only       | `HEAD /users`     |

```go
func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.List},
        {Method: "POST", Path: "/users", Handler: c.Create},
        {Method: "PUT", Path: "/users/:id", Handler: c.Update},
        {Method: "PATCH", Path: "/users/:id", Handler: c.Patch},
        {Method: "DELETE", Path: "/users/:id", Handler: c.Delete},
    }
}
```

## Route Parameters

### Path Parameters

Use `:name` syntax for dynamic segments:

```go
// Route: /users/:id
func (c *UserController) Show(ctx types.Context) any {
    id := ctx.Param("id")  // Extract parameter
    return c.service.GetUser(id)
}
```

### Multiple Parameters

```go
func (c *PostController) Routes() types.Routes {
    return types.Routes{
        // /users/:userId/posts/:postId
        {Method: "GET", Path: "/users/:userId/posts/:postId", Handler: c.Show},
    }
}

func (c *PostController) Show(ctx types.Context) any {
    userId := ctx.Param("userId")
    postId := ctx.Param("postId")
    return c.service.GetUserPost(userId, postId)
}
```

## Query Parameters

Access query string parameters:

```go
// GET /users?page=1&limit=10&search=john
func (c *UserController) List(ctx types.Context) any {
    page := ctx.Query("page")           // "1"
    limit := ctx.Query("limit")         // "10"
    search := ctx.Query("search")       // "john"

    // With defaults
    pageNum := ctx.QueryDefault("page", "1")

    return c.service.ListUsers(page, limit, search)
}
```

## Route Groups

Organize routes hierarchically using nested routes:

```go
func (c *AppController) Routes() types.Routes {
    return types.Routes{
        {
            Path: "/api",
            Children: types.Routes{
                {
                    Path: "/v1",
                    Children: types.Routes{
                        {Method: "GET", Path: "/users", Handler: c.ListUsers},
                        {Method: "GET", Path: "/users/:id", Handler: c.GetUser},
                    },
                },
            },
        },
    }
}
// Results in: /api/v1/users, /api/v1/users/:id
```

## Route Middleware

Apply middleware to specific routes:

```go
func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.List},  // Public
        {
            Method:      "POST",
            Path:        "/users",
            Handler:     c.Create,
            Middlewares: types.Middlewares{&AuthMiddleware{}},  // Protected
        },
        {
            Path:        "/admin",
            Middlewares: types.Middlewares{&AdminMiddleware{}},
            Children: types.Routes{
                {Method: "GET", Path: "/users", Handler: c.AdminList},
                {Method: "DELETE", Path: "/users/:id", Handler: c.AdminDelete},
            },
        },
    }
}
```

## Wildcard Routes

Match any path segment with `*`:

```go
func (c *FileController) Routes() types.Routes {
    return types.Routes{
        // Matches /files/anything/here
        {Method: "GET", Path: "/files/*", Handler: c.Serve},
    }
}

func (c *FileController) Serve(ctx types.Context) any {
    // Handle file serving
}
```

## Route Options

Each route can have several options:

```go
type Route struct {
    Method      Method       // HTTP method (GET, POST, etc.)
    Path        string       // URL path pattern
    Handler     any          // Handler function
    Middlewares Middlewares  // Route-specific middleware
    Children    Routes       // Nested routes
}
```

## CLI Routes (Commands)

For CLI applications, routes define commands:

```go
func (c *CliController) Routes() types.Routes {
    return types.Routes{
        {Path: "greet", Handler: c.Greet},         // goose cli greet
        {Path: "users/list", Handler: c.ListUsers}, // goose cli users/list
        {
            Path: "db",
            Children: types.Routes{
                {Path: "migrate", Handler: c.Migrate},  // goose cli db/migrate
                {Path: "seed", Handler: c.Seed},        // goose cli db/seed
            },
        },
    }
}
```

## Route Resolution

Goose resolves routes in this order:

1. Exact path match
2. Pattern match with parameters
3. Wildcard match
4. 404 Not Found

```go
Routes:
- GET /users           (exact)
- GET /users/:id       (parameter)
- GET /files/*         (wildcard)

Request: GET /users/123
Match:   GET /users/:id (id = "123")
```

## Route Caching

Goose caches resolved routes for performance. The cache is automatically invalidated when routes change.

## Best Practices

### 1. RESTful Conventions

Follow REST conventions for resource routes:

```go
func (c *PostController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/posts", Handler: c.Index},        // List
        {Method: "GET", Path: "/posts/:id", Handler: c.Show},      // Show one
        {Method: "POST", Path: "/posts", Handler: c.Create},       // Create
        {Method: "PUT", Path: "/posts/:id", Handler: c.Update},    // Update
        {Method: "DELETE", Path: "/posts/:id", Handler: c.Delete}, // Delete
    }
}
```

### 2. Version Your API

```go
{
    Path: "/api",
    Children: types.Routes{
        {
            Path: "/v1",
            Children: v1Routes,
        },
        {
            Path: "/v2",
            Children: v2Routes,
        },
    },
}
```

### 3. Group by Feature

```go
{
    Path: "/admin",
    Middlewares: types.Middlewares{&AdminAuthMiddleware{}},
    Children: adminRoutes,
}
```

## Next Steps

- [Controllers](controllers.md) - Building controllers
- [Middleware](middleware.md) - Request interceptors
- [Requests](requests.md) - Handling request data
