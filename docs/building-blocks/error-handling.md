# Error Handling

Learn how to handle errors gracefully in your Goose applications.

## Basic Error Handling

### In Controllers

Handlers return a `types.Output`. Translate errors into the semantic `output.*` helpers:

```go
type ShowUserDto struct {
    ID string `param:"id"`
}

func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetUserByID(dto.ID)
    if err != nil {
        return output.NotFound("User not found")
    }
    return output.JSON(user)
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
func (c *UserController) Create(dto *CreateUserDto) types.Output {
    user, err := c.service.CreateUser(dto)
    if err != nil {
        if errors.Is(err, ErrEmailInUse) {
            return output.Conflict("Email already in use")
        }
        return output.InternalServerError("Failed to create user")
    }
    return output.Created(user)
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
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetUserByID(dto.ID)
    if err != nil {
        return handleError(c.log, err)
    }
    return output.JSON(user)
}

func handleError(log types.Log, err error) types.Output {
    switch {
    case errors.Is(err, ErrNotFound):
        return output.NotFound("Resource not found")
    case errors.Is(err, ErrUnauthorized):
        return output.Unauthorized("Authentication required")
    case errors.Is(err, ErrForbidden):
        return output.Forbidden("Access denied")
    case errors.Is(err, ErrValidation):
        return output.UnprocessableEntity(err.Error(), nil)
    default:
        log.Error("Unhandled error", "error", err)
        return output.InternalServerError("An unexpected error occurred")
    }
}
```

### Built-in Error Helpers

You don't need to write your own error responders — the `io/output` package already provides them, each with the correct status code and a consistent `{"success": false, "message": ...}` envelope:

```go
output.BadRequest("...")           // 400
output.Unauthorized("...")         // 401
output.Forbidden("...")            // 403
output.NotFound("...")             // 404
output.Conflict("...")             // 409
output.UnprocessableEntity(msg, errs) // 422
output.InternalServerError("...")  // 500

// Usage
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetUserByID(dto.ID)
    if err != nil {
        return output.NotFound("User not found")
    }
    return output.JSON(user)
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

## Prefer Errors Over Panics

Handlers and services should return errors rather than panic. Convert a returned error into a response with the `output.*` helpers so the client always gets a well-formed body:

```go
func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user, err := c.service.GetUserByID(dto.ID)
    if err != nil {
        c.log.Error("Failed to load user", "id", dto.ID, "error", err)
        return output.InternalServerError("An unexpected error occurred")
    }
    return output.JSON(user)
}
```

Background workers (queues, cron) recover from panics in their jobs automatically, but request handlers do not — guard risky operations and return an error instead of letting a panic escape.

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
func (c *Controller) Handle(dto *HandleDto) types.Output {
    if err != nil {
        c.log.Error("Database error", "error", err)
        return output.InternalServerError("An error occurred")
    }
    return output.JSON(result)
}

// ❌ Bad: Expose internal details
func (c *Controller) Handle(dto *HandleDto) types.Output {
    if err != nil {
        return output.InternalServerError(err.Error())  // May expose DB details
    }
    return output.JSON(result)
}
```

## Next Steps

- [Logging](logging.md) - Log errors and events
- [Middleware](middleware.md) - Error handling middleware
- [Responses](responses.md) - Error response formats
