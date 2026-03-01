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

func (m *RoleMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Get user role from context (set by auth middleware)
    userRole := ctx.Get("user_role").(string)

    // Check if user has required role
    for _, role := range m.allowedRoles {
        if userRole == role {
            return next()
        }
    }

    return ctx.Status(403).JSON(map[string]string{
        "error": "Insufficient permissions",
    })
}
```

### Apply to Routes

```go
func (c *AdminController) Routes() types.Routes {
    return types.Routes{
        // User routes
        {Method: "GET", Path: "/profile", Handler: c.Profile,
            Middlewares: []any{&AuthMiddleware{}}},

        // Editor routes
        {Method: "POST", Path: "/posts", Handler: c.CreatePost,
            Middlewares: []any{&AuthMiddleware{}, RequireRole("editor", "admin")}},

        // Admin only routes
        {Method: "GET", Path: "/admin/users", Handler: c.ListUsers,
            Middlewares: []any{&AuthMiddleware{}, RequireRole("admin")}},
        {Method: "DELETE", Path: "/admin/users/:id", Handler: c.DeleteUser,
            Middlewares: []any{&AuthMiddleware{}, RequireRole("admin")}},
    }
}
```

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

func (m *PermissionMiddleware) Handle(ctx types.Context, next types.Next) any {
    userRole := ctx.Get("user_role").(string)

    if !HasPermission(userRole, m.requiredPermission) {
        return ctx.Status(403).JSON(map[string]string{
            "error": "Permission denied",
        })
    }

    return next()
}
```

### Apply Permissions

```go
func (c *PostController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/posts", Handler: c.Index,
            Middlewares: []any{&AuthMiddleware{}, RequirePermission(PermissionReadPosts)}},
        {Method: "POST", Path: "/posts", Handler: c.Create,
            Middlewares: []any{&AuthMiddleware{}, RequirePermission(PermissionWritePosts)}},
        {Method: "DELETE", Path: "/posts/:id", Handler: c.Delete,
            Middlewares: []any{&AuthMiddleware{}, RequirePermission(PermissionDeletePosts)}},
    }
}
```

## Resource-Based Authorization

### Ownership Check

```go
type PostController struct {
    postService *PostService `inject:""`
}

func (c *PostController) Update(ctx types.Context) any {
    postID := ctx.Param("id")
    userID := ctx.Get("user_id").(string)
    userRole := ctx.Get("user_role").(string)

    // Get the post
    post, err := c.postService.GetByID(postID)
    if err != nil {
        return ctx.Status(404).JSON(map[string]string{
            "error": "Post not found",
        })
    }

    // Check ownership or admin
    if post.AuthorID != userID && userRole != "admin" {
        return ctx.Status(403).JSON(map[string]string{
            "error": "You can only edit your own posts",
        })
    }

    // Proceed with update...
    var dto UpdatePostDTO
    ctx.Bind(&dto)
    return c.postService.Update(postID, dto)
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

func (c *PostController) Update(ctx types.Context) any {
    user := ctx.Get("user").(*User)
    post, _ := c.postService.GetByID(ctx.Param("id"))

    if !c.policy.CanUpdate(user, post) {
        return ctx.Status(403).JSON(map[string]string{
            "error": "Cannot update this post",
        })
    }

    // Proceed...
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

func (g *AdminGuard) Handle(ctx types.Context, next types.Next) any {
    // Check authentication
    claims, err := g.authService.ValidateToken(ctx.Header("Authorization"))
    if err != nil {
        return ctx.Status(401).JSON(map[string]string{"error": "Unauthorized"})
    }

    // Check admin role
    if claims.Role != "admin" {
        return ctx.Status(403).JSON(map[string]string{"error": "Admin access required"})
    }

    ctx.Set("user_id", claims.UserID)
    ctx.Set("user_role", claims.Role)

    return next()
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
