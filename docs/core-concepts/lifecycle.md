# Application Lifecycle

Understanding the Goose application lifecycle helps you hook into key stages of your application's execution.

## Lifecycle Overview

```
┌────────────────────────────────────────────────────────────────┐
│                    Application Start                            │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  1. Platform Creation                                           │
│     - Create platform (API/Web/CLI)                            │
│     - Load platform configuration                               │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  2. Container Initialization                                    │
│     - Create IoC container                                      │
│     - Register core services                                    │
│     - Run initializers                                          │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  3. Module Traversal                                            │
│     - Resolve module imports                                    │
│     - Register declarations                                     │
│     - Publish exports                                           │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  4. Boot Hooks                                                  │
│     - Execute OnBoot hooks                                      │
│     - Initialize routes                                         │
│     - Start background services                                 │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  5. Platform Boot                                               │
│     - Start HTTP server (API/Web)                              │
│     - Execute commands (CLI)                                    │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  6. Application Running                                         │
│     - Handle requests                                           │
│     - Process jobs                                              │
│     - Run cron tasks                                            │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  7. Shutdown                                                    │
│     - Stop accepting requests                                   │
│     - Complete pending operations                               │
│     - Execute cleanup hooks                                     │
└────────────────────────────────────────────────────────────────┘
```

## Lifecycle Hooks

### Boot Hook

Execute code after the application is fully initialized but before it starts running.

Implement the `types.Bootable` interface:

```go
type MyService struct {
    db    *sql.Db   `inject:""`
    cache *cache.Cache `inject:""`
}

// Boot is called during the boot phase
func (s *MyService) Boot(k types.Kernel) error {
    // Warm up cache
    // ...

    // Verify database connection
    // ...

    return nil
}
```

### Shutdown Hook

Execute cleanup code when the application shuts down.

Implement the `types.Shutdownable` interface:

```go
func (s *MyService) Shutdown(k types.Kernel) error {
    // Close connections
    // Flush buffers
    // Save state
    return nil
}
```

### Hook Interfaces

```go
// types/hooks.go

// Configurable is called during module registration, before declarations are created.
type Configurable interface {
    Configure(c Container) error
}

// Bootable is implemented by services that need initialization after setup
type Bootable interface {
    Boot(k Kernel) error
}

// Shutdownable is implemented by services that need cleanup on shutdown
type Shutdownable interface {
    Shutdown(k Kernel) error
}
```

## Using Initializers

Initializers run before module traversal:

```go
func main() {
    initializers := []func(types.Container) error{
        loadConfiguration,
        setupLogging,
        validateEnvironment,
    }

    stop, err := goose.Start(goose.API(platform, module, initializers))
    if err != nil {
        panic(err)
    }
    defer stop()
}

func loadConfiguration(container types.Container) error {
    cfg, err := config.NewConfig("./config")
    if err != nil {
        return err
    }
    return container.Register(func() *config.Config {
        return cfg
    }, "", true)
}

func setupLogging(container types.Container) error {
    logger := log.NewLog(
        log.AppLogChannel("app"),
        log.NewLogger(
            []types.Modifier{},
            formatters.NewLine(),
            processors.NewConsole(),
        ),
    )
    container.Register(func() types.Log {
        return logger
    }, "", true)
    return nil
}

func validateEnvironment(container types.Container) error {
    required := []string{"APP_NAME", "DB_DATABASE"}
    env := env.NewEnv()

    for _, key := range required {
        if env.Get(key) == "" {
            return fmt.Errorf("missing required env: %s", key)
        }
    }
    return nil
}
```

## Request Lifecycle

For each HTTP request:

```
┌────────────────────────────────────────────────────────────────┐
│                     Incoming Request                            │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  1. Request Parsing                                             │
│     - Parse HTTP method, path, headers                         │
│     - Create request context                                    │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  2. Route Matching                                              │
│     - Find matching route                                       │
│     - Extract path parameters                                   │
│     - Collect middleware stack                                  │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  3. Middleware Execution                                        │
│     - Execute middleware chain                                  │
│     - Authentication, logging, etc.                            │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  4. Handler Execution                                           │
│     - Resolve controller                                        │
│     - Call handler method                                       │
│     - Execute business logic                                    │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│  5. Response Serialization                                      │
│     - Serialize return value                                    │
│     - Set headers                                               │
│     - Send response                                             │
└────────────────────────────────────────────────────────────────┘
```

## Graceful Shutdown

Goose handles graceful shutdown automatically:

```go
func main() {
    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }

    // The stop function is called when:
    // - SIGINT received (Ctrl+C)
    // - SIGTERM received (kill)
    // - Application error

    defer func() {
        if err := stop(); err != nil {
            log.Println("Shutdown error:", err)
        }
    }()
}
```

### Custom Shutdown Handling

```go
func main() {
    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }

    // Wait for shutdown signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    fmt.Println("Shutting down...")

    // Custom cleanup
    cleanupResources()

    // Stop the application
    if err := stop(); err != nil {
        log.Fatal(err)
    }
}
```

## Multi-Instance Lifecycle

For multi-platform applications:

```go
func main() {
    stop, err := goose.Start(
        goose.API(apiPlatform, apiModule, nil),
        goose.Web(webPlatform, webModule, nil),
        goose.CLI(cliPlatform, cliModule, nil),
    )
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

- API and Web platforms run concurrently
- CLI platform runs when `cli` argument is passed
- All share the same dependency container

## Best Practices

### 1. Fail Fast

Validate critical requirements during boot:

```go
func (s *DatabaseService) OnBoot(kernel types.Kernel) error {
    // Fail immediately if database is unreachable
    if err := s.db.Ping(); err != nil {
        return fmt.Errorf("database not available: %w", err)
    }
    return nil
}
```

### 2. Clean Shutdown

Always clean up resources:

```go
func (s *QueueService) OnShutdown() error {
    // Wait for in-progress jobs
    s.queue.WaitForCompletion(30 * time.Second)

    // Close connections
    return s.connection.Close()
}
```

### 3. Health Checks

Implement health checks for production:

```go
func (c *HealthController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/health", Handler: c.Check},
        {Method: "GET", Path: "/ready", Handler: c.Ready},
    }
}

func (c *HealthController) Check(ctx types.Context) any {
    return map[string]string{"status": "ok"}
}
```

## Next Steps

- [Multi-Platform Support](multi-platform.md) - Running multiple platforms
- [Middleware](../building-blocks/middleware.md) - Request interceptors
- [Error Handling](../building-blocks/error-handling.md) - Handling errors
