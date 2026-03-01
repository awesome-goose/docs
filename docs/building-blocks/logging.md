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
import "github.com/awesome-goose/goose/log"

func setupLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("app"),
        log.NewLogger(
            []types.Modifier{},
            log.NewLineFormatter(),
            log.NewConsoleProcessor(),
        ),
    )
}
```

### With Formatters

```go
// Line format (human readable)
log.NewLineFormatter()
// Output: [2024-01-15 10:30:45] app.INFO: User created {"user_id": "123"}

// JSON format (machine readable)
log.NewJsonFormatter()
// Output: {"datetime": "2024-01-15T10:30:45Z", "level": "info", "message": "User created", "user_id": "123"}

// Syslog format
log.NewSyslogFormatter()
```

### With Processors

```go
// Console output
log.NewConsoleProcessor()

// File output
log.NewFileProcessor("/var/log/app.log")

// Syslog output
log.NewSyslogProcessor("app", syslog.LOG_LOCAL0)
```

### With Modifiers

```go
log.NewLogger(
    []types.Modifier{
        log.NewColorModifier(),      // Colorized output
        log.NewStackTraceModifier(), // Include stack traces
        log.NewSystemInfoModifier(), // Include system info
        log.NewUuidModifier(),       // Add UUID to each log
    },
    log.NewLineFormatter(),
    log.NewConsoleProcessor(),
)
```

## Multiple Channels

Configure different logging channels:

```go
logger := log.NewLog(
    log.AppLogChannel("console"),
    log.NewLogger(
        []types.Modifier{log.NewColorModifier()},
        log.NewLineFormatter(),
        log.NewConsoleProcessor(),
    ),
)

// Add file channel
logger.Add("file",
    log.NewLogger(
        []types.Modifier{},
        log.NewJsonFormatter(),
        log.NewFileProcessor("/var/log/app.log"),
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

For production, configure log rotation:

```go
// Using file processor with rotation
log.NewFileProcessor("/var/log/app.log",
    log.WithMaxSize(100),      // 100 MB max
    log.WithMaxBackups(10),    // Keep 10 backups
    log.WithMaxAge(30),        // 30 days retention
    log.WithCompress(true),    // Compress old logs
)
```

## Production Configuration

```go
func NewProductionLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("production"),
        // File logging with JSON format
        log.NewLogger(
            []types.Modifier{
                log.NewUuidModifier(),
                log.NewSystemInfoModifier(),
            },
            log.NewJsonFormatter(),
            log.NewFileProcessor("/var/log/app.log"),
        ),
    )
}
```

## Development Configuration

```go
func NewDevelopmentLogger() types.Log {
    return log.NewLog(
        log.AppLogChannel("development"),
        // Console logging with colors
        log.NewLogger(
            []types.Modifier{
                log.NewColorModifier(),
                log.NewStackTraceModifier(),
            },
            log.NewLineFormatter(),
            log.NewConsoleProcessor(),
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
