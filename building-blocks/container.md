<!--
// Awesome Goose
// Container
// The core part of this container package is cloned from https://github.com/golobby/container
//
// Container is a dependency injection container that allows you to bind abstractions to concrete implementations.
// It supports both singleton and transient bindings, and allows you to resolve dependencies by calling functions or
// resolving interfaces.
//
// It is designed to be simple and easy to use, while providing powerful features for managing dependencies
// in your application. -->

<!-- // TODO:
// Features:
//
// Binding
// - bind
// - bindif
// - singleton
// - singletonIf
// - scoped (request scoped)
// - scopedIf
// - tagging/grouped binds (a set of interfaces bind & resolved together)
// - primitives (bind primitives like number, strings etc using token)
// - directly bind a struct type, without interface or name
//
// Resolve
// - resolve
// - resolve optional resolve with tag
// - directly resolve type
//
// Events
// - resolving
// - rebinding
// - call a method for an binding
// - core platform types should (e.g controllers, event listeners, middleware, and more) should embed container -->

# Goose Framework: Service Container Documentation

The Goose service container is a core component of the framework, responsible for managing the lifecycle of services and their dependencies. It provides a powerful, in-memory registry for resolving and injecting dependencies seamlessly across the application.

## Introduction

The service container:

- Maintains an **in-memory registry** of services.
- Acts as the **main entry point** for the framework.
- Is responsible for **creating, resolving, and managing the lifecycle** of providers and their dependencies.

## Principles

The Goose service container is designed with the following principles in mind:

1. **Minimal Configuration**:

   - Developers should not need to explicitly inject dependencies into their providers.
   - Dependencies are automatically resolved and injected by the container.

2. **Automatic Resolution**:

   - The container automatically resolves dependencies for services, including interfaces, structs, functions, and other primitives.

3. **Global Accessibility**:
   - The container is designed to work across the entire framework, enabling seamless integration with requests, responses, middlewares, controllers, and other Goose primitives.

## Features

### 1. **Service Registration**

The container supports multiple ways to register services:

- **Register a Service**:
  - The container can auto-resolve embedded dependencies or use a `New` method (if defined) as a factory method to create a new instance of the service.
  - If a `New` method is defined, embedded dependencies are not auto-resolved.

```go
// example of a service with a New method
type MyService struct {
    Dependency1 *Dependency1
    Dependency2 *Dependency2
}
func (s *MyService) New() *MyService {
    return &MyService{
        Dependency1: NewDependency1(),
        Dependency2: NewDependency2(),
    }
}
```

- **Register a Function Factory**:
  - A function can be registered as a factory for a service. The return value of the function is used as the service instance.

```go
container.register(func() *MyService {
    return &MyService{
        Dependency1: NewDependency1(),
        Dependency2: NewDependency2(),
    }
})
```

- **Register Interfaces to Implementations**:
  - Interfaces can be mapped to their concrete implementations, allowing for flexible dependency injection.

```go
container.registerInterface((*MyInterface)(nil), &MyImplementation{})
```

- **Request-Scoped Services**:
  - Services with parameters like `*request.Request` are automatically registered as request-scoped services, meaning a new instance is created for each request.

```go
type MyService struct {
    *request.Request
}

// OR
type MyService struct {
}
container.registerRequestScoped(&MyService{})
```

### 2. **Service Lifecycle Hooks**

Services can implement optional lifecycle hooks to handle specific events:

- **OnRegister**:

  - Called when the service is registered with the container.

- **OnUnregister**:

  - Called when the service is unregistered from the container.

- **OnResolve**:

  - Called when the service is resolve from the container.

- **HasNew**:
  - Not really a hook. If the service implements this interface, the 'New' it will be used as a factory method to create a new instance of the service. Embedded dependencies are not auto-resolved in this case.

### 3. **Singletons and Scopes**

- All registered services are **singletons** by default.
- Services can also be registered as **request-scoped**, ensuring a new instance is created for each request.
- Services can also be registered as **fresh**, ensuring a new instance is created for each resolution.

## Core Methods

The container exposes the following core methods:

### 1. **register()**

Registers a service, function factory, or interface-to-implementation mapping with the container.

Example:

```go
container.register(MyService{})
container.register(func() MyService {
    return NewMyService()
})
container.registerInterface((*MyInterface)(nil), MyImplementation{})
```

### 2. **resolve()**

Resolves a service or dependency from the container.

Example:

```go
service := container.resolve(MyService{})
```

### 3. **call()**

Invokes a method on a type, resolving any missing arguments from the container.

- Accepts a type, method name, and optional arguments.
- Any method argument not provided in the arguments list is resolved from the container.
- The method is invoked with the combined arguments, and the result is returned.

Example:

```go
result := container.call(MyService{}, "DoSomething", arg1, arg2)
```

## Global Capabilities

The Goose service container is designed to handle all Go primitives, including:

- Interfaces
- Structs
- Functions
- Other Goose primitives (e.g., requests, responses, middlewares, controllers)

This ensures the container can be used across the framework to manage dependencies seamlessly.

## Notes

- Services are **singletons** by default unless explicitly registered as request-scoped.
- Services with `*request.Request` or `*response.Response` parameters are automatically registered as **request-scoped**.
- The container ensures minimal configuration and automatic resolution of dependencies, reducing boilerplate code for developers.

## Summary

The Goose service container is a powerful and flexible tool for managing dependencies and services within the framework. By adhering to its principles of minimal configuration, automatic resolution, and global accessibility, developers can focus on building robust applications without worrying about dependency management.
