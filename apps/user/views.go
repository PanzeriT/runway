package user

import (
	"fmt"
	"net/http"
)

func (a *UserApp) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Runway.")
}
