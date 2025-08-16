package router

import (
	"log/slog"
	"net/http"
	"strings"
)

// HandlerFunc represents a custom handler function
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Route represents a single route
type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// Router represents our custom router
type Router struct {
	routes []Route
}

// New creates a new router instance
func New() *Router {
	return &Router{
		routes: make([]Route, 0),
	}
}

// GET registers a GET route
func (r *Router) GET(path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  "GET",
		Path:    path,
		Handler: handler,
	})
}

// POST registers a POST route
func (r *Router) POST(path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  "POST",
		Path:    path,
		Handler: handler,
	})
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	slog.Info("Received request", req.Method, req.URL.Path)
	// Find matching route
	for _, route := range r.routes {
		if r.matchRoute(route, req) {
			route.Handler(w, req)
			return
		}
	}

	// No route found - return 404
	http.NotFound(w, req)
}

// matchRoute checks if a route matches the request
func (r *Router) matchRoute(route Route, req *http.Request) bool {
	// Check HTTP method
	if route.Method != req.Method {
		return false
	}

	// Check path
	return r.matchPath(route.Path, req.URL.Path)
}

// matchPath checks if the route path matches the request path
func (r *Router) matchPath(routePath, requestPath string) bool {
	// Clean paths (remove trailing slashes except for root)
	routePath = r.cleanPath(routePath)
	requestPath = r.cleanPath(requestPath)

	return routePath == requestPath
}

// cleanPath normalizes a path
func (r *Router) cleanPath(path string) string {
	if path == "" {
		return "/"
	}

	// Remove trailing slash unless it's root
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}
