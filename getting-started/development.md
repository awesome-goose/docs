# Goose

## Docs

### Getting Started

### First App

`See (My First App - Quote of the Day)[./quotes/index.md] for a tutorial style introduction to Goose application.`

#### Introduction

Goose provides an exceptional and unified development experience across web, API, and CLI applications.

#### Watch Changes

Goose supports live reloading of your application code, allowing you to see changes in real-time without having to restart your application. This feature is particularly useful during development, as it speeds up the development process and allows you to iterate quickly.
You can enable live reloading by using the `--watch` flag when running your application:

```bash
$ goose run --watch
```

This command will start your application and watch for changes in your code. When a change is detected, Goose will automatically reload your application, allowing you to see the changes immediately.

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

The command would also register the module in the Goose application, making it available for use. You can then start adding your models, controllers, routes, and other components to this module.

##### Directory Structure

A Goose module has no rigid structure, but it is recommended to follow a structure that makes sense for your module. A common structure includes the following directories:

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
$ goose migration create --name=create_users_table --module=my-module
```

This command will create a new migration file in the `migrations` directory with the specified name. You can then add your migration code to this file.

```go
package migrations

import (
    "github.com/awesome-goose/framework/migration"
)

type CreateUsersTable struct{}

func (cut CreateUsersTable) Up() error {
    // Code to create the users table
    return migration.CreateTable("users", func(table migration.Table) {
        table.String("name")
        table.String("email").Unique()
    })
}

func (cut CreateUsersTable) Down() error {
    // Code to drop the users table
    return migration.DropTable("users")
}
```

You can then run the migration using the Goose CLI:

```bash
$ goose migration up --name=create_users_table --module=my-module
```

To revert the migration, you can use the following command:

```bash
$ goose migration down  --name=create_users_table --module=my-module
```

`--name` is optional, if not provided, Goose will run all pending migrations.

##### Seeds

Seeds are used to populate the database with initial data. They are typically used to create default records or test data that can be used during development or testing. Goose provides a seeding system that allows you to create and run seed files easily.

```bash
$ goose seed create --name=seed_users --module=my-module
```

This command will create a new seed file in the `seeds` directory with the specified name. You can then add your seed code to this file.

```go
package seeds

import (
    "github.com/awesome-goose/framework/seed"
)

type SeedUsers struct{}

func (su *SeedUsers) Up() error {
    // Code to seed the users table with initial data
    return seed.Create("users", func(s seed.Seeder) {
        s.Insert("name", "John Doe")
        s.Insert("email", "john.doe@goose.com")
    })
}

func (su *SeedUsers) Down() error {
    // Code to remove the seeded data
    return seed.Truncate("users")
}
```

You can then run the seed using the Goose CLI:

```bash
$ goose seed up --name=seed_users --module=my-module
```

To revert the seed, you can use the following command:

```bash
$ goose seed down --name=seed_users --module=my-module
```

`--name` is optional, if not provided, Goose will run all pending seeds.

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


func (u *User) TableName() string {
    return "users" // Specify the table name for the model
}

// in services.controller, you can use the model like this
models.User{
    Name:  "John Doe",
    Email: "i@i.com",
}

models.User.New(&User{})
models.User.Save()
models.User.Save()
```

##### DTOs

A Data Transfer Object (DTO) is used to transfer data between different layers of your application.
DTOs are typically used to encapsulate data that is sent to or received from the client,
and they help to decouple the internal representation of your data from the external representation.

```go
package dtos

import (
    "github.com/awesome-goose/framework/validation"
)

type UserDTO struct {
    Name  string `json:"name" validate:"required"` //declarative validation
    Email string `json:"email" validate:"required,email"`
}

// imperative validation
func (dto *UserDTO) Validate() error {
    // Implement validation logic here
    return nil // Return nil if validation passes
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
    "github.com/awesome-goose/framework/testing"
)
```

To run all tests in your module, you can use the following command:

```bash
$ goose test --module=my-module
```

This will run all the tests defined in the `tests` directory of your module.

` --module` is optional, if not provided, Goose will run all tests in the current application.
