package router

import (
	"strings"
)

type node struct {
	part     string // e.g. "users", ":id", "*filepath"
	children []*node
	handler  HandlerFunc
	isParam  bool // true if part starts with ":"
	isWild   bool // true if part starts with "*"
}

func (n *node) insert(parts []string, handler HandlerFunc) {
	// check if it's the last element or the root page
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		n.handler = handler
		return
	}

	// Check if there is a child node for the next part
	part := parts[0]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:    part,
			isParam: len(part) > 0 && part[0] == ':',
			isWild:  len(part) > 0 && part[0] == '*',
		}
		// TODO len(part) is just a patch for now
		n.children = append(n.children, child)
	}

	child.insert(parts[1:], handler)
}

func (n *node) search(parts []string, params map[string]string) HandlerFunc {
	if len(parts) == 0 {
		return n.handler
	}

	part := parts[0]

	for _, child := range n.children {
		if child.part == part || child.isParam || child.isWild {
			if child.isParam {
				params[child.part[1:]] = part // store param without ":"
			}
			if child.isWild {
				params[child.part[1:]] = strings.Join(parts, "/") // capture rest
				return child.handler
			}
			handler := child.search(parts[1:], params)
			if handler != nil {
				return handler
			}
		}
	}
	return nil
}

func (n *node) getSubRoutes(base string) []string {
	routes := []string{}

	path := base + "/" + n.part
	if n.handler != nil {
		routes = append(routes, strings.TrimRight(path, "/"))
	}

	for _, child := range n.children {
		routes = append(routes, child.getSubRoutes(path)...)
	}

	return routes
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isParam || child.isWild {
			return child
		}
	}
	return nil
}
