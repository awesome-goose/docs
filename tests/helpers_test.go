package tests

import (
	"os"
	"path/filepath"
)

// Test helpers for docs tests

// createTestDocsDir creates a temporary docs directory with sample markdown files
func createTestDocsDir(baseDir string) (string, error) {
	docsDir := filepath.Join(baseDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return "", err
	}

	// Create index.md
	indexContent := `# Welcome to Goose

This is the main documentation page.

## Getting Started

Follow these steps to get started with Goose.
`
	if err := os.WriteFile(filepath.Join(docsDir, "index.md"), []byte(indexContent), 0644); err != nil {
		return "", err
	}

	return docsDir, nil
}

// createTestDocsWithSections creates a docs directory with multiple sections
func createTestDocsWithSections(baseDir string) (string, error) {
	docsDir := filepath.Join(baseDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return "", err
	}

	// Create index.md
	indexContent := `# Welcome to Goose

This is the main documentation page.
`
	if err := os.WriteFile(filepath.Join(docsDir, "index.md"), []byte(indexContent), 0644); err != nil {
		return "", err
	}

	// Create getting-started section
	gettingStartedDir := filepath.Join(docsDir, "getting-started")
	if err := os.MkdirAll(gettingStartedDir, 0755); err != nil {
		return "", err
	}

	installContent := `# Installation

Install Goose using go get:

` + "```bash" + `
go get github.com/awesome-goose/goose
` + "```" + `
`
	if err := os.WriteFile(filepath.Join(gettingStartedDir, "installation.md"), []byte(installContent), 0644); err != nil {
		return "", err
	}

	quickStartContent := `# Quick Start

Get started quickly with Goose.

1. Create a new project
2. Add your modules
3. Run the application
`
	if err := os.WriteFile(filepath.Join(gettingStartedDir, "quick-start.md"), []byte(quickStartContent), 0644); err != nil {
		return "", err
	}

	// Create core-concepts section
	coreConceptsDir := filepath.Join(docsDir, "core-concepts")
	if err := os.MkdirAll(coreConceptsDir, 0755); err != nil {
		return "", err
	}

	modulesContent := `# Modules

Modules are the building blocks of a Goose application.

## Creating a Module

Create a module by defining its components.
`
	if err := os.WriteFile(filepath.Join(coreConceptsDir, "modules.md"), []byte(modulesContent), 0644); err != nil {
		return "", err
	}

	return docsDir, nil
}

// createOutputDir creates an output directory for test builds
func createOutputDir(baseDir string) (string, error) {
	outputDir := filepath.Join(baseDir, "dist")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	return outputDir, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// readFileContent reads the content of a file
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// countFilesInDir counts the number of files in a directory (non-recursive)
func countFilesInDir(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// countFilesRecursive counts all files in a directory recursively
func countFilesRecursive(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}
