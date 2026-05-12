# Logging

Goose provides a flexible logging system with multiple formatters, processors, and output destinations.

## Basic Usage

Inject the logger and use it:

```go
type UserService struct {
    log types.Log `inject:""`
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    s.log.Info("Creating user", "email", dto.Email)

    user, err := s.repo.Create(dto)
    if err != nil {
        s.log.Error("Failed to create user", "error", err)
        return nil, err
    }

    s.log.Info("User created", "user_id", user.ID)
    return user, nil
}
```

## Log Levels

Goose supports standard log levels:

```go
// Debug - Detailed debugging information
s.log.Debug("Processing request", "request_id", id)

// Info - General information
s.log.Info("User logged in", "user_id", userID)

// Notice - Normal but significant events
s.log.Notice("Cache cleared", "cache_keys", count)

// Warning - Warning conditions
s.log.Warning("Rate limit approaching", "current", count, "limit", limit)

// Error - Error conditions
s.log.Error("Database connection failed", "error", err)

// Critical - Critical conditions
s.log.Critical("Service unavailable", "service", name)

// Alert - Action must be taken immediately
s.log.Alert("Security breach detected", "ip", ip)

// Emergency - System is unusable
s.log.Emergency("System crash", "error", err)
```

## Configuring the Logger

### Basic Setup

```go
import (
    "github.com/awesome-goose/goose/log"
    "github.com/awesome-goose/goose/log/formatters"
    "github.com/awesome-goose/goose/log/processors"
)

func setupLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("app"),
        log.NewLogger(
            []types.Modifier{},
            formatters.NewLine(),
            processors.NewConsole(),
        ),
    )
}
```

### With Formatters

```go
import "github.com/awesome-goose/goose/log/formatters"

// Line format (human readable)
formatters.NewLine()
// Output: 2024-01-15 10:30:45 [app] INFO: User created | extra: [...]

// JSON format (machine readable)
formatters.NewJSON()
// Output: {"datetime":"...","channel":"app","level":"info","message":"User created",...}

// Syslog format
formatters.NewSyslog()
```

### With Processors

```go
import "github.com/awesome-goose/goose/log/processors"

// Console output
processors.NewConsole()

// File output (directory path - files are bucketed by hour)
processors.NewFileProcessor("./storage/logs")
// Creates: ./storage/logs/2024/01/15/10.txt

// Syslog output (returns (proc, error) on Unix; no-op on Windows)
proc, err := processors.NewSyslog("myapp")
```

### With Modifiers

```go
import "github.com/awesome-goose/goose/log/modifiers"

log.NewLogger(
    []types.Modifier{
        modifiers.NewColorTagsModifier(), // Colorized output (Line formatter)
        modifiers.NewStackTrace(),        // Include stack traces on errors
        modifiers.NewSystemInfo(),        // Hostname / PID
        modifiers.NewUUID(),              // Per-record UUID
    },
    formatters.NewLine(),
    processors.NewConsole(),
)
```

## Multiple Channels

Configure different logging channels:

```go
import (
    "github.com/awesome-goose/goose/log"
    "github.com/awesome-goose/goose/log/formatters"
    "github.com/awesome-goose/goose/log/modifiers"
    "github.com/awesome-goose/goose/log/processors"
)

logger := log.NewLog(
    log.AppLogChannel("console"),
    log.NewLogger(
        []types.Modifier{modifiers.NewColorTagsModifier()},
        formatters.NewLine(),
        processors.NewConsole(),
    ),
)

// Add file channel
logger.Add("file",
    log.NewLogger(
        []types.Modifier{},
        formatters.NewJSON(),
        processors.NewFileProcessor("./storage/logs"),
    ),
)

// Use specific channel
logger.Use("file").Info("This goes to file")
logger.Use("console").Info("This goes to console")
```

## Structured Logging

Always use structured logging with key-value pairs:

```go
// ✅ Good: Structured logging
s.log.Info("Order processed",
    "order_id", order.ID,
    "customer_id", order.CustomerID,
    "total", order.Total,
    "items", len(order.Items),
)

// ❌ Bad: String formatting
s.log.Info(fmt.Sprintf("Order %s processed for customer %s, total: $%.2f",
    order.ID, order.CustomerID, order.Total))
```

## Request Logging

Log HTTP requests:

```go
type RequestLoggerMiddleware struct {
    log types.Log `inject:""`
}

func (m *RequestLoggerMiddleware) Handle(ctx types.Context) error {
    start := time.Now()

    // Log request start
    m.log.Info("Request started",
        "method", ctx.Method(),
        "path", ctx.Path(),
        "ip", ctx.ClientIP(),
    )

    // After handler (using deferred logging)
    defer func() {
        duration := time.Since(start)
        m.log.Info("Request completed",
            "method", ctx.Method(),
            "path", ctx.Path(),
            "duration_ms", duration.Milliseconds(),
        )
    }()

    return nil
}
```

## Error Logging

Log errors with context:

```go
func (s *PaymentService) ProcessPayment(orderID string, amount float64) error {
    err := s.gateway.Charge(amount)
    if err != nil {
        s.log.Error("Payment processing failed",
            "order_id", orderID,
            "amount", amount,
            "error", err,
            "gateway", s.gateway.Name(),
        )
        return err
    }

    s.log.Info("Payment processed successfully",
        "order_id", orderID,
        "amount", amount,
    )
    return nil
}
```

## Log Rotation

The built-in `processors.NewFileProcessor(dir)` writes records into a
date-bucketed directory tree (`/dir/YYYY/MM/DD/HH.txt`), which gives you
implicit time-based rotation. There is no `WithMaxSize`/`WithMaxBackups`/
`WithCompress` option on the framework processor — for size-based rotation
or compression, configure logrotate/journald at the OS level, or pipe to a
log shipper.

## Production Configuration

```go
import (
    "github.com/awesome-goose/goose/log"
    "github.com/awesome-goose/goose/log/formatters"
    "github.com/awesome-goose/goose/log/modifiers"
    "github.com/awesome-goose/goose/log/processors"
    "github.com/awesome-goose/goose/types"
)

func NewProductionLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("production"),
        log.NewLogger(
            []types.Modifier{
                modifiers.NewUUID(),
                modifiers.NewSystemInfo(),
            },
            formatters.NewJSON(),
            processors.NewFileProcessor("/var/log/app"),
        ),
    )
}
```

## Development Configuration

```go
func NewDevelopmentLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("development"),
        log.NewLogger(
            []types.Modifier{
                modifiers.NewColorTagsModifier(),
                modifiers.NewStackTrace(),
            },
            formatters.NewLine(),
            processors.NewConsole(),
        ),
    )
}
```

## Best Practices

### 1. Use Appropriate Log Levels

```go
// Debug: For development details
s.log.Debug("SQL query executed", "query", sql, "duration", dur)

// Info: For normal operations
s.log.Info("User registered", "user_id", id)

// Warning: For recoverable issues
s.log.Warning("Cache miss", "key", key)

// Error: For errors that need attention
s.log.Error("API call failed", "error", err)
```

### 2. Include Context

```go
// ✅ Good: Rich context
s.log.Error("Failed to send email",
    "user_id", userID,
    "email", email,
    "template", template,
    "error", err,
    "attempt", attemptNumber,
)

// ❌ Bad: No context
s.log.Error("Email failed")
```

### 3. Don't Log Sensitive Data

```go
// ❌ Bad: Logging sensitive data
s.log.Info("User login", "password", password)

// ✅ Good: Mask or omit
s.log.Info("User login", "user_id", userID)
```

### 4. Use Consistent Keys

```go
// ✅ Good: Consistent naming
s.log.Info("Operation", "user_id", id)
s.log.Info("Another op", "user_id", id)

// ❌ Bad: Inconsistent naming
s.log.Info("Operation", "userId", id)
s.log.Info("Another op", "user", id)
```

## Next Steps

- [Error Handling](error-handling.md) - Handle errors properly
- [Environment Variables](environment.md) - Configure logging
- [Middleware](middleware.md) - Request logging
