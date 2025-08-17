package apps

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/panzerit/runway/router"
)

type Base struct {
	Templates   map[string]*template.Template
	Initialized bool
}

func (b *Base) Initialize(ctx context.Context) error {
	log.Println("Initializing UserApp...")

	// Initialize database connections, etc.
	b.Initialized = true
	return nil
}

func (b *Base) Routes() map[string]router.HandlerFunc {
	return map[string]router.HandlerFunc{}
}

func (b *Base) Shutdown(ctx context.Context) error {
	log.Println("Shutting down UserApp...")
	// Clean up resources
	b.Initialized = false
	return nil
}

func (b *Base) HealthCheck() error {
	if !b.Initialized {
		return fmt.Errorf("UserApp not initialized")
	}
	return nil
}

// LoadTemplates combines layouts, pages, and partials into a map of templates
func (b *Base) LoadTemplates() error {
	return nil
	slog.Info("Loading templates...")
	// Get all template files
	layoutFiles, _ := filepath.Glob("apps/root/templates/layouts/*.html")
	partialFiles, _ := filepath.Glob("apps/root/templates/partials/*.html")
	pageFiles, _ := filepath.Glob("apps/root/templates/pages/*.html")

	b.Templates = make(map[string]*template.Template, len(pageFiles))

	// Combine base files (layouts + partials)
	baseFiles := append(layoutFiles, partialFiles...)

	// Create a template for each page
	for _, pageFile := range pageFiles {
		// Extract page name from filename
		pageName := filepath.Base(pageFile)
		pageName = pageName[:len(pageName)-len(filepath.Ext(pageName))]

		// Combine all files: layouts + partials + current page
		allFiles := append(baseFiles, pageFile)

		// Parse all files together
		tmpl, err := template.ParseFiles(allFiles...)
		if err != nil {
			return err
		}

		_ = tmpl
		slog.Info("  - loaded template", "page", pageName)
		// a.Templates[pageName] = tmpl
		b.Templates[pageName] = template.New("<h1>" + pageName + "</h1>") // Placeholder for actual template parsing

	}

	return nil
}
