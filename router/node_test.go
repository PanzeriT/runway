package router

import (
	"fmt"
	"net/http"
	"testing"
)

func routePrintingHandlerFunc(route string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Route: %s", route)
	}
}

// TODO: There is still a mess with the slahses. Clean-up needed.
func TestNode(t *testing.T) {
	r := New()

	route1 := "test/with/a/route/"
	r.GET(route1, routePrintingHandlerFunc(route1))
	route2 := "test/with/another/route/"
	r.GET(route2, routePrintingHandlerFunc(route2))

	routesFromNode := r.GetRoutes()
	if routesFromNode[0]+"/" != "//"+route1 {
		t.Errorf("cannot find route1: %s; found: %s", route1, routesFromNode[0])
	}
	if routesFromNode[1]+"/" != "//"+route2 {
		t.Errorf("cannot find route2: %s; found: %s", route2, routesFromNode[1])
	}
}
