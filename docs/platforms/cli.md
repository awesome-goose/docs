# CLI Platform

Build command-line applications with the Goose CLI platform.

## Overview

The CLI platform parses the process arguments, routes them to a command handler, and prints the handler's output. Positional arguments form the command **path**; `--flags` are bound to your DTO.

## Quick Start

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/cli"
)

func main() {
    platform := cli.NewPlatform(
        cli.WithName("myapp"),
        cli.WithVersion("0.0.0"),
    )

    module := &app.AppModule{}

    stop, err := goose.Start(goose.CLI(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Configuration

```go
platform := cli.NewPlatform(
    cli.WithName("myapp"),
    cli.WithVersion("0.0.0"),
    cli.WithDescription("My CLI application"),
)
```

## Defining Commands

Commands are routes registered with `router.Cli`. The handler is a `[]any{Controller{}, "Method"}` tuple, takes a DTO, and returns a `types.Output` — build these with the `output` console helpers.

### Basic Command

```go
import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/modules/router"
    "github.com/awesome-goose/goose/types"
)

type GreetController struct{}

type GreetDto struct {
    Name string `flag:"name"`
}

func (c *GreetController) Greet(dto *GreetDto) types.Output {
    name := dto.Name
    if name == "" {
        name = "World"
    }
    return output.ConsoleSuccess(fmt.Sprintf("Hello, %s!", name))
}
```

Register it:

```go
var ROUTES = router.ForRoutes(
    router.Cli("/greet", []any{GreetController{}, "Greet"}),
)
```

**Usage:**

```bash
myapp greet --name=John
# Output: Hello, John!
```

### Flags

Every `--flag` (and `-f` shorthand) binds to a DTO field via the `flag` tag. A bare `--force` with no value binds as `"true"`, so `bool` fields work:

```go
type CreateDto struct {
    Name  string `flag:"name"`
    Force bool   `flag:"force"`
    Count int    `flag:"count"`
}

func (c *Controller) Create(dto *CreateDto) types.Output {
    if dto.Force {
        return output.ConsoleSuccess(fmt.Sprintf("Force creating: %s", dto.Name))
    }
    return output.ConsoleSuccess(fmt.Sprintf("Creating: %s", dto.Name))
}
```

**Usage:**

```bash
myapp create --name=myproject
myapp create -n myproject -f
myapp create --name=myproject --count=3 --force
```

### Positional Arguments

Positional arguments make up the command path. To accept one as input, declare it as a route parameter and bind it with the `param` tag:

```go
// Route: /user/show/:id
type ShowDto struct {
    ID string `param:"id"`
}

func (c *UserController) Show(dto *ShowDto) types.Output {
    user, err := c.service.GetByID(dto.ID)
    if err != nil {
        return output.ConsoleError(fmt.Sprintf("User %s not found", dto.ID))
    }
    return output.ConsoleSuccess(fmt.Sprintf("User: %s", user.Email))
}
```

```go
var ROUTES = router.ForRoutes(
    router.Cli("/user/show/:id", []any{UserController{}, "Show"}),
)
```

**Usage:** `myapp user show 123`

### Subcommands

Nest commands with `router.Group` or slash-separated paths:

```go
type UserController struct {
    service *UserService `inject:""`
}

var ROUTES = router.ForRoutes(
    router.Group("/user",
        router.Cli("/list", []any{UserController{}, "List"}),
        router.Cli("/create", []any{UserController{}, "Create"}),
        router.Cli("/delete/:id", []any{UserController{}, "Delete"}),
    ),
)

func (c *UserController) List(dto *EmptyDto) types.Output {
    users := c.service.GetAll()
    rows := make([][]string, 0, len(users))
    for _, u := range users {
        rows = append(rows, []string{u.ID, u.Email})
    }
    return output.Table([]string{"ID", "Email"}, rows)
}

type CreateUserDto struct {
    Email string `flag:"email"`
    Name  string `flag:"name"`
}

func (c *UserController) Create(dto *CreateUserDto) types.Output {
    user, err := c.service.Create(dto)
    if err != nil {
        return output.ConsoleError(fmt.Sprintf("Error: %v", err))
    }
    return output.ConsoleSuccess(fmt.Sprintf("Created user: %s", user.ID))
}
```

**Usage:**

```bash
myapp user list
myapp user create --email=user@example.com --name=John
myapp user delete 123
```

## Output

Return an `output` console value from the handler rather than printing directly:

```go
func (c *Controller) Run(dto *EmptyDto) types.Output {
    return output.Line("Hello, World!")
}
```

The console helpers include:

| Helper                          | Purpose                          |
| ------------------------------- | -------------------------------- |
| `output.Line(msg)`              | A plain line                     |
| `output.ConsoleSuccess(msg)`    | Green success message            |
| `output.ConsoleError(msg)`      | Red error message (exit code 1)  |
| `output.Info` / `output.Warning`| Cyan / yellow status messages    |
| `output.Table(headers, rows)`   | A formatted table                |
| `output.List(items)`            | A bulleted list                  |
| `output.Box(title, lines)`      | A boxed panel                    |
| `output.ProgressBar(cur, total, width)` | A progress bar           |

### Formatted Tables

```go
func (c *UserController) List(dto *EmptyDto) types.Output {
    users := c.service.GetAll()

    rows := make([][]string, 0, len(users))
    for _, u := range users {
        rows = append(rows, []string{u.ID, u.Email, u.Name})
    }

    return output.Table([]string{"ID", "Email", "Name"}, rows)
}
```

## Interactive Input

Read from stdin inside the handler, then return an output value:

```go
import "bufio"

func (c *Controller) CreateInteractive(dto *EmptyDto) types.Output {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter email: ")
    email, _ := reader.ReadString('\n')
    email = strings.TrimSpace(email)

    fmt.Print("Enter name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    // Create user...
    return output.ConsoleSuccess(fmt.Sprintf("Created user with email: %s", email))
}
```

## Progress Indicators

For live progress, print in the loop and return the final result:

```go
func (c *Controller) Process(dto *EmptyDto) types.Output {
    items := []string{"item1", "item2", "item3"}
    total := len(items)

    for i := range items {
        // Process item...
        time.Sleep(time.Second)
        percent := float64(i+1) / float64(total) * 100
        fmt.Printf("\rProgress: %.0f%% (%d/%d)", percent, i+1, total)
    }
    fmt.Println()

    return output.ConsoleSuccess("Complete!")
}
```

## Error Handling & Exit Codes

`output.ConsoleError` prints in red and sets a non-zero exit code. For a custom code, use `output.ConsoleWithCode`:

```go
type RunDto struct {
    ID string `param:"id"`
}

func (c *Controller) Run(dto *RunDto) types.Output {
    if dto.ID == "" {
        return output.ConsoleError("user ID is required") // exit code 1
    }

    user, err := c.service.GetByID(dto.ID)
    if err != nil {
        return output.ConsoleWithCode(fmt.Sprintf("failed to get user: %v", err), 2)
    }

    return output.ConsoleSuccess(fmt.Sprintf("User: %s", user.Email))
}
```

## Common Patterns

### CRUD Commands

```go
var ROUTES = router.ForRoutes(
    router.Group("/users",
        router.Cli("/", []any{UserController{}, "List"}),
        router.Cli("/show/:id", []any{UserController{}, "Show"}),
        router.Cli("/create", []any{UserController{}, "Create"}),
        router.Cli("/update/:id", []any{UserController{}, "Update"}),
        router.Cli("/delete/:id", []any{UserController{}, "Delete"}),
    ),
)
```

### Database Commands

```go
type DatabaseController struct {
    db *gorm.DB `inject:""`
}

var ROUTES = router.ForRoutes(
    router.Group("/db",
        router.Cli("/migrate", []any{DatabaseController{}, "Migrate"}),
        router.Cli("/seed", []any{DatabaseController{}, "Seed"}),
        router.Cli("/reset", []any{DatabaseController{}, "Reset"}),
    ),
)

func (c *DatabaseController) Migrate(dto *EmptyDto) types.Output {
    c.db.AutoMigrate(&User{}, &Post{})
    return output.ConsoleSuccess("Migrations complete")
}
```

## Best Practices

1. **Bind flags with `flag` tags** rather than parsing arguments by hand
2. **Use consistent flag naming** across commands
3. **Return `output.ConsoleError`** for failures so the exit code is non-zero
4. **Show progress** for long operations
5. **Group related subcommands** under a shared path

## Next Steps

- [CLI Tutorial](../tutorials/cli-tool.md) - Build a CLI tool
- [Services](../building-blocks/services.md) - Business logic
- [Configuration](../getting-started/configuration.md) - App configuration
