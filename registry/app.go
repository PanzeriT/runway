package registry

import (
	"context"

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
