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
    web.WithTimeout(30),                // Request timeout (seconds)
    web.WithName("My Web App"),
    web.WithVersion("0.0.0"),
    web.WithAuthor("Your Name"),
    web.WithDescription("Web app description"),
)
```

### Available Options

```go
type Config struct {
    Name        string  // App name
    Version     string  // App version
    Author      string  // Author
    Description string  // Description
    Host        string  // Listen address
    Port        int     // Port
    Timeout     int     // Request timeout
}
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

Return data that your serializer/renderer will process:

```go
func (c *HomeController) Index(ctx types.Context) any {
    return map[string]interface{}{
        "_template": "pages/home.html",
        "Title":     "Home",
        "Message":   "Welcome to our site!",
    }
}
```

### With Layout

```go
func (c *HomeController) Index(ctx types.Context) any {
    return map[string]interface{}{
        "_template": "pages/home.html",
        "_layout":   "base/layout.html",
        "Title":     "Home",
        "User":      c.getCurrentUser(ctx),
    }
}
```

### With Data from Service

```go
type UserController struct {
    service *UserService `inject:""`
}

func (c *UserController) List(ctx types.Context) any {
    users := c.service.GetAll()

    return map[string]interface{}{
        "_template": "pages/users/list.html",
        "Title":     "Users",
        "Users":     users,
    }
}

func (c *UserController) Show(ctx types.Context) any {
    user := c.service.GetByID(ctx.Request().Params()["id"])
    if user == nil {
        return map[string]interface{}{
            "_template": "pages/404.html",
            "_status":   404,
        }
    }

    return map[string]interface{}{
        "_template": "pages/users/show.html",
        "Title":     user.Name,
        "User":      user,
    }
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
import (
    "encoding/json"
    "net/url"
)

type CreateUserForm struct {
    Email string
    Name  string
}

func (c *UserController) Create(ctx types.Context) any {
    body, err := ctx.Request().Body()
    if err != nil {
        return map[string]interface{}{
            "_template": "pages/users/create.html",
            "Errors":    map[string]string{"_form": "Failed to read form"},
        }
    }

    values, _ := url.ParseQuery(string(body))
    form := CreateUserForm{
        Email: values.Get("email"),
        Name:  values.Get("name"),
    }

    errors := c.validate(form)
    if len(errors) > 0 {
        return map[string]interface{}{
            "_template": "pages/users/create.html",
            "Form":      form,
            "Errors":    errors,
        }
    }

    c.service.Create(form)
    return map[string]interface{}{
        "_redirect": "/users",
    }
}
```

## Sessions

Session management is typically handled via cookies and context values:

```go
func (c *AuthController) Login(ctx types.Context) any {
    // Validate credentials...

    // Store user info in context for the request
    ctx.SetValue("user_id", user.ID)
    ctx.SetValue("user_name", user.Name)

    // For persistent sessions, set a cookie via raw response
    // or use a session middleware

    return map[string]interface{}{
        "_redirect": "/dashboard",
    }
}

func (c *AuthController) Logout(ctx types.Context) any {
    // Clear session cookie via raw response
    return map[string]interface{}{
        "_redirect": "/",
    }
}
```

### Session Middleware

```go
type SessionMiddleware struct{}

func (m *SessionMiddleware) Handle(ctx types.Context) error {
    // Read session from cookie/header
    // Store in context
    ctx.SetValue("current_user_id", userID)
    return nil
}
```

## Flash Messages

Flash messages can be implemented via session/cookie storage:

```go
func (c *UserController) Create(ctx types.Context) any {
    // Create user...

    // Store flash message (implement via session storage)
    return map[string]interface{}{
        "_redirect": "/users",
        "_flash":    map[string]string{"success": "User created successfully!"},
    }
}

func (c *UserController) List(ctx types.Context) any {
    // Retrieve flash from session
    flashMsg := "" // Read from session/cookie

    return map[string]interface{}{
        "_template": "pages/users/list.html",
        "Users":     c.service.GetAll(),
        "FlashMsg":  flashMsg,
    }
}
```

## Static Files

Static file serving is handled by your HTTP server configuration.
Goose does not currently have built-in `WithStatic` option.
Configure static file serving at the server level or via middleware.

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
return map[string]interface{}{
    "_redirect": "/users",
}

// With status code
return map[string]interface{}{
    "_redirect": "/login",
    "_status":   302,
}
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
