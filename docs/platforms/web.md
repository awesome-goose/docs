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

## Templates

### Template Structure

Templates live under `app/templates` by default:

```
app/templates/
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

Partials in the `partials/` directory are auto-loaded and can be referenced by name.

### Base Layout

```html
<!-- app/templates/base/layout.html -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}} - MyApp</title>
    <link rel="stylesheet" href="/static/css/app.css" />
  </head>
  <body>
    {{template "header" .}}

    <main>{{template "content" .}}</main>

    {{template "footer" .}}

    <script src="/static/js/app.js"></script>
  </body>
</html>
```

### Page Template

```html
<!-- app/templates/pages/home.html -->
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

Return `output.View(templatePath, data)`. The path is resolved under the templates directory and wrapped in the default layout (`base/layout.html`).

### Basic Rendering

```go
func (c *HomeController) Index(dto *EmptyDto) types.Output {
    return output.View("pages/home.html", map[string]any{
        "Title":   "Home",
        "Message": "Welcome to our site!",
    })
}
```

### Overriding the Layout

```go
func (c *HomeController) Index(dto *EmptyDto) types.Output {
    return output.View("pages/home.html", map[string]any{
        "Title": "Home",
        "User":  c.getCurrentUser(),
    }, output.WithLayout("base/admin.html"))
}
```

### With Data from a Service

```go
type UserController struct {
    service *UserService `inject:""`
}

type ShowUserDto struct {
    ID string `param:"id"`
}

func (c *UserController) List(dto *EmptyDto) types.Output {
    return output.View("pages/users/list.html", map[string]any{
        "Title": "Users",
        "Users": c.service.GetAll(),
    })
}

func (c *UserController) Show(dto *ShowUserDto) types.Output {
    user := c.service.GetByID(dto.ID)
    if user == nil {
        return output.View("pages/404.html", nil, output.WithHTMLCode(404))
    }

    return output.View("pages/users/show.html", map[string]any{
        "Title": user.Name,
        "User":  user,
    })
}
```

## Forms

### Form Markup

```html
<!-- app/templates/pages/users/create.html -->
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

### Processing a Form

Bind form fields with the `form` tag. On error, re-render the form; on success, redirect:

```go
type CreateUserForm struct {
    Email string `form:"email"`
    Name  string `form:"name"`
}

func (c *UserController) Create(dto *CreateUserForm) types.Output {
    if errs := c.validate(dto); len(errs) > 0 {
        return output.View("pages/users/create.html", map[string]any{
            "Form":   dto,
            "Errors": errs,
        }, output.WithHTMLCode(422))
    }

    c.service.Create(dto)
    return output.Redirect("/users").WithSuccess("User created successfully!")
}
```

## Sessions

Middleware works with `types.Context` and can stash per-request values:

```go
type SessionMiddleware struct{}

func (m *SessionMiddleware) Handle(ctx types.Context) error {
    // Read the session id from a cookie, look up the user, then:
    ctx.SetValue("current_user_id", userID)
    return nil
}
```

Handlers receive those values through a `context`-tagged DTO field. On login, set a cookie via a response header and redirect:

```go
type LoginDto struct {
    Email    string `form:"email"`
    Password string `form:"password"`
}

func (c *AuthController) Login(dto *LoginDto) types.Output {
    user, err := c.authService.Authenticate(dto.Email, dto.Password)
    if err != nil {
        return output.View("pages/login.html", map[string]any{
            "Error": "Invalid credentials",
        }, output.WithHTMLCode(401))
    }

    cookie := &http.Cookie{Name: "session", Value: user.SessionToken, Path: "/", HttpOnly: true}
    return output.Redirect("/dashboard",
        output.WithRedirectHeaders(map[string]string{"Set-Cookie": cookie.String()}),
    )
}

func (c *AuthController) Logout(dto *EmptyDto) types.Output {
    // Expire the session cookie
    cookie := &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1}
    return output.Redirect("/",
        output.WithRedirectHeaders(map[string]string{"Set-Cookie": cookie.String()}),
    )
}
```

## Flash Messages

Attach flash data to a redirect with the chainable helpers:

```go
func (c *UserController) Create(dto *CreateUserForm) types.Output {
    c.service.Create(dto)
    return output.Redirect("/users").
        WithSuccess("User created successfully!")
    // also: .WithError(...), .WithWarning(...), .WithInfo(...),
    //       .WithErrors(map[string][]string{...}), .WithOldInput()
}
```

Rendering flash on the next page requires session middleware to persist it across the redirect and expose it to the template.

## Static Files

Static file serving is handled at the HTTP-server or middleware level — Goose does not currently ship a `WithStatic` option.

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
// Simple redirect (302)
return output.Redirect("/users")

// Permanent redirect (301)
return output.RedirectPermanent("/new-home")

// Back to the previous page
return output.Back("/fallback")
```

See [Responses → Redirects](../building-blocks/responses.md#redirects) for the full set.

## Routes

```go
var ROUTES = router.ForRoutes(
    router.Get("/", []any{HomeController{}, "Index"}),
    router.Get("/about", []any{HomeController{}, "About"}),
    router.Get("/contact", []any{HomeController{}, "Contact"}),
    router.Post("/contact", []any{HomeController{}, "SubmitContact"}),
)
```

## Best Practices

1. **Use layouts** for consistent page structure
2. **Implement CSRF protection** for forms
3. **Use sessions securely** with proper secrets
4. **Escape output** in templates (Go's `html/template` does this by default)
5. **Handle form validation** gracefully by re-rendering with errors
6. **Use flash messages** for user feedback
7. **Organize templates** by feature

## Next Steps

- [Templates](../building-blocks/templates.md) - Template system
- [Sessions](../building-blocks/session.md) - Session management
- [CSRF Protection](../security/csrf.md) - Form security
