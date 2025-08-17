package router

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

// Router represents our custom router
type Router struct {
	tree map[method]*node
}

// New creates a new router instance
func New() *Router {
	r := &Router{
		tree: make(map[method]*node),
	}

	// initialize the trees for all allowed HTTP methods
	allowedMethods := []method{GET, POST} // TODO: add PUT, DELETE, PATCH, HEAD, OPTIONS
	for _, m := range allowedMethods {
		r.tree[m] = &node{}
	}

	return r
}

func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.Handle(GET, pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.Handle(POST, pattern, handler)
}

func (r *Router) Handle(method method, path string, handler HandlerFunc) {
	slog.Info("registering route", "method", method, "path", path)
	parts := strings.Split(strings.Trim(path, "/"), "/")

	r.tree[method].insert(parts, handler)

	routes := r.tree[method].getSubRoutes("/")
	slog.Info("currently registered routes", "routes", routes)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	parts := []string{}
	if path != "" {
		parts = strings.Split(path, "/")
	}

	root := r.tree[method(req.Method)]
	params := make(map[string]string)
	handler := root.search(parts, params)

	if handler != nil {
		// wrap request with params in context
		ctx := context.WithValue(req.Context(), "params", params)
		handler(w, req.WithContext(ctx))
		return
	}

	http.NotFound(w, req)
}

func (r *Router) GetRoutes() []string {
	return r.tree[GET].getSubRoutes("")
}
