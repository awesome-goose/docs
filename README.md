<p align="center">
  <strong>🪿 Goose Framework</strong><br>
  <sub>Modular • Scalable • Multi-Platform</sub>
</p>

---

# Goose Documentation Builder

A CLI tool built with the Goose framework to build beautiful, developer-friendly documentation websites from markdown files.

---

## Features

- **Beautiful Design** - Clean, modern interface with easy navigation
- **Dark Mode** - Toggle between light and dark themes
- **Fast Search** - Full-text search with keyboard shortcuts (⌘K)
- **Syntax Highlighting** - Code blocks with Dracula theme
- **Responsive** - Mobile-friendly design
- **GitHub Pages Ready** - Build static sites for easy hosting
- **Table of Contents** - Auto-generated TOC for each page
- **Copy Code** - One-click code block copying
- **Goose Framework** - Built using the Goose CLI platform pattern

---

## Project Structure

```
docs/
├── main.go              # Goose CLI entry point
├── app/
│   ├── app.module.go    # Module definition
│   ├── app.controller.go # CLI command handlers
│   ├── app.routes.go    # CLI route definitions
│   ├── app.dtos.go      # Command DTOs with flag tags
│   └── app.service.go   # Business logic service
├── builder/
│   ├── builder.go       # Core build logic
│   ├── templates.go     # HTML templates
│   ├── styles.go        # CSS styles
│   └── scripts.go       # JavaScript
└── docs/                # Markdown documentation files
```

---

## Installation

```bash
cd docs
go mod tidy
go build -o goose-docs .
```

---

## Usage

### Build Documentation

```bash
# Build with default options
./goose-docs build --input docs --output dist

# Build with custom options
./goose-docs build --input ./docs --output ./public --title "My Docs" --base-url /my-project/
```

### Options

| Option       | Default               | Description                          |
| ------------ | --------------------- | ------------------------------------ |
| `--input`    | `docs`                | Input directory containing .md files |
| `--output`   | `dist`                | Output directory for static site     |
| `--title`    | `Goose Documentation` | Site title                           |
| `--base-url` | `/`                   | Base URL (for GitHub Pages subpaths) |

### Development Server

```bash
# Start local server
./goose-docs serve --dir dist --port 3000

# Custom port
./goose-docs serve --port 8080
```

---

## Documentation Directory Structure

Your markdown files should be organized like this:

```
docs/
├── index.md                    # Home page content
├── getting-started/
│   ├── introduction.md
│   ├── installation.md
│   └── quick-start.md
├── core-concepts/
│   ├── architecture.md
│   └── modules.md
└── ...
```

Sections are automatically detected from folder names.

---

## GitHub Pages Deployment

1. Build with the repository base URL:

   ```bash
   ./goose-docs build --input docs --output dist --base-url /your-repo-name/
   ```

2. Commit the `dist` folder or configure GitHub Actions to build automatically.

3. In your repo settings, set GitHub Pages to serve from the `dist` folder.

---

## Running Tests

```bash
# Run all tests
go test ./tests/...

# Run with verbose output
go test ./tests/... -v
```

---

## Code Coverage

```bash
# Coverage for all goose packages
go test ./tests/... -coverpkg=./...

```

---

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

---

## License

MIT License - see LICENSE file for details.
