package user

import (
	"log"

	"github.com/panzerit/runway/apps"
	"github.com/panzerit/runway/registry"
)

// Register this app during package initialization
func init() {
	if err := registry.Register(&UserApp{}); err != nil {
		log.Fatal("Failed to register UserApp:", err)
	}
}

type UserApp struct {
	apps.Base
}

func (a *UserApp) Name() string {
	return "user"
}
