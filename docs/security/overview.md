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
type AuthMiddleware struct {
    userService *UserService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context, next types.Next) any {
    token := ctx.Header("Authorization")
    if token == "" {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Unauthorized",
        })
    }

    user, err := m.userService.ValidateToken(token)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid token",
        })
    }

    ctx.Set("user", user)
    return next()
}
```

### Apply to Routes

```go
func (c *Controller) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/public", Handler: c.Public},
        {Method: "GET", Path: "/private", Handler: c.Private, Middlewares: []any{&AuthMiddleware{}}},
    }
}
```

## Security Best Practices

### 1. Environment Variables

Never hardcode secrets:

```go
// Good
jwtSecret := env.String("JWT_SECRET", "")
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

func (c *Controller) Create(ctx types.Context) any {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Status(400).JSON(errors.ValidationError(err))
    }
    // Proceed with validated data
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
    return ctx.Status(401).JSON(map[string]string{
        "error": "Invalid credentials",
    })
}

// Bad - reveals too much
if err != nil {
    return ctx.Status(401).JSON(map[string]string{
        "error": "User not found in database table users", // Reveals info!
    })
}
```

## Security Headers

Add security headers:

```go
type SecurityMiddleware struct{}

func (m *SecurityMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Prevent XSS
    ctx.SetHeader("X-XSS-Protection", "1; mode=block")

    // Prevent clickjacking
    ctx.SetHeader("X-Frame-Options", "DENY")

    // Prevent MIME sniffing
    ctx.SetHeader("X-Content-Type-Options", "nosniff")

    // HSTS
    ctx.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

    return next()
}
```

## HTTPS

Always use HTTPS in production:

```go
// In production, enable HTTPS
platform := api.NewPlatform(
    api.WithHost("0.0.0.0"),
    api.WithPort(443),
    api.WithTLS("/path/to/cert.pem", "/path/to/key.pem"),
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
