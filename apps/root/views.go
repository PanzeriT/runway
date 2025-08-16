package root

import (
	"fmt"
	"net/http"
)

func (a *RootApp) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Runway.")
}
