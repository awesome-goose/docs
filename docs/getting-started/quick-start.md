# Quick Start

Build your first Goose application in under 5 minutes.

## Create Your First API

### Step 1: Create the Application

```bash
goose app --name=myapi --template=api
```

This creates a new API application with the following structure:

```
myapi/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    └── app.dtos.go
```

### Step 2: Install Dependencies

```bash
cd myapi
go mod tidy
```

### Step 3: Run the Application

```bash
go run main.go
```

Your API is now running at `http://localhost:8080`!

### Step 4: Test the API

Open a new terminal and test your API:

```bash
curl http://localhost:8080/
```

Response:

```json
{
  "success": true,
  "data": { "message": "Welcome to Goose API!" }
}
```

## Understanding the Generated Code

### main.go

```go
package main

import (
    "myapi/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    // Create the API platform with configuration
    platform := api.NewPlatform(
        api.WithHost("localhost"),
        api.WithPort(8080),
    )

    // Create the application module
    module := &app.AppModule{}

    // Start the application
    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

### app/app.module.go

```go
package app

import "github.com/awesome-goose/goose/types"

type AppModule struct{}

// Imports returns modules this module depends on (including its routes)
func (m *AppModule) Imports() []types.Module {
    return []types.Module{ROUTES}
}

// Exports returns components available to other modules
func (m *AppModule) Exports() []any {
    return []any{}
}

// Declarations returns components this module provides
func (m *AppModule) Declarations() []any {
    return []any{
        &AppService{},
    }
}
```

### app/app.controller.go

```go
package app

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type AppController struct {
    service *AppService `inject:""`
}

func (c *AppController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetWelcomeMessage())
}
```

### app/app.routes.go

```go
package app

import "github.com/awesome-goose/goose/modules/router"

var ROUTES = router.ForRoutes(
    router.Get("/", []any{AppController{}, "Index"}),
)
```

### app/app.dtos.go

```go
package app

// EmptyDto is used by handlers that take no request input.
type EmptyDto struct{}
```

## Adding a New Endpoint

Let's add a `/hello/:name` endpoint:

### 1. Add the Handler Method

In `app/app.controller.go`:

```go
type HelloDto struct {
    Name string `param:"name"`
}

func (c *AppController) Hello(dto *HelloDto) types.Output {
    return output.JSON(map[string]string{
        "message": "Hello, " + dto.Name + "!",
    })
}
```

### 2. Register the Route

In `app/app.routes.go`:

```go
var ROUTES = router.ForRoutes(
    router.Get("/", []any{AppController{}, "Index"}),
    router.Get("/hello/:name", []any{AppController{}, "Hello"}),
)
```

### 3. Test It

```bash
curl http://localhost:8080/hello/World
```

Response:

```json
{
  "success": true,
  "data": { "message": "Hello, World!" }
}
```

## Creating Different Application Types

### Web Application

```bash
goose app --name=myweb --template=web
```

### CLI Application

```bash
goose app --name=mycli --template=cli
```

### Multi-Platform Application

```bash
goose app --name=mymulti --template=multi
```

This creates an application that combines API, Web, and CLI platforms.

## Adding Modules

Generate a new module in your application:

```bash
# Navigate to your project
cd myapi

# Generate a plain module
goose g module --name=users --type=plain

# Or generate a resource module with CRUD operations
goose g module --name=posts --type=resource
```

## Next Steps

Now that you have a running application:

- [Directory Structure](directory-structure.md) - Understand project layout
- [Configuration](configuration.md) - Configure your application
- [Routing](../building-blocks/routing.md) - Learn about routing
- [Controllers](../building-blocks/controllers.md) - Build controllers
- [Modules](../core-concepts/modules.md) - Create custom modules
