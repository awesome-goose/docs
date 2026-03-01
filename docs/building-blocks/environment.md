# Environment Variables

Goose provides a clean way to manage environment configuration through the `Env` package.

## Basic Usage

### Loading Environment Variables

Environment variables are automatically loaded from:

1. OS environment variables
2. `.env` files

```go
import "github.com/awesome-goose/goose/env"

// Create new Env instance (auto-loads .env)
e := env.NewEnv()

// Get values
appName := e.Get("APP_NAME")
port := e.GetInt("PORT")
debug := e.GetBool("DEBUG")
```

## .env File

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
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=myapp
DB_PASSWORD=secret

# Cache
CACHE_DRIVER=memory
CACHE_TTL=3600

# Logging
LOG_CHANNEL=console
LOG_LEVEL=debug

# Security
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=your-encryption-key
```

## Getting Values

### String Values

```go
// Get value (empty string if not found)
appName := e.Get("APP_NAME")

// Get with default
appEnv := e.GetWithDefault("APP_ENV", "development")
```

### Integer Values

```go
// Get as integer
port := e.GetInt("PORT")  // Returns 0 if not found or invalid

// Usage
server.Listen(fmt.Sprintf(":%d", port))
```

### Boolean Values

```go
// Get as boolean
debug := e.GetBool("DEBUG")  // Returns false if not found

// Valid truthy values: "true", "1", "yes", "on"
// Valid falsy values: "false", "0", "no", "off"
```

### Float Values

```go
// Get as float64
rate := e.GetFloat("RATE_LIMIT")
```

## Setting Values

Override values at runtime:

```go
e.Set("APP_ENV", "testing")
```

## Environment-Specific Files

Create environment-specific `.env` files:

```
.env                  # Default values
.env.development      # Development overrides
.env.production       # Production overrides
.env.test             # Test overrides
.env.local            # Local overrides (gitignored)
```

### Loading Order

1. `.env` (base configuration)
2. `.env.{APP_ENV}` (environment-specific)
3. `.env.local` (local overrides)
4. OS environment variables (highest priority)

## Environment Sources

### Custom Sources

```go
import (
    "github.com/awesome-goose/goose/env"
    "github.com/awesome-goose/goose/env/sources"
)

e := &env.Env{}
e.FromSources(
    sources.NewOsEnvSource(),    // OS environment
    sources.NewFileEnvSource(),   // .env files
)
```

### Creating Custom Sources

```go
type VaultEnvSource struct {
    client *vault.Client
}

func (s *VaultEnvSource) Load(e types.EnvStore) {
    secrets := s.client.GetSecrets("app/config")
    for key, value := range secrets {
        e.Set(key, value)
    }
}

// Use it
e.FromSources(
    sources.NewOsEnvSource(),
    &VaultEnvSource{client: vaultClient},
)
```

## Using in Applications

### In Initializers

```go
func main() {
    initializers := []func(types.Container) error{
        func(c types.Container) error {
            e := env.NewEnv()
            c.Register(func() *env.Env {
                return e
            }, "", true)
            return nil
        },
    }

    // ...
}
```

### In Services

```go
type ConfigService struct {
    env *env.Env `inject:""`
}

func (s *ConfigService) GetDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        Driver:   s.env.Get("DB_DRIVER"),
        Host:     s.env.GetWithDefault("DB_HOST", "localhost"),
        Port:     s.env.GetInt("DB_PORT"),
        Database: s.env.Get("DB_DATABASE"),
        Username: s.env.Get("DB_USERNAME"),
        Password: s.env.Get("DB_PASSWORD"),
    }
}
```

## Validating Environment

Validate required variables on startup:

```go
func ValidateEnvironment(e *env.Env) error {
    required := []string{
        "APP_NAME",
        "DB_DATABASE",
        "JWT_SECRET",
    }

    var missing []string
    for _, key := range required {
        if e.Get(key) == "" {
            missing = append(missing, key)
        }
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required env vars: %v", missing)
    }

    return nil
}

// Use in initializer
func(c types.Container) error {
    e := env.NewEnv()
    if err := ValidateEnvironment(e); err != nil {
        return err
    }
    // ...
}
```

## Best Practices

### 1. Never Commit Secrets

```bash
# .gitignore
.env
.env.local
.env.*.local
```

### 2. Use .env.example

Create a template for required variables:

```env
# .env.example (committed to repo)
APP_NAME=
APP_ENV=development
DB_DATABASE=
JWT_SECRET=
```

### 3. Validate Early

```go
func main() {
    e := env.NewEnv()

    // Fail fast if missing required vars
    if e.Get("DB_DATABASE") == "" {
        panic("DB_DATABASE is required")
    }

    // ...
}
```

### 4. Use Typed Getters

```go
// ✅ Good: Use typed getters
port := e.GetInt("PORT")
debug := e.GetBool("DEBUG")

// ❌ Bad: Manual conversion
portStr := e.Get("PORT")
port, _ := strconv.Atoi(portStr)
```

### 5. Provide Sensible Defaults

```go
// ✅ Good: Defaults for non-critical settings
logLevel := e.GetWithDefault("LOG_LEVEL", "info")
timeout := e.GetWithDefault("TIMEOUT", "30")

// ❌ Bad: No defaults for configuration
logLevel := e.Get("LOG_LEVEL")  // Empty if not set
```

### 6. Group Related Variables

```env
# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=myapp

# Cache
CACHE_DRIVER=redis
CACHE_HOST=localhost
CACHE_PORT=6379

# Mail
MAIL_DRIVER=smtp
MAIL_HOST=smtp.example.com
MAIL_PORT=587
```

## Security Considerations

1. **Never log environment variables**
2. **Use separate secrets for each environment**
3. **Rotate secrets regularly**
4. **Use secret management services in production**

```go
// ❌ Never do this
log.Println("JWT_SECRET:", e.Get("JWT_SECRET"))

// ✅ Log that config was loaded, not values
log.Println("Configuration loaded successfully")
```

## Next Steps

- [Configuration](../getting-started/configuration.md) - YAML configuration
- [Logging](logging.md) - Configure logging
- [Security](../security/overview.md) - Security best practices
