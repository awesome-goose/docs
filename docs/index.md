# Goose Framework Documentation

Welcome to the official documentation for **Goose** — an awesome, all-in-one, modular, and scalable Go framework for building modern web, API, and CLI applications.

---

## Why Goose?

- **🚀 Fast** - Built on Go for maximum performance and efficiency
- **📦 Modular** - First-class module system for organized, reusable code
- **🔧 Flexible** - Build APIs, web apps, or CLI tools with the same patterns
- **📝 Declarative** - Convention over configuration for predictable behavior
- **🧪 Testable** - Built-in testing support for reliable applications
- **🌐 Multi-Platform** - Single codebase for API, Web, and CLI platforms

---

## Table of Contents

### 📚 Getting Started

| Document                                                      | Description                             |
| ------------------------------------------------------------- | --------------------------------------- |
| [Introduction](getting-started/introduction.md)               | What is Goose and why use it            |
| [Installation](getting-started/installation.md)               | How to install Goose and prerequisites  |
| [Quick Start](getting-started/quick-start.md)                 | Build your first application in minutes |
| [Configuration](getting-started/configuration.md)             | Configure your Goose application        |
| [Directory Structure](getting-started/directory-structure.md) | Understanding project layout            |

### 🏗️ Core Concepts

| Document                                                      | Description                         |
| ------------------------------------------------------------- | ----------------------------------- |
| [Architecture Overview](core-concepts/architecture.md)        | Understanding Goose architecture    |
| [Modules](core-concepts/modules.md)                           | Building modular applications       |
| [Dependency Injection](core-concepts/dependency-injection.md) | IoC container and service injection |
| [Lifecycle](core-concepts/lifecycle.md)                       | Application lifecycle and hooks     |
| [Multi-Platform Support](core-concepts/multi-platform.md)     | API, Web, and CLI platforms         |

### 🧱 Building Blocks

| Document                                                | Description                         |
| ------------------------------------------------------- | ----------------------------------- |
| [Routing](building-blocks/routing.md)                   | Defining routes and URL handling    |
| [Controllers](building-blocks/controllers.md)           | Request handlers and business logic |
| [Middleware](building-blocks/middleware.md)             | Request/response interceptors       |
| [Requests](building-blocks/requests.md)                 | Handling incoming HTTP requests     |
| [Responses](building-blocks/responses.md)               | Sending HTTP responses              |
| [DTOs & Validation](building-blocks/validation.md)      | Data transfer and validation        |
| [Services](building-blocks/services.md)                 | Business logic layer                |
| [Error Handling](building-blocks/error-handling.md)     | Managing errors gracefully          |
| [Logging](building-blocks/logging.md)                   | Application logging                 |
| [Environment Variables](building-blocks/environment.md) | Managing configuration              |

### 🗄️ Database

| Document                                  | Description                       |
| ----------------------------------------- | --------------------------------- |
| [Overview](database/overview.md)          | Database support in Goose         |
| [SQL Databases](database/sql.md)          | Working with relational databases |
| [Key-Value Store](database/kv.md)         | Redis-compatible KV operations    |
| [Migrations](database/migrations.md)      | Database schema management        |
| [Entities & Models](database/entities.md) | Defining data models              |

### ⚡ Advanced Features

| Document                            | Description                    |
| ----------------------------------- | ------------------------------ |
| [Caching](advanced/caching.md)      | Application caching strategies |
| [Queues & Jobs](advanced/queues.md) | Background job processing      |
| [Cron Jobs](advanced/cron.md)       | Scheduled task execution       |
| [Events](advanced/events.md)        | Event-driven architecture      |

### 🔒 Security

| Document                                     | Description                   |
| -------------------------------------------- | ----------------------------- |
| [Overview](security/overview.md)             | Security best practices       |
| [Authentication](security/authentication.md) | User authentication           |
| [Authorization](security/authorization.md)   | Access control                |
| [CORS](security/cors.md)                     | Cross-Origin Resource Sharing |
| [Rate Limiting](security/rate-limiting.md)   | Protecting your API           |

### 🖥️ Platforms

| Document                         | Description                 |
| -------------------------------- | --------------------------- |
| [API Platform](platforms/api.md) | Building RESTful APIs       |
| [Web Platform](platforms/web.md) | Building web applications   |
| [SPA Platform](platforms/spa.md) | Single-page apps + API as one service |
| [CLI Platform](platforms/cli.md) | Building command-line tools |

### 🛠️ Goose CLI

| Document                                        | Description                    |
| ----------------------------------------------- | ------------------------------ |
| [Overview](cli/overview.md)                     | CLI tool introduction          |
| [Installation](cli/installation.md)             | Installing the CLI             |
| [Creating Applications](cli/creating-apps.md)   | Scaffolding new projects       |
| [Generating Modules](cli/generating-modules.md) | Adding modules to projects     |
| [Command Reference](cli/reference.md)           | Complete command documentation |

### 🧪 Testing

| Document                                      | Description                     |
| --------------------------------------------- | ------------------------------- |
| [Overview](testing/overview.md)               | Testing strategies              |
| [Unit Testing](testing/unit.md)               | Testing components in isolation |
| [Integration Testing](testing/integration.md) | Testing component interactions  |
| [HTTP Testing](testing/http.md)               | Testing API endpoints           |
| [Mocking](testing/mocking.md)                 | Test doubles and mocks          |

### 🚀 Deployment

| Document                                         | Description               |
| ------------------------------------------------ | ------------------------- |
| [Overview](deployment/overview.md)               | Deployment strategies     |
| [Docker](deployment/docker.md)                   | Containerized deployments |
| [Production Checklist](deployment/production.md) | Going to production       |

### 📖 Reference

| Document                                            | Description                |
| --------------------------------------------------- | -------------------------- |
| [API Reference](reference/api.md)                   | Complete API documentation |
| [Configuration Options](reference/configuration.md) | All configuration options  |
| [Error Codes](reference/errors.md)                  | Error code reference       |

### 👥 Community

| Document                                  | Description                |
| ----------------------------------------- | -------------------------- |
| [Contributing](community/contributing.md) | How to contribute          |
| [Changelog](community/changelog.md)       | Version history            |
| [Roadmap](community/roadmap.md)           | Future plans               |
| [FAQ](community/faq.md)                   | Frequently asked questions |

### 📘 Tutorials

| Document                                        | Description               |
| ----------------------------------------------- | ------------------------- |
| [Tutorials Index](tutorials/index.md)           | All tutorials             |
| [Sample Applications](tutorials/sample-apps.md) | Complete working examples |
| [Building a REST API](tutorials/rest-api.md)    | Complete API tutorial     |
| [Building a Web App](tutorials/web-app.md)      | Complete web tutorial     |
| [Building a CLI Tool](tutorials/cli-tool.md)    | Complete CLI tutorial     |

---

## Quick Start

```bash
# Install Goose CLI
go install github.com/awesome-goose/cli@latest

# Create a new API application
goose app --name=myapi --template=api

# Navigate and run
cd myapi
go mod tidy
go run main.go
```

Your API is now running at `http://localhost:8080`! 🎉

> 💡 **Looking for complete examples?** Check out the [Sample Applications](tutorials/sample-apps.md) for working reference implementations of API, Web, CLI, and Multi-platform apps.

---

## Principles

Goose follows these core principles:

1. **Structs are first-class citizens** - Define your application structure declaratively
2. **Tags are powerful** - Use struct tags for validation, injection, and configuration
3. **Declarative over imperative** - Convention over configuration
4. **Go idioms first** - Stick to Go conventions and standard library
5. **Modular by design** - Everything is a module

---

## Version

This documentation is for Goose **v1.0**.

---

_Built with ❤️ by the Goose team_
