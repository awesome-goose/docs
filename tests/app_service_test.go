package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awesome-goose/docs/app"
	test "github.com/awesome-goose/goose/testing"
)

func TestAppService(t *testing.T) {
	test.NewSuiteRunner(t, &AppServiceSuite{}).Run()
}

type AppServiceSuite struct {
	test.Suite
	service *app.AppService
	tempDir string
}

func (s *AppServiceSuite) SetupTest() {
	s.service = &app.AppService{}
	tempDir, err := os.MkdirTemp("", "docs-app-service-test-*")
	if err != nil {
		s.T.T().Fatalf("Failed to create temp directory: %v", err)
	}
	s.tempDir = tempDir
}

func (s *AppServiceSuite) TeardownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *AppServiceSuite) TestBuild_WithValidInput() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir, err := createOutputDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	indexPath := filepath.Join(outputDir, "index.html")
	s.T.Expect(fileExists(indexPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestBuild_WithSections() {
	docsDir, err := createTestDocsWithSections(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	s.T.Expect(fileExists(filepath.Join(outputDir, "index.html"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputDir, "getting-started", "installation.html"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputDir, "getting-started", "quick-start.html"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(outputDir, "core-concepts", "modules.html"))).ToBeTrue()
}

func (s *AppServiceSuite) TestBuild_CreatesAssetsDirectory() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	assetsDir := filepath.Join(outputDir, "assets")
	s.T.Expect(dirExists(assetsDir)).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(assetsDir, "styles.css"))).ToBeTrue()
	s.T.Expect(fileExists(filepath.Join(assetsDir, "app.js"))).ToBeTrue()
}

func (s *AppServiceSuite) TestBuild_CreatesSearchIndex() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	searchIndexPath := filepath.Join(outputDir, "search-index.json")
	s.T.Expect(fileExists(searchIndexPath)).ToBeTrue()
}

func (s *AppServiceSuite) TestBuild_WithCustomBaseURL() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/awesome-goose/")
	s.T.Expect(err).ToBeNil()

	indexContent, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(indexContent).ToContainString("/awesome-goose/")
}

func (s *AppServiceSuite) TestBuild_WithCustomTitle() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "My Custom Docs", "/")
	s.T.Expect(err).ToBeNil()

	indexContent, err := readFileContent(filepath.Join(outputDir, "index.html"))
	s.T.Expect(err).ToBeNil()
	s.T.Expect(indexContent).ToContainString("My Custom Docs")
}

func (s *AppServiceSuite) TestBuild_WithEmptyInput() {
	docsDir := filepath.Join(s.tempDir, "empty-docs")
	err := os.MkdirAll(docsDir, 0755)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()
}

func (s *AppServiceSuite) TestBuild_OutputDirectoryCreated() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "new", "nested", "output")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	s.T.Expect(dirExists(outputDir)).ToBeTrue()
}

func (s *AppServiceSuite) TestBuild_MarkdownConvertedToHTML() {
	docsDir, err := createTestDocsDir(s.tempDir)
	s.T.Expect(err).ToBeNil()

	outputDir := filepath.Join(s.tempDir, "dist")

	err = s.service.Build(docsDir, outputDir, "Test Docs", "/")
	s.T.Expect(err).ToBeNil()

	indexPath := filepath.Join(outputDir, "index.html")
	s.T.Expect(fileExists(indexPath)).ToBeTrue()

	content, err := readFileContent(indexPath)
	s.T.Expect(err).ToBeNil()
	s.T.Expect(content).ToContainString("<!DOCTYPE html>")
	s.T.Expect(content).ToContainString("<html")
}
