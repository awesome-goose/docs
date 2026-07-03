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

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type UsersController struct {
    service *UsersService `inject:""`
}

type EmptyDto struct{}

func (c *UsersController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}
```

**users.routes.go:**

```go
package users

import "github.com/awesome-goose/goose/modules/router"

var ROUTES = router.ForRoutes(
    router.Get("/users", []any{UsersController{}, "Index"}),
)
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

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type ProductsController struct {
    service *ProductsService `inject:""`
}

type ProductIDDto struct {
    ID string `param:"id"`
}

type UpdateProductDTO struct {
    ID    string  `param:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func (c *ProductsController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}

func (c *ProductsController) Show(dto *ProductIDDto) types.Output {
    return output.JSON(c.service.GetByID(dto.ID))
}

func (c *ProductsController) Create(dto *CreateProductDTO) types.Output {
    return output.Created(c.service.Create(dto))
}

func (c *ProductsController) Update(dto *UpdateProductDTO) types.Output {
    return output.JSON(c.service.Update(dto.ID, dto))
}

func (c *ProductsController) Delete(dto *ProductIDDto) types.Output {
    return output.NoContent()
}
```

**products.routes.go:**

```go
package products

import "github.com/awesome-goose/goose/modules/router"

var ROUTES = router.ForRoutes(
    router.Get("/products", []any{ProductsController{}, "Index"}),
    router.Get("/products/:id", []any{ProductsController{}, "Show"}),
    router.Post("/products", []any{ProductsController{}, "Create"}),
    router.Put("/products/:id", []any{ProductsController{}, "Update"}),
    router.Delete("/products/:id", []any{ProductsController{}, "Delete"}),
)
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
