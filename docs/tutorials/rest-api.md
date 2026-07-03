# Tutorial: Build a REST API

Learn to build a complete REST API with Goose from scratch.

## What We'll Build

A simple blog API with:

- User authentication
- CRUD for blog posts
- Comments on posts
- Request validation

## Prerequisites

- Go 1.21+
- Goose CLI installed

## Step 1: Create the Project

```bash
goose app --name=blog-api --template=api
cd blog-api
go mod tidy
```

## Step 2: Define Entities

Create `app/entities/user.go`:

```go
package entities

import "time"

type User struct {
    ID        string     `json:"id" gorm:"primaryKey"`
    Email     string     `json:"email" gorm:"uniqueIndex"`
    Password  string     `json:"-"`
    Name      string     `json:"name"`
    CreatedAt *time.Time `json:"created_at"`
}
```

Create `app/entities/post.go`:

```go
package entities

import "time"

type Post struct {
    ID        string     `json:"id" gorm:"primaryKey"`
    Title     string     `json:"title"`
    Content   string     `json:"content"`
    AuthorID  string     `json:"author_id"`
    Author    *User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
    Published bool       `json:"published" gorm:"default:false"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}
```

## Step 3: Create Auth Module

Generate the module:

```bash
goose g module --name=auth --type=plain
```

Update `app/auth/auth.service.go`:

```go
package auth

import (
    "blog-api/app/entities"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type AuthService struct {
    db        *gorm.DB `inject:""`
    jwtSecret string
}

func (s *AuthService) SetSecret(secret string) {
    s.jwtSecret = secret
}

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func (s *AuthService) Register(email, password, name string) (*entities.User, error) {
    // Check if user exists
    var existing entities.User
    if s.db.Where("email = ?", email).First(&existing).Error == nil {
        return nil, errors.New("email already registered")
    }

    // Hash password
    hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

    user := &entities.User{
        ID:       generateID(),
        Email:    email,
        Password: string(hashed),
        Name:     name,
    }

    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
    var user entities.User
    if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

Update `app/auth/auth.controller.go`:

```go
package auth

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type AuthController struct {
    service *AuthService `inject:""`
}

type RegisterDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

type LoginDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// The DTO is bound from the request body automatically.
func (c *AuthController) Register(dto *RegisterDTO) types.Output {
    user, err := c.service.Register(dto.Email, dto.Password, dto.Name)
    if err != nil {
        return output.BadRequest(err.Error())
    }
    return output.Created(user)
}

func (c *AuthController) Login(dto *LoginDTO) types.Output {
    token, err := c.service.Login(dto.Email, dto.Password)
    if err != nil {
        return output.Unauthorized(err.Error())
    }
    return output.JSON(map[string]string{"token": token})
}
```

Update `app/auth/auth.routes.go`:

```go
package auth

import "github.com/awesome-goose/goose/modules/router"

var ROUTES = router.ForRoutes(
    router.Post("/auth/register", []any{AuthController{}, "Register"}),
    router.Post("/auth/login", []any{AuthController{}, "Login"}),
)
```

Make sure the auth module imports these routes:

```go
func (m *AuthModule) Imports() []types.Module {
    return []types.Module{ROUTES}
}
```

## Step 4: Create Posts Module

Generate the module:

```bash
goose g module --name=posts --type=resource
```

Update `app/posts/posts.service.go`:

```go
package posts

import (
    "blog-api/app/entities"

    "gorm.io/gorm"
)

type PostsService struct {
    db *gorm.DB `inject:""`
}

func (s *PostsService) GetAll() []entities.Post {
    var posts []entities.Post
    s.db.Where("published = ?", true).Preload("Author").Find(&posts)
    return posts
}

func (s *PostsService) GetByID(id string) *entities.Post {
    var post entities.Post
    if err := s.db.Preload("Author").First(&post, "id = ?", id).Error; err != nil {
        return nil
    }
    return &post
}

func (s *PostsService) GetByAuthor(authorID string) []entities.Post {
    var posts []entities.Post
    s.db.Where("author_id = ?", authorID).Find(&posts)
    return posts
}

func (s *PostsService) Create(authorID string, dto CreatePostDTO) (*entities.Post, error) {
    post := &entities.Post{
        ID:       generateID(),
        Title:    dto.Title,
        Content:  dto.Content,
        AuthorID: authorID,
    }

    if err := s.db.Create(post).Error; err != nil {
        return nil, err
    }

    return post, nil
}

func (s *PostsService) Update(id string, dto UpdatePostDTO) (*entities.Post, error) {
    var post entities.Post
    if err := s.db.First(&post, "id = ?", id).Error; err != nil {
        return nil, err
    }

    post.Title = dto.Title
    post.Content = dto.Content
    post.Published = dto.Published

    if err := s.db.Save(&post).Error; err != nil {
        return nil, err
    }

    return &post, nil
}

func (s *PostsService) Delete(id string) error {
    return s.db.Delete(&entities.Post{}, "id = ?", id).Error
}
```

Update `app/posts/posts.controller.go`:

```go
package posts

import (
    "github.com/awesome-goose/goose/io/output"
    "github.com/awesome-goose/goose/types"
)

type PostsController struct {
    service *PostsService `inject:""`
}

// EmptyDto is used by handlers that take no request input.
type EmptyDto struct{}

type PostIDDto struct {
    ID string `param:"id"`
}

type CreatePostDTO struct {
    Title   string `json:"title" validate:"required"`
    Content string `json:"content" validate:"required"`
    // Populated by AuthMiddleware via ctx.SetValue("user_id", ...)
    AuthorID string `context:"user_id"`
}

type UpdatePostDTO struct {
    ID        string `param:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    Published bool   `json:"published"`
}

type MyPostsDto struct {
    AuthorID string `context:"user_id"`
}

func (c *PostsController) Index(dto *EmptyDto) types.Output {
    return output.JSON(c.service.GetAll())
}

func (c *PostsController) Show(dto *PostIDDto) types.Output {
    post := c.service.GetByID(dto.ID)
    if post == nil {
        return output.NotFound("Post not found")
    }
    return output.JSON(post)
}

func (c *PostsController) Create(dto *CreatePostDTO) types.Output {
    post, err := c.service.Create(dto.AuthorID, *dto)
    if err != nil {
        return output.InternalServerError(err.Error())
    }
    return output.Created(post)
}

func (c *PostsController) Update(dto *UpdatePostDTO) types.Output {
    post, err := c.service.Update(dto.ID, *dto)
    if err != nil {
        return output.InternalServerError(err.Error())
    }
    return output.JSON(post)
}

func (c *PostsController) Delete(dto *PostIDDto) types.Output {
    if err := c.service.Delete(dto.ID); err != nil {
        return output.InternalServerError(err.Error())
    }
    return output.NoContent()
}

func (c *PostsController) MyPosts(dto *MyPostsDto) types.Output {
    return output.JSON(c.service.GetByAuthor(dto.AuthorID))
}
```

> The `CreatePostDTO` and `UpdatePostDTO` service calls take the DTO by value, so update their signatures to `Create(authorID string, dto CreatePostDTO)` and `Update(id string, dto UpdatePostDTO)` as shown in Step 4 (they read only `Title`, `Content`, and `Published`).

## Step 5: Create Auth Middleware

Create `app/auth/auth.middleware.go`:

Middleware implements `Handle(ctx types.Context) error`. Returning an `error` aborts the request; returning `nil` continues to the handler. On success it stashes the authenticated user id with `ctx.SetValue`, which protected handlers read through a `context`-tagged DTO field.

```go
package auth

import (
    "errors"
    "strings"

    "github.com/awesome-goose/goose/types"
)

type AuthMiddleware struct {
    service *AuthService `inject:""`
}

func (m *AuthMiddleware) Handle(ctx types.Context) error {
    authHeader := ""
    if h := ctx.Request().Headers()["Authorization"]; len(h) > 0 {
        authHeader = h[0]
    }
    if authHeader == "" {
        return errors.New("missing authorization")
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    claims, err := m.service.ValidateToken(token)
    if err != nil {
        return errors.New("invalid token")
    }

    ctx.SetValue("user_id", claims.UserID)
    ctx.SetValue("user_email", claims.Email)

    return nil
}
```

## Step 6: Update Routes

Update `app/posts/posts.routes.go`:

```go
package posts

import (
    "blog-api/app/auth"

    "github.com/awesome-goose/goose/modules/router"
)

var ROUTES = router.ForRoutes(
    // Public routes
    router.Get("/posts", []any{PostsController{}, "Index"}),
    router.Get("/posts/:id", []any{PostsController{}, "Show"}),

    // Protected routes — AuthMiddleware runs before the handler
    router.Post("/posts", []any{PostsController{}, "Create"}, &auth.AuthMiddleware{}),
    router.Put("/posts/:id", []any{PostsController{}, "Update"}, &auth.AuthMiddleware{}),
    router.Delete("/posts/:id", []any{PostsController{}, "Delete"}, &auth.AuthMiddleware{}),
    router.Get("/my-posts", []any{PostsController{}, "MyPosts"}, &auth.AuthMiddleware{}),
)
```

And import it from the posts module:

```go
func (m *PostsModule) Imports() []types.Module {
    return []types.Module{ROUTES}
}
```

## Step 7: Update App Module

Update `app/app.module.go`:

```go
package app

import (
    "blog-api/app/auth"
    "blog-api/app/posts"

    "github.com/awesome-goose/goose/types"
)

type AppModule struct{}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        &auth.AuthModule{},
        &posts.PostsModule{},
    }
}

func (m *AppModule) Exports() []any {
    return []any{}
}

func (m *AppModule) Declarations() []any {
    return []any{
        &MigrationService{},
    }
}
```

Create `app/migrations.service.go`:

The `Boot` hook runs once at startup (it satisfies the `types.Bootable` interface):

```go
package app

import (
    "blog-api/app/entities"

    "github.com/awesome-goose/goose/types"
    "gorm.io/gorm"
)

type MigrationService struct {
    db *gorm.DB `inject:""`
}

func (s *MigrationService) Boot(k types.Kernel) error {
    return s.db.AutoMigrate(
        &entities.User{},
        &entities.Post{},
    )
}
```

## Step 8: Run the API

```bash
go run main.go
```

## Step 9: Test the API

### Register a user

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John"}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Create a post (with token)

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title":"My First Post","content":"Hello World!"}'
```

### Get all posts

```bash
curl http://localhost:8080/posts
```

## Summary

You've built a REST API with:

- User authentication (JWT)
- CRUD operations for posts
- Route protection with middleware
- Database migrations
- Input validation

## Next Steps

- Add comments feature
- Implement pagination
- Add search functionality
- Add API rate limiting
