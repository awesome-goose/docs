# Directory Structure

Understanding the Goose project structure helps you organize your code effectively.

## API Application Structure

```
myapi/
├── .env                    # Environment variables
├── .gitignore             # Git ignore rules
├── go.mod                 # Go module definition
├── go.sum                 # Go dependencies checksum
├── main.go                # Application entry point
├── config/                # Configuration files (optional)
│   ├── app.yaml
│   └── database.yaml
├── app/                   # Main application module
│   ├── app.module.go      # Module definition
│   ├── app.controller.go  # Request handlers
│   ├── app.service.go     # Business logic
│   ├── app.routes.go      # Route definitions
│   ├── app.dtos.go        # Data transfer objects
│   ├── jobs/              # Background jobs (optional)
│   │   └── sample.job.go
│   └── queries/           # Database queries (optional)
│       └── sample.query.go
└── modules/               # Additional modules
    └── users/
        ├── users.module.go
        ├── users.controller.go
        ├── users.service.go
        ├── users.routes.go
        └── users.dtos.go
```

## Web Application Structure

```
myweb/
├── .env
├── go.mod
├── main.go
├── app/
│   ├── app.module.go
│   ├── app.controller.go
│   ├── app.service.go
│   ├── app.routes.go
│   ├── app.dtos.go
│   └── templates/         # HTML templates
│       ├── base/
│       │   └── layout.html
│       ├── pages/
│       │   └── home.html
│       └── partials/
│           ├── header.html
│           └── footer.html
├── static/                # Static assets
│   ├── css/
│   ├── js/
│   └── images/
└── modules/
```

## CLI Application Structure

```
mycli/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go      # Command definitions
    └── app.dtos.go
```

## Multi-Platform Structure

```
mymulti/
├── .env
├── go.mod
├── main.go
└── app/
    ├── api/               # API platform
    │   ├── api.module.go
    │   ├── api.controller.go
    │   ├── api.routes.go
    │   └── api.service.go
    ├── web/               # Web platform
    │   ├── web.module.go
    │   ├── web.controller.go
    │   ├── web.routes.go
    │   ├── web.service.go
    │   └── templates/
    ├── cli/               # CLI platform
    │   ├── cli.module.go
    │   ├── cli.controller.go
    │   ├── cli.routes.go
    │   └── cli.service.go
    └── shared/            # Shared components
        ├── shared.module.go
        ├── shared.service.go
        └── entities/
```

## Module Structure

### Plain Module

A basic module for simple functionality:

```
users/
├── users.module.go        # Module definition
├── users.controller.go    # Request handlers
├── users.service.go       # Business logic
├── users.routes.go        # Route definitions
└── users.dtos.go          # Data transfer objects
```

### Resource Module

A module with database entity:

```
products/
├── products.module.go
├── products.controller.go
├── products.service.go
├── products.routes.go
├── products.dtos.go
├── products.entity.go     # Database entity
├── migrations/            # Database migrations
│   └── 001_create_products.go
└── seeds/                 # Seed data
    └── products.seed.go
```

## File Naming Conventions

| File              | Purpose               | Example              |
| ----------------- | --------------------- | -------------------- |
| `*.module.go`     | Module definition     | `app.module.go`      |
| `*.controller.go` | Request handlers      | `app.controller.go`  |
| `*.service.go`    | Business logic        | `app.service.go`     |
| `*.routes.go`     | Route definitions     | `app.routes.go`      |
| `*.dtos.go`       | Data transfer objects | `app.dtos.go`        |
| `*.entity.go`     | Database entities     | `user.entity.go`     |
| `*.middleware.go` | Middleware            | `auth.middleware.go` |
| `*.job.go`        | Background jobs       | `email.job.go`       |

## Key Files Explained

### main.go

The entry point of your application:

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

### app.module.go

Defines the module structure:

```go
package app

import "github.com/awesome-goose/goose/types"

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        // Import other modules
    }
}

func (m *AppModule) Exports() []any {
    return []any{
        // Export services for other modules
    }
}

func (m *AppModule) Declarations() []any {
    return []any{
        &AppController{},
        &AppService{},
    }
}
```

### .env

Environment configuration:

```env
APP_NAME=myapp
APP_ENV=development
HOST=localhost
PORT=8080
```

## Organizing Large Applications

For large applications, organize by feature:

```
app/
├── app.module.go          # Root module
├── auth/                  # Authentication feature
│   ├── auth.module.go
│   ├── auth.controller.go
│   ├── auth.service.go
│   └── auth.middleware.go
├── users/                 # Users feature
├── products/              # Products feature
├── orders/                # Orders feature
└── shared/                # Shared utilities
    ├── shared.module.go
    ├── validators/
    └── helpers/
```

## Next Steps

- [Modules](../core-concepts/modules.md) - Deep dive into modules
- [Controllers](../building-blocks/controllers.md) - Controller patterns
- [Services](../building-blocks/services.md) - Service layer
