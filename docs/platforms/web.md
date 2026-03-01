# Web Platform

Build server-rendered web applications with the Goose Web platform.

## Overview

The Web platform serves HTML pages with template rendering support.

## Quick Start

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/web"
)

func main() {
    platform := web.NewPlatform(
        web.WithHost("localhost"),
        web.WithPort(3000),
        web.WithTemplates("./templates"),
    )

    module := &app.AppModule{}

    stop, err := goose.Start(goose.Web(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## Configuration

### Platform Options

```go
platform := web.NewPlatform(
    web.WithHost("0.0.0.0"),
    web.WithPort(3000),
    web.WithTemplates("./templates"),
    web.WithStatic("./public", "/static"),
    web.WithSessionSecret("secret-key"),
)
```

## Templates

### Template Structure

```
templates/
├── base/
│   └── layout.html
├── pages/
│   ├── home.html
│   ├── about.html
│   └── users/
│       ├── list.html
│       └── show.html
└── partials/
    ├── header.html
    ├── footer.html
    └── nav.html
```

### Base Layout

```html
<!-- templates/base/layout.html -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}} - MyApp</title>
    <link rel="stylesheet" href="/static/css/app.css" />
  </head>
  <body>
    {{template "partials/header.html" .}}

    <main>{{template "content" .}}</main>

    {{template "partials/footer.html" .}}

    <script src="/static/js/app.js"></script>
  </body>
</html>
```

### Page Template

```html
<!-- templates/pages/home.html -->
{{define "content"}}
<div class="container">
  <h1>Welcome, {{.User.Name}}</h1>

  <ul>
    {{range .Items}}
    <li>{{.Name}} - {{.Price}}</li>
    {{end}}
  </ul>
</div>
{{end}}
```

## Rendering Views

### Basic Rendering

```go
func (c *HomeController) Index(ctx types.Context) any {
    return ctx.View("pages/home.html", map[string]interface{}{
        "Title":   "Home",
        "Message": "Welcome to our site!",
    })
}
```

### With Layout

```go
func (c *HomeController) Index(ctx types.Context) any {
    return ctx.View("pages/home.html", map[string]interface{}{
        "Title": "Home",
        "User":  c.getCurrentUser(ctx),
    }).Layout("base/layout.html")
}
```

### With Data from Service

```go
type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) List(ctx types.Context) any {
    users := c.service.GetAll()

    return ctx.View("pages/users/list.html", map[string]interface{}{
        "Title": "Users",
        "Users": users,
    })
}

func (c *UserController) Show(ctx types.Context) any {
    user := c.service.GetByID(ctx.Param("id"))
    if user == nil {
        return ctx.Status(404).View("pages/404.html", nil)
    }

    return ctx.View("pages/users/show.html", map[string]interface{}{
        "Title": user.Name,
        "User":  user,
    })
}
```

## Forms

### Form Handling

```html
<!-- templates/pages/users/create.html -->
{{define "content"}}
<form method="POST" action="/users">
  <input type="hidden" name="_csrf" value="{{.CSRFToken}}" />

  <label for="email">Email:</label>
  <input type="email" name="email" id="email" value="{{.Form.Email}}" />
  {{if .Errors.Email}}<span class="error">{{.Errors.Email}}</span>{{end}}

  <label for="name">Name:</label>
  <input type="text" name="name" id="name" value="{{.Form.Name}}" />
  {{if .Errors.Name}}<span class="error">{{.Errors.Name}}</span>{{end}}

  <button type="submit">Create User</button>
</form>
{{end}}
```

### Processing Form

```go
type CreateUserForm struct {
    Email string `form:"email"`
    Name  string `form:"name"`
}

func (c *UserController) Create(ctx types.Context) any {
    var form CreateUserForm
    ctx.Bind(&form)

    errors := c.validate(form)
    if len(errors) > 0 {
        return ctx.View("pages/users/create.html", map[string]interface{}{
            "Form":   form,
            "Errors": errors,
        })
    }

    c.service.Create(form)
    return ctx.Redirect("/users")
}
```

## Sessions

### Session Management

```go
func (c *AuthController) Login(ctx types.Context) any {
    // Validate credentials...

    // Set session
    ctx.Session().Set("user_id", user.ID)
    ctx.Session().Set("user_name", user.Name)

    return ctx.Redirect("/dashboard")
}

func (c *AuthController) Logout(ctx types.Context) any {
    ctx.Session().Clear()
    return ctx.Redirect("/")
}
```

### Session Middleware

```go
type SessionMiddleware struct{}

func (m *SessionMiddleware) Handle(ctx types.Context, next types.Next) any {
    userID := ctx.Session().Get("user_id")
    if userID != nil {
        ctx.Set("current_user_id", userID)
    }
    return next()
}
```

## Flash Messages

```go
func (c *UserController) Create(ctx types.Context) any {
    // Create user...

    ctx.Flash("success", "User created successfully!")
    return ctx.Redirect("/users")
}

func (c *UserController) List(ctx types.Context) any {
    return ctx.View("pages/users/list.html", map[string]interface{}{
        "Users":    c.service.GetAll(),
        "FlashMsg": ctx.Flash("success"),
    })
}
```

## Static Files

### Configuration

```go
platform := web.NewPlatform(
    web.WithStatic("./public", "/static"),
)
```

### Directory Structure

```
public/
├── css/
│   └── app.css
├── js/
│   └── app.js
└── images/
    └── logo.png
```

### Usage in Templates

```html
<link rel="stylesheet" href="/static/css/app.css" />
<script src="/static/js/app.js"></script>
<img src="/static/images/logo.png" alt="Logo" />
```

## Redirects

```go
// Simple redirect
return ctx.Redirect("/users")

// With status code
return ctx.Redirect("/login", 302)

// Back to previous page
return ctx.Back()
```

## Routes

```go
func (c *HomeController) Routes() types.Routes {
    return types.Routes{
        {Method: "GET", Path: "/", Handler: c.Index},
        {Method: "GET", Path: "/about", Handler: c.About},
        {Method: "GET", Path: "/contact", Handler: c.Contact},
        {Method: "POST", Path: "/contact", Handler: c.SubmitContact},
    }
}
```

## Best Practices

1. **Use layouts** for consistent page structure
2. **Implement CSRF protection** for forms
3. **Use sessions securely** with proper secrets
4. **Escape output** in templates
5. **Handle form validation** gracefully
6. **Use flash messages** for user feedback
7. **Organize templates** by feature

## Next Steps

- [Templates](../building-blocks/templates.md) - Template system
- [Sessions](../building-blocks/session.md) - Session management
- [CSRF Protection](../security/csrf.md) - Form security
