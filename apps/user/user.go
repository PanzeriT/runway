package users

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/panzerit/runway/registry"
)

type UserApp struct {
	initialized bool
}

func (u *UserApp) Name() string { return "users" }

func (u *UserApp) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/":       u.listUsers,
		"/create": u.createUser,
	}
}

func (u *UserApp) Initialize(ctx context.Context) error {
	log.Println("Initializing UserApp...")
	// Initialize database connections, etc.
	u.initialized = true
	return nil
}

func (u *UserApp) Shutdown(ctx context.Context) error {
	log.Println("Shutting down UserApp...")
	// Clean up resources
	u.initialized = false
	return nil
}

func (u *UserApp) HealthCheck() error {
	if !u.initialized {
		return fmt.Errorf("UserApp not initialized")
	}
	return nil
}

func (u *UserApp) listUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User list from UserApp")
}

func (u *UserApp) createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create user in UserApp")
}

// Register this sub-app during package initialization
func init() {
	if err := registry.Register(&UserApp{}); err != nil {
		log.Fatal("Failed to register UserApp:", err)
	}
}
