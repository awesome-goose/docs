package app

import (
	"github.com/awesome-goose/docs/builder"
)

type AppService struct{}

// Build builds the documentation site from markdown files
func (s *AppService) Build(input, output, title, baseURL string) error {
	config := builder.Config{
		InputDir:  input,
		OutputDir: output,
		Title:     title,
		BaseURL:   baseURL,
	}
	return builder.Build(config)
}

// Serve starts a local development server
func (s *AppService) Serve(dir, port string) error {
	return builder.Serve(dir, port)
}
