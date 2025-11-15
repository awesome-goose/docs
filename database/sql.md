# Goose

## Docs

### Getting Started

### Working with SQL Database

#### Introduction

Goose provides a robust SQL database management system that allows you to easily handle SQL database operations such as migrations, seeding, and model management. This section will guide you through the process of working with SQL databases in Goose applications.

Not all apps require a SQL database. Goose is a modular platform, so SQL database support is provided through official modules that you can add to your application as needed. This allows you to keep your application lightweight and only include the components that are necessary for your specific use case.

#### Import SQL Database Module

To add SQL database support to your Goose application, you need to import the official SQL database module. You can do this by using the Goose CLI to add the module to your application.

```bash
$ goose module add sql # add to the root module
$ goose module add sql --module=my-module # add to a specific module
```

This command will download and install the SQL database module, and it will also update your application's configuration to include the SQL database settings.

The above command would add the official SQL database module from `github.com/awesome-goose/sql` to your application. You can then start using the SQL database features provided by the module.

To import the SQL database module manually without using the CLI, you can add the following code to your application's main module file or to any specific module:

```go
import (
    "github.com/awesome-goose/sql"
)

type MyModule struct {}

func (m *MyModule) NewMyModule() *MyModule {
    return &MyModule{}
}

func (m *MyModule) Imports() {
    return []contracts.Module{
        ...
        database.Module{},
        ...
    }
}
```

#### Configure SQL Database

To configure the SQL database connection, you can update the configuration file located at `config/sql.yaml` of the root module or any specific module where the SQL database module is imported.

```yaml
sql:
  driver: { { envOr "SQL_DRIVER" "mysql" } }
  host: { { envOr "SQL_HOST" "localhost" } }
  port: { { envOr "SQL_PORT" 3306 } }
  database: { { envOr "SQL_DATABASE" "my_database" } }
  username: { { envOr "SQL_USERNAME" "root" } }
  password: { { envOr "SQL_PASSWORD" "secret" } }
```

Or set the relevant environment variables in the root `.env` file:

```
SQL_CONNECTION=mysql
SQL_HOST=localhost
SQL_PORT=3306
SQL_DATABASE=my_database
SQL_USERNAME=root
SQL_PASSWORD=secret
```

##### Add Migrations

Migrations are used to manage changes to the database schema over time. They allow you to version control your database schema and apply changes in a controlled manner. Goose provides a migration system that allows you to create, apply, and rollback migrations easily.

```bash
$ goose migration create --name=create_users_table --module=my-module
```

This command will create a new migration file in the `migrations` directory with the specified name. You can then add your migration code to this file.

```go
package migrations

import (
    "github.com/awesome-goose/sql/migration"
)

type CreateUsersTable struct{
    migration *migration.Migration
}

func (cut CreateUsersTable) NewCreateUsersTable(migration *migration.Migration) *CreateUsersTable {
    return &CreateUsersTable{
        migration: migration,
    }
}

func (cut CreateUsersTable) Up() error {
    return migration.CreateTable("users", func(table migration.Table) {
        table.String("name")
        table.String("email").Unique()
    })
}

func (cut CreateUsersTable) Down() error {
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

Seeds are used to populate the database with initial data. They are typically used to create default records or test data that can be used during development or testing. Similar to migrations, Goose provides a seeding system that allows you to create and run seed files easily.

```bash
$ goose seed create --name=seed_users --module=my-module
```

This command will create a new seed file in the `seeds` directory with the specified name. You can then add your seed code to this file.

```go
package seed

import (
    "github.com/awesome-goose/sql/seed"
)

type CreateUsers struct{
    seed *seed.Seed
}

func (cut CreateUsers) NewCreateUsers(seed *seed.Seed) *CreateUsers {
    return &CreateUsers{
        seed: seed,
    }
}

func (cut CreateUsers) Up() error {
    return seed.Create("users", func(s seed.Seeder) {
        s.Insert("name", "John Doe")
        s.Insert("email", "john.doe@goose.com")
    })
}

func (cut CreateUsers) Down() error {
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

Models should not be confused with DTOs (Data Transfer Objects). Models are used to interact with the database, while DTOs are used to transfer data between different layers of your application.

```go
package models

import (
    "github.com/awesome-goose/sql/model"
)

type User struct {
    model model.Model

    Name  string `json:"name"`
    Email string `json:"email"`
}

func NewUser() *User {
    return &User{}
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
    "github.com/awesome-goose/sql/validation"
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
    "github.com/awesome-goose/sql/controller"
    "github.com/awesome-goose/sql/dtos"
    "github.com/awesome-goose/sql/models"
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
    "github.com/awesome-goose/sql/router"
    "github.com/awesome-goose/sql/controllers"
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
    "github.com/awesome-goose/sql/models"
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
    "github.com/awesome-goose/sql/middleware"
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
    "github.com/awesome-goose/sql/assets"
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
    "github.com/awesome-goose/sql/view"
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
    "github.com/awesome-goose/sql/testing"
)
```

To run all tests in your module, you can use the following command:

```bash
$ goose test --module=my-module
```

This will run all the tests defined in the `tests` directory of your module.

` --module` is optional, if not provided, Goose will run all tests in the current application.
