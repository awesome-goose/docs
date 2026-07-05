# SPA Platform

Serve a single-page application and its JSON API as one service with the Goose SPA platform.

## Overview

The SPA platform runs one HTTP server that does two jobs:

- Requests under the API prefix (default `/api`) are routed to your Goose modules and answered as JSON — exactly like the API platform.
- Every other request is served from a static directory (default `public/`) containing your built frontend, with an `index.html` fallback for client-side routes.

This is the natural fit for React, Vue, Svelte, or Angular apps that ship with their own Go backend: one binary, one port, no CORS.

## Quick Start

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/spa"
)

func main() {
    platform := spa.NewPlatform(
        spa.WithPort(8080),
        spa.WithStaticDir("public"),
        spa.WithAPIPrefix("/api"),
    )

    module := &app.AppModule{}

    stop, err := goose.Start(goose.SPA(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Configuration

### Platform Options

```go
platform := spa.NewPlatform(
    spa.WithHost("0.0.0.0"),            // Listen address
    spa.WithPort(8080),                  // Port
    spa.WithTimeout(30),                 // Request timeout (seconds)
    spa.WithName("My SPA"),              // App name
    spa.WithVersion("0.0.0"),            // App version
    spa.WithStaticDir("public"),         // Directory with the built frontend
    spa.WithIndexFile("index.html"),     // Fallback file for client-side routes
    spa.WithAPIPrefix("/api"),           // Prefix routed to your modules
)
```

### Available Options

```go
type Config struct {
    Name        string  // App name
    Version     string  // App version
    Author      string  // Author
    Description string  // Description
    Host        string  // Listen address
    Port        int     // Port
    Timeout     int     // Request timeout

    StaticDir   string  // Static assets directory (default "public")
    IndexFile   string  // SPA entry file within StaticDir (default "index.html")
    APIPrefix   string  // URL prefix routed to the kernel (default "/api")
}
```

A relative `StaticDir` resolves against the working directory of the running
binary, so a deployed binary sitting next to its `public/` folder works
without configuration.

## Request Lifecycle

For each request the platform decides in order:

1. **API route** — the path equals the API prefix or starts with it: the prefix is stripped and the request goes through the normal Goose router → middleware → controller pipeline. Errors surface as HTTP 500.
2. **Method guard** — non-`GET`/`HEAD` requests outside the API prefix get `405 Method Not Allowed`.
3. **Static file** — if the cleaned path exists inside `StaticDir`, it is served with correct MIME types, range support, and conditional-GET handling.
4. **Index fallback** — if the path has no file extension and the client accepts HTML, `IndexFile` is served with `Cache-Control: no-cache` so client-side routes like `/dashboard/settings` load your SPA.
5. **404** — everything else, including missing assets like `/logo.png` (which never fall back to `index.html`).

Path traversal is blocked: request paths are cleaned and rooted before touching the filesystem.

## Routes Are Prefix-Free

Declare routes exactly as you would in an API app — the platform adds the prefix:

```go
var ROUTES = router.ForRoutes(
    router.Get("/", []any{AppController{}, "Health"}),        // GET /api
    router.Get("/users", []any{UserController{}, "List"}),    // GET /api/users
    router.Get("/users/:id", []any{UserController{}, "Get"}), // GET /api/users/:id
)
```

Controllers return JSON with `output.JSON(...)`, identical to the [API platform](api.md).

## Development Workflow

Run two processes during development:

```bash
make dev-backend    # go run main.go — Goose on :8080
make dev-frontend   # framework dev server with hot reload
```

Point the frontend dev server's proxy at the backend so `/api` calls hit real endpoints (the CLI's spa template pre-configures this):

```js
// vite.config.js
server: {
  proxy: { '/api': 'http://localhost:8080' },
}
```

```json
// Angular proxy.conf.json
{ "/api": { "target": "http://localhost:8080", "secure": false } }
```

## Production

Build the frontend into `StaticDir`, compile the binary, and ship both:

```bash
make dist
# dist/
# ├── myapp        # Go binary
# ├── .env
# └── public/      # Built frontend (index.html + hashed assets)

cd dist && ./myapp
```

The Goose service serves the whole app — no separate web server needed.

## Scaffolding

The Goose CLI generates a complete SPA project — Go backend, your choice of frontend, and a Makefile wiring them together:

```bash
goose app --name=myspa --template=spa --framework=react   # or vue | svelte | ng
```

See [Creating Apps](../cli/creating-apps.md) for the generated structure.

## Multi-Platform

A SPA instance composes with other platforms like any server instance:

```go
stop, err := goose.Start(
    goose.SPA(spaPlatform, spaModule, initializers),
    goose.CLI(cliPlatform, cliModule, initializers),
)
```

See [Multi-Platform Applications](../core-concepts/multi-platform.md).
