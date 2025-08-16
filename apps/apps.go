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

func (a *Base) Initialize(ctx context.Context) error {
	log.Println("Initializing UserApp...")

	// Initialize database connections, etc.
	a.Initialized = true
	return nil
}

func (a *Base) Routes() map[string]router.HandlerFunc {
	return map[string]router.HandlerFunc{}
}

func (a *Base) Shutdown(ctx context.Context) error {
	log.Println("Shutting down UserApp...")
	// Clean up resources
	a.Initialized = false
	return nil
}

func (a *Base) HealthCheck() error {
	if !a.Initialized {
		return fmt.Errorf("UserApp not initialized")
	}
	return nil
}

// LoadTemplates combines layouts, pages, and partials into a map of templates
func (a *Base) LoadTemplates() error {
	slog.Info("Loading templates...")
	// Get all template files
	layoutFiles, _ := filepath.Glob("apps/root/templates/layouts/*.html")
	partialFiles, _ := filepath.Glob("apps/root/templates/partials/*.html")
	pageFiles, _ := filepath.Glob("apps/root/templates/pages/*.html")

	a.Templates = make(map[string]*template.Template, len(pageFiles))

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

		slog.Info("  - loaded template", "page", pageName)
		a.Templates[pageName] = tmpl

	}

	return nil
}
