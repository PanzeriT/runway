package registry

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/panzerit/runway/router"
)

// App defines the interface that all applications must implement
type App interface {
	Name() string
	Routes() map[string]router.HandlerFunc
	Initialize(ctx context.Context) error
	LoadTemplates() error
	Shutdown(ctx context.Context) error
	HealthCheck() error
}

// Registry manages all registered applications
type Registry struct {
	mu       sync.RWMutex
	apps     map[string]App
	handlers map[string]router.HandlerFunc
}

// Global registry instance
var (
	globalRegistry *Registry
	registryOnce   sync.Once
)

// GetRegistry returns the singleton registry instance
func GetRegistry() *Registry {
	registryOnce.Do(func() {
		globalRegistry = &Registry{
			apps:     make(map[string]App),
			handlers: make(map[string]router.HandlerFunc),
		}
	})
	return globalRegistry
}

// Register adds a app to the registry (called during init)
func (r *Registry) Register(app App) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := app.Name()
	if _, exists := r.apps[name]; exists {
		return fmt.Errorf("app %s already registered", name)
	}

	r.apps[name] = app
	log.Printf("Registered app: %s", name)
	return nil
}

// Initialize all registered apps
func (r *Registry) InitializeAll(ctx context.Context, router *router.Router) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, app := range r.apps {
		if err := app.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize app %s: %w", name, err)
		}

		if err := app.LoadTemplates(); err != nil {
			return fmt.Errorf("failed to load templates for app %s: %w", name, err)
		}

		// Register routes
		routes := app.Routes()
		for pattern, handler := range routes {
			fullPattern := fmt.Sprintf("/%s%s", name, pattern)
			r.handlers[fullPattern] = handler
			router.GET(fullPattern, handler)
			log.Printf("Registered route: %s", fullPattern)
		}
	}

	return nil
}

// GetHandler returns the handler for a given pattern
func (r *Registry) GetHandler(pattern string) (router.HandlerFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[pattern]
	return handler, exists
}

// ListApps returns all registered app names
func (r *Registry) ListApps() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.apps))
	for name := range r.apps {
		names = append(names, name)
	}
	return names
}

// GetApp returns a app by name (for health checks, etc.)
func (r *Registry) GetApp(name string) (App, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	app, exists := r.apps[name]
	return app, exists
}

// Shutdown gracefully shuts down all apps
func (r *Registry) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var lastErr error
	for name, app := range r.apps {
		if err := app.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down app %s: %v", name, err)
			lastErr = err
		}
	}
	return lastErr
}

// Convenience function for global registration
func Register(app App) error {
	return GetRegistry().Register(app)
}
