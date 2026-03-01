# CLI Platform

Build command-line applications with the Goose CLI platform.

## Overview

The CLI platform handles command-line argument parsing and executing commands.

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
        cli.WithVersion("1.0.0"),
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
    cli.WithVersion("1.0.0"),
    cli.WithDescription("My CLI application"),
)
```

## Defining Commands

### Basic Command

```go
type GreetController struct{}

func (c *GreetController) Routes() types.Routes {
    return types.Routes{
        {Path: "greet", Handler: c.Greet},
    }
}

func (c *GreetController) Greet(ctx types.Context) any {
    name := ctx.Arg(0)
    if name == "" {
        name = "World"
    }
    return fmt.Sprintf("Hello, %s!", name)
}
```

**Usage:**

```bash
myapp greet John
# Output: Hello, John!
```

### Command with Flags

```go
func (c *Controller) Routes() types.Routes {
    return types.Routes{
        {
            Path:    "create",
            Handler: c.Create,
            Meta: map[string]interface{}{
                "flags": []cli.Flag{
                    {Name: "name", Short: "n", Required: true},
                    {Name: "force", Short: "f", Type: "bool"},
                },
            },
        },
    }
}

func (c *Controller) Create(ctx types.Context) any {
    name := ctx.Flag("name")
    force := ctx.FlagBool("force")

    if force {
        return fmt.Sprintf("Force creating: %s", name)
    }
    return fmt.Sprintf("Creating: %s", name)
}
```

**Usage:**

```bash
myapp create --name=myproject
myapp create -n myproject -f
```

### Subcommands

```go
type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Path: "user/list", Handler: c.List},
        {Path: "user/create", Handler: c.Create},
        {Path: "user/delete", Handler: c.Delete},
    }
}

func (c *UserController) List(ctx types.Context) any {
    users := c.service.GetAll()

    for _, user := range users {
        fmt.Printf("%s - %s\n", user.ID, user.Email)
    }

    return nil
}

func (c *UserController) Create(ctx types.Context) any {
    email := ctx.Flag("email")
    name := ctx.Flag("name")

    user, err := c.service.Create(CreateUserDTO{
        Email: email,
        Name:  name,
    })

    if err != nil {
        return fmt.Sprintf("Error: %v", err)
    }

    return fmt.Sprintf("Created user: %s", user.ID)
}
```

**Usage:**

```bash
myapp user list
myapp user create --email=user@example.com --name=John
myapp user delete 123
```

## Context Methods

### Arguments

```go
func (c *Controller) Run(ctx types.Context) any {
    // Get argument by position
    first := ctx.Arg(0)
    second := ctx.Arg(1)

    // Get all arguments
    args := ctx.Args()

    return nil
}
```

### Flags

```go
func (c *Controller) Run(ctx types.Context) any {
    // String flag
    name := ctx.Flag("name")

    // Bool flag
    verbose := ctx.FlagBool("verbose")

    // Int flag
    count := ctx.FlagInt("count")

    // With defaults
    port := ctx.FlagInt("port")
    if port == 0 {
        port = 8080
    }

    return nil
}
```

## Output

### Console Output

```go
func (c *Controller) Run(ctx types.Context) any {
    // Return string to print
    return "Hello, World!"
}

func (c *Controller) RunComplex(ctx types.Context) any {
    // Print directly
    fmt.Println("Processing...")

    // Return formatted output
    return fmt.Sprintf("Completed %d items", count)
}
```

### Formatted Tables

```go
func (c *UserController) List(ctx types.Context) any {
    users := c.service.GetAll()

    // Print header
    fmt.Printf("%-36s %-30s %-20s\n", "ID", "EMAIL", "NAME")
    fmt.Println(strings.Repeat("-", 86))

    // Print rows
    for _, user := range users {
        fmt.Printf("%-36s %-30s %-20s\n", user.ID, user.Email, user.Name)
    }

    return nil
}
```

## Interactive Input

```go
import "bufio"

func (c *Controller) CreateInteractive(ctx types.Context) any {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter email: ")
    email, _ := reader.ReadString('\n')
    email = strings.TrimSpace(email)

    fmt.Print("Enter name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    // Create user...
    return fmt.Sprintf("Created user with email: %s", email)
}
```

## Progress Indicators

```go
func (c *Controller) Process(ctx types.Context) any {
    items := []string{"item1", "item2", "item3"}
    total := len(items)

    for i, item := range items {
        // Process item...
        time.Sleep(time.Second)

        // Show progress
        percent := float64(i+1) / float64(total) * 100
        fmt.Printf("\rProgress: %.0f%% (%d/%d)", percent, i+1, total)
    }

    fmt.Println("\nComplete!")
    return nil
}
```

## Error Handling

```go
func (c *Controller) Run(ctx types.Context) any {
    id := ctx.Arg(0)
    if id == "" {
        return errors.New("user ID is required")
    }

    user, err := c.service.GetByID(id)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }

    return fmt.Sprintf("User: %s", user.Email)
}
```

## Exit Codes

```go
func (c *Controller) Run(ctx types.Context) any {
    err := c.doSomething()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    return "Success"
}
```

## Common Patterns

### CRUD Commands

```go
func (c *UserController) Routes() types.Routes {
    return types.Routes{
        {Path: "users", Handler: c.List},
        {Path: "users/show", Handler: c.Show},
        {Path: "users/create", Handler: c.Create},
        {Path: "users/update", Handler: c.Update},
        {Path: "users/delete", Handler: c.Delete},
    }
}
```

### Database Commands

```go
type DatabaseController struct {
    db *gorm.DB `inject:""`
}

func (c *DatabaseController) Routes() types.Routes {
    return types.Routes{
        {Path: "db/migrate", Handler: c.Migrate},
        {Path: "db/seed", Handler: c.Seed},
        {Path: "db/reset", Handler: c.Reset},
    }
}

func (c *DatabaseController) Migrate(ctx types.Context) any {
    c.db.AutoMigrate(&User{}, &Post{})
    return "Migrations complete"
}
```

### Help Text

Commands automatically get help:

```bash
myapp --help
myapp user --help
myapp user create --help
```

## Best Practices

1. **Provide clear help text** for all commands
2. **Use consistent flag naming** across commands
3. **Handle errors gracefully** with clear messages
4. **Show progress** for long operations
5. **Use exit codes** appropriately
6. **Support both flags and arguments** where sensible

## Next Steps

- [CLI Tutorial](../tutorials/cli-tool.md) - Build a CLI tool
- [Services](../building-blocks/services.md) - Business logic
- [Configuration](../getting-started/configuration.md) - App configuration
