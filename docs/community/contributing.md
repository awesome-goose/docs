# Contributing to Goose

Thank you for your interest in contributing to Goose!

## Ways to Contribute

- **Report bugs** - File issues for bugs you encounter
- **Suggest features** - Propose new functionality
- **Submit PRs** - Fix bugs or add features
- **Improve docs** - Fix typos, add examples
- **Help others** - Answer questions in discussions

## Getting Started

### Prerequisites

- Go 1.21+
- Git
- A GitHub account

### Fork and Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/goose.git
cd goose
git remote add upstream https://github.com/awesome-goose/goose.git
```

### Set Up Development Environment

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
go vet ./...
```

## Development Workflow

### Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/my-new-feature
```

### Make Changes

1. Write your code
2. Add/update tests
3. Update documentation if needed
4. Run tests and linter

### Commit Guidelines

Use conventional commits:

```
feat: add new validation rules
fix: resolve routing conflict
docs: update installation guide
test: add tests for auth middleware
refactor: simplify container logic
```

### Submit Pull Request

```bash
# Push your branch
git push origin feature/my-new-feature
```

Then create a PR on GitHub with:

- Clear title and description
- Reference any related issues
- Screenshots for UI changes

## Code Style

### Go Style

Follow standard Go conventions:

```go
// Good: clear naming
func (s *UserService) GetByEmail(email string) (*User, error) {
    // ...
}

// Good: error handling
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Good: concise comments
// GetByEmail retrieves a user by their email address.
```

### Formatting

```bash
# Format code
go fmt ./...

# Check formatting
gofmt -d .
```

### Linting

```bash
# Run go vet
go vet ./...

# Install and run golangci-lint
golangci-lint run
```

## Testing

### Write Tests

```go
func TestFeature(t *testing.T) {
    // Arrange
    input := "test"

    // Act
    result := Feature(input)

    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific tests
go test -run TestFeature ./...
```

## Documentation

### Code Documentation

Add comments for exported functions:

```go
// NewModule creates a new cache module with the given options.
// It returns an error if the configuration is invalid.
func NewModule(opts ...Option) (*Module, error) {
    // ...
}
```

### Update Docs

For documentation changes:

1. Edit relevant `.md` files in `/docs`
2. Preview changes locally
3. Submit PR

## Issue Guidelines

### Bug Reports

Include:

- Go version
- Goose version
- Steps to reproduce
- Expected vs actual behavior
- Error messages/logs

### Feature Requests

Include:

- Use case description
- Proposed solution
- Alternatives considered

## Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** thoroughly
5. **Submit** PR
6. **Address** review feedback
7. **Get merged** 🎉

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help maintain a welcoming community

## Questions?

- Open a [Discussion](https://github.com/awesome-goose/goose/discussions)
- Check existing issues
- Read the documentation

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.

Thank you for contributing! 🙏
