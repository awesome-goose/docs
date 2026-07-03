# Authentication

Implement user authentication in Goose applications.

## Overview

Authentication verifies user identity. Goose provides patterns for:

- JWT (JSON Web Tokens)
- Session-based authentication
- API keys
- OAuth2

## JWT Authentication

### Setup

```go
import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
    userService *UserService `inject:""`
    env         types.Env    `inject:""`
}

func (s *AuthService) secret() string {
    return s.env.Get("JWT_SECRET")
}
```

### Generate Token

```go
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(user *User) (string, error) {
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secret()))
}
```

### Validate Token

```go
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(s.secret()), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

### Auth Middleware

Middleware implements `Handle(ctx types.Context) error`. Returning an `error` aborts the request; returning `nil` continues to the handler. Stash the authenticated identity with `ctx.SetValue` so handlers can read it via a `context`-tagged DTO field.

```go
import "errors"

type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    // Get token from header
    authHeader := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        authHeader = h[0]
    }
    if authHeader == "" {
        return errors.New("missing authorization header")
    }

    // Parse Bearer token
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        return errors.New("invalid authorization format")
    }

    // Validate token
    claims, err := m.authService.ValidateToken(tokenString)
    if err != nil {
        return errors.New("invalid token")
    }

    // Set user in context for downstream handlers
    ctx.SetValue("user_id", claims.UserID)
    ctx.SetValue("user_email", claims.Email)
    ctx.SetValue("user_role", claims.Role)

    return nil
}
```

### Login Endpoint

```go
type AuthController struct {
    authService *AuthService `inject:""`
    userService *UserService `inject:""`
}

type LoginDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token     string `json:"token"`
    ExpiresIn int64  `json:"expires_in"`
    User      *User  `json:"user"`
}

func (c *AuthController) Login(dto *LoginDTO) types.Output {
    // Find user
    user, err := c.userService.GetByEmail(dto.Email)
    if err != nil {
        return output.Unauthorized("Invalid credentials")
    }

    // Check password
    if !CheckPassword(dto.Password, user.Password) {
        return output.Unauthorized("Invalid credentials")
    }

    // Generate token
    token, err := c.authService.GenerateToken(user)
    if err != nil {
        return output.InternalServerError("Failed to generate token")
    }

    return output.JSON(LoginResponse{
        Token:     token,
        ExpiresIn: 86400, // 24 hours
        User:      user,
    })
}
```

### Registration

```go
type RegisterDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

func (c *AuthController) Register(dto *RegisterDTO) types.Output {
    // Check if user exists
    existing, _ := c.userService.GetByEmail(dto.Email)
    if existing != nil {
        return output.Conflict("Email already registered")
    }

    // Hash password
    hashedPassword, _ := HashPassword(dto.Password)

    // Create user
    user, err := c.userService.Create(&User{
        Email:    dto.Email,
        Password: hashedPassword,
        Name:     dto.Name,
    })
    if err != nil {
        return output.InternalServerError("Failed to create user")
    }

    // Generate token
    token, _ := c.authService.GenerateToken(user)

    return output.Created(LoginResponse{
        Token:     token,
        ExpiresIn: 86400,
        User:      user,
    })
}
```

### Routes

```go
var ROUTES = router.ForRoutes(
    router.Post("/auth/login", []any{AuthController{}, "Login"}),
    router.Post("/auth/register", []any{AuthController{}, "Register"}),
    router.Post("/auth/logout", []any{AuthController{}, "Logout"}),
    router.Get("/auth/me", []any{AuthController{}, "Me"}, &AuthMiddleware{}),
)
```

## Session Authentication

### Setup Session Store

```go
import "github.com/awesome-goose/goose/modules/kv"

type SessionService struct {
    kv *kv.KV `inject:""`
}

func (s *SessionService) Create(userID string) (string, error) {
    sessionID := uuid.New().String()
    return sessionID, s.kv.Set("session:"+sessionID, map[string]any{
        "user_id":    userID,
        "created_at": time.Now(),
    }, 24*time.Hour)
}

func (s *SessionService) Get(sessionID string) (*Session, error) {
    v, err := s.kv.Get("session:" + sessionID)
    if err != nil {
        return nil, err
    }
    // kv stores JSON-decoded values; coerce as needed.
    raw, _ := v.(map[string]any)
    return &Session{UserID: raw["user_id"].(string)}, nil
}

func (s *SessionService) Destroy(sessionID string) error {
    _, err := s.kv.Del("session:" + sessionID)
    return err
}
```

### Session Middleware

```go
type SessionMiddleware struct {
    sessionService *SessionService `inject:""`
}

func (m *SessionMiddleware) Handle(ctx types.Context) error {
    // Read the session id from the Cookie header
    cookieHeader := ""
    if h := ctx.Request().Headers()["Cookie"]; len(h) > 0 {
        cookieHeader = h[0]
    }
    sessionID := parseCookie(cookieHeader, "session_id") // small helper
    if sessionID == "" {
        return errors.New("not authenticated")
    }

    // Validate session
    session, err := m.sessionService.Get(sessionID)
    if err != nil {
        return errors.New("invalid session")
    }

    // Set user in context
    ctx.SetValue("user_id", session.UserID)

    return nil
}
```

## API Key Authentication

```go
type APIKeyMiddleware struct {
    apiKeyService *APIKeyService `inject:""`
}

func (m *APIKeyMiddleware) Handle(ctx types.Context) error {
    apiKey := ""
    if h := ctx.Request().Headers()["X-API-Key"]; len(h) > 0 {
        apiKey = h[0]
    }
    if apiKey == "" {
        return errors.New("missing API key")
    }

    key, err := m.apiKeyService.Validate(apiKey)
    if err != nil {
        return errors.New("invalid API key")
    }

    ctx.SetValue("api_key", key)
    ctx.SetValue("user_id", key.UserID)

    return nil
}
```

## Password Utilities

```go
import "golang.org/x/crypto/bcrypt"

// Hash password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// Check password
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## Refresh Tokens

```go
type RefreshTokenService struct {
    kv          *kv.KV       `inject:""`
    authService *AuthService `inject:""`
    userService *UserService `inject:""`
}

func (s *RefreshTokenService) Generate(userID string) (string, error) {
    refreshToken := uuid.New().String()
    return refreshToken, s.kv.Set("refresh:"+refreshToken, userID, 7*24*time.Hour)
}

func (s *RefreshTokenService) Refresh(refreshToken string) (*TokenPair, error) {
    v, err := s.kv.Get("refresh:" + refreshToken)
    if err != nil {
        return nil, fmt.Errorf("invalid refresh token")
    }
    userID, _ := v.(string)

    // Invalidate old refresh token
    _, _ = s.kv.Del("refresh:" + refreshToken)

    // Generate new tokens
    user, _ := s.userService.GetByID(userID)
    accessToken, _ := s.authService.GenerateToken(user)
    newRefreshToken, _ := s.Generate(userID)

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
    }, nil
}
```

## Best Practices

1. **Use HTTPS** for all authentication endpoints
2. **Hash passwords** with bcrypt or argon2
3. **Set appropriate token expiry** (short for access, longer for refresh)
4. **Implement rate limiting** on login endpoints
5. **Log authentication events** for security monitoring
6. **Use secure cookies** (HttpOnly, Secure, SameSite)
7. **Validate token on every request**

## Next Steps

- [Authorization](authorization.md) - Access control
- [CORS](cors.md) - Cross-origin configuration
- [Rate Limiting](rate-limiting.md) - Prevent brute force
