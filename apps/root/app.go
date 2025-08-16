package root

import (
	"log"

	"github.com/panzerit/runway/apps"
	"github.com/panzerit/runway/registry"
)

// Register this app during package initialization
func init() {
	if err := registry.Register(&RootApp{}); err != nil {
		log.Fatal("Failed to register RootApp:", err)
	}
}

type RootApp struct {
	apps.Base
}

func (a *RootApp) Name() string {
	return "root"
}
