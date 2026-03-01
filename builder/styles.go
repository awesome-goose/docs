package builder

const cssStyles = `/* ========================================
   Goose Documentation - Modern CSS Styles
   ======================================== */

/* CSS Custom Properties (Design Tokens) */
:root {
  /* Colors - Light Mode */
  --color-bg: #ffffff;
  --color-bg-secondary: #f8fafc;
  --color-bg-tertiary: #f1f5f9;
  --color-text: #1e293b;
  --color-text-secondary: #475569;
  --color-text-muted: #64748b;
  --color-border: #e2e8f0;
  --color-border-light: #f1f5f9;
  --color-primary: #3b82f6;
  --color-primary-hover: #2563eb;
  --color-primary-light: #dbeafe;
  --color-accent: #f59e0b;
  --color-success: #10b981;
  --color-code-bg: #1e1e1e;
  --color-code-text: #d4d4d4;
  
  /* Sidebar */
  --sidebar-width: 280px;
  --sidebar-bg: var(--color-bg);
  
  /* TOC */
  --toc-width: 220px;
  
  /* Typography */
  --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  --font-mono: 'SF Mono', SFMono-Regular, ui-monospace, 'DejaVu Sans Mono', Menlo, Consolas, monospace;
  
  /* Spacing */
  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 1.5rem;
  --space-xl: 2rem;
  --space-2xl: 3rem;
  
  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 250ms ease;
  
  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  
  /* Border Radius */
  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-full: 9999px;
}

/* Dark Mode */
[data-theme="dark"] {
  --color-bg: #0f172a;
  --color-bg-secondary: #1e293b;
  --color-bg-tertiary: #334155;
  --color-text: #f1f5f9;
  --color-text-secondary: #cbd5e1;
  --color-text-muted: #94a3b8;
  --color-border: #334155;
  --color-border-light: #1e293b;
  --color-primary: #60a5fa;
  --color-primary-hover: #93c5fd;
  --color-primary-light: #1e3a5f;
  --sidebar-bg: #0f172a;
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
}

/* Reset & Base */
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  font-size: 16px;
  scroll-behavior: smooth;
  -webkit-text-size-adjust: 100%;
}

body {
  font-family: var(--font-sans);
  font-size: 1rem;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-bg);
  display: flex;
  min-height: 100vh;
  overflow-x: hidden;
}

/* ========================================
   Sidebar Styles
   ======================================== */

.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: var(--sidebar-width);
  background: var(--sidebar-bg);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: transform var(--transition-normal);
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-lg);
  border-bottom: 1px solid var(--color-border);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  text-decoration: none;
  color: var(--color-text);
  font-weight: 700;
  font-size: 1.25rem;
}

.logo-icon {
  font-size: 1.5rem;
}

.logo:hover {
  color: var(--color-primary);
}

/* Theme Toggle */
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-toggle:hover {
  background: var(--color-bg-tertiary);
  color: var(--color-text);
}

.moon-icon { display: none; }
[data-theme="dark"] .sun-icon { display: none; }
[data-theme="dark"] .moon-icon { display: block; }

/* Search */
.search-container {
  padding: var(--space-md) var(--space-lg);
  position: relative;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: var(--space-sm);
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: var(--space-sm) var(--space-md);
  padding-left: 2.25rem;
  padding-right: 3rem;
  font-size: 0.875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-secondary);
  color: var(--color-text);
  transition: all var(--transition-fast);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-shortcut {
  position: absolute;
  right: var(--space-sm);
  padding: 2px 6px;
  font-size: 0.75rem;
  font-family: var(--font-sans);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-md) 0;
  scrollbar-width: thin;
  scrollbar-color: var(--color-border) transparent;
}

.sidebar-nav::-webkit-scrollbar {
  width: 6px;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: var(--radius-full);
}

.nav-section {
  margin-bottom: var(--space-xs);
}

.nav-section-header {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  width: 100%;
  padding: var(--space-sm) var(--space-lg);
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text);
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: all var(--transition-fast);
}

.nav-section-header:hover {
  background: var(--color-bg-secondary);
}

.nav-section-icon {
  font-size: 1rem;
}

.nav-section-title {
  flex: 1;
}

.nav-section-arrow {
  transition: transform var(--transition-fast);
  color: var(--color-text-muted);
}

.nav-section.active .nav-section-arrow {
  transform: rotate(90deg);
}

.nav-section-items {
  display: none;
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-section.active .nav-section-items {
  display: block;
}

.nav-link {
  display: block;
  padding: var(--space-xs) var(--space-lg);
  padding-left: calc(var(--space-lg) + 1.75rem);
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  border-left: 2px solid transparent;
}

.nav-link:hover {
  color: var(--color-text);
  background: var(--color-bg-secondary);
}

.nav-link.active {
  color: var(--color-primary);
  background: var(--color-primary-light);
  border-left-color: var(--color-primary);
  font-weight: 500;
}

/* Sidebar Footer */
.sidebar-footer {
  padding: var(--space-md) var(--space-lg);
  border-top: 1px solid var(--color-border);
}

.sidebar-footer-link {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-sm) var(--space-md);
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  text-decoration: none;
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.sidebar-footer-link:hover {
  background: var(--color-bg-secondary);
  color: var(--color-text);
}

/* ========================================
   Main Content
   ======================================== */

.main-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  display: flex;
  min-height: 100vh;
}

.content {
  flex: 1;
  max-width: 900px;
  padding: var(--space-xl) var(--space-2xl);
  margin: 0 auto;
}

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin-bottom: var(--space-xl);
}

.breadcrumb a {
  color: var(--color-text-muted);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.breadcrumb a:hover {
  color: var(--color-primary);
}

.breadcrumb-separator {
  color: var(--color-border);
}

.breadcrumb-current {
  color: var(--color-text);
  font-weight: 500;
}

/* ========================================
   Prose (Content Styling)
   ======================================== */

.prose {
  color: var(--color-text);
  max-width: none;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  color: var(--color-text);
  font-weight: 700;
  line-height: 1.3;
  margin-top: var(--space-xl);
  margin-bottom: var(--space-md);
  scroll-margin-top: var(--space-xl);
}

.prose h1 { font-size: 2.25rem; margin-top: 0; }
.prose h2 { font-size: 1.75rem; border-bottom: 1px solid var(--color-border); padding-bottom: var(--space-sm); }
.prose h3 { font-size: 1.375rem; }
.prose h4 { font-size: 1.125rem; }

.prose p {
  margin-bottom: var(--space-md);
}

.prose a {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
  transition: color var(--transition-fast);
}

.prose a:hover {
  color: var(--color-primary-hover);
  text-decoration: underline;
}

.prose strong {
  font-weight: 600;
  color: var(--color-text);
}

.prose ul, .prose ol {
  margin-bottom: var(--space-md);
  padding-left: var(--space-lg);
}

.prose li {
  margin-bottom: var(--space-xs);
}

.prose li::marker {
  color: var(--color-text-muted);
}

/* Code */
.prose code {
  font-family: var(--font-mono);
  font-size: 0.875em;
  background: var(--color-bg-tertiary);
  padding: 0.125rem 0.375rem;
  border-radius: var(--radius-sm);
  color: var(--color-primary);
}

.prose pre {
  background: var(--color-code-bg);
  border-radius: var(--radius-lg);
  padding: var(--space-lg);
  margin-bottom: var(--space-md);
  overflow-x: auto;
  font-size: 0.875rem;
  line-height: 1.7;
  border: 1px solid var(--color-border);
}

.prose pre code {
  background: transparent;
  padding: 0;
  color: var(--color-code-text);
  font-size: inherit;
}

/* Tables */
.prose table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: var(--space-lg);
  font-size: 0.875rem;
}

.prose th, .prose td {
  padding: var(--space-sm) var(--space-md);
  text-align: left;
  border: 1px solid var(--color-border);
}

.prose th {
  background: var(--color-bg-secondary);
  font-weight: 600;
  color: var(--color-text);
}

.prose tr:nth-child(even) {
  background: var(--color-bg-secondary);
}

/* Blockquotes */
.prose blockquote {
  border-left: 4px solid var(--color-primary);
  background: var(--color-bg-secondary);
  padding: var(--space-md) var(--space-lg);
  margin-bottom: var(--space-md);
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  color: var(--color-text-secondary);
}

.prose blockquote p:last-child {
  margin-bottom: 0;
}

/* Horizontal Rule */
.prose hr {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: var(--space-xl) 0;
}

/* Images */
.prose img {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-md);
  margin: var(--space-md) 0;
}

/* Task Lists */
.prose input[type="checkbox"] {
  margin-right: var(--space-sm);
}

/* ========================================
   Table of Contents (TOC)
   ======================================== */

.toc {
  position: sticky;
  top: var(--space-xl);
  width: var(--toc-width);
  max-height: calc(100vh - var(--space-2xl));
  overflow-y: auto;
  padding: var(--space-xl) var(--space-lg);
  flex-shrink: 0;
  display: none;
}

@media (min-width: 1280px) {
  .toc {
    display: block;
  }
}

.toc-header {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  margin-bottom: var(--space-md);
}

.toc-nav a {
  display: block;
  padding: var(--space-xs) 0;
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  text-decoration: none;
  border-left: 2px solid transparent;
  padding-left: var(--space-md);
  margin-left: calc(-1 * var(--space-md));
  transition: all var(--transition-fast);
}

.toc-nav a:hover {
  color: var(--color-text);
}

.toc-nav a.active {
  color: var(--color-primary);
  border-left-color: var(--color-primary);
}

.toc-nav a.toc-h3 {
  padding-left: calc(var(--space-md) + var(--space-md));
}

/* ========================================
   Page Navigation
   ======================================== */

.page-nav {
  margin-top: var(--space-2xl);
  padding-top: var(--space-xl);
  border-top: 1px solid var(--color-border);
}

.page-nav-info {
  text-align: center;
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.page-nav-edit a {
  color: var(--color-primary);
  text-decoration: none;
}

.page-nav-edit a:hover {
  text-decoration: underline;
}

/* ========================================
   Home Page Styles
   ======================================== */

.home-content {
  max-width: 1100px;
}

/* Hero */
.hero {
  text-align: center;
  padding: var(--space-2xl) 0;
  margin-bottom: var(--space-2xl);
}

.hero-icon {
  font-size: 5rem;
  margin-bottom: var(--space-lg);
}

.hero-title {
  font-size: 3rem;
  font-weight: 800;
  color: var(--color-text);
  margin-bottom: var(--space-md);
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--color-text-secondary);
  max-width: 600px;
  margin: 0 auto var(--space-xl);
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  gap: var(--space-md);
  justify-content: center;
  flex-wrap: wrap;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-sm) var(--space-xl);
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  cursor: pointer;
  border: none;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.btn-secondary {
  background: var(--color-bg-secondary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-primary);
}

/* Features Grid */
.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--space-lg);
  margin-bottom: var(--space-2xl);
}

.feature-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  transition: all var(--transition-fast);
}

.feature-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: var(--space-md);
}

.feature-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--space-sm);
}

.feature-description {
  font-size: 0.9375rem;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

/* Quick Links */
.quick-links {
  margin-bottom: var(--space-2xl);
}

.quick-links-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: var(--space-lg);
  text-align: center;
}

.quick-links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: var(--space-md);
}

.quick-link-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  text-decoration: none;
  transition: all var(--transition-fast);
}

.quick-link-card:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  transform: translateY(-2px);
}

.quick-link-icon {
  font-size: 1.5rem;
}

.quick-link-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text);
  text-align: center;
}

/* ========================================
   Search Modal
   ======================================== */

.search-modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: none;
  align-items: flex-start;
  justify-content: center;
  padding-top: 15vh;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.search-modal.active {
  display: flex;
}

.search-modal-content {
  width: 100%;
  max-width: 600px;
  background: var(--color-bg);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  margin: 0 var(--space-md);
}

.search-modal-header {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  border-bottom: 1px solid var(--color-border);
}

.search-modal-header .search-icon {
  color: var(--color-text-muted);
}

.search-modal-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 1rem;
  color: var(--color-text);
  outline: none;
}

.search-modal-input::placeholder {
  color: var(--color-text-muted);
}

.search-modal-close {
  background: transparent;
  border: none;
  cursor: pointer;
}

.search-modal-close kbd {
  padding: 4px 8px;
  font-size: 0.75rem;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
}

.search-modal-results {
  max-height: 400px;
  overflow-y: auto;
  padding: var(--space-sm) 0;
}

.search-result-item {
  display: block;
  padding: var(--space-md);
  text-decoration: none;
  transition: background var(--transition-fast);
}

.search-result-item:hover,
.search-result-item.selected {
  background: var(--color-bg-secondary);
}

.search-result-title {
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--space-xs);
}

.search-result-section {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-xs);
}

.search-result-excerpt {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.search-no-results {
  padding: var(--space-xl);
  text-align: center;
  color: var(--color-text-muted);
}

.search-modal-footer {
  display: flex;
  gap: var(--space-lg);
  justify-content: center;
  padding: var(--space-sm) var(--space-md);
  border-top: 1px solid var(--color-border);
  background: var(--color-bg-secondary);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.search-modal-footer kbd {
  padding: 2px 6px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  margin-right: var(--space-xs);
}

/* ========================================
   Mobile Styles
   ======================================== */

.mobile-menu-toggle {
  display: none;
  position: fixed;
  top: var(--space-md);
  left: var(--space-md);
  z-index: 101;
  width: 44px;
  height: 44px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  align-items: center;
  justify-content: center;
  color: var(--color-text);
  box-shadow: var(--shadow-sm);
}

.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

@media (max-width: 1024px) {
  .mobile-menu-toggle {
    display: flex;
  }

  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .sidebar-overlay.active {
    display: block;
  }

  .main-content {
    margin-left: 0;
  }

  .content {
    padding: var(--space-xl) var(--space-lg);
    padding-top: calc(var(--space-xl) + 60px);
  }

  .toc {
    display: none;
  }
}

@media (max-width: 640px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .prose h1 { font-size: 1.75rem; }
  .prose h2 { font-size: 1.5rem; }
  .prose h3 { font-size: 1.25rem; }

  .features {
    grid-template-columns: 1fr;
  }

  .hero-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .btn {
    width: 100%;
  }
}

/* ========================================
   Syntax Highlighting (Dracula Theme)
   ======================================== */

/* Chroma syntax highlighting overrides */
.chroma {
  background: var(--color-code-bg);
  color: #f8f8f2;
}

.chroma .err { color: #ff5555; }
.chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
.chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
.chroma .hl { display: block; width: 100%; background-color: #44475a; }
.chroma .lnt { margin-right: 0.4em; padding: 0 0.4em 0 0.4em; color: #6272a4; }
.chroma .ln { margin-right: 0.4em; padding: 0 0.4em 0 0.4em; color: #6272a4; }
.chroma .k { color: #ff79c6; }
.chroma .kc { color: #ff79c6; }
.chroma .kd { color: #8be9fd; font-style: italic; }
.chroma .kn { color: #ff79c6; }
.chroma .kp { color: #ff79c6; }
.chroma .kr { color: #ff79c6; }
.chroma .kt { color: #8be9fd; }
.chroma .na { color: #50fa7b; }
.chroma .nb { color: #8be9fd; font-style: italic; }
.chroma .nc { color: #50fa7b; }
.chroma .no { color: #bd93f9; }
.chroma .nd { color: #50fa7b; }
.chroma .ne { color: #50fa7b; }
.chroma .nf { color: #50fa7b; }
.chroma .nl { color: #8be9fd; font-style: italic; }
.chroma .nn { color: #f8f8f2; }
.chroma .nt { color: #ff79c6; }
.chroma .nv { color: #8be9fd; font-style: italic; }
.chroma .s { color: #f1fa8c; }
.chroma .sa { color: #f1fa8c; }
.chroma .sb { color: #f1fa8c; }
.chroma .sc { color: #f1fa8c; }
.chroma .dl { color: #f1fa8c; }
.chroma .sd { color: #f1fa8c; }
.chroma .s2 { color: #f1fa8c; }
.chroma .se { color: #f1fa8c; }
.chroma .sh { color: #f1fa8c; }
.chroma .si { color: #f1fa8c; }
.chroma .sx { color: #f1fa8c; }
.chroma .sr { color: #f1fa8c; }
.chroma .s1 { color: #f1fa8c; }
.chroma .ss { color: #f1fa8c; }
.chroma .m { color: #bd93f9; }
.chroma .mb { color: #bd93f9; }
.chroma .mf { color: #bd93f9; }
.chroma .mh { color: #bd93f9; }
.chroma .mi { color: #bd93f9; }
.chroma .il { color: #bd93f9; }
.chroma .mo { color: #bd93f9; }
.chroma .o { color: #ff79c6; }
.chroma .ow { color: #ff79c6; }
.chroma .p { color: #f8f8f2; }
.chroma .c { color: #6272a4; }
.chroma .ch { color: #6272a4; }
.chroma .cm { color: #6272a4; }
.chroma .c1 { color: #6272a4; }
.chroma .cs { color: #6272a4; }
.chroma .cp { color: #ff79c6; }
.chroma .cpf { color: #ff79c6; }
.chroma .gd { color: #ff5555; }
.chroma .ge { font-style: italic; }
.chroma .gi { color: #50fa7b; }
.chroma .gs { font-weight: bold; }
.chroma .gu { color: #6272a4; font-weight: bold; }

/* Print Styles */
@media print {
  .sidebar,
  .toc,
  .mobile-menu-toggle,
  .search-container,
  .page-nav {
    display: none;
  }

  .main-content {
    margin-left: 0;
  }

  .content {
    max-width: 100%;
    padding: 0;
  }
}
`
