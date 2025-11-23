# Goose Service Container

The Goose service container is a robust, production-ready dependency injection system for Go applications. It manages the lifecycle of services and their dependencies, providing automatic resolution, injection, and lifecycle hooks for advanced use cases.

## Key Features

- **In-memory registry** for services and dependencies
- **Automatic dependency injection** for structs and functions
- **Singleton by default** (with extensibility for other scopes)
- **Lifecycle hooks** (`OnRegister`, `OnResolve`) for service events
- **Reflection caching** for high performance
- **Support for unexported fields** (using `unsafe`)
- **Flexible registration** for structs, factory functions, and interface-to-implementation bindings
- **Method invocation** with automatic argument resolution

## Principles

1. **Minimal Configuration**
   - No manual wiring required; dependencies are resolved automatically.
2. **Automatic Resolution**
   - Supports interfaces, structs, functions, and primitives.
3. **Global Accessibility**
   - Works across the entire framework, including controllers, middlewares, and more.

## Usage

### Service Registration

You can register services as either structs, factory functions, or bind interfaces to implementations. All services are singletons by default.

#### Struct Registration: Constructor Method Discovery

When registering a struct, the container will look for a constructor method named either `New` or `New<ServiceName>` (e.g., `NewLogger` for a `Logger` service). If such a method exists, it will be used as the factory for the service. Otherwise, dependencies are auto-injected into fields.

```go
// Register a struct (uses New or New<ServiceName> method if present, else auto-injects dependencies)
container.Register(Logger{}, "", true) // Will use Logger.NewLogger() if present
container.Register(MyService{}, "", true) // Will use MyService.New() or MyService.NewMyService() if present
```

#### Factory Function Registration

```go
// Register a factory function
container.Register(func() *MyService {
    return &MyService{ /* ... */ }
}, "", true)
```

#### Interface to Implementation Registration

```go
// Register an interface to a concrete implementation (factory returns interface)
container.Register(func() MyInterface {
    return &MyImplementation{}
}, "", true)
```

### Service Resolution

Resolve a service by providing a pointer to the desired type (interface or struct):

```go
myService := &MyService{}
err := container.Resolve(myService, "")

myInterface := new(MyInterface) // or var myInterface MyInterface = nil
err := container.Resolve(&myInterface, "")
```

If the service implements `OnResolve`, it will be called after resolution.

### Lifecycle Hooks

Services can implement the following optional hooks:

```go
func (s *MyService) OnRegister() {
    // Called after registration
}

func (s *MyService) OnResolve() {
    // Called after resolution
}
```

### Method Invocation

Invoke a method on a service, with automatic resolution of missing arguments:

```go
result, err := container.Call(myService, "DoSomething", arg1, arg2)
```

### Automatic Dependency Injection

If a struct has a `New` or `New<ServiceName>` method, it is used as a factory. Otherwise, all fields (including unexported) are auto-injected:

```go
// With New method
func (s *MyService) New(dep1 *Dep1, dep2 *Dep2) *MyService {
    return &MyService{dep1: dep1, dep2: dep2}
}

// With New<ServiceName> method
func (l *Logger) NewLogger(cfg *Config) *Logger {
    return &Logger{config: cfg}
}

// Without New/New<ServiceName> method
// Dependencies are injected automatically
```

### Reflection Caching

The container caches type lookups for performance, ensuring efficient resolution even in large applications.

### Unexported Field Injection

Unexported fields are injected using Go's `unsafe` package. This is powerful but should be used with caution.

## Error Handling

All errors are returned as custom Goose error types, providing clear diagnostics for invalid registration, resolution, or constructor issues.

## Extensibility

The container is designed for extensibility. You can add support for transient/request-scoped services, custom hooks, or additional features as needed.

## Example

```go
type MyService struct {
    dep1 *Dep1
    dep2 *Dep2
}

// Register interface to implementation
container.Register(func() MyInterface {
    return &MyImplementation{}
}, "", true)

myInterface := new(MyInterface)
err := container.Resolve(&myInterface, "")

func (s *MyService) OnRegister() {
    // Custom registration logic
}

func (s *MyService) OnResolve() {
    // Custom resolution logic
}
```

## Notes

- All services are singletons by default.
- When registering a struct, the container will look for both `New` and `New<ServiceName>` methods (e.g., `NewLogger` for a `Logger` service) and use them as constructors if present.
- Unexported field injection uses `unsafe` and may break in future Go versions.
- Lifecycle hooks are optional and provide advanced customization.
- Reflection caching is automatic and requires no configuration.
- Factory functions returning interfaces are registered under the interface type, enabling interface-to-implementation mapping.

## Summary

The Goose service container is a powerful, flexible, and production-ready solution for dependency management in Go. It provides automatic injection, lifecycle hooks, interface binding, and high performance for modern applications.
