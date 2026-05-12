# Configuration

Goose applications can be configured using environment variables and YAML configuration files.

## Environment Variables

### .env File

Create a `.env` file in your project root:

```env
# Application
APP_NAME=myapp
APP_ENV=development
APP_DEBUG=true

# Server
HOST=localhost
PORT=8080

# Database
DB_DRIVER=sqlite
DB_DATABASE=./data/app.db

# Logging
LOG_CHANNEL=console
LOG_LEVEL=debug
```

### Loading Environment Variables

Environment variables are automatically loaded when your application starts. The Goose `Env` package handles this:

```go
import "github.com/awesome-goose/goose/env"

// Create a new Env instance
e := env.NewEnv()

// Get values
appName := e.Get("APP_NAME")
port := e.GetInt("PORT")
debug := e.GetBool("APP_DEBUG")

// Get with default
env := e.GetWithDefault("APP_ENV", "development")
```

### Environment-Specific Files

You can create environment-specific `.env` files:

- `.env` - Default configuration
- `.env.development` - Development overrides
- `.env.production` - Production overrides
- `.env.test` - Test overrides

## YAML Configuration

For complex configuration, use YAML files in a `config/` directory.

### Directory Structure

```
config/
├── app.yaml
├── database.yaml
├── cache.yaml
└── logging.yaml
```

### Example: config/app.yaml

```yaml
name: myapp
version: 0.0.0
timezone: UTC

server:
  host: localhost
  port: 8080
  timeout: 30s
```

### Example: config/database.yaml

```yaml
default: sqlite

connections:
  sqlite:
    driver: sqlite
    database: ./data/app.db

  postgres:
    driver: postgres
    host: localhost
    port: 5432
    database: myapp
    username: user
    password: secret
```

### Accessing Configuration

```go
import "github.com/awesome-goose/goose/config"

// In your service or controller
type MyService struct {
    config *config.Config `inject:""`
}

func (s *MyService) GetAppName() string {
    return s.config.Get("app.name").(string)
}

func (s *MyService) GetDBConfig() map[string]any {
    return s.config.Get("database.connections.sqlite").(map[string]any)
}
```

## Platform Configuration

Each platform has its own configuration options.

### API Platform

```go
import "github.com/awesome-goose/goose/platforms/api"

platform := api.NewPlatform(
    api.WithName("myapi"),
    api.WithHost("0.0.0.0"),
    api.WithPort(8080),
)
```

### Web Platform

```go
import "github.com/awesome-goose/goose/platforms/web"

platform := web.NewPlatform(
    web.WithName("myweb"),
    web.WithHost("0.0.0.0"),
    web.WithPort(3000),
    web.WithTemplatesDir("./templates"),
)
```

### CLI Platform

```go
import "github.com/awesome-goose/goose/platforms/cli"

platform := cli.NewPlatform(
    cli.WithName("mycli"),
    cli.WithVersion("0.0.0"),
)
```

## Module Configuration

Modules can have their own configuration.

### Cache Module

```go
import "github.com/awesome-goose/goose/modules/cache"

cacheModule := cache.NewModule(
    cache.WithGroup("app"),
    cache.WithDefaultTTL(5 * time.Minute),
)
```

### SQL Module

```go
import "github.com/awesome-goose/goose/modules/sql"

sqlModule := sql.Root(&sql.Config{
    Dialect: "sqlite",
    Name:    "./data/app.db",
    Sync:    true,
})
```

### Queue Module

```go
import "github.com/awesome-goose/goose/modules/queues"

cfg := queues.DefaultConfig()
queues.WithQueue("default")(cfg)
queues.WithPollInterval(500 * time.Millisecond)(cfg)

queueModule := queues.Root(cfg, queueHandlers...)
```

Per-handler worker count is configured via
`handler.WithMinWorkers(n)` / `WithMaxWorkers(n)` — there is no
`WithWorkerCount` on the module.

## Configuration Best Practices

### 1. Use Environment Variables for Secrets

Never commit secrets to your repository:

```env
# .env (gitignored)
DB_PASSWORD=my-secret-password
API_KEY=sk-1234567890
```

### 2. Environment-Based Configuration

```go
func newPlatform() types.Platform {
    env := os.Getenv("APP_ENV")

    if env == "production" {
        return api.NewPlatform(
            api.WithHost("0.0.0.0"),
            api.WithPort(80),
        )
    }

    return api.NewPlatform(
        api.WithHost("localhost"),
        api.WithPort(8080),
    )
}
```

### 3. Validate Configuration Early

```go
func main() {
    requiredEnvVars := []string{"DB_DATABASE", "APP_NAME"}

    e := env.NewEnv()
    for _, v := range requiredEnvVars {
        if e.Get(v) == "" {
            panic("Missing required env var: " + v)
        }
    }

    // Continue with application startup
}
```

### 4. Use Defaults

```go
port := e.GetWithDefault("PORT", "8080")
timeout := e.GetWithDefault("TIMEOUT", "30")
```

## Next Steps

- [Directory Structure](directory-structure.md) - Project layout
- [Environment Variables](../building-blocks/environment.md) - Detailed env guide
- [Logging](../building-blocks/logging.md) - Configure logging
