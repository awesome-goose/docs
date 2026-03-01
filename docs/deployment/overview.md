# Deployment Overview

Deploy your Goose application to production environments.

## Deployment Options

| Platform                        | Best For               | Setup Complexity |
| ------------------------------- | ---------------------- | ---------------- |
| [Docker](docker.md)             | Containers, Kubernetes | Medium           |
| [Traditional Server](server.md) | VPS, bare metal        | Simple           |
| Cloud Platforms                 | AWS, GCP, Azure        | Varies           |

## Building for Production

### Build Binary

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o myapp main.go

# Build for different platforms
GOOS=darwin GOARCH=amd64 go build -o myapp-mac main.go
GOOS=windows GOARCH=amd64 go build -o myapp.exe main.go
```

### Build Flags

```bash
# Strip debug info (smaller binary)
go build -ldflags="-s -w" -o myapp main.go

# Embed version info
go build -ldflags="-X main.Version=1.2.3" -o myapp main.go
```

## Environment Configuration

### Production Environment Variables

```env
# Application
APP_NAME=myapp
APP_ENV=production
APP_DEBUG=false
HOST=0.0.0.0
PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=db.example.com
DB_PORT=5432
DB_NAME=myapp_prod
DB_USER=myapp
DB_PASSWORD=secret

# Redis
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=secret

# Security
JWT_SECRET=your-production-secret
```

### Load Environment

```go
// main.go
func main() {
    // Load .env file if exists
    if _, err := os.Stat(".env"); err == nil {
        env.Load(".env")
    }

    // Or use environment-specific files
    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" {
        appEnv = "production"
    }
    env.Load(".env." + appEnv)
}
```

## Quick Docker Deployment

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary
COPY --from=builder /app/myapp .

# Copy config (if needed)
COPY .env.production .env

EXPOSE 8080

CMD ["./myapp"]
```

### Build and Run

```bash
# Build image
docker build -t myapp:latest .

# Run container
docker run -d \
    -p 8080:8080 \
    -e APP_ENV=production \
    --name myapp \
    myapp:latest
```

## Health Checks

### Health Endpoint

```go
type HealthController struct {
    db    *gorm.DB  `inject:""`
    redis *kv.Client `inject:""`
}

func (c *HealthController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/health", Handler: c.Health},
        {Method: "GET", Path: "/health/ready", Handler: c.Ready},
        {Method: "GET", Path: "/health/live", Handler: c.Live},
    }
}

func (c *HealthController) Health(ctx types.Context) any {
    return map[string]string{"status": "healthy"}
}

func (c *HealthController) Ready(ctx types.Context) any {
    // Check all dependencies
    checks := map[string]string{}

    // Database
    if err := c.checkDB(); err != nil {
        checks["database"] = "unhealthy"
    } else {
        checks["database"] = "healthy"
    }

    // Redis
    if err := c.checkRedis(); err != nil {
        checks["redis"] = "unhealthy"
    } else {
        checks["redis"] = "healthy"
    }

    // Determine overall status
    for _, status := range checks {
        if status == "unhealthy" {
            return ctx.Status(503).JSON(checks)
        }
    }

    return checks
}

func (c *HealthController) Live(ctx types.Context) any {
    return map[string]string{"status": "alive"}
}
```

## Graceful Shutdown

```go
func main() {
    // Start application
    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        log.Fatal(err)
    }

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := stop(); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited")
}
```

## Logging in Production

```go
type LogConfig struct {
    Level  string
    Format string // json or text
    Output string // stdout, file
}

func setupLogging() {
    config := LogConfig{
        Level:  env.String("LOG_LEVEL", "info"),
        Format: env.String("LOG_FORMAT", "json"),
        Output: env.String("LOG_OUTPUT", "stdout"),
    }

    // Configure logger based on environment
    if env.String("APP_ENV", "") == "production" {
        // Use JSON format for production (easier to parse)
        log.SetFormatter("json")
        log.SetLevel("info")
    }
}
```

## Production Checklist

### Before Deployment

- [ ] Set `APP_ENV=production`
- [ ] Disable debug mode
- [ ] Use production database
- [ ] Set secure secrets
- [ ] Configure HTTPS
- [ ] Set up health checks
- [ ] Configure logging
- [ ] Set up monitoring

### Security

- [ ] Use strong JWT secrets
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS properly
- [ ] Set security headers
- [ ] Enable rate limiting
- [ ] Validate all input

### Performance

- [ ] Configure connection pools
- [ ] Set up caching
- [ ] Enable compression
- [ ] Optimize database queries
- [ ] Use CDN for static assets

### Monitoring

- [ ] Set up error tracking
- [ ] Configure log aggregation
- [ ] Monitor health endpoints
- [ ] Set up alerts
- [ ] Track metrics

## Process Management

### Using systemd

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Goose Application
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=always
RestartSec=5
Environment=APP_ENV=production

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable myapp
sudo systemctl start myapp

# Check status
sudo systemctl status myapp
```

## Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Next Steps

- [Docker Deployment](docker.md) - Container deployment
- [Production Best Practices](production.md) - Optimization tips
- [Monitoring](monitoring.md) - Set up observability
