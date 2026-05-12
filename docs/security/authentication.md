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

```go
type AuthMiddleware struct {
    authService *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Get token from header
    authHeader := ctx.Header("Authorization")
    if authHeader == "" {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Missing authorization header",
        })
    }

    // Parse Bearer token
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid authorization format",
        })
    }

    // Validate token
    claims, err := m.authService.ValidateToken(tokenString)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid token",
        })
    }

    // Set user in context
    ctx.Set("user_id", claims.UserID)
    ctx.Set("user_email", claims.Email)
    ctx.Set("user_role", claims.Role)

    return next()
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

func (c *AuthController) Login(ctx types.Context) any {
    var dto LoginDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": "Invalid request",
        })
    }

    // Find user
    user, err := c.userService.GetByEmail(dto.Email)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid credentials",
        })
    }

    // Check password
    if !CheckPassword(dto.Password, user.Password) {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid credentials",
        })
    }

    // Generate token
    token, err := c.authService.GenerateToken(user)
    if err != nil {
        return ctx.Status(500).JSON(map[string]string{
            "error": "Failed to generate token",
        })
    }

    return LoginResponse{
        Token:     token,
        ExpiresIn: 86400, // 24 hours
        User:      user,
    }
}
```

### Registration

```go
type RegisterDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

func (c *AuthController) Register(ctx types.Context) any {
    var dto RegisterDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.Status(400).JSON(map[string]string{
            "error": "Invalid request",
        })
    }

    // Check if user exists
    existing, _ := c.userService.GetByEmail(dto.Email)
    if existing != nil {
        return ctx.Status(409).JSON(map[string]string{
            "error": "Email already registered",
        })
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
        return ctx.Status(500).JSON(map[string]string{
            "error": "Failed to create user",
        })
    }

    // Generate token
    token, _ := c.authService.GenerateToken(user)

    return ctx.Status(201).JSON(LoginResponse{
        Token:     token,
        ExpiresIn: 86400,
        User:      user,
    })
}
```

### Routes

```go
func (c *AuthController) Routes() types.Routes {
    return types.Routes{
        {Method: "POST", Path: "/auth/login", Handler: c.Login},
        {Method: "POST", Path: "/auth/register", Handler: c.Register},
        {Method: "POST", Path: "/auth/logout", Handler: c.Logout},
        {Method: "GET", Path: "/auth/me", Handler: c.Me, Middlewares: []any{&AuthMiddleware{}}},
    }
}
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

func (m *SessionMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Get session from cookie
    sessionID := ctx.Cookie("session_id")
    if sessionID == "" {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Not authenticated",
        })
    }

    // Validate session
    session, err := m.sessionService.Get(sessionID)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid session",
        })
    }

    // Set user in context
    ctx.Set("user_id", session.UserID)

    return next()
}
```

## API Key Authentication

```go
type APIKeyMiddleware struct {
    apiKeyService *APIKeyService `inject:""`
}

func (m *APIKeyMiddleware) Handle(ctx types.Context, next types.Next) any {
    apiKey := ctx.Header("X-API-Key")
    if apiKey == "" {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Missing API key",
        })
    }

    key, err := m.apiKeyService.Validate(apiKey)
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{
            "error": "Invalid API key",
        })
    }

    ctx.Set("api_key", key)
    ctx.Set("user_id", key.UserID)

    return next()
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
