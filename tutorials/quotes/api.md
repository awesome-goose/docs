# Goose

## Docs

### Getting Started

### Quote of the Day

#### Introduction

We would build a simple "Quote of the Day" application using Goose. This application will demonstrate the basic concepts of Goose, including modules, migrations, seeds, models, controllers, and routing.

#### Prerequisites

Before we start, make sure you have the following prerequisites installed:

- Go (version 1.18 or later)
- Goose CLI (installed via `go install github.com/awesome-goose/cli@latest`)
- A database (we will use SQLite for this example, but you can use any database supported by Goose)

#### Step 1: Create a New Goose Application

First, we need to create a new Goose application. Open your terminal and run the following command:

```bash
$ goose init --platform=web --name=quote-of-the-day
```

This command will create a new directory called `quote-of-the-day` with the basic structure of a Goose application.

#### Bare Goose

A goose application initially is a bare-bones application that includes the essential components needed to get started. It provides a solid foundation for building your application, allowing you to focus on adding features and functionality without worrying about the underlying structure.

#### Modules

Goose is a module first platform, meaning that it is built around the concept of modules. Each module represents a specific functionality or feature of your application, allowing you to easily add, remove, or modify features as needed.

Critical components of application such as database, cache, queue, jobs are not included in the bare Goose application. Instead, you can add these components as modules, allowing you to customize your application based on your specific requirements. These components are provided as official modules.

Custom modules can also be created to encapsulate specific functionality or features, making it easy to share and reuse code across different applications. See ([Modules](./modules/index.md)) for more information on how to create and use modules in Goose.

##### Init Module

To create your first module, you can use the Goose CLI to generate a new module. This will create a new directory with the necessary files and structure for your module.

```bash
$ goose module create --name=my-module
```

This command will create a new module named `my-module` in the current directory. You can then add your code and functionality to this module.

##### Directory Structure

A Goose module has no regid structure, but it is recommended to follow a structure that makes sense for your module. A common structure includes the following directories:

```
my-module/
├── migrations/       # Database migrations for the module, if it would be involved in database operations
├── seeds/            # Database seed files for the module, if it would be involved in
├── models/           # Database models for the module, if it would be involved in database operations
├── dtos/             # Configuration files for the module
├── controllers/      # Controllers for handling requests and responses
├── routes/           # Routing configuration for the module
├── services/         # Business logic and services for the module
├── middlewares/      # Middleware functions for the module
├── utils/            # Utility functions for the module
├── config/           # Configuration files for the module
├── assets/           # Static assets for the module, if applicable
├── views/            # View templates for the module, if applicable
├── tests/            # Test files for the module
└── README.md         # (Optional) Documentation for the module
```

When you generate a module using the Goose CLI, it will create a basic structure for you, which you can then customize based on your needs. Depending on application platform (cli, web or api), the structure would vary slightly, but the core concepts remain the same.

##### Migrations

Migrations are used to manage changes to the database schema over time. They allow you to version control your database schema and apply changes in a controlled manner. Goose provides a migration system that allows you to create, apply, and rollback migrations easily.

```bash
$ goose migration create --name=create_users_table
```

This command will create a new migration file in the `migrations` directory with the specified name. You can then add your migration code to this file.

```go
package migrations

import (
    "github.com/awesome-goose/framework/migration"
)

func Up() error {
    // Code to create the users table
    return migration.CreateTable("users", func(table migration.Table) {
        table.String("name")
        table.String("email").Unique()
    })
}

func Down() error {
    // Code to drop the users table
    return migration.DropTable("users")
}
```

##### Seeds

Seeds are used to populate the database with initial data. They are typically used to create default records or test data that can be used during development or testing. Goose provides a seeding system that allows you to create and run seed files easily.

```bash
$ goose seed create --name=seed_users
```

This command will create a new seed file in the `seeds` directory with the specified name. You can then add your seed code to this file.

```go
package seeds
import (
    "github.com/awesome-goose/framework/seed"
)
func Up() error {
    // Code to seed the users table with initial data
    return seed.Create("users", func(s seed.Seeder) {
        s.Create("users", map[string]interface{}{
            "name":  "John Doe",
            "email": "
```

##### Model

A model in Goose represents a database entity. It is used to define the structure of your data and how it interacts with the database. Models are typically defined in the `models` directory of your module.

```go
package models
import (
    "github.com/awesome-goose/framework/model"
)
type User struct {
    model.Model
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

##### DTOs

A Data Transfer Object (DTO) is used to transfer data between different layers of your application.
DTOs are typically used to encapsulate data that is sent to or received from the client,
and they help to decouple the internal representation of your data from the external representation.

```go
package dtos
type UserDTO struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

##### Controller

A controller in Goose is responsible for handling incoming requests and returning responses. It acts as the intermediary between the client and the application logic. Controllers are typically defined in the `controllers` directory of your module.

```go
package controllers
import (
    "github.com/awesome-goose/framework/controller"
    "github.com/awesome-goose/framework/dtos"
    "github.com/awesome-goose/framework/models"
)
type UserController struct {
    controller.BaseController
}
func (c *UserController) Create(ctx controller.Context) {
    var userDTO dtos.UserDTO
    if err := ctx.Bind(&userDTO); err != nil {
        ctx.JSON(400, "Invalid request")
        return
    }

    user := models.User{
        Name:  userDTO.Name,
        Email: userDTO.Email,
    }

    if err := user.Save(); err != nil {
        ctx.JSON(500, "Failed to create user")
        return
    }

    ctx.JSON(201, user)
}
```

#### Routing

Routing in Goose is handled through the `routes` directory of your module. You can define your routes and associate them with controllers to handle incoming requests.

```go
package routes
import (
    "github.com/awesome-goose/framework/router"
    "github.com/awesome-goose/framework/controllers"
)
func RegisterRoutes(r *router.Router) {
    userController := &controllers.UserController{}
    r.POST("/users", userController.Create)
}
```

#### Services

Services in Goose are used to encapsulate business logic and provide a clean interface for interacting with your application. They are typically defined in the `services` directory of your module.

```go
package services
import (
    "github.com/awesome-goose/framework/models"
)
type UserService struct {}
func (s *UserService) CreateUser(name, email string) (*models.User, error) {
    user := &models.User{
        Name:  name,
        Email: email,
    }
    if err := user.Save(); err != nil {
        return nil, err
    }
    return user, nil
}
```

#### Middlewares

Middlewares in Goose are used to intercept requests and responses, allowing you to perform actions such as authentication, logging, or modifying the request/response. They are typically defined in the `middlewares` directory of your module.

```go
package middlewares
import (
    "github.com/awesome-goose/framework/middleware"
)
func AuthMiddleware(next middleware.Handler) middleware.Handler {
    return func(ctx middleware.Context) {
        // Perform authentication logic
        if !isAuthenticated(ctx) {
            ctx.JSON(401, "Unauthorized")
            return
        }
        next(ctx)
    }
}
```

#### Utils

Utils in Goose are used to encapsulate utility functions that can be reused across your application. They are typically defined in the `utils` directory of your module.

```go
package utils
import "strings"
func ToLowerCase(s string) string {
    return strings.ToLower(s)
}
```

#### Config

Configuration in Goose is managed through the `config` directory of your module. You can define your configuration settings in a YAML file, which can be loaded and accessed throughout your application.

```yaml
# config/app.yaml
name: my-app
environment: development
port: 8080
log_level: info
```

#### [Web] Assets

Assets in Goose are used to serve static files such as images, CSS, and JavaScript. They are typically defined in the `assets` directory of your module.

```go
package assets
import (
    "github.com/awesome-goose/framework/assets"
)
func RegisterAssets() {
    assets.Register("/static", "assets/static")
}
```

#### [Web] Views

Views in Goose are used to render HTML templates for web applications. They are typically defined in the `views` directory of your module.

```go
package views
import (
    "github.com/awesome-goose/framework/view"
)
func RegisterViews() {
    view.Register("templates", "views/templates")
}
```

#### Tests

Tests in Goose are used to ensure the correctness of your application. They are typically defined in the `tests` directory of your module. You can use the built-in testing framework or any other testing framework of your choice.

```go
package tests
import (
    "testing"
    "github.com/awesome-goose/framework/models"
)
func TestUserModel(t *testing.T) {
    user := &models.User{
        Name:  "John Doe",
        Email: "
```
