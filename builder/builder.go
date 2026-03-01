package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Config holds the build configuration
type Config struct {
	InputDir  string
	OutputDir string
	Title     string
	BaseURL   string
}

// Page represents a documentation page
type Page struct {
	Title       string `json:"title"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Content     string `json:"content,omitempty"`
	HTMLContent template.HTML
	Section     string `json:"section"`
	Order       int    `json:"order"`
}

// Section represents a documentation section
type Section struct {
	Name  string
	Slug  string
	Pages []Page
	Icon  string
}

// NavItem represents a navigation item
type NavItem struct {
	Name     string
	Slug     string
	Pages    []Page
	Icon     string
	IsActive bool
}

// SearchIndex represents the search index
type SearchIndex struct {
	Pages []SearchPage `json:"pages"`
}

// SearchPage represents a page in the search index
type SearchPage struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Section string `json:"section"`
	Content string `json:"content"`
}

// Build builds the documentation site
func Build(config Config) error {
	// Normalize config values with defaults
	if config.InputDir == "" {
		config.InputDir = "docs"
	}
	if config.OutputDir == "" {
		config.OutputDir = "dist"
	}
	if config.Title == "" {
		config.Title = "Goose Documentation"
	}
	if config.BaseURL == "" {
		config.BaseURL = "/"
	}
	// Ensure BaseURL ends with /
	if !strings.HasSuffix(config.BaseURL, "/") {
		config.BaseURL = config.BaseURL + "/"
	}

	// Initialize goldmark with syntax highlighting
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	// Collect all markdown files
	pages, err := collectPages(config.InputDir, config.BaseURL, md)
	if err != nil {
		return fmt.Errorf("failed to collect pages: %w", err)
	}

	// Organize pages into sections
	sections := organizeSections(pages)

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate HTML pages
	for i := range pages {
		if err := generatePage(&pages[i], sections, config); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", pages[i].Path, err)
		}
	}

	// Generate index page
	if err := generateIndexPage(sections, config); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate search index
	if err := generateSearchIndex(pages, config); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	// Copy static assets
	if err := copyAssets(config.OutputDir); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	return nil
}

func collectPages(inputDir, baseURL string, md goldmark.Markdown) ([]Page, error) {
	var pages []Page

	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Convert markdown to HTML
		var buf bytes.Buffer
		if err := md.Convert(content, &buf); err != nil {
			return fmt.Errorf("failed to convert markdown %s: %w", path, err)
		}

		// Extract title from content
		title := extractTitle(string(content))
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(path), ".md")
			title = toTitle(strings.ReplaceAll(title, "-", " "))
		}

		// Calculate relative path
		relPath, _ := filepath.Rel(inputDir, path)
		section := filepath.Dir(relPath)
		if section == "." {
			section = "root"
		}

		// Generate URL
		url := strings.TrimSuffix(relPath, ".md") + ".html"
		url = filepath.ToSlash(url)
		if baseURL != "/" {
			url = strings.TrimSuffix(baseURL, "/") + "/" + url
		}

		// Determine order based on filename or position
		order := getPageOrder(filepath.Base(path))

		pages = append(pages, Page{
			Title:       title,
			Path:        relPath,
			URL:         url,
			Content:     string(content),
			HTMLContent: template.HTML(buf.String()),
			Section:     section,
			Order:       order,
		})

		return nil
	})

	return pages, err
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func toTitle(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func getPageOrder(filename string) int {
	orderMap := map[string]int{
		"index.md":               0,
		"introduction.md":        1,
		"overview.md":            2,
		"installation.md":        3,
		"quick-start.md":         4,
		"configuration.md":       5,
		"directory-structure.md": 6,
	}
	if order, ok := orderMap[filename]; ok {
		return order
	}
	return 100
}

func organizeSections(pages []Page) []Section {
	sectionMap := make(map[string]*Section)
	sectionOrder := map[string]int{
		"root":            0,
		"getting-started": 1,
		"core-concepts":   2,
		"building-blocks": 3,
		"platforms":       4,
		"database":        5,
		"security":        6,
		"testing":         7,
		"advanced":        8,
		"deployment":      9,
		"cli":             10,
		"tutorials":       11,
		"community":       12,
	}

	sectionIcons := map[string]string{
		"root":            "📚",
		"getting-started": "🚀",
		"core-concepts":   "🏗️",
		"building-blocks": "🧱",
		"platforms":       "🖥️",
		"database":        "🗄️",
		"security":        "🔒",
		"testing":         "🧪",
		"advanced":        "⚡",
		"deployment":      "🚢",
		"cli":             "💻",
		"tutorials":       "📖",
		"community":       "👥",
	}

	for _, page := range pages {
		if _, ok := sectionMap[page.Section]; !ok {
			name := toTitle(strings.ReplaceAll(page.Section, "-", " "))
			if page.Section == "root" {
				name = "Overview"
			}
			icon := sectionIcons[page.Section]
			if icon == "" {
				icon = "📄"
			}
			sectionMap[page.Section] = &Section{
				Name:  name,
				Slug:  page.Section,
				Icon:  icon,
				Pages: []Page{},
			}
		}
		sectionMap[page.Section].Pages = append(sectionMap[page.Section].Pages, page)
	}

	var sections []Section
	for _, section := range sectionMap {
		sort.Slice(section.Pages, func(i, j int) bool {
			if section.Pages[i].Order != section.Pages[j].Order {
				return section.Pages[i].Order < section.Pages[j].Order
			}
			return section.Pages[i].Title < section.Pages[j].Title
		})
		sections = append(sections, *section)
	}

	sort.Slice(sections, func(i, j int) bool {
		oi := sectionOrder[sections[i].Slug]
		oj := sectionOrder[sections[j].Slug]
		if oi != oj {
			return oi < oj
		}
		return sections[i].Name < sections[j].Name
	})

	return sections
}

func generatePage(page *Page, sections []Section, config Config) error {
	var navItems []NavItem
	for _, section := range sections {
		isActive := section.Slug == page.Section
		navItems = append(navItems, NavItem{
			Name:     section.Name,
			Slug:     section.Slug,
			Pages:    section.Pages,
			Icon:     section.Icon,
			IsActive: isActive,
		})
	}

	data := struct {
		Page     *Page
		Sections []NavItem
		Config   Config
	}{
		Page:     page,
		Sections: navItems,
		Config:   config,
	}

	tmpl, err := template.New("page").Funcs(template.FuncMap{
		"isActive": func(p Page, current *Page) bool {
			return p.Path == current.Path
		},
	}).Parse(pageTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputPath := filepath.Join(config.OutputDir, strings.TrimSuffix(page.Path, ".md")+".html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func generateIndexPage(sections []Section, config Config) error {
	var navItems []NavItem
	for _, section := range sections {
		navItems = append(navItems, NavItem{
			Name:     section.Name,
			Slug:     section.Slug,
			Pages:    section.Pages,
			Icon:     section.Icon,
			IsActive: section.Slug == "root",
		})
	}

	data := struct {
		Sections []NavItem
		Config   Config
	}{
		Sections: navItems,
		Config:   config,
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	outputPath := filepath.Join(config.OutputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func generateSearchIndex(pages []Page, config Config) error {
	var searchPages []SearchPage
	re := regexp.MustCompile(`<[^>]*>`)

	for _, page := range pages {
		plainContent := re.ReplaceAllString(string(page.HTMLContent), " ")
		plainContent = strings.Join(strings.Fields(plainContent), " ")
		if len(plainContent) > 500 {
			plainContent = plainContent[:500] + "..."
		}

		searchPages = append(searchPages, SearchPage{
			Title:   page.Title,
			URL:     page.URL,
			Section: page.Section,
			Content: plainContent,
		})
	}

	index := SearchIndex{Pages: searchPages}
	jsonData, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	outputPath := filepath.Join(config.OutputDir, "search-index.json")
	return os.WriteFile(outputPath, jsonData, 0644)
}

func copyAssets(outputDir string) error {
	cssPath := filepath.Join(outputDir, "assets", "styles.css")
	if err := os.MkdirAll(filepath.Dir(cssPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(cssPath, []byte(cssStyles), 0644); err != nil {
		return err
	}

	jsPath := filepath.Join(outputDir, "assets", "app.js")
	if err := os.WriteFile(jsPath, []byte(jsScript), 0644); err != nil {
		return err
	}

	return nil
}

// Serve starts a local development server
func Serve(dir, port string) error {
	fmt.Printf("🪿 Starting development server at http://localhost:%s\n", port)
	fmt.Println("   Press Ctrl+C to stop")

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	return http.ListenAndServe(":"+port, nil)
}
