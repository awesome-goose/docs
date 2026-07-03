# DTOs & Validation

Data Transfer Objects (DTOs) and validation ensure data integrity in your Goose applications.

## Data Transfer Objects

DTOs define the shape of data for requests and responses:

```go
// Request DTO
type CreateUserDTO struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// Response DTO
type UserResponseDTO struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // Password excluded from response
}
```

## Validation Tags

Use struct tags for declarative validation:

### Required Fields

```go
type CreateUserDTO struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required"`
}
```

### String Validation

```go
type ProfileDTO struct {
    Username    string `json:"username" validate:"required,min=3,max=20"`
    Bio         string `json:"bio" validate:"max=500"`
    Website     string `json:"website" validate:"omitempty,url"`
    TwitterHandle string `json:"twitter" validate:"omitempty,startswith=@"`
}
```

### Numeric Validation

```go
type ProductDTO struct {
    Price    float64 `json:"price" validate:"required,gt=0"`
    Quantity int     `json:"quantity" validate:"required,gte=0,lte=1000"`
    Rating   float64 `json:"rating" validate:"gte=0,lte=5"`
}
```

### Email & URL

```go
type ContactDTO struct {
    Email   string `json:"email" validate:"required,email"`
    Website string `json:"website" validate:"omitempty,url"`
}
```

### Enum/OneOf

```go
type OrderDTO struct {
    Status   string `json:"status" validate:"oneof=pending processing shipped delivered"`
    Priority string `json:"priority" validate:"oneof=low medium high"`
}
```

### Nested Validation

```go
type OrderDTO struct {
    Customer CustomerDTO   `json:"customer" validate:"required"`
    Items    []OrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CustomerDTO struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

type OrderItemDTO struct {
    ProductID string `json:"product_id" validate:"required,uuid"`
    Quantity  int    `json:"quantity" validate:"required,gt=0"`
}
```

## Common Validation Tags

| Tag          | Description             | Example                      |
| ------------ | ----------------------- | ---------------------------- |
| `required`   | Field must be present   | `validate:"required"`        |
| `email`      | Valid email format      | `validate:"email"`           |
| `url`        | Valid URL format        | `validate:"url"`             |
| `uuid`       | Valid UUID format       | `validate:"uuid"`            |
| `min`        | Minimum length/value    | `validate:"min=3"`           |
| `max`        | Maximum length/value    | `validate:"max=100"`         |
| `len`        | Exact length            | `validate:"len=10"`          |
| `gt`         | Greater than            | `validate:"gt=0"`            |
| `gte`        | Greater than or equal   | `validate:"gte=0"`           |
| `lt`         | Less than               | `validate:"lt=100"`          |
| `lte`        | Less than or equal      | `validate:"lte=100"`         |
| `oneof`      | Must be one of values   | `validate:"oneof=a b c"`     |
| `contains`   | Must contain substring  | `validate:"contains=@"`      |
| `startswith` | Must start with         | `validate:"startswith=http"` |
| `endswith`   | Must end with           | `validate:"endswith=.com"`   |
| `omitempty`  | Skip if empty           | `validate:"omitempty,email"` |
| `dive`       | Validate array elements | `validate:"dive,required"`   |

## Using Validation in Controllers

The DTO is **bound automatically** from the request before your handler runs — there's no manual `Bind` step. To check it, inject a `types.Validator` and validate the populated DTO, returning `output.UnprocessableEntity` on failure:

```go
type UserController struct {
    service   *UserService    `inject:""`
    validator types.Validator `inject:""`
}

func (c *UserController) Create(dto *CreateUserDTO) types.Output {
    // dto is already populated from the request body.
    if err := c.validator.Validate(dto); err != nil {
        return output.UnprocessableEntity("Validation failed", formatValidationErrors(err))
    }

    return output.Created(c.service.CreateUser(dto))
}
```

## Custom Validation

### Custom Validate Method

Implement complex validation logic:

```go
type CreateOrderDTO struct {
    CustomerID string    `json:"customer_id" validate:"required"`
    Items      []ItemDTO `json:"items" validate:"required,min=1,dive"`
    ScheduledAt *time.Time `json:"scheduled_at"`
}

// Custom validation method
func (dto *CreateOrderDTO) Validate() error {
    // Check scheduled date is in the future
    if dto.ScheduledAt != nil && dto.ScheduledAt.Before(time.Now()) {
        return fmt.Errorf("scheduled_at must be in the future")
    }

    // Check for duplicate items
    seen := make(map[string]bool)
    for _, item := range dto.Items {
        if seen[item.ProductID] {
            return fmt.Errorf("duplicate product: %s", item.ProductID)
        }
        seen[item.ProductID] = true
    }

    return nil
}
```

### Using Custom Validation

```go
func (c *OrderController) Create(dto *CreateOrderDTO) types.Output {
    // Tag-based validation
    if err := c.validator.Validate(dto); err != nil {
        return output.UnprocessableEntity("Validation failed", err.Error())
    }

    // Custom validation
    if err := dto.Validate(); err != nil {
        return output.UnprocessableEntity(err.Error(), nil)
    }

    return output.Created(c.service.CreateOrder(dto))
}
```

## Validation Error Formatting

Format validation errors for API responses:

```go
func formatValidationErrors(err error) []map[string]string {
    var errors []map[string]string

    for _, e := range err.(validator.ValidationErrors) {
        errors = append(errors, map[string]string{
            "field":   e.Field(),
            "tag":     e.Tag(),
            "message": getErrorMessage(e),
        })
    }

    return errors
}

func getErrorMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return fmt.Sprintf("%s is required", e.Field())
    case "email":
        return fmt.Sprintf("%s must be a valid email", e.Field())
    case "min":
        return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
    case "max":
        return fmt.Sprintf("%s must be at most %s", e.Field(), e.Param())
    default:
        return fmt.Sprintf("%s is invalid", e.Field())
    }
}
```

## DTO Patterns

### Request/Response DTOs

```go
// Request DTOs (input)
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// Response DTOs (output)
type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    // Password excluded
}

// Conversion
func ToUserResponse(user *User) *UserResponse {
    return &UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }
}
```

### Query DTOs

```go
type ListUsersQuery struct {
    Page    int    `query:"page" validate:"gte=1"`
    Limit   int    `query:"limit" validate:"gte=1,lte=100"`
    Search  string `query:"search" validate:"max=100"`
    SortBy  string `query:"sort_by" validate:"oneof=name email created_at"`
    Order   string `query:"order" validate:"oneof=asc desc"`
}

func (c *UserController) List(dto *ListUsersQuery) types.Output {
    // Query params are bound automatically; apply defaults for omitted values.
    if dto.Page == 0 {
        dto.Page = 1
    }
    if dto.Limit == 0 {
        dto.Limit = 10
    }
    if dto.Order == "" {
        dto.Order = "asc"
    }

    if err := c.validator.Validate(dto); err != nil {
        return output.UnprocessableEntity("Validation failed", formatValidationErrors(err))
    }

    return output.JSON(c.service.ListUsers(dto))
}
```

## Best Practices

### 1. Separate Request and Response DTOs

```go
// ✅ Good: Separate DTOs
type CreateUserRequest struct { ... }  // Input
type UserResponse struct { ... }       // Output

// ❌ Bad: Same DTO for both
type UserDTO struct { ... }  // Used for input and output
```

### 2. Keep DTOs Simple

```go
// ✅ Good: Flat structure
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ❌ Bad: Unnecessary nesting
type CreateUserRequest struct {
    User struct {
        Info struct {
            Name string `json:"name"`
        } `json:"info"`
    } `json:"user"`
}
```

### 3. Use Specific Types

```go
// ✅ Good: Specific types
type UpdateOrderRequest struct {
    Status OrderStatus `json:"status" validate:"required"`
}

type OrderStatus string
const (
    OrderStatusPending OrderStatus = "pending"
    OrderStatusShipped OrderStatus = "shipped"
)

// ❌ Bad: Generic string
type UpdateOrderRequest struct {
    Status string `json:"status"`
}
```

## Next Steps

- [Requests](requests.md) - Handling requests
- [Responses](responses.md) - Sending responses
- [Error Handling](error-handling.md) - Managing errors
