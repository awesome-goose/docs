package tests

import (
	"testing"

	"github.com/awesome-goose/docs/app"
	test "github.com/awesome-goose/goose/testing"
)

func TestAppController(t *testing.T) {
	test.NewSuiteRunner(t, &AppControllerSuite{}).Run()
}

type AppControllerSuite struct {
	test.Suite
}

func (s *AppControllerSuite) TestBuildDto_Initialization() {
	dto := &app.BuildDto{
		Input:   "docs",
		Output:  "dist",
		Title:   "My Docs",
		BaseURL: "/my-docs/",
	}

	s.T.Expect(dto.Input).ToEqual("docs")
	s.T.Expect(dto.Output).ToEqual("dist")
	s.T.Expect(dto.Title).ToEqual("My Docs")
	s.T.Expect(dto.BaseURL).ToEqual("/my-docs/")
}

func (s *AppControllerSuite) TestBuildDto_EmptyValues() {
	dto := &app.BuildDto{}

	s.T.Expect(dto.Input).ToEqual("")
	s.T.Expect(dto.Output).ToEqual("")
	s.T.Expect(dto.Title).ToEqual("")
	s.T.Expect(dto.BaseURL).ToEqual("")
}

func (s *AppControllerSuite) TestBuildDto_WithDefaultLikeValues() {
	dto := &app.BuildDto{
		Input:   "docs",
		Output:  "dist",
		Title:   "Goose Documentation",
		BaseURL: "/",
	}

	s.T.Expect(dto.Input).ToEqual("docs")
	s.T.Expect(dto.Output).ToEqual("dist")
	s.T.Expect(dto.Title).ToEqual("Goose Documentation")
	s.T.Expect(dto.BaseURL).ToEqual("/")
}

func (s *AppControllerSuite) TestServeDto_Initialization() {
	dto := &app.ServeDto{
		Dir:  "dist",
		Port: "8080",
	}

	s.T.Expect(dto.Dir).ToEqual("dist")
	s.T.Expect(dto.Port).ToEqual("8080")
}

func (s *AppControllerSuite) TestServeDto_EmptyValues() {
	dto := &app.ServeDto{}

	s.T.Expect(dto.Dir).ToEqual("")
	s.T.Expect(dto.Port).ToEqual("")
}

func (s *AppControllerSuite) TestServeDto_WithDefaultLikeValues() {
	dto := &app.ServeDto{
		Dir:  "dist",
		Port: "3000",
	}

	s.T.Expect(dto.Dir).ToEqual("dist")
	s.T.Expect(dto.Port).ToEqual("3000")
}

func (s *AppControllerSuite) TestHelpDto_Initialization() {
	dto := &app.HelpDto{}
	s.T.Expect(dto).Not().ToBeNil()
}

func (s *AppControllerSuite) TestBuildDto_WithCustomBaseURL() {
	dto := &app.BuildDto{
		Input:   "./my-docs",
		Output:  "./public",
		Title:   "Custom Documentation",
		BaseURL: "/awesome-goose/docs/",
	}

	s.T.Expect(dto.Input).ToEqual("./my-docs")
	s.T.Expect(dto.Output).ToEqual("./public")
	s.T.Expect(dto.Title).ToEqual("Custom Documentation")
	s.T.Expect(dto.BaseURL).ToEqual("/awesome-goose/docs/")
}

func (s *AppControllerSuite) TestServeDto_WithCustomPort() {
	dto := &app.ServeDto{
		Dir:  "build",
		Port: "9000",
	}

	s.T.Expect(dto.Dir).ToEqual("build")
	s.T.Expect(dto.Port).ToEqual("9000")
}

func (s *AppControllerSuite) TestBuildDto_AllFieldsSet() {
	dto := &app.BuildDto{
		Input:   "content/docs",
		Output:  "public/docs",
		Title:   "API Reference",
		BaseURL: "/api/",
	}

	s.T.Expect(dto.Input).Not().ToBeEmpty()
	s.T.Expect(dto.Output).Not().ToBeEmpty()
	s.T.Expect(dto.Title).Not().ToBeEmpty()
	s.T.Expect(dto.BaseURL).Not().ToBeEmpty()
}

func (s *AppControllerSuite) TestServeDto_AllFieldsSet() {
	dto := &app.ServeDto{
		Dir:  "docs-output",
		Port: "4000",
	}

	s.T.Expect(dto.Dir).Not().ToBeEmpty()
	s.T.Expect(dto.Port).Not().ToBeEmpty()
}
