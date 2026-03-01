# Generating Modules

Use the Goose CLI to generate modules within existing applications.

## Basic Usage

```bash
goose g module --name=<module-name> --type=<type>
# or
goose generate module --name=<module-name> --type=<type>
```

## Module Types

### Plain Module

A simple module with basic components:

```bash
goose g module --name=users --type=plain
```

**Generated structure:**

```
app/users/
├── users.module.go      # Module definition
├── users.controller.go  # Request handlers
├── users.service.go     # Business logic
├── users.routes.go      # Route definitions
└── users.dtos.go        # Data transfer objects
```

**Best for:**

- Simple features
- API endpoints
- Services without database entities

### Resource Module

A full CRUD module with database entity:

```bash
goose g module --name=products --type=resource
```

**Generated structure:**

```
app/products/
├── products.module.go      # Module definition
├── products.controller.go  # CRUD handlers
├── products.service.go     # Business logic
├── products.routes.go      # RESTful routes
├── products.dtos.go        # Request/Response DTOs
└── products.entity.go      # Database entity
```

**For web apps (additional):**

```
app/products/
└── templates/
    └── pages/
        ├── list.html      # List view
        └── show.html      # Detail view
```

**Best for:**

- Database-backed resources
- CRUD operations
- RESTful APIs

## Command Options

```bash
goose g module --name=<name> --type=<type> [--template=<platform>]
```

| Flag         | Description                         | Required | Default       |
| ------------ | ----------------------------------- | -------- | ------------- |
| `--name`     | Module name                         | Yes      | -             |
| `--type`     | Module type (`plain` or `resource`) | Yes      | -             |
| `--template` | Platform type (`api`, `web`, `cli`) | No       | Auto-detected |

## Auto-Detection

The CLI automatically detects your project type from `main.go`:

```bash
# In an API project
goose g module --name=users --type=resource
# Generates API-style module

# In a Web project
goose g module --name=users --type=resource
# Generates Web-style module with templates
```

### Explicit Template

Override auto-detection:

```bash
goose g module --name=users --type=resource --template=api
```

## Generated Code

### Plain Module Files

**users.module.go:**

```go
package users

import "github.com/awesome-goose/goose/types"

type UsersModule struct{}

func (m *UsersModule) Imports() []types.Module {
    return []types.Module{}
}

func (m *UsersModule) Exports() []any {
    return []any{
        &UsersService{},
    }
}

func (m *UsersModule) Declarations() []any {
    return []any{
        &UsersController{},
        &UsersService{},
    }
}
```

**users.controller.go:**

```go
package users

import "github.com/awesome-goose/goose/types"

type UsersController struct {
    service *UsersService `inject:""`
}

func (c *UsersController) Index(ctx types.Context) any {
    return c.service.GetAll()
}
```

**users.routes.go:**

```go
package users

import "github.com/awesome-goose/goose/types"

func (c *UsersController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/users", Handler: c.Index},
    }
}
```

### Resource Module Files

**products.entity.go:**

```go
package products

import "time"

type Product struct {
    ID          string     `json:"id" gorm:"primaryKey"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    Price       float64    `json:"price"`
    CreatedAt   *time.Time `json:"created_at"`
    UpdatedAt   *time.Time `json:"updated_at"`
}
```

**products.controller.go:**

```go
package products

import "github.com/awesome-goose/goose/types"

type ProductsController struct {
    service *ProductsService `inject:""`
}

func (c *ProductsController) Index(ctx types.Context) any {
    return c.service.GetAll()
}

func (c *ProductsController) Show(ctx types.Context) any {
    return c.service.GetByID(ctx.Param("id"))
}

func (c *ProductsController) Create(ctx types.Context) any {
    var dto CreateProductDTO
    ctx.Bind(&dto)
    return c.service.Create(dto)
}

func (c *ProductsController) Update(ctx types.Context) any {
    var dto UpdateProductDTO
    ctx.Bind(&dto)
    return c.service.Update(ctx.Param("id"), dto)
}

func (c *ProductsController) Delete(ctx types.Context) any {
    return c.service.Delete(ctx.Param("id"))
}
```

**products.routes.go:**

```go
package products

import "github.com/awesome-goose/goose/types"

func (c *ProductsController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/products", Handler: c.Index},
        {Method: "GET", Path: "/products/:id", Handler: c.Show},
        {Method: "POST", Path: "/products", Handler: c.Create},
        {Method: "PUT", Path: "/products/:id", Handler: c.Update},
        {Method: "DELETE", Path: "/products/:id", Handler: c.Delete},
    }
}
```

## Registering Modules

After generating, import in your root module:

```go
// app/app.module.go
package app

import (
    "myapi/app/users"
    "myapi/app/products"
    "github.com/awesome-goose/goose/types"
)

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        &users.UsersModule{},
        &products.ProductsModule{},
    }
}
```

## Examples

### Create Multiple Modules

```bash
# Generate several modules
goose g module --name=users --type=resource
goose g module --name=products --type=resource
goose g module --name=orders --type=resource
goose g module --name=payments --type=plain
goose g module --name=notifications --type=plain
```

### Web Application Modules

```bash
# For a web app
goose g module --name=posts --type=resource --template=web

# Includes templates:
# app/posts/templates/pages/list.html
# app/posts/templates/pages/show.html
```

### CLI Commands Module

```bash
# For a CLI app
goose g module --name=database --type=plain --template=cli

# Routes become commands:
# goose cli database/migrate
# goose cli database/seed
```

## Best Practices

1. **Use singular names for resource modules**: `user` not `users`
2. **Use plain modules for utilities**: notifications, auth, etc.
3. **Register modules immediately** in your root module
4. **Run `go mod tidy`** after generating

## Next Steps

- [Modules](../core-concepts/modules.md) - Module concepts
- [Controllers](../building-blocks/controllers.md) - Customize controllers
- [Routing](../building-blocks/routing.md) - Define routes
