# Security Overview

Goose provides built-in security features to protect your applications from common vulnerabilities.

## Security Features

| Feature         | Description                        | Documentation                       |
| --------------- | ---------------------------------- | ----------------------------------- |
| Authentication  | Verify user identity               | [Authentication](authentication.md) |
| Authorization   | Control access to resources        | [Authorization](authorization.md)   |
| CORS            | Cross-Origin Resource Sharing      | [CORS](cors.md)                     |
| CSRF Protection | Prevent cross-site request forgery | [CSRF](csrf.md)                     |
| Rate Limiting   | Prevent abuse and DoS              | [Rate Limiting](rate-limiting.md)   |

## Quick Security Setup

### Basic Authentication

```go
import "errors"

type AuthMiddleware struct {
    userService *UserService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    token := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        token = h[0]
    }
    if token == "" {
        return errors.New("unauthorized")
    }

    user, err := m.userService.ValidateToken(token)
    if err != nil {
        return errors.New("invalid token")
    }

    ctx.SetValue("user", user)
    return nil
}
```

### Apply to Routes

```go
var ROUTES = router.ForRoutes(
    router.Get("/public", []any{Controller{}, "Public"}),
    router.Get("/private", []any{Controller{}, "Private"}, &AuthMiddleware{}),
)
```

## Security Best Practices

### 1. Environment Variables

Never hardcode secrets:

```go
// Good (e is a *env.Env injected as types.Env)
jwtSecret := e.Get("JWT_SECRET")
if jwtSecret == "" {
    panic("JWT_SECRET is required")
}

// Bad
jwtSecret := "my-secret-key" // Never do this!
```

### 2. Input Validation

Always validate user input:

```go
type CreateUserDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,max=100"`
}

func (c *Controller) Create(dto *CreateUserDTO) types.Output {
    // dto is bound automatically; validate it with an injected types.Validator
    if err := c.validator.Validate(dto); err != nil {
        return output.UnprocessableEntity("Validation failed", err)
    }
    // Proceed with validated data
    return output.Created(c.service.Create(dto))
}
```

### 3. Password Hashing

Never store plain passwords:

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 4. SQL Injection Prevention

Use parameterized queries:

```go
// Good - parameterized
db.Where("email = ?", email).First(&user)

// Bad - string concatenation
db.Raw("SELECT * FROM users WHERE email = '" + email + "'") // Vulnerable!
```

### 5. Error Messages

Don't expose sensitive information:

```go
// Good - generic message
if err != nil {
    return output.Unauthorized("Invalid credentials")
}

// Bad - reveals too much
if err != nil {
    return output.Unauthorized("User not found in database table users") // Reveals info!
}
```

## Security Headers

Add security headers:

```go
type SecurityMiddleware struct{}

func (m *SecurityMiddleware) Handle(ctx types.Context) error {
    resp := ctx.Response()

    // Prevent XSS
    resp.SetHeader("X-XSS-Protection", "1; mode=block")

    // Prevent clickjacking
    resp.SetHeader("X-Frame-Options", "DENY")

    // Prevent MIME sniffing
    resp.SetHeader("X-Content-Type-Options", "nosniff")

    // HSTS
    resp.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

    return nil
}
```

## HTTPS

Always use HTTPS in production. The platform listens over plain HTTP, so terminate TLS at a reverse proxy (nginx, Caddy, or your cloud load balancer) in front of the app:

```go
platform := api.NewPlatform(
    api.WithHost("0.0.0.0"),
    api.WithPort(8080), // proxy forwards HTTPS traffic here
)
```

## Security Checklist

### Authentication

- [ ] Use strong password requirements
- [ ] Implement account lockout
- [ ] Use secure session management
- [ ] Implement proper logout

### Authorization

- [ ] Implement role-based access
- [ ] Validate permissions on every request
- [ ] Use principle of least privilege

### Data Protection

- [ ] Hash passwords with bcrypt
- [ ] Encrypt sensitive data
- [ ] Use HTTPS everywhere
- [ ] Sanitize all input

### Infrastructure

- [ ] Keep dependencies updated
- [ ] Use security headers
- [ ] Implement rate limiting
- [ ] Log security events

## Common Vulnerabilities

### Injection Attacks

- Use parameterized queries
- Validate and sanitize input
- Use ORM methods

### Broken Authentication

- Use secure password hashing
- Implement MFA where possible
- Secure session management

### XSS (Cross-Site Scripting)

- Escape output by default
- Use Content Security Policy
- Validate input

### CSRF (Cross-Site Request Forgery)

- Use CSRF tokens
- Validate origin headers
- Use SameSite cookies

## Next Steps

- [Authentication](authentication.md) - User authentication
- [Authorization](authorization.md) - Access control
- [CORS](cors.md) - Cross-origin requests
- [Rate Limiting](rate-limiting.md) - Abuse prevention
