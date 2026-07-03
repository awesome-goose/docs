# Authorization

Control access to resources based on user roles and permissions.

## Overview

Authorization determines what an authenticated user can do. Common patterns:

- **Role-Based Access Control (RBAC)** - Users have roles with permissions
- **Attribute-Based Access Control (ABAC)** - Policies based on attributes
- **Resource-Based Access Control** - Permissions tied to specific resources

## Role-Based Access Control

### Define Roles

```go
const (
    RoleUser    = "user"
    RoleEditor  = "editor"
    RoleAdmin   = "admin"
)

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

### Role Middleware

```go
type RoleMiddleware struct {
    allowedRoles []string
}

func RequireRole(roles ...string) *RoleMiddleware {
    return &RoleMiddleware{allowedRoles: roles}
}

func (m *RoleMiddleware) Handle(ctx types.Context) error {
    // Get user role from context (set by auth middleware)
    userRole, _ := ctx.GetValue("user_role").(string)

    // Check if user has required role
    for _, role := range m.allowedRoles {
        if userRole == role {
            return nil // allowed — continue
        }
    }

    return errors.New("insufficient permissions")
}
```

### Apply to Routes

```go
var ROUTES = router.ForRoutes(
    // User routes
    router.Get("/profile", []any{AdminController{}, "Profile"}, &AuthMiddleware{}),

    // Editor routes
    router.Post("/posts", []any{AdminController{}, "CreatePost"},
        &AuthMiddleware{}, RequireRole("editor", "admin")),

    // Admin only routes
    router.Get("/admin/users", []any{AdminController{}, "ListUsers"},
        &AuthMiddleware{}, RequireRole("admin")),
    router.Delete("/admin/users/:id", []any{AdminController{}, "DeleteUser"},
        &AuthMiddleware{}, RequireRole("admin")),
)
```

Middleware runs left to right — `AuthMiddleware` populates `user_role`, then `RequireRole` checks it.

## Permission-Based Control

### Define Permissions

```go
const (
    PermissionReadPosts   = "posts:read"
    PermissionWritePosts  = "posts:write"
    PermissionDeletePosts = "posts:delete"
    PermissionReadUsers   = "users:read"
    PermissionWriteUsers  = "users:write"
    PermissionDeleteUsers = "users:delete"
)

// Role to permissions mapping
var rolePermissions = map[string][]string{
    "user": {
        PermissionReadPosts,
    },
    "editor": {
        PermissionReadPosts,
        PermissionWritePosts,
    },
    "admin": {
        PermissionReadPosts,
        PermissionWritePosts,
        PermissionDeletePosts,
        PermissionReadUsers,
        PermissionWriteUsers,
        PermissionDeleteUsers,
    },
}
```

### Permission Check

```go
func HasPermission(role string, permission string) bool {
    permissions, ok := rolePermissions[role]
    if !ok {
        return false
    }

    for _, p := range permissions {
        if p == permission {
            return true
        }
    }
    return false
}
```

### Permission Middleware

```go
type PermissionMiddleware struct {
    requiredPermission string
}

func RequirePermission(permission string) *PermissionMiddleware {
    return &PermissionMiddleware{requiredPermission: permission}
}

func (m *PermissionMiddleware) Handle(ctx types.Context) error {
    userRole, _ := ctx.GetValue("user_role").(string)

    if !HasPermission(userRole, m.requiredPermission) {
        return errors.New("permission denied")
    }

    return nil
}
```

### Apply Permissions

```go
var ROUTES = router.ForRoutes(
    router.Get("/posts", []any{PostController{}, "Index"},
        &AuthMiddleware{}, RequirePermission(PermissionReadPosts)),
    router.Post("/posts", []any{PostController{}, "Create"},
        &AuthMiddleware{}, RequirePermission(PermissionWritePosts)),
    router.Delete("/posts/:id", []any{PostController{}, "Delete"},
        &AuthMiddleware{}, RequirePermission(PermissionDeletePosts)),
)
```

## Resource-Based Authorization

### Ownership Check

```go
type PostController struct {
    postService *PostService `inject:""`
}

type UpdatePostDTO struct {
    ID       string `param:"id"`
    UserID   string `context:"user_id"`
    UserRole string `context:"user_role"`
    Title    string `json:"title"`
    Content  string `json:"content"`
}

func (c *PostController) Update(dto *UpdatePostDTO) types.Output {
    // Get the post
    post, err := c.postService.GetByID(dto.ID)
    if err != nil {
        return output.NotFound("Post not found")
    }

    // Check ownership or admin
    if post.AuthorID != dto.UserID && dto.UserRole != "admin" {
        return output.Forbidden("You can only edit your own posts")
    }

    // Proceed with update...
    return output.JSON(c.postService.Update(dto.ID, dto))
}
```

### Resource Authorization Service

```go
type AuthorizationService struct{}

func (s *AuthorizationService) CanEdit(user *User, resource interface{}) bool {
    switch r := resource.(type) {
    case *Post:
        return r.AuthorID == user.ID || user.Role == "admin"
    case *Comment:
        return r.AuthorID == user.ID || user.Role == "admin" || user.Role == "moderator"
    default:
        return false
    }
}

func (s *AuthorizationService) CanDelete(user *User, resource interface{}) bool {
    if user.Role == "admin" {
        return true
    }

    switch r := resource.(type) {
    case *Post:
        return r.AuthorID == user.ID
    case *Comment:
        return r.AuthorID == user.ID
    default:
        return false
    }
}
```

## Policy Pattern

### Define Policies

```go
type Policy interface {
    CanView(user *User, resource interface{}) bool
    CanCreate(user *User) bool
    CanUpdate(user *User, resource interface{}) bool
    CanDelete(user *User, resource interface{}) bool
}

type PostPolicy struct{}

func (p *PostPolicy) CanView(user *User, resource interface{}) bool {
    post := resource.(*Post)

    // Published posts are viewable by all
    if post.Published {
        return true
    }

    // Drafts only by author or admin
    return post.AuthorID == user.ID || user.Role == "admin"
}

func (p *PostPolicy) CanCreate(user *User) bool {
    return HasPermission(user.Role, PermissionWritePosts)
}

func (p *PostPolicy) CanUpdate(user *User, resource interface{}) bool {
    post := resource.(*Post)
    return post.AuthorID == user.ID || user.Role == "admin"
}

func (p *PostPolicy) CanDelete(user *User, resource interface{}) bool {
    post := resource.(*Post)
    return post.AuthorID == user.ID || user.Role == "admin"
}
```

### Use Policies

```go
type PostController struct {
    postService *PostService `inject:""`
    policy      *PostPolicy
}

type UpdatePostDto struct {
    ID   string `param:"id"`
    User *User  `context:"user"`
}

func (c *PostController) Update(dto *UpdatePostDto) types.Output {
    post, _ := c.postService.GetByID(dto.ID)

    if !c.policy.CanUpdate(dto.User, post) {
        return output.Forbidden("Cannot update this post")
    }

    // Proceed...
    return output.JSON(post)
}
```

## Scoped Queries

Filter data based on user:

```go
type PostService struct {
    db *gorm.DB `inject:""`
}

// Get posts user can see
func (s *PostService) GetVisiblePosts(user *User) []Post {
    var posts []Post

    query := s.db.Model(&Post{})

    if user.Role == "admin" {
        // Admin sees all
        query.Find(&posts)
    } else {
        // Others see published or own posts
        query.Where("published = ? OR author_id = ?", true, user.ID).Find(&posts)
    }

    return posts
}
```

## Guard Pattern

Group authorization logic:

```go
type AdminGuard struct {
    authService *AuthService `inject:""`
}

func (g *AdminGuard) Handle(ctx types.Context) error {
    // Check authentication
    authHeader := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        authHeader = h[0]
    }
    claims, err := g.authService.ValidateToken(authHeader)
    if err != nil {
        return errors.New("unauthorized")
    }

    // Check admin role
    if claims.Role != "admin" {
        return errors.New("admin access required")
    }

    ctx.SetValue("user_id", claims.UserID)
    ctx.SetValue("user_role", claims.Role)

    return nil
}
```

## Best Practices

1. **Fail secure** - Deny by default
2. **Check authorization on every request**
3. **Use middleware** for consistent enforcement
4. **Combine roles with resource ownership**
5. **Log authorization failures**
6. **Test authorization thoroughly**
7. **Keep authorization logic centralized**

## Common Patterns

### Admin Override

```go
func (s *AuthService) CanPerformAction(user *User, action string, resource interface{}) bool {
    // Admin can do anything
    if user.Role == "admin" {
        return true
    }

    // Check specific permissions
    return s.checkPermission(user, action, resource)
}
```

### Hierarchical Roles

```go
var roleHierarchy = map[string]int{
    "user":   1,
    "editor": 2,
    "admin":  3,
}

func HasRoleOrHigher(userRole, requiredRole string) bool {
    return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
```

## Next Steps

- [Authentication](authentication.md) - User identity
- [Middleware](../building-blocks/middleware.md) - Request processing
- [Guards](../building-blocks/guards.md) - Route protection
