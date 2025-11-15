# Goose

## Docs

### Getting Started

#### Introduction

Goose provides an exceptional and unified development experience across web, API, and CLI applications. Its modular architecture allows developers to build applications as a collection of independent modules, each encapsulating specific functionality. This modularity promotes code reusability, maintainability, and scalability.

Following are the key concepts and components involved in building a Goose application.

#### Boostrap a new application

To create a new Goose application, you can use the Goose CLI to generate a new application. This will create a new directory with the necessary files and structure for your application.

```bash
$ goose init --platform=<api | web | cli> --name=my-app
```

This command will create a new Goose application named `my-app` in the current directory. You can then navigate to the newly created directory and start building your application.

```bash
$ cd ./my-app
```

#### Run & Watch Changes

Goose supports live reloading of your application code, allowing you to see changes in real-time without having to restart your application. This feature is particularly useful during development, as it speeds up the development process and allows you to iterate quickly.
You can enable live reloading by using the `--watch` flag when running your application:

```bash
$ goose run # only run
$ goose run --watch # run with live reloading
```

This command will start your application and watch for changes in your code. When a change is detected, Goose will automatically reload your application, allowing you to see the changes immediately.

#### Add Modules

A goose application initially is a bare-bones application that includes the essential components needed to get started. It provides a solid foundation for building your application, allowing you to focus on adding features and functionality without worrying about the underlying structure.

Goose is a module-first platform, meaning that it is built around the concept of modules. Each module represents a specific functionality or feature of your application, allowing you to easily add, remove, or modify features as needed.

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
... included by default for all platforms
├── dtos/             # Configuration files for the module
├── controllers/      # Controllers for handling requests and responses
├── routes/           # Routing configuration for the module
├── services/         # Business logic and services for the module
├── middlewares/      # Middleware functions for the module
├── utils/            # Utility functions for the module
├── config/           # Configuration files for the module
├── tests/            # Test files for the module
├── README.md         # Documentation for the module

... only for web applications i.e `--platform=web`
├── assets/           # Static assets for the module, if applicable
├── views/            # View templates for the module, if applicable


... included if database operations would be involved using the `with=db` flag
├── migrations/       # Database migrations for the module, if it would be involved in database operations
├── seeds/            # Database seed files for the module, if it would be involved in
└── models/           # Database models for the module, if it would be involved in database operations
```

You can create a module with specific components using the `--with` flag. For example, to create a module with database migrations and seed files, you can use the following command:

```bash
$ goose module create --name=my-module --with=db
```

To better illustrate how to build a Goose module, below are explanations and examples of each component within a module.

##### DTOs

A Data Transfer Object (DTO) is used to transfer data between different layers of your application.
DTOs are typically used to encapsulate data that is sent to or received from the client,
and they help to decouple the internal representation of your data from the external representation.

```go
package dtos

import (
    "github.com/awesome-goose/validation"
)

type UserDTO struct {
    // declarative validation using tags
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// imperative validation using Validate method
func (dto *UserDTO) Validate() error {
    return nil
}
```

DTOs often include validation logic to ensure that the data being transferred meets certain criteria. This can be done using tags or by implementing a `Validate` method on the DTO struct.

See ([Validation](./building-blocks/validation.md)) for more information on how to validate DTOs in Goose.

##### Controller

A controller in Goose is responsible for handling incoming requests and returning responses. It acts as the intermediary between the client and the application logic. Controllers are typically defined in the `controllers` directory of your module.

```go
package controllers

import (
    "github.com/awesome-goose/controller"
    "github.com/awesome-goose/contracts"
)

type UserController struct {
    controller controller.Controller
}

func NewUserController(controller *controller.Controller) *UserController {
    return &UserController{
        controller: controller,
    }
}

func (c *UserController) Create(ctx controller.Context) contracts.Output {
    var userDTO dtos.UserDTO
    if err := ctx.Bind(&userDTO); err != nil {
        return errors.New("Failed to create user")
    }

    user := &models.User{
        Name:  userDTO.Name,
        Email: userDTO.Email,
    }

    if err := user.Save(); err != nil {
        return errors.New("Failed to create user")
    }

    controller.JSON(user)
}
```

#### Routing

Routing in Goose is handled through the `routes` directory of your module. You can define your routes and associate them with controllers to handle incoming requests.

```go
package app

import (
    "github.com/awesome-goose/contracts"
)

type Router struct {
    userController *UserController
}

func NewRouter(userController *UserController) *Router {
    return &Router{
        userController: userController,
    }
}

func (r *Router) Routes() (contracts.Routes, error) {
    return contracts.Routes{
        {
            Validators:  contracts.Validators{},
            Middlewares: contracts.Middlewares{},

            Path:  "users",
            Alias: []string{"people", "customers"},

            Children: contracts.Routes{
                {Path: "create", Handler: r.userController.Create},

                // Or
                {Path: "create", Handler: []any{r.userController, "Create"}},
            },
        },
    }, nil
}
```

#### Services

Services in Goose are used to encapsulate business logic and provide a clean interface for interacting with your application. They are typically defined in the `services` directory of your module.

Services are generic and can be used for various purposes, such as making request to other services, processing data, or performing complex operations.

```go
import (
    "github.com/awesome-goose/contracts"
)

type Service struct {
    contracts.Service
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) SomehtingComplec() (any, error) {

    return nil, nil
}
```

#### Middlewares

Middlewares in Goose are used to intercept requests and responses, allowing you to perform actions such as authentication, logging, or modifying the request/response. They are typically defined in the `middlewares` directory of your module.

```go
package middlewares

import (
    "github.com/awesome-goose/framework/middleware"
)

type AuthMiddleware struct {
    middleware.Middleware
}

func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{}
}

func (m *AuthMiddleware) Before(next middleware.Handler) error{
    if !isAuthenticated(ctx) {
        return errors.New("Unauthorized")
    }

    return nil
}


func (m *AuthMiddleware) After(next middleware.Handler) error {
    return nil
}

```

#### Utils

Utils in Goose are used to encapsulate utility functions that can be reused across your application. They are typically defined in the `utils` directory of your module.

It consists of:

- functions.go files
- values.go files

`functions.go` contains utility or helper functions used within the module.
``values.go` contains constant or variable values or configuration settings used within the module.

#### Config

Configuration in Goose is managed through the `config` directory of your module. You can define your configuration settings in a YAML file, which can be loaded and accessed throughout your application.

There is a default `app.yaml` file created in the `config` directory when you bootstrap a new Goose application. It contains basic configuration settings for your application, such as the application name, environment, port, and log level.

```yaml
# app.yaml
name: my-app
environment: development
port: 8080
log_level: info

database:
  host: localhost
  port: 5432
  username: user
  password: password
  dbname: mydb
```

Every module can have its own configuration file, allowing you to manage settings specific to that module. You can create additional configuration files as needed for your module.

Each module config is in the `config.yaml` at the root of the directory of the module.

```yaml
# my-module/config.yaml
sql:
  driver: { { envOr "SQL_DRIVER" "mysql" } }
  host: { { envOr "SQL_HOST" "localhost" } }
  port: { { envOr "SQL_PORT" 3306 } }
  database: { { envOr "SQL_DATABASE" "my_database" } }
  username: { { envOr "SQL_USERNAME" "root" } }
  password: { { envOr "SQL_PASSWORD" "secret" } }
```

Configs are hierachical, meaning that module configs overrides the application-level configs, when there's a conflict.

Configs can be accessed using the configuration service provided by Goose.

```go
package services

import (
    "github.com/awesome-goose/platform/config"
)

type MyService struct {
    configService *config.Service
}

func (s MyService) NewMyService(configService *config.Service) *MyService {
    return &MyService{
        configService: configService,
    }
}

func (s *MyService) GetDatabaseHost() string {
    return s.configService.GetString("sql.host")
}

```

See ([Configuration](./building-blocks/configuration.md)) for more information on how to manage and extend configuration in Goose.

##### Environment Variables

In addition to configs, Goose also supports environment variables, which can be used to override configuration settings. You can set environment variables in your system or in a `.env` file at the root of your application.

Environment variables take precedence over configuration settings defined in YAML files.

Environment variables can be accessed directly (though discouraged) as follows:

```go
package services

import (
    "github.com/awesome-goose/platform/env"
)

type MyService struct {
    envService *env.EnvService
}

func (s MyService) NewMyService(envService *env.EnvService) *MyService {
    return &MyService{
        envService: envService,
    }
}

func (s *MyService) GetDatabaseHost() string {
    return s.envService.GetString("SQL_HOST")
}
```

See ([Environment Variables](./building-blocks/environment-variables.md)) for more information on how to manage and use environment variables in Goose.

#### Logs

Goose provides a built-in logging service that allows you to log messages at different levels (debug, info, warn, error). The logging service is configurable, allowing you to customize the log format, output destination, and log level.

```go
package services

import (
    "github.com/awesome-goose/platform/logging"
)

type MyService struct {
    logger *logging.Logger
}

func (s MyService) NewMyService(logger *logging.Logger) *MyService {
    return &MyService{
        logger: logger,
    }
}

func (s *MyService) DoSomething() {
    s.logger.Info("Doing something...")
}
```

See ([Logging](./building-blocks/logging.md)) for more information on how to use and configure logging in Goose.

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

See ([Testing](./building-blocks/testing.md)) for more information on how to write and run tests in Goose.

<!--
addition docs:
- validation
- extending config, env, logs
- tests -->

### Deployment

#### Introduction

Goose can be deployed in various ways to suit your needs. Here are some common deployment methods:

1. **CLI**: Deploy Goose applications using the command line interface.
2. **Docker**: Use Docker to containerize and deploy Goose applications.
3. **Goose Cloud**: Deploy your applications to Goose Cloud, a managed service for Goose applications.

#### Run Your App Directly

To run your Goose application on your local machine, private server, use the following command:

```bash
$ goose run
```

This command will package your application and prepare it for deployment. You can then run the packaged application on any server that has Goose installed.

#### Docker Deployment

To deploy your Goose application using Docker, a dockerfile is provided by default. You can build and run your application in a Docker container with the following commands:

```bash
$ docker build -t my-app .
$ docker run -d -p 8080:8080 my-app
```

This will build your application into a Docker image and run it in a container, exposing port 8080.

#### Goose Cloud Deployment

To deploy your Goose application to Goose Cloud, you can use the following command:

```bash
$ goose deploy --platform=cloud --name=my-app
```
