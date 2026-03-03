# Docker Deployment

Deploy your Goose application using Docker containers.

## Basic Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server main.go

# Runtime stage
FROM alpine:3.19

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files if needed
# COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/public ./public

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
CMD ["./server"]
```

## Multi-Stage Build (Optimized)

```dockerfile
# Dependencies stage
FROM golang:1.21-alpine AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Build stage
FROM deps AS builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server main.go

# Test stage (optional)
FROM builder AS tester

RUN go test -v ./...

# Runtime stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
```

## Docker Compose

### Basic Setup

```yaml
# docker-compose.yml
version: "3.8"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

### With Nginx

```yaml
version: "3.8"

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - app
    restart: unless-stopped

  app:
    build: .
    expose:
      - "8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    deploy:
      replicas: 2

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

### Development Compose

```yaml
# docker-compose.dev.yml
version: "3.8"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - .:/app
      - go_modules:/go/pkg/mod
    environment:
      - APP_ENV=development
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: myapp_dev
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_dev:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  go_modules:
  postgres_dev:
```

## Environment Files

### .env.production

```env
APP_ENV=production
APP_DEBUG=false
HOST=0.0.0.0
PORT=8080

DB_DRIVER=postgres
DB_HOST=postgres
DB_PORT=5432
DB_NAME=myapp
DB_USER=myapp
DB_PASSWORD=secret

REDIS_HOST=redis
REDIS_PORT=6379

JWT_SECRET=your-production-secret
```

### Using Environment Files

```yaml
services:
  app:
    build: .
    env_file:
      - .env.production
```

## Building Images

### Build Commands

```bash
# Build image
docker build -t myapp:latest .

# Build with tag
docker build -t myapp:0.0.0 .

# Build with build args
docker build --build-arg VERSION=0.0.0 -t myapp:0.0.0 .

# Build for specific platform
docker build --platform linux/amd64 -t myapp:latest .
```

### Push to Registry

```bash
# Tag for registry
docker tag myapp:latest registry.example.com/myapp:latest

# Push to registry
docker push registry.example.com/myapp:latest

# Push to Docker Hub
docker tag myapp:latest username/myapp:latest
docker push username/myapp:latest
```

## Running Containers

### Basic Run

```bash
# Run container
docker run -d \
    --name myapp \
    -p 8080:8080 \
    -e APP_ENV=production \
    myapp:latest

# Run with environment file
docker run -d \
    --name myapp \
    -p 8080:8080 \
    --env-file .env.production \
    myapp:latest
```

### With Docker Compose

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d app

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Health Checks

### In Dockerfile

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### In Docker Compose

```yaml
services:
  app:
    build: .
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

## Resource Limits

```yaml
services:
  app:
    build: .
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 512M
        reservations:
          cpus: "0.25"
          memory: 128M
```

## Logging

### Docker Logging

```yaml
services:
  app:
    build: .
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### View Logs

```bash
# View logs
docker logs myapp

# Follow logs
docker logs -f myapp

# Compose logs
docker-compose logs -f app
```

## Networking

### Custom Network

```yaml
version: "3.8"

services:
  app:
    networks:
      - frontend
      - backend

  postgres:
    networks:
      - backend

  nginx:
    networks:
      - frontend

networks:
  frontend:
  backend:
```

## Secrets Management

### Using Docker Secrets

```yaml
version: "3.8"

services:
  app:
    build: .
    secrets:
      - db_password
      - jwt_secret
    environment:
      DB_PASSWORD_FILE: /run/secrets/db_password
      JWT_SECRET_FILE: /run/secrets/jwt_secret

secrets:
  db_password:
    file: ./secrets/db_password.txt
  jwt_secret:
    file: ./secrets/jwt_secret.txt
```

## Production Tips

1. **Use multi-stage builds** to minimize image size
2. **Run as non-root user** for security
3. **Include health checks** for orchestration
4. **Set resource limits** to prevent runaway containers
5. **Use `.dockerignore`** to exclude unnecessary files
6. **Pin image versions** for reproducibility
7. **Scan images** for vulnerabilities

### .dockerignore

```
.git
.gitignore
*.md
Dockerfile*
docker-compose*
.env*
coverage.out
*.test
tmp/
```

## Next Steps

- [Kubernetes Deployment](kubernetes.md) - Container orchestration
- [Production Best Practices](production.md) - Optimization
- [Monitoring](monitoring.md) - Observability
