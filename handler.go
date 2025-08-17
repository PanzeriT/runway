package runway

import (
	"fmt"
	"net/http"

	"github.com/panzerit/runway/router"
)

func (a *Runway) Routes(w http.ResponseWriter, r *http.Request) {
	for _, method := range router.AllowedMethods {
		fmt.Fprintf(w, "\n\nRoutes for %s:\n", method)
		for i, route := range a.Router.GetRoutes(method) {
			fmt.Fprintf(w, "%d: %s\n", i, route)
		}
	}

	for _, method := range router.AllowedMethods {
		fmt.Fprintf(w, "\n\nNodes for %s:\n", method)
		for i, node := range a.Router.GetNodes(method) {
			fmt.Fprintf(w, "%d: %#v\n", i, *node)
		}
	}
}
