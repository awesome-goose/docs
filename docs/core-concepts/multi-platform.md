# Multi-Platform Support

Goose supports building applications for multiple platforms: API, Web, and CLI. You can run them separately or combine them in a single application.

## Platform Types

### API Platform

RESTful API applications that return JSON responses.

```go
import "github.com/awesome-goose/goose/platforms/api"

platform := api.NewPlatform(
    api.WithName("myapi"),
    api.WithHost("0.0.0.0"),
    api.WithPort(8080),
)
```

### Web Platform

Server-rendered web applications with HTML templates.

```go
import "github.com/awesome-goose/goose/platforms/web"

platform := web.NewPlatform(
    web.WithName("myweb"),
    web.WithHost("0.0.0.0"),
    web.WithPort(3000),
    web.WithTemplatesDir("./templates"),
)
```

### CLI Platform

Command-line tools and utilities.

```go
import "github.com/awesome-goose/goose/platforms/cli"

platform := cli.NewPlatform(
    cli.WithName("mycli"),
    cli.WithVersion("0.0.0"),
)
```

## Single Platform Applications

### API Application

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    platform := api.NewPlatform()
    module := &app.AppModule{}

    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

### Web Application

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/web"
)

func main() {
    platform := web.NewPlatform()
    module := &app.AppModule{}

    stop, err := goose.Start(goose.Web(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

### CLI Application

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/cli"
)

func main() {
    platform := cli.NewPlatform()
    module := &app.AppModule{}

    stop, err := goose.Start(goose.CLI(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Multi-Platform Application

Run multiple platforms from a single codebase:

```go
package main

import (
    apiApp "myapp/app/api"
    webApp "myapp/app/web"
    cliApp "myapp/app/cli"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
    "github.com/awesome-goose/goose/platforms/web"
    "github.com/awesome-goose/goose/platforms/cli"
)

func main() {
    // Create platforms
    apiPlatform := api.NewPlatform(
        api.WithPort(8080),
    )
    webPlatform := web.NewPlatform(
        web.WithPort(3000),
    )
    cliPlatform := cli.NewPlatform()

    // Create modules
    apiModule := &apiApp.ApiModule{}
    webModule := &webApp.WebModule{}
    cliModule := &cliApp.CliModule{}

    // Start all platforms
    stop, err := goose.Start(
        goose.API(apiPlatform, apiModule, nil),
        goose.Web(webPlatform, webModule, nil),
        goose.CLI(cliPlatform, cliModule, nil),
    )
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Running Multi-Platform Apps

### Server Mode (Default)

When running without arguments, API and Web platforms start concurrently:

```bash
go run main.go
# Starts API on :8080 and Web on :3000
```

### CLI Mode

To run CLI commands, pass `cli` as the first argument:

```bash
go run main.go cli <command>
# Runs the CLI with the specified command
```

## Shared Code

Share code between platforms using a shared module:

### Project Structure

```
app/
├── api/
│   ├── api.module.go
│   ├── api.controller.go
│   └── api.routes.go
├── web/
│   ├── web.module.go
│   ├── web.controller.go
│   ├── web.routes.go
│   └── templates/
├── cli/
│   ├── cli.module.go
│   ├── cli.controller.go
│   └── cli.routes.go
└── shared/
    ├── shared.module.go
    ├── shared.service.go
    ├── entities/
    │   └── user.entity.go
    └── repositories/
        └── user.repository.go
```

### Shared Module

```go
// app/shared/shared.module.go
package shared

import "github.com/awesome-goose/goose/types"

type SharedModule struct{}

func (m *SharedModule) Imports() []types.Module {
    return []types.Module{}
}

func (m *SharedModule) Exports() []any {
    return []any{
        &UserService{},
        &ProductService{},
    }
}

func (m *SharedModule) Declarations() []any {
    return []any{
        &UserService{},
        &ProductService{},
        &UserRepository{},
        &ProductRepository{},
    }
}
```

### Using Shared Module

```go
// app/api/api.module.go
package api

import (
    "myapp/app/shared"
    "github.com/awesome-goose/goose/types"
)

type ApiModule struct{}

func (m *ApiModule) Imports() []types.Module {
    return []types.Module{
        &shared.SharedModule{},  // Import shared services
    }
}

func (m *ApiModule) Declarations() []any {
    return []any{
        &ApiController{},
    }
}
```

## Platform-Specific Responses

### API Responses (JSON)

Each platform binds a DTO and returns a `types.Output`; only the `output.*` helper differs.

```go
type GetUserDto struct {
    ID string `param:"id"`
}

func (c *ApiController) GetUser(dto *GetUserDto) types.Output {
    user, err := c.service.GetUser(dto.ID)
    if err != nil {
        return output.NotFound(err.Error())
    }
    return output.JSON(user)
}
```

### Web Responses (HTML)

```go
func (c *WebController) GetUser(dto *GetUserDto) types.Output {
    user, err := c.service.GetUser(dto.ID)
    if err != nil {
        return output.View("pages/error.html", map[string]any{
            "message": err.Error(),
        }, output.WithHTMLCode(404))
    }
    return output.View("pages/users/show.html", map[string]any{
        "user": user,
    })
}
```

### CLI Responses (Console)

The CLI reads the id as a route parameter (`/user/:id`):

```go
func (c *CliController) GetUser(dto *GetUserDto) types.Output {
    user, err := c.service.GetUser(dto.ID)
    if err != nil {
        return output.ConsoleError(fmt.Sprintf("Error: %s", err.Error()))
    }
    return output.ConsoleSuccess(fmt.Sprintf("User: %s (%s)", user.Name, user.Email))
}
```

## Use Cases

### Typical Multi-Platform Scenarios

1. **Admin Tools**: API for frontend + CLI for admin tasks
2. **Full-Stack**: API + Web for server-rendered pages
3. **DevOps**: CLI for deployment tools + API for monitoring
4. **Migration**: Gradually moving from Web to API

### Example: E-Commerce

```go
// API: Mobile app clients
goose.API(apiPlatform, apiModule, nil)

// SPA: Storefront single-page app + its JSON API as one service
goose.SPA(spaPlatform, spaModule, nil)

// Web: Server-rendered store pages
goose.Web(webPlatform, webModule, nil)

// CLI: Inventory management, reports
goose.CLI(cliPlatform, cliModule, nil)
```

SPA instances run as server instances alongside API and Web — see the
[SPA Platform](../platforms/spa.md) guide.

## Best Practices

### 1. Keep Controllers Thin

Platform-specific controllers should only handle request/response:

```go
// API Controller — dto is bound automatically
func (c *ApiUserController) Create(dto *CreateUserDTO) types.Output {
    return output.Created(c.service.CreateUser(dto)) // Delegate to shared service
}

// Web Controller — same service, different response
func (c *WebUserController) Create(dto *CreateUserDTO) types.Output {
    user, err := c.service.CreateUser(dto) // Same service
    if err != nil {
        return output.Redirect("/users/new").WithError(err.Error())
    }
    return output.Redirect("/users/" + user.ID)
}
```

### 2. Centralize Business Logic

Put all business logic in shared services:

```go
// shared/user.service.go
func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    // Validation
    // Business rules
    // Database operations
    // All platforms use this
}
```

### 3. Use Appropriate Response Types

Each platform has its own response handling:

```go
// API: JSON
return output.JSON(user)

// Web: render a template
return output.View("pages/users/show.html", data)

// CLI: formatted console output
return output.ConsoleSuccess(fmt.Sprintf("Created user: %s", user.Name))
```

## Next Steps

- [API Platform](../platforms/api.md) - Detailed API docs
- [Web Platform](../platforms/web.md) - Detailed Web docs
- [CLI Platform](../platforms/cli.md) - Detailed CLI docs
