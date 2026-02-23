# Multi-Platform Applications

Goose supports running multiple platforms (API, Web, CLI) in a single application. This allows you to share code, services, and resources across different interfaces while maintaining separate configurations for each platform.

## Overview

A multi-platform application can have:

- **Multiple API servers** - Running on different ports concurrently
- **Multiple Web servers** - Running on different ports concurrently
- **One CLI interface** - Invoked via command-line arguments

## Quick Start

### Create a Multi-Platform App

```bash
goose app --name=myapp --template=multi
cd myapp
go mod tidy
```

### Run the Application

```bash
# Start API and Web servers concurrently
go run main.go

# Run CLI commands
go run main.go cli help
go run main.go cli info
```

## Basic Usage

```go
package main

import (
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
    "github.com/awesome-goose/goose/platforms/cli"
    "github.com/awesome-goose/goose/platforms/web"
    "github.com/awesome-goose/goose/types"
)

func main() {
    stop, err := goose.Start(
        // API platform - runs on port 8080
        goose.API(
            api.NewPlatform(api.WithName("my-api"), api.WithPort(8080)),
            &app.ApiModule{},
            commonInitializers,
        ),
        // Web platform - runs on port 3000
        goose.Web(
            web.NewPlatform(web.WithName("my-web"), web.WithPort(3000)),
            &app.WebModule{},
            commonInitializers,
        ),
        // CLI platform - invoked with `cli` argument
        goose.CLI(
            cli.NewPlatform(cli.WithName("my-cli")),
            &cliApp.CliModule{},
            commonInitializers,
        ),
    )
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Platform Instances

### goose.API()

Creates an API platform instance. Multiple API instances can run concurrently on different ports.

```go
goose.API(
    api.NewPlatform(
        api.WithName("users-api"),
        api.WithPort(8080),
    ),
    &UsersApiModule{},
    initializers,
)
```

### goose.Web()

Creates a Web platform instance. Multiple Web instances can run concurrently on different ports.

```go
goose.Web(
    web.NewPlatform(
        web.WithName("admin-web"),
        web.WithPort(3000),
    ),
    &AdminWebModule{},
    initializers,
)
```

### goose.CLI()

Creates a CLI platform instance. Only one CLI instance is allowed per application.

```go
goose.CLI(
    cli.NewPlatform(
        cli.WithName("admin-cli"),
    ),
    &AdminCliModule{},
    initializers,
)
```

## Execution Behavior

### Server Mode (Default)

When running `go run main.go`:

- All API and Web platform instances start concurrently
- Each runs in its own goroutine
- The application waits for shutdown signals (SIGINT, SIGTERM)
- Graceful shutdown is handled automatically

### CLI Mode

When running `go run main.go cli <command>`:

- Only the CLI platform instance is executed
- Runs in the main goroutine (blocking)
- Server instances are not started
- Exits after command completion

## Project Structure

A typical multi-platform app structure:

```
myapp/
├── main.go                 # Application entry point
├── go.mod
├── .env
└── app/
    ├── api/                    # API platform
    │   ├── api.module.go       # API platform module
    │   ├── api.controller.go   # API request handlers
    │   ├── api.service.go      # API business logic
    │   └── api.routes.go       # API route definitions
    │
    ├── web/                    # Web platform
    │   ├── web.module.go       # Web platform module
    │   ├── web.controller.go   # Web request handlers
    │   ├── web.service.go      # Web business logic
    │   └── web.routes.go       # Web route definitions
    │
    ├── cli/                    # CLI platform
    │   ├── cli.module.go       # CLI platform module
    │   ├── cli.controller.go   # CLI command handlers
    │   ├── cli.service.go      # CLI business logic
    │   └── cli.routes.go       # CLI command definitions
    │
    └── shared/                 # Shared across all platforms
        ├── shared.module.go
        └── shared.service.go
```

## Sharing Code

### Shared Modules

Create a shared module for services used across platforms:

```go
// app/shared/shared.module.go
package shared

type SharedModule struct{}

func (m *SharedModule) Imports() []types.Module {
    return []types.Module{}
}

func (m *SharedModule) Exports() []any {
    return []any{&UserService{}, &AuthService{}}
}

func (m *SharedModule) Declarations() []any {
    return []any{&UserService{}, &AuthService{}}
}
```

Import the shared module in platform-specific modules:

```go
// app/api.module.go
func (m *ApiModule) Imports() []types.Module {
    return []types.Module{
        &shared.SharedModule{},
        API_ROUTES,
    }
}
```

### Shared Initializers

Use the same initializers across platforms for common services:

```go
var sharedInitializers = []func(container types.Container) error{
    func(container types.Container) error {
        return container.Register(
            func() types.Log { return log.NewLog(...) },
            "", true,
        )
    },
    func(container types.Container) error {
        return container.Register(
            func() *database.DB { return database.Connect(...) },
            "", true,
        )
    },
}
```

## Error Handling

The framework handles errors gracefully:

- If any instance fails to start, all others are stopped
- Runtime errors in one instance trigger shutdown of all instances
- The stop function returns the first error encountered

```go
stop, err := goose.Start(instances...)
if err != nil {
    // Handle startup error
    log.Fatal(err)
}
defer func() {
    if err := stop(); err != nil {
        // Handle shutdown error
        log.Printf("Shutdown error: %v", err)
    }
}()
```

## Best Practices

1. **Separate Concerns**: Keep platform-specific code in separate directories
2. **Share Wisely**: Only share code that truly needs to be shared
3. **Use Shared Initializers**: Register common services once in shared initializers
4. **Configure Ports**: Ensure API and Web servers use different ports
5. **Graceful Shutdown**: Always call the stop function to ensure clean shutdown

## Migration from Single Platform

To convert an existing single-platform app to multi-platform:

1. Move platform-specific code to appropriate directories
2. Create shared module for common services
3. Update `main.go` to use `goose.Start()`
4. Add platform-specific modules and routes
5. Test each platform independently
