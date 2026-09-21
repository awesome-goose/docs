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

The options below are all off by default; an app that sets none of them behaves exactly as before.

```go
platform := spa.NewPlatform(
    spa.WithHost("0.0.0.0"),                                  // required in a container: the default is localhost
    spa.WithHandler("/ws", wsHandler),                        // mount an http.Handler: WebSocket, SSE, webhooks
    spa.WithMiddleware(requestID, accessLog),                 // wrap the whole pipeline
    spa.WithBodyLimit(25 << 20),                              // 25 MiB request bodies
    spa.WithCORS(spa.CORS{AllowedOrigins: []string{"https://app.example"}}),
    spa.WithSecurityHeaders(headers),                         // see below
    spa.WithCompression(),                                    // gzip + precompressed .br/.gz files
    spa.WithRecovery(func(r *http.Request, v any, stack []byte) { /* log it */ }),
    spa.WithWriteTimeout(30 * time.Second),                   // 0 disables the server-wide write timeout
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

    WriteTimeout time.Duration // Server write timeout when set with WithWriteTimeout (0 = none)
    BodyLimit    int64         // Request body cap in bytes (0 = none)
    CORS         *CORS         // nil = no CORS headers
    SecurityHeaders *SecurityHeaders // nil = none
    Recover      bool          // Panic recovery (WithRecovery)
    Compression  bool          // gzip and precompressed siblings (WithCompression)
    Middlewares  []Middleware  // Wrap the whole pipeline
    Handlers     []Mount       // Mounted http.Handlers
}
```

A relative `StaticDir` resolves against the working directory of the running
binary, so a deployed binary sitting next to its `public/` folder works
without configuration.

## Request Lifecycle

Options run as a pipeline around routing. From the outside in: **recovery → security headers → CORS → compression → body limit → your middleware (first is outermost) → routing.** CORS preflights are answered before your middleware runs, so an auth middleware never sees them.

For each request the platform then decides in order:

0. **Mounted handler** — a path registered with `WithHandler` (exact, or a subtree when it ends in `/`; the longest match wins) is served by its `http.Handler` with the original path and any method. It beats the API prefix and the static files.
1. **API route** — the path equals the API prefix or starts with it: the prefix is stripped and the request goes through the normal Goose router → middleware → controller pipeline. Errors surface as HTTP 500.
2. **Method guard** — non-`GET`/`HEAD` requests outside the API prefix get `405 Method Not Allowed`.
3. **Static file** — if the cleaned path exists inside `StaticDir`, it is served with correct MIME types, range support, and conditional-GET handling.
4. **Index fallback** — if the path has no file extension and the client accepts HTML, `IndexFile` is served with `Cache-Control: no-cache` so client-side routes like `/dashboard/settings` load your SPA.
5. **404** — everything else, including missing assets like `/logo.png` (which never fall back to `index.html`).

Path traversal is blocked: request paths are cleaned and rooted before touching the filesystem.

## WebSockets and Server-Sent Events

Mount a plain `http.Handler`. It receives the real request and a writer that still supports `Flush`, `Hijack` and per-request deadlines, whatever other options are on:

```go
spa.NewPlatform(
    spa.WithCompression(),
    spa.WithHandler("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil) // gorilla/websocket asserts w.(http.Hijacker)
        // ...
    })),
    spa.WithHandler("/events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/event-stream")
        rc := http.NewResponseController(w)
        rc.SetWriteDeadline(time.Time{}) // this stream may outlive the server's write timeout
        for msg := range messages(r.Context()) {
            fmt.Fprintf(w, "data: %s\n\n", msg)
            rc.Flush()
        }
    })),
)
```

Every writer the platform puts in front of your handler forwards `Flush`, `Hijack` and `Unwrap`. Write your own middleware the same way (or embed the writer and add `Unwrap() http.ResponseWriter`), or `http.NewResponseController` cannot reach the connection.

## Timeouts and Long Responses

The server write timeout defaults to 30 seconds and applies to the whole response, so any response that takes longer (a stream, a large download) is cut. Either raise or remove it for every route with `WithWriteTimeout` (`0` disables it), or leave the default and let only the handlers that need it clear their own deadline with `http.NewResponseController(w).SetWriteDeadline(time.Time{})`. `WithWriteTimeout` takes precedence over the write half of `WithTimeout`, in either option order.

## CORS

Because the SPA and the API share an origin, most apps need no CORS. When a separate frontend origin calls the API, list it:

```go
spa.WithCORS(spa.CORS{
    AllowedOrigins:   []string{"https://app.example"},
    AllowCredentials: true,
    MaxAge:           10 * time.Minute,
})
```

Only listed origins get `Access-Control-Allow-Origin`. A preflight from an unlisted origin, or asking for a method or header that is not allowed, is answered `403`. `"*"` allows any origin but cannot be combined with `AllowCredentials`; `Boot` returns an error if you try. Methods default to GET, HEAD, POST, PUT, PATCH and DELETE, and headers to `Content-Type` and `Authorization`.

## Security Headers, Body Limit and Recovery

`WithSecurityHeaders(spa.DefaultSecurityHeaders())` sends `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin` and `X-Frame-Options: SAMEORIGIN`. It sends no `Content-Security-Policy`; a policy depends on your assets, so set `ContentSecurityPolicy` yourself. A handler can still override any of them.

`WithBodyLimit(n)` answers `413` when the declared size is over `n` and otherwise cuts the body at `n` bytes (reading it returns `*http.MaxBytesError`).

`WithRecovery` turns a panic into `500 {"message":"internal server error"}`. The panic value never reaches the client; your callback gets it and the stack. If the response has already started, the connection is aborted instead, because no status can be changed. Without this option a panic propagates to `net/http` as before.

## Compression

`WithCompression` gzips responses whose type is worth compressing (text, JSON, JavaScript, SVG, WebAssembly and similar) for clients that accept it, and adds `Vary: Accept-Encoding`. It never compresses an event stream, a WebSocket upgrade, a range request or a response that already has a `Content-Encoding`, and `Flush` on a compressed stream flushes the gzip block.

Static files also honour precompressed siblings: if `app.js.br` (or `app.js.gz`) exists next to `app.js` and the client accepts that encoding, it is sent with the original file's `Content-Type`. Brotli is preferred. Ship them from your build (`brotli -k dist/*.js`) to avoid compressing on every request.

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
