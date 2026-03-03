# Installation

This guide covers installing Goose and its prerequisites.

## Prerequisites

### Go Version

Goose requires **Go 1.21** or later.

Check your Go version:

```bash
go version
```

If you need to install or upgrade Go, visit [golang.org/dl](https://golang.org/dl/).

### Database (Optional)

For database features, you'll need one of:

- **SQLite** - No additional setup required
- **PostgreSQL** - Version 12+
- **MySQL** - Version 8+

## Installing the Goose CLI

The Goose CLI helps you scaffold applications and generate modules.

### Using Go Install (Recommended)

```bash
go install github.com/awesome-goose/cli@latest
```

### Verify Installation

```bash
goose --version
```

You should see output like:

```
goose version 0.0.0
```

### Build from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/awesome-goose/cli.git

# Navigate to the directory
cd cli

# Build and install
go install .
```

## Installing the Goose Framework

The framework is installed automatically when you create a new project. For manual installation in an existing project:

```bash
go get github.com/awesome-goose/goose
```

## Project Setup

### Create a New Project

Using the CLI (recommended):

```bash
goose app --name=myapp --template=api
```

### Manual Setup

1. Create project directory:

```bash
mkdir myapp && cd myapp
```

2. Initialize Go module:

```bash
go mod init myapp
```

3. Install Goose:

```bash
go get github.com/awesome-goose/goose
```

4. Create `main.go`:

```go
package main

import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/platforms/api"
)

func main() {
    platform := api.NewPlatform()
    module := &app.AppModule{}

    stop, err := goose.Start(goose.API(platform, module, nil))
    if err != nil {
        panic(err)
    }
    defer stop()
}
```

## IDE Setup

### VS Code

Install the Go extension for the best experience:

1. Open VS Code
2. Go to Extensions (Cmd/Ctrl + Shift + X)
3. Search for "Go" by the Go Team
4. Click Install

Recommended settings for `.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "gofmt",
  "editor.formatOnSave": true
}
```

### GoLand

GoLand provides excellent Go support out of the box. No additional setup required.

## Troubleshooting

### Command Not Found

If `goose` command is not found, ensure your Go bin directory is in your PATH:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

Then reload your shell:

```bash
source ~/.bashrc  # or ~/.zshrc
```

### Module Not Found

If you get module errors, try:

```bash
go mod tidy
```

### Permission Denied

On Unix systems, you may need to make the binary executable:

```bash
chmod +x $(go env GOPATH)/bin/goose
```

## Next Steps

Now that Goose is installed:

- [Quick Start](quick-start.md) - Create your first application
- [Configuration](configuration.md) - Learn about configuration options
