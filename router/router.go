package router

import (
	"context"
	"net/http"
	"strings"
)

var AllowedMethods = []method{GET, POST} // TODO: add PUT, DELETE, PATCH, HEAD, OPTIONS

type Router struct {
	tree map[method]*node
}

// New creates a new router instance
func New() *Router {
	r := &Router{
		tree: make(map[method]*node),
	}

	// initialize the trees for all allowed HTTP methods
	for _, m := range AllowedMethods {
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
	parts := strings.Split(strings.Trim(path, "/"), "/")
	r.tree[method].insert(parts, handler)
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

func (r *Router) GetRoutes(method method) []string {
	paths := []string{}
	root := r.tree[method]

	var collect func(n *node, path string)
	collect = func(n *node, path string) {
		if n == nil {
			return
		}
		if n.handler != nil {
			paths = append(paths, path)
		}
		for _, child := range n.children {
			collect(child, path+"/"+child.part)
		}
	}
	collect(root, "")
	return paths
}

func (r *Router) GetNodes(method method) []*node {
	nodes := []*node{}
	root := r.tree[method]

	// Recursive function to flatten the tree
	var collect func(n *node)
	collect = func(n *node) {
		if n == nil {
			return
		}
		nodes = append(nodes, n)
		for _, child := range n.children {
			collect(child)
		}
	}
	collect(root)

	return nodes
}
