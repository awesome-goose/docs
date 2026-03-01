package app

import (
	"fmt"

	"github.com/awesome-goose/goose/io/output"
	"github.com/awesome-goose/goose/types"
)

type AppController struct {
	appService *AppService `inject:""`
	log        types.Log   `inject:""`
}

// Build builds the documentation site from markdown files
func (c *AppController) Build(dto *BuildDto) types.Output {
	c.log.Info("Building documentation site", map[string]any{
		"input":    dto.Input,
		"output":   dto.Output,
		"title":    dto.Title,
		"base_url": dto.BaseURL,
	})

	err := c.appService.Build(dto.Input, dto.Output, dto.Title, dto.BaseURL)
	if err != nil {
		return output.ConsoleError(fmt.Sprintf("Failed to build: %v", err))
	}

	return output.ConsoleSuccess(fmt.Sprintf("✅ Documentation built successfully!\n   Output: %s/", dto.Output))
}

// Serve starts a local development server
func (c *AppController) Serve(dto *ServeDto) types.Output {
	c.log.Info("Starting development server", map[string]any{
		"dir":  dto.Dir,
		"port": dto.Port,
	})

	err := c.appService.Serve(dto.Dir, dto.Port)
	if err != nil {
		return output.ConsoleError(fmt.Sprintf("Server error: %v", err))
	}

	return output.ConsoleSuccess("Server stopped")
}

// Help displays usage information
func (c *AppController) Help(dto *HelpDto) types.Output {
	helpText := `
🪿 Goose Documentation Builder

A CLI tool to build beautiful documentation websites from markdown files.

Usage:
  goose-docs <command> [options]

Commands:
  build     Build the documentation site from markdown files
  serve     Start a local development server
  help      Show this help message

Build Options:
  --input     Input directory containing markdown files (default: "docs")
  --output    Output directory for the static site (default: "dist")
  --title     Site title (default: "Goose Documentation")
  --base-url  Base URL for the site (default: "/")

Serve Options:
  --dir       Directory to serve (default: "dist")
  --port      Port to serve on (default: "3000")

Examples:
  goose-docs build
  goose-docs build --input ./docs --output ./public --base-url /awesome-goose/
  goose-docs serve --port 8080
`
	return output.ConsoleSuccess(helpText)
}
