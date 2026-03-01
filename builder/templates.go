package builder

const pageTemplate = `<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Page.Title}} - {{.Config.Title}}">
    <title>{{.Page.Title}} | {{.Config.Title}}</title>
    <link rel="stylesheet" href="{{.Config.BaseURL}}assets/styles.css">
    <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🪿</text></svg>">
</head>
<body>
    <!-- Mobile Menu Toggle -->
    <button class="mobile-menu-toggle" aria-label="Toggle navigation menu">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
    </button>

    <!-- Sidebar -->
    <aside class="sidebar">
        <div class="sidebar-header">
            <a href="{{.Config.BaseURL}}index.html" class="logo">
                <span class="logo-icon">🪿</span>
                <span class="logo-text">Goose</span>
            </a>
            <button class="theme-toggle" aria-label="Toggle dark mode">
                <svg class="sun-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="5"></circle>
                    <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"></path>
                </svg>
                <svg class="moon-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                </svg>
            </button>
        </div>

        <!-- Search -->
        <div class="search-container">
            <div class="search-input-wrapper">
                <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-input" placeholder="Search docs..." aria-label="Search documentation">
                <kbd class="search-shortcut">⌘K</kbd>
            </div>
            <div class="search-results"></div>
        </div>

        <!-- Navigation -->
        <nav class="sidebar-nav">
            {{range .Sections}}
            <div class="nav-section {{if .IsActive}}active{{end}}">
                <button class="nav-section-header" aria-expanded="{{if .IsActive}}true{{else}}false{{end}}">
                    <span class="nav-section-icon">{{.Icon}}</span>
                    <span class="nav-section-title">{{.Name}}</span>
                    <svg class="nav-section-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 18l6-6-6-6"></path>
                    </svg>
                </button>
                <ul class="nav-section-items">
                    {{range .Pages}}
                    <li>
                        <a href="{{.URL}}" class="nav-link {{if isActive . $.Page}}active{{end}}">
                            {{.Title}}
                        </a>
                    </li>
                    {{end}}
                </ul>
            </div>
            {{end}}
        </nav>

        <!-- Sidebar Footer -->
        <div class="sidebar-footer">
            <a href="https://github.com/awesome-goose/goose" target="_blank" rel="noopener" class="sidebar-footer-link">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                </svg>
                <span>GitHub</span>
            </a>
        </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
        <article class="content">
            <nav class="breadcrumb">
                <a href="{{.Config.BaseURL}}index.html">Home</a>
                {{if ne .Page.Section "root"}}
                <span class="breadcrumb-separator">/</span>
                <span>{{.Page.Section}}</span>
                {{end}}
                <span class="breadcrumb-separator">/</span>
                <span class="breadcrumb-current">{{.Page.Title}}</span>
            </nav>

            <div class="prose">
                {{.Page.HTMLContent}}
            </div>

            <!-- Page Navigation -->
            <nav class="page-nav">
                <div class="page-nav-info">
                    <span class="page-nav-edit">
                        Found an issue? <a href="https://github.com/awesome-goose/goose/edit/main/docs/{{.Page.Path}}" target="_blank" rel="noopener">Edit this page on GitHub</a>
                    </span>
                </div>
            </nav>
        </article>

        <!-- Table of Contents -->
        <aside class="toc">
            <div class="toc-header">On this page</div>
            <nav class="toc-nav">
                <!-- Populated by JavaScript -->
            </nav>
        </aside>
    </main>

    <!-- Overlay for mobile -->
    <div class="sidebar-overlay"></div>

    <!-- Search Modal -->
    <div class="search-modal" role="dialog" aria-modal="true" aria-label="Search documentation">
        <div class="search-modal-content">
            <div class="search-modal-header">
                <svg class="search-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-modal-input" placeholder="Search documentation..." autofocus>
                <button class="search-modal-close" aria-label="Close search">
                    <kbd>ESC</kbd>
                </button>
            </div>
            <div class="search-modal-results"></div>
            <div class="search-modal-footer">
                <span><kbd>↑↓</kbd> Navigate</span>
                <span><kbd>↵</kbd> Select</span>
                <span><kbd>ESC</kbd> Close</span>
            </div>
        </div>
    </div>

    <script src="{{.Config.BaseURL}}assets/app.js"></script>
</body>
</html>`

const indexTemplate = `<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Config.Title}} - Modern, modular Go framework for web, API, and CLI applications">
    <title>{{.Config.Title}}</title>
    <link rel="stylesheet" href="{{.Config.BaseURL}}assets/styles.css">
    <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🪿</text></svg>">
</head>
<body>
    <!-- Mobile Menu Toggle -->
    <button class="mobile-menu-toggle" aria-label="Toggle navigation menu">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
    </button>

    <!-- Sidebar -->
    <aside class="sidebar">
        <div class="sidebar-header">
            <a href="{{.Config.BaseURL}}index.html" class="logo">
                <span class="logo-icon">🪿</span>
                <span class="logo-text">Goose</span>
            </a>
            <button class="theme-toggle" aria-label="Toggle dark mode">
                <svg class="sun-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="5"></circle>
                    <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"></path>
                </svg>
                <svg class="moon-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                </svg>
            </button>
        </div>

        <!-- Search -->
        <div class="search-container">
            <div class="search-input-wrapper">
                <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-input" placeholder="Search docs..." aria-label="Search documentation">
                <kbd class="search-shortcut">⌘K</kbd>
            </div>
            <div class="search-results"></div>
        </div>

        <!-- Navigation -->
        <nav class="sidebar-nav">
            {{range .Sections}}
            <div class="nav-section {{if eq .Slug "root"}}active{{end}}">
                <button class="nav-section-header" aria-expanded="{{if eq .Slug "root"}}true{{else}}false{{end}}">
                    <span class="nav-section-icon">{{.Icon}}</span>
                    <span class="nav-section-title">{{.Name}}</span>
                    <svg class="nav-section-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M9 18l6-6-6-6"></path>
                    </svg>
                </button>
                <ul class="nav-section-items">
                    {{range .Pages}}
                    <li>
                        <a href="{{.URL}}" class="nav-link">
                            {{.Title}}
                        </a>
                    </li>
                    {{end}}
                </ul>
            </div>
            {{end}}
        </nav>

        <!-- Sidebar Footer -->
        <div class="sidebar-footer">
            <a href="https://github.com/awesome-goose/goose" target="_blank" rel="noopener" class="sidebar-footer-link">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                </svg>
                <span>GitHub</span>
            </a>
        </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
        <article class="content home-content">
            <!-- Hero Section -->
            <div class="hero">
                <div class="hero-icon">🪿</div>
                <h1 class="hero-title">Goose Framework</h1>
                <p class="hero-subtitle">A modern, modular Go framework for building scalable web applications, RESTful APIs, and command-line tools.</p>
                <div class="hero-actions">
                    <a href="{{.Config.BaseURL}}getting-started/quick-start.html" class="btn btn-primary">Get Started</a>
                    <a href="{{.Config.BaseURL}}getting-started/introduction.html" class="btn btn-secondary">Learn More</a>
                </div>
            </div>

            <!-- Features Grid -->
            <div class="features">
                <div class="feature-card">
                    <div class="feature-icon">🚀</div>
                    <h3 class="feature-title">High Performance</h3>
                    <p class="feature-description">Built on Go's fast runtime for maximum performance with minimal overhead.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">📦</div>
                    <h3 class="feature-title">Modular Architecture</h3>
                    <p class="feature-description">Everything is a module. Share, reuse, and maintain code easily.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">🖥️</div>
                    <h3 class="feature-title">Multi-Platform</h3>
                    <p class="feature-description">Build APIs, web apps, or CLI tools with consistent patterns.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">💉</div>
                    <h3 class="feature-title">Dependency Injection</h3>
                    <p class="feature-description">Built-in IoC container for clean dependency management.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">🧪</div>
                    <h3 class="feature-title">Testing Support</h3>
                    <p class="feature-description">Comprehensive testing utilities for reliable applications.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">📝</div>
                    <h3 class="feature-title">Declarative Design</h3>
                    <p class="feature-description">Convention over configuration for predictable behavior.</p>
                </div>
            </div>

            <!-- Quick Links -->
            <div class="quick-links">
                <h2 class="quick-links-title">Quick Links</h2>
                <div class="quick-links-grid">
                    {{range .Sections}}
                    {{if ne .Slug "root"}}
                    <a href="{{if .Pages}}{{(index .Pages 0).URL}}{{else}}#{{end}}" class="quick-link-card">
                        <span class="quick-link-icon">{{.Icon}}</span>
                        <span class="quick-link-title">{{.Name}}</span>
                    </a>
                    {{end}}
                    {{end}}
                </div>
            </div>
        </article>
    </main>

    <!-- Overlay for mobile -->
    <div class="sidebar-overlay"></div>

    <!-- Search Modal -->
    <div class="search-modal" role="dialog" aria-modal="true" aria-label="Search documentation">
        <div class="search-modal-content">
            <div class="search-modal-header">
                <svg class="search-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-modal-input" placeholder="Search documentation..." autofocus>
                <button class="search-modal-close" aria-label="Close search">
                    <kbd>ESC</kbd>
                </button>
            </div>
            <div class="search-modal-results"></div>
            <div class="search-modal-footer">
                <span><kbd>↑↓</kbd> Navigate</span>
                <span><kbd>↵</kbd> Select</span>
                <span><kbd>ESC</kbd> Close</span>
            </div>
        </div>
    </div>

    <script src="{{.Config.BaseURL}}assets/app.js"></script>
</body>
</html>`
