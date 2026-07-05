# Sample Applications

Explore complete, working examples of Goose applications to accelerate your learning and development.

## Sandbox Repository

The **Goose Sandbox** contains reference implementations for all supported platform types. These examples demonstrate real-world patterns and best practices.

📦 **Repository:** [github.com/awesome-goose/sandbox](https://github.com/awesome-goose/sandbox)

---

## Available Examples

| Application                                                        | Description                     | Platform        |
| ------------------------------------------------------------------ | ------------------------------- | --------------- |
| [api/](https://github.com/awesome-goose/sandbox/tree/main/api)     | REST API with JSON responses    | API             |
| [web/](https://github.com/awesome-goose/sandbox/tree/main/web)     | Server-rendered web application | Web             |
| [cli/](https://github.com/awesome-goose/sandbox/tree/main/cli)     | Command-line interface tool     | CLI             |
| [multi/](https://github.com/awesome-goose/sandbox/tree/main/multi) | Multi-platform application      | API + Web + CLI |
| [spa/](https://github.com/awesome-goose/sandbox/tree/main/spa)     | Single-page app (Angular)       | SPA             |

---

## API Example

A REST API application demonstrating:

- JSON response handling
- RESTful routing patterns
- Resource modules with CRUD operations
- Dependency injection
- Structured logging

```bash
cd sandbox/api
go mod tidy
go run main.go
# Server starts at http://localhost:8080
```

**Endpoints:**

- `GET /` - Health check
- `GET /user` - List all users
- `GET /user/:id` - Get user by ID

---

## Web Example

A server-rendered web application demonstrating:

- HTML template rendering
- Form handling
- Static file serving
- Session management patterns

```bash
cd sandbox/web
go mod tidy
go run main.go
# Server starts at http://localhost:8080
```

---

## CLI Example

A command-line application demonstrating:

- Command registration
- Argument and flag handling
- Output formatting
- Interactive prompts

```bash
cd sandbox/cli
go mod tidy
go run main.go
```

---

## Multi-Platform Example

A single codebase serving multiple platforms:

- Shared business logic across platforms
- API endpoints (port 8080)
- Web interface (port 3000)
- CLI commands

```bash
cd sandbox/multi
go mod tidy
go run main.go
```

This example shows how to:

- Structure code for multi-platform support
- Share services and modules between platforms
- Configure platform-specific settings

---

## SPA Example

A single-page application demonstrating:

- One service serving both `/api` JSON routes and an Angular frontend
- `index.html` fallback for client-side routing
- An in-memory resource module (no database required)
- A Makefile-driven install / dev / build / dist workflow

```bash
cd sandbox/spa
make install   # go mod tidy + npm install
make run       # ng build -> public/, then server at http://localhost:8080
```

During development, run `make dev` to get the Go backend and the Angular
dev server (with `/api` proxying and hot reload) together.

---

## Project Structure

Each sandbox application follows the standard Goose module structure:

```
app/
├── main.go              # Application entry point
├── config.yaml          # Configuration file
├── go.mod               # Go module definition
├── app/
│   ├── app.module.go    # Root module
│   ├── app.controller.go
│   ├── app.service.go
│   ├── app.routes.go
│   └── app.dtos.go
└── [feature]/           # Feature modules
    ├── [feature].module.go
    ├── [feature].controller.go
    ├── [feature].service.go
    ├── [feature].routes.go
    └── [feature].dtos.go
```

---

## Running Examples Locally

1. Clone the sandbox repository:

   ```bash
   git clone https://github.com/awesome-goose/sandbox.git
   cd sandbox
   ```

2. Choose an example:

   ```bash
   cd api  # or web, cli, multi, spa
   ```

3. Install dependencies and run:
   ```bash
   go mod tidy
   go run main.go
   ```

---

## Using as Templates

The sandbox examples can serve as starting points for your own applications:

1. Copy the example that matches your use case
2. Update `go.mod` with your module name
3. Modify the application code as needed

Or use the Goose CLI to scaffold a new project:

```bash
# Create from template
goose app --name=myproject --template=api
goose app --name=myproject --template=web
goose app --name=myproject --template=cli
goose app --name=myproject --template=multi
goose app --name=myproject --template=spa --framework=react
```

---

## Next Steps

- [Quick Start](../getting-started/quick-start.md) - Build your first app
- [Architecture](../core-concepts/architecture.md) - Understand Goose design
- [Modules](../core-concepts/modules.md) - Learn module patterns
- [Building a REST API](rest-api.md) - Step-by-step tutorial
