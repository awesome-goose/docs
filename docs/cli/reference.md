# CLI Command Reference

Complete reference for all Goose CLI commands.

## Global Flags

These flags work with any command:

| Flag              | Description           |
| ----------------- | --------------------- |
| `--help`, `-h`    | Show help for command |
| `--version`, `-v` | Show CLI version      |

## Commands

### goose app

Create a new Goose application.

**Syntax:**

```bash
goose app --name=<name> --template=<template> [options]
```

**Required Flags:**
| Flag | Description | Values |
|------|-------------|--------|
| `--name` | Application name | Any valid Go package name |
| `--template` | Application template | `api`, `web`, `cli`, `multi`, `spa` |

**Optional Flags:**
| Flag | Description | Default |
|------|-------------|---------|
| `--framework` | Frontend framework (required with `--template=spa`) | `react`, `vue`, `svelte`, `ng` |
| `--path` | Output directory | Current directory |

**Examples:**

```bash
# Create an API application
goose app --name=myapi --template=api

# Create a web application
goose app --name=myweb --template=web

# Create a CLI application
goose app --name=mycli --template=cli

# Create a multi-platform application
goose app --name=mymulti --template=multi

# Create a single-page application (Go API + frontend as one service)
goose app --name=myspa --template=spa --framework=react

# Create in specific directory
goose app --name=myapp --template=api --path=/projects
```

**Output:**
Creates a directory with the project structure.

---

### goose generate module

Generate a new module in an existing project.

**Syntax:**

```bash
goose generate module --name=<name> --type=<type> [options]
goose g module --name=<name> --type=<type> [options]
```

**Required Flags:**
| Flag | Description | Values |
|------|-------------|--------|
| `--name` | Module name | Any valid Go package name |
| `--type` | Module type | `plain`, `resource` |

**Optional Flags:**
| Flag | Description | Default |
|------|-------------|---------|
| `--template` | Platform template | Auto-detected from project |

**Examples:**

```bash
# Generate a plain module
goose g module --name=auth --type=plain

# Generate a resource module
goose g module --name=products --type=resource

# Specify platform explicitly
goose g module --name=users --type=resource --template=api
goose g module --name=pages --type=plain --template=web
goose g module --name=commands --type=plain --template=cli
```

**Output:**
Creates module directory under `app/` with generated files.

---

### goose --version

Display CLI version information.

**Syntax:**

```bash
goose --version
goose -v
```

**Output:**

```
goose version 0.0.0
```

---

### goose --help

Display help information.

**Syntax:**

```bash
goose --help
goose -h
goose <command> --help
```

**Output:**
Displays available commands and flags.

## Template Reference

### API Template Structure

```
app-name/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    ├── app.dtos.go
    ├── jobs/
    │   └── sample.job.go
    └── queries/
        └── sample.query.go
```

### Web Template Structure

```
app-name/
├── .env
├── go.mod
├── main.go
└── app/
    ├── app.module.go
    ├── app.controller.go
    ├── app.service.go
    ├── app.routes.go
    ├── app.dtos.go
    └── templates/
        ├── base/
        │   └── layout.html
        ├── pages/
        │   └── home.html
        └── partials/
            ├── header.html
            └── footer.html
```

### CLI Template Structure

```
app-name/
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

### Multi Template Structure

```
app-name/
├── .env
├── go.mod
├── main.go
└── app/
    ├── api/
    │   ├── api.module.go
    │   ├── api.controller.go
    │   ├── api.routes.go
    │   └── api.service.go
    ├── web/
    │   ├── web.module.go
    │   ├── web.controller.go
    │   ├── web.routes.go
    │   ├── web.service.go
    │   └── templates/
    ├── cli/
    │   ├── cli.module.go
    │   ├── cli.controller.go
    │   ├── cli.routes.go
    │   └── cli.service.go
    └── shared/
        ├── shared.module.go
        ├── shared.service.go
        └── entities/
```

## Module Types Reference

### Plain Module

Generated files:

- `module.go` - Module definition
- `controller.go` - Request handlers
- `service.go` - Business logic
- `routes.go` - Route definitions
- `dtos.go` - Data transfer objects

### Resource Module

Generated files (all from plain module, plus):

- `entity.go` - Database entity definition
- `templates/` (web only) - HTML templates

## Exit Codes

| Code | Meaning           |
| ---- | ----------------- |
| 0    | Success           |
| 1    | General error     |
| 2    | Invalid arguments |

## Environment Variables

The CLI respects these environment variables:

| Variable     | Description                         |
| ------------ | ----------------------------------- |
| `GOOSE_PATH` | Default output path for `goose app` |

## Troubleshooting

### Common Issues

**Module detection fails:**

```bash
# Specify template explicitly
goose g module --name=users --type=resource --template=api
```

**Permission denied:**

```bash
# Check directory permissions
chmod 755 ./app
```

**Go mod errors after generation:**

```bash
# Tidy dependencies
go mod tidy
```

## See Also

- [Creating Applications](creating-apps.md)
- [Generating Modules](generating-modules.md)
- [Installation](installation.md)
