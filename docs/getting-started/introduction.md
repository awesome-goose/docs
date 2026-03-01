# Introduction to Goose

## What is Goose?

Goose is a modern, modular Go framework designed for building scalable web applications, RESTful APIs, and command-line tools. It provides a unified development experience across all platforms while maintaining Go's simplicity and performance.

## Key Features

### 🚀 High Performance

Built on Go's fast runtime, Goose delivers exceptional performance for demanding applications. The framework adds minimal overhead while providing powerful abstractions.

### 📦 Modular Architecture

Everything in Goose is a module. This design promotes:

- **Code Reusability** - Share modules across projects
- **Maintainability** - Isolated, focused components
- **Scalability** - Add or remove features easily
- **Testability** - Test modules in isolation

### 🔧 Multi-Platform Support

Build for multiple platforms with consistent patterns:

- **API Platform** - RESTful APIs with JSON responses
- **Web Platform** - Server-rendered HTML applications
- **CLI Platform** - Command-line tools and utilities

### 📝 Declarative Design

Goose favors declarative over imperative code:

```go
// Define your module declaratively
type UserModule struct{}

func (m *UserModule) Declarations() []any {
    return []any{
        &UserController{},
        &UserService{},
    }
}

func (m *UserModule) Exports() []any {
    return []any{&UserService{}}
}
```

### 💉 Dependency Injection

Built-in IoC container for clean dependency management:

```go
type UserController struct {
    service *UserService `inject:""`
    log     types.Log    `inject:""`
}
```

### 🧪 Testing Support

Comprehensive testing utilities for all components:

```go
func TestUserController(t *testing.T) {
    container := goose.NewTestContainer()
    controller := container.Resolve(&UserController{})
    // Test your controller
}
```

## Who Should Use Goose?

Goose is ideal for:

- **Go developers** looking for a structured framework
- **Teams** building microservices or monolithic applications
- **Projects** requiring multiple interfaces (API + Web + CLI)
- **Anyone** who values clean, modular code

## Comparison with Other Frameworks

| Feature                      | Goose | Gin | Echo | Fiber |
| ---------------------------- | ----- | --- | ---- | ----- |
| Multi-Platform (API/Web/CLI) | ✅    | ❌  | ❌   | ❌    |
| Built-in DI Container        | ✅    | ❌  | ❌   | ❌    |
| Module System                | ✅    | ❌  | ❌   | ❌    |
| CLI Generator                | ✅    | ❌  | ❌   | ❌    |
| Job Queues                   | ✅    | ❌  | ❌   | ❌    |
| Cron Jobs                    | ✅    | ❌  | ❌   | ❌    |
| Caching                      | ✅    | ❌  | ❌   | ❌    |

## Philosophy

Goose follows these core principles:

1. **Structs are first-class citizens**
   - Define your application using Go structs
   - Use struct tags for configuration and behavior

2. **Favor declarative over imperative**
   - Describe what you want, not how to do it
   - Let the framework handle the wiring

3. **Go idioms first**
   - Use standard library conventions
   - Keep things simple and readable

4. **Modules over monoliths**
   - Break applications into focused modules
   - Each module handles one responsibility

## Getting Help

- **Documentation** - You're reading it!
- **GitHub Issues** - Report bugs and request features
- **Community** - Join our Discord/Slack (coming soon)

## Next Steps

Ready to start? Continue to:

- [Installation](installation.md) - Install Goose
- [Quick Start](quick-start.md) - Build your first app
