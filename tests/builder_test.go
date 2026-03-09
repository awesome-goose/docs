package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/docs/builder"
	test "github.com/awesome-goose/goose/testing"
)

func TestBuilder(t *testing.T) {
	test.NewSuiteRunner(t, &BuilderSuite{}).Run()
}

type BuilderSuite struct {
	test.Suite
	tempDir string
}

func (s *BuilderSuite) SetupTest() {
	tempDir, err := os.MkdirTemp("", "docs-builder-test-*")
	if err != nil {
		s.T.T().Fatalf("Failed to create temp directory: %v", err)
	}
	s.tempDir = tempDir
}

func (s *BuilderSuite) TeardownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *BuilderSuite) TestBuild_WithDefaultConfig() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()
}

func (s *BuilderSuite) TestBuild_CreatesHTMLFiles() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	s.T.Expect(fileExists(filepath.Join(outputDir, "index.html"))).ToBeTrue()
}

func (s *BuilderSuite) TestBuild_CreatesAssets() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	s.T.Expect(fileExists(filepath.Join(outputDir, "assets", "styles.css"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputDir, "assets", "app.js"))).ToBeTrue()
}

func (s *BuilderSuite) TestBuild_CreatesSearchIndex() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	searchIndexPath := filepath.Join(outputDir, "search-index.json")
	s.T.Expect(fileExists(searchIndexPath)).ToBeTrue()
}

func (s *BuilderSuite) TestBuild_SearchIndexIsValidJSON() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "search-index.json"))
	s.T.Expect(err).ToBeNil()

	var data map[string]any
	err = json.Unmarshal([]byte(content), &data)
	s.T.Expect(err).ToBeNil()
}

func (s *BuilderSuite) TestBuild_WithMultipleSections() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	s.T.Expect(dirExists(filepath.Join(outputDir, "getting-started"))).ToBeTrue()
	s.T.Expect(dirExists(filepath.Join(outputDir, "core-concepts"))).ToBeTrue()
}

func (s *BuilderSuite) TestBuild_ProcessesMarkdownCorrectly() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	installContent, err := readFileContent(filepath.Join(outputDir, "getting-started", "installation.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(installContent).ToContainString("<html")
}

func (s *BuilderSuite) TestBuild_WithBaseURL() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/awesome-goose/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(content).ToContainString("/awesome-goose/")
}

func (s *BuilderSuite) TestBuild_HTMLStructure() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()

	s.T.Expect(content).ToContainString("<!DOCTYPE html>")
	s.T.Expect(content).ToContainString("<head>")
	s.T.Expect(content).ToContainString("<body>")
}

func (s *BuilderSuite) TestBuild_IncludesNavigation() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()

	s.T.Expect(content).ToContainString("sidebar")
}

func (s *BuilderSuite) TestBuild_TitleInOutput() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "My Custom Documentation",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(content).ToContainString("My Custom Documentation")
}

func (s *BuilderSuite) TestConfig_Initialization() {
	config := builder.Config{
		InputDir:  "docs",
		OutputDir: "dist",
		Title:     "My Docs",
		BaseURL:   "/",
	}

	s.T.Expect(config.InputDir).ToEqual("docs")
	s.T.Expect(config.OutputDir).ToEqual("dist")
	s.T.Expect(config.Title).ToEqual("My Docs")
	s.T.Expect(config.BaseURL).ToEqual("/")
}

func (s *BuilderSuite) TestPage_Initialization() {
	page := builder.Page{
		Title:   "Test Page",
		Path:    "test.md",
		URL:     "test.html",
		Content: "# Test",
		Section: "root",
		Order:   0,
	}

	s.T.Expect(page.Title).ToEqual("Test Page")
	s.T.Expect(page.Path).ToEqual("test.md")
	s.T.Expect(page.URL).ToEqual("test.html")
	s.T.Expect(page.Section).ToEqual("root")
	s.T.Expect(page.Order).ToEqual(0)
}

func (s *BuilderSuite) TestSection_Initialization() {
	section := builder.Section{
		Name: "Getting Started",
		Slug: "getting-started",
		Icon: "🚀",
		Pages: []builder.Page{
			{Title: "Installation", Path: "getting-started/installation.md"},
		},
	}

	s.T.Expect(section.Name).ToEqual("Getting Started")
	s.T.Expect(section.Slug).ToEqual("getting-started")
	s.T.Expect(section.Icon).ToEqual("🚀")
	s.T.Expect(len(section.Pages)).ToEqual(1)
}

func (s *BuilderSuite) TestBuild_EmptyInputDirectory() {
	docsDir := filepath.Join(s.tempDir, "empty-docs")
	err := os.MkdirAll(docsDir, 0755)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()
}

func (s *BuilderSuite) TestBuild_CreatesNestedOutputDirectories() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "nested", "output", "dir")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	s.T.Expect(dirExists(outputDir)).ToBeTrue()
}

func (s *BuilderSuite) TestBuild_CodeBlockSyntaxHighlighting() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "getting-started", "installation.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(content).ToContainString("<pre")
}

func (s *BuilderSuite) TestBuild_SearchIndexContainsPages() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	config := builder.Config{
		InputDir:  docsDir,
		OutputDir: outputDir,
		Title:     "Test Docs",
		BaseURL:   "/",
	}

	err = builder.Build(config)
	s.T.Expect(err).ToBeNil()

	content, err := readFileContent(filepath.Join(outputDir, "search-index.json"))
	s.T.Expect(err).ToBeNil()

	s.T.Expect(content).ToContainString("pages")
	s.T.Expect(content).ToContainString("title")
}
