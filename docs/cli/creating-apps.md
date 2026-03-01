# Creating Applications

Use the Goose CLI to scaffold new applications with proper structure and boilerplate.

## Basic Usage

```bash
goose app --name=<app-name> --template=<template>
```

## Available Templates

### API Template

Create a RESTful API application:

```bash
goose app --name=myapi --template=api
```

**Generated structure:**

```
myapi/
├── .env                    # Environment configuration
├── go.mod                  # Go module file
├── main.go                 # Application entry point
└── app/
    ├── app.module.go       # Root module
    ├── app.controller.go   # Main controller
    ├── app.service.go      # Main service
    ├── app.routes.go       # Route definitions
    ├── app.dtos.go         # Data transfer objects
    ├── jobs/               # Background jobs
    │   └── sample.job.go
    └── queries/            # Database queries
        └── sample.query.go
```

**Features included:**

- JSON API server
- Database support (SQLite default)
- Background job examples
- Query examples

### Web Template

Create a server-rendered web application:

```bash
goose app --name=myweb --template=web
```

**Generated structure:**

```
myweb/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    ├── app.dtos.go
    └── templates/          # HTML templates
        ├── base/
        │   └── layout.html
        ├── pages/
        │   └── home.html
        └── partials/
            ├── header.html
            └── footer.html
```

**Features included:**

- HTML template rendering
- Layout system
- Partial templates
- Static file serving

### CLI Template

Create a command-line application:

```bash
goose app --name=mycli --template=cli
```

**Generated structure:**

```
mycli/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go       # Command definitions
    └── app.dtos.go
```

**Features included:**

- Command routing
- Argument parsing
- Console output

### Multi-Platform Template

Create an application with API, Web, and CLI:

```bash
goose app --name=mymulti --template=multi
```

**Generated structure:**

```
mymulti/
├── .env
├── go.mod
├── main.go
└── app/
    ├── api/                # API platform
    │   ├── api.module.go
    │   ├── api.controller.go
    │   ├── api.routes.go
    │   └── api.service.go
    ├── web/                # Web platform
    │   ├── web.module.go
    │   ├── web.controller.go
    │   ├── web.routes.go
    │   ├── web.service.go
    │   └── templates/
    ├── cli/                # CLI platform
    │   ├── cli.module.go
    │   ├── cli.controller.go
    │   ├── cli.routes.go
    │   └── cli.service.go
    └── shared/             # Shared code
        ├── shared.module.go
        ├── shared.service.go
        └── entities/
```

## Command Options

```bash
goose app --name=<name> --template=<template> [--path=<directory>]
```

| Flag         | Description      | Required | Default           |
| ------------ | ---------------- | -------- | ----------------- |
| `--name`     | Application name | Yes      | -                 |
| `--template` | Template type    | Yes      | -                 |
| `--path`     | Output directory | No       | Current directory |

### Examples

```bash
# Create in current directory
goose app --name=myapi --template=api

# Create in specific directory
goose app --name=myapi --template=api --path=/home/user/projects

# Create web app
goose app --name=blog --template=web

# Create CLI tool
goose app --name=tools --template=cli

# Create multi-platform app
goose app --name=platform --template=multi
```

## After Creating

### 1. Navigate to Project

```bash
cd myapi
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Edit `.env` file:

```env
APP_NAME=myapi
APP_ENV=development
HOST=localhost
PORT=8080
```

### 4. Run the Application

```bash
go run main.go
```

### 5. Test It

**API:**

```bash
curl http://localhost:8080/
```

**Web:**
Open `http://localhost:3000` in browser

**CLI:**

```bash
go run main.go cli <command>
```

## Generated Code Overview

### main.go (API Template)

```go
package main

import (
    "myapi/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    platform := api.NewPlatform(
        api.WithHost("localhost"),
        api.WithPort(8080),
    )

    module := &app.AppModule{}

    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

### app.module.go

```go
package app

import "github.com/awesome-goose/goose/types"

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{}
}

func (m *AppModule) Exports() []any {
    return []any{}
}

func (m *AppModule) Declarations() []any {
    return []any{
        &AppController{},
        &AppService{},
    }
}
```

## Customizing Templates

After scaffolding, customize:

1. **Update `.env`** with your configuration
2. **Modify `app.controller.go`** for your endpoints
3. **Add modules** with `goose g module`
4. **Configure database** if needed

## Next Steps

- [Generating Modules](generating-modules.md) - Add features
- [Directory Structure](../getting-started/directory-structure.md) - Understand layout
- [Quick Start](../getting-started/quick-start.md) - Build your first feature
