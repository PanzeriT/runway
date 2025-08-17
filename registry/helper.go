package registry

import (
	"reflect"
	"strings"
)

func InferName(obj any) string {
	fullName := reflect.TypeOf(obj).String()

	// Remove the package name (and the 'App', if it has been used in the name)
	name := strings.TrimSuffix(strings.Split(fullName, ".")[1], "App")

	lowerName := strings.ToLower(name)

	return lowerName
}

func InferPath(obj any) string {
	path := InferName(obj)

	// Handle special case for root app
	if path == "root" {
		return ""
	}

	return path
}
