# Error Handling

Learn how to handle errors gracefully in your Goose applications.

## Basic Error Handling

### In Controllers

```go
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetUserByID(ctx.Param("id"))
    if err != nil {
        return ctx.Error(404, "User not found")
    }
    return user
}
```

### In Services

```go
func (s *UserService) GetUserByID(id string) (*User, error) {
    var user User
    err := s.db.First(&user, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    return &user, nil
}
```

## Custom Errors

### Defining Custom Errors

```go
package errors

import "errors"

var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrValidation    = errors.New("validation error")
    ErrDuplicate     = errors.New("duplicate entry")
    ErrInternalError = errors.New("internal error")
)

// Domain-specific errors
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidPassword = errors.New("invalid password")
    ErrEmailInUse      = errors.New("email already in use")
)
```

### Using Custom Errors

```go
func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    // Check for duplicate
    existing, err := s.repo.FindByEmail(dto.Email)
    if err == nil && existing != nil {
        return nil, ErrEmailInUse
    }

    // Create user...
}

// In controller
func (c *UserController) Create(ctx types.Context) any {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        if errors.Is(err, ErrEmailInUse) {
            return ctx.Error(409, "Email already in use")
        }
        return ctx.Error(500, "Failed to create user")
    }
    return user
}
```

## Error Types

### Structured Errors

```go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(code, message string, details any) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

// Usage
func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
    if exists := s.repo.EmailExists(dto.Email); exists {
        return nil, NewAppError(
            "EMAIL_EXISTS",
            "Email address is already registered",
            map[string]string{"email": dto.Email},
        )
    }
    // ...
}
```

### Validation Errors

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
    if len(e) == 0 {
        return "validation error"
    }
    return e[0].Message
}

// Usage
func ValidateUser(dto CreateUserDTO) ValidationErrors {
    var errors ValidationErrors

    if dto.Name == "" {
        errors = append(errors, ValidationError{
            Field:   "name",
            Message: "Name is required",
        })
    }

    if !isValidEmail(dto.Email) {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: "Invalid email format",
        })
    }

    return errors
}
```

## HTTP Error Responses

### Standard Error Response

```go
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetUserByID(ctx.Param("id"))
    if err != nil {
        return c.handleError(ctx, err)
    }
    return user
}

func (c *UserController) handleError(ctx types.Context, err error) any {
    switch {
    case errors.Is(err, ErrNotFound):
        return ctx.Error(404, map[string]any{
            "error": "not_found",
            "message": "Resource not found",
        })
    case errors.Is(err, ErrUnauthorized):
        return ctx.Error(401, map[string]any{
            "error": "unauthorized",
            "message": "Authentication required",
        })
    case errors.Is(err, ErrForbidden):
        return ctx.Error(403, map[string]any{
            "error": "forbidden",
            "message": "Access denied",
        })
    case errors.Is(err, ErrValidation):
        return ctx.Error(422, map[string]any{
            "error": "validation_failed",
            "message": err.Error(),
        })
    default:
        c.log.Error("Unhandled error", "error", err)
        return ctx.Error(500, map[string]any{
            "error": "internal_error",
            "message": "An unexpected error occurred",
        })
    }
}
```

### Error Helper Functions

```go
func NotFound(ctx types.Context, resource string) any {
    return ctx.Error(404, map[string]any{
        "error": "not_found",
        "message": fmt.Sprintf("%s not found", resource),
    })
}

func BadRequest(ctx types.Context, message string) any {
    return ctx.Error(400, map[string]any{
        "error": "bad_request",
        "message": message,
    })
}

func Unauthorized(ctx types.Context) any {
    return ctx.Error(401, map[string]any{
        "error": "unauthorized",
        "message": "Authentication required",
    })
}

func InternalError(ctx types.Context) any {
    return ctx.Error(500, map[string]any{
        "error": "internal_error",
        "message": "An unexpected error occurred",
    })
}

// Usage
func (c *UserController) Show(ctx types.Context) any {
    user, err := c.service.GetUserByID(ctx.Param("id"))
    if err != nil {
        return NotFound(ctx, "User")
    }
    return user
}
```

## Error Wrapping

Use Go's error wrapping to preserve context:

```go
func (s *UserService) GetUserByID(id string) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}

func (s *OrderService) CreateOrder(dto CreateOrderDTO) (*Order, error) {
    user, err := s.userService.GetUserByID(dto.UserID)
    if err != nil {
        return nil, fmt.Errorf("create order failed: %w", err)
    }
    // ...
}

// Unwrap to check specific errors
if errors.Is(err, gorm.ErrRecordNotFound) {
    // Handle not found
}
```

## Panic Recovery

Goose automatically recovers from panics, but you can add custom handling:

```go
type RecoveryMiddleware struct {
    log types.Log `inject:""`
}

func (m *RecoveryMiddleware) Handle(ctx types.Context) error {
    defer func() {
        if r := recover(); r != nil {
            m.log.Error("Panic recovered",
                "panic", r,
                "stack", string(debug.Stack()))
            ctx.SetStatus(500)
            ctx.SetResponse(map[string]string{
                "error": "internal_error",
                "message": "An unexpected error occurred",
            })
        }
    }()
    return nil
}
```

## Logging Errors

Always log errors with context:

```go
func (s *UserService) DeleteUser(id string) error {
    err := s.repo.Delete(id)
    if err != nil {
        s.log.Error("Failed to delete user",
            "user_id", id,
            "error", err,
        )
        return fmt.Errorf("delete user failed: %w", err)
    }

    s.log.Info("User deleted", "user_id", id)
    return nil
}
```

## Best Practices

### 1. Don't Swallow Errors

```go
// ✅ Good: Return errors
func (s *Service) DoSomething() error {
    if err := s.operation(); err != nil {
        return err
    }
    return nil
}

// ❌ Bad: Swallowed error
func (s *Service) DoSomething() {
    _ = s.operation()  // Error ignored
}
```

### 2. Add Context to Errors

```go
// ✅ Good: Error with context
return fmt.Errorf("failed to process order %s: %w", orderID, err)

// ❌ Bad: No context
return err
```

### 3. Use Appropriate HTTP Status Codes

| Status | When to Use          |
| ------ | -------------------- |
| 400    | Invalid request body |
| 401    | Not authenticated    |
| 403    | Not authorized       |
| 404    | Resource not found   |
| 409    | Conflict (duplicate) |
| 422    | Validation error     |
| 500    | Server error         |

### 4. Don't Expose Internal Errors

```go
// ✅ Good: Generic message to client
func (c *Controller) Handle(ctx types.Context) any {
    if err != nil {
        c.log.Error("Database error", "error", err)
        return ctx.Error(500, "An error occurred")
    }
}

// ❌ Bad: Expose internal details
func (c *Controller) Handle(ctx types.Context) any {
    if err != nil {
        return ctx.Error(500, err.Error())  // May expose DB details
    }
}
```

## Next Steps

- [Logging](logging.md) - Log errors and events
- [Middleware](middleware.md) - Error handling middleware
- [Responses](responses.md) - Error response formats
