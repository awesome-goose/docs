# CLI Installation

Install the Goose CLI to scaffold applications and generate code.

## Prerequisites

- **Go 1.21+** - Required for installation and running generated projects

Verify Go is installed:

```bash
go version
```

## Installation Methods

### Using Go Install (Recommended)

```bash
go install github.com/awesome-goose/cli@latest
```

This installs the `goose` binary to your `$GOPATH/bin` directory.

### Verify Installation

```bash
goose --version
```

Expected output:

```
goose version 1.0.0
```

### Building from Source

Clone and build:

```bash
# Clone the repository
git clone https://github.com/awesome-goose/cli.git
cd cli

# Build the binary
go build -o goose .

# Install globally
go install .
```

## PATH Configuration

Ensure `$GOPATH/bin` is in your PATH.

### macOS/Linux (bash)

Add to `~/.bashrc` or `~/.bash_profile`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Reload:

```bash
source ~/.bashrc
```

### macOS/Linux (zsh)

Add to `~/.zshrc`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Reload:

```bash
source ~/.zshrc
```

### Windows (PowerShell)

```powershell
$env:Path += ";$(go env GOPATH)\bin"
```

For permanent changes, add to your system environment variables.

## Troubleshooting

### Command Not Found

If `goose` is not found:

1. Check Go bin directory:

```bash
ls $(go env GOPATH)/bin
```

2. Verify PATH:

```bash
echo $PATH | grep "$(go env GOPATH)/bin"
```

3. Reinstall:

```bash
go install github.com/awesome-goose/cli@latest
```

### Permission Denied (Unix)

Make the binary executable:

```bash
chmod +x $(go env GOPATH)/bin/goose
```

### Old Version

Update to latest:

```bash
go install github.com/awesome-goose/cli@latest
```

Or specify a version:

```bash
go install github.com/awesome-goose/cli@v1.2.0
```

## Verifying Installation

Run the help command:

```bash
goose --help
```

Expected output:

```
Goose CLI - A command-line tool for Goose framework

Usage:
  goose <command> [flags]

Commands:
  app       Create a new application
  generate  Generate code (alias: g)

Flags:
  --version  Show version
  --help     Show help

Use "goose <command> --help" for more information about a command.
```

## Updating

To update to the latest version:

```bash
go install github.com/awesome-goose/cli@latest
```

## Uninstalling

Remove the binary:

```bash
rm $(go env GOPATH)/bin/goose
```

## Next Steps

- [Creating Applications](creating-apps.md) - Create your first app
- [Generating Modules](generating-modules.md) - Add modules
- [Command Reference](reference.md) - Full documentation
