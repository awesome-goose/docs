# Goose CLI Overview

The Goose CLI is a command-line tool for scaffolding and managing Goose framework applications.

## Features

- **Application Scaffolding** - Create new API, Web, SPA, or CLI applications
- **Module Generation** - Add modules to existing projects
- **Code Generation** - Generate controllers, services, entities
- **Project Management** - Manage application configuration

## Quick Start

```bash
# Install the CLI
go install github.com/awesome-goose/cli@latest

# Create a new API application
goose app --name=myapi --template=api

# Navigate to project
cd myapi

# Run the application
go mod tidy && go run main.go
```

## Available Commands

| Command                 | Description                       |
| ----------------------- | --------------------------------- |
| `goose app`             | Create a new application          |
| `goose g module`        | Generate a new module             |
| `goose generate module` | Generate a new module (full form) |
| `goose --version`       | Show CLI version                  |
| `goose --help`          | Show help information             |

## Command Structure

```
goose <command> [subcommand] [flags]
```

### Examples

```bash
# Create applications
goose app --name=myapi --template=api
goose app --name=myweb --template=web
goose app --name=mycli --template=cli
goose app --name=mymulti --template=multi
goose app --name=myspa --template=spa --framework=react

# Generate modules
goose g module --name=users --type=plain
goose g module --name=products --type=resource
goose generate module --name=orders --type=resource --template=api
```

## Use Cases

### Starting a New Project

```bash
# Create a REST API
goose app --name=blog-api --template=api

# Project structure is created:
# blog-api/
# ├── .env
# ├── go.mod
# ├── main.go
# └── app/
#     ├── app.module.go
#     ├── app.controller.go
#     ├── app.service.go
#     ├── app.routes.go
#     └── app.dtos.go
```

### Adding Features

```bash
# Add a users module
cd blog-api
goose g module --name=users --type=resource

# Module is created:
# app/users/
# ├── users.module.go
# ├── users.controller.go
# ├── users.service.go
# ├── users.routes.go
# ├── users.dtos.go
# └── users.entity.go
```

### Multi-Platform Project

```bash
# Create a project with API, Web, and CLI
goose app --name=myproject --template=multi

# Includes separate modules for each platform:
# app/
# ├── api/
# ├── web/
# ├── cli/
# └── shared/
```

## Configuration

The CLI uses sensible defaults but can be customized:

```bash
# Specify output directory
goose app --name=myapp --template=api --path=/projects

# Auto-detect platform from existing project
goose g module --name=users --type=resource  # Detects from main.go
```

## Next Steps

- [Installation](installation.md) - Install the CLI
- [Creating Applications](creating-apps.md) - Create new projects
- [Generating Modules](generating-modules.md) - Add modules
- [Command Reference](reference.md) - Full command documentation
