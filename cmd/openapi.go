package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Generate OpenAPI spec (bonus feature)
func (r *Router) GenerateOpenAPIJSON() (any, error) {
	spec := map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title":   "My API",
			"version": "1.0.0",
		},
		"paths": make(map[string]any),
	}

	paths := spec["paths"].(map[string]any)

	for _, route := range r.routes {
		if paths[route.Path] == nil {
			paths[route.Path] = make(map[string]any)
		}

		pathItem := paths[route.Path].(map[string]any)
		pathItem[strings.ToLower(route.Method)] = map[string]any{
			"summary": fmt.Sprintf("%s %s", route.Method, route.Path),
			"responses": map[string]any{
				"200": map[string]any{
					"description": "Success",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": r.typeToSchema(route.Handler.outputType),
						},
					},
				},
			},
		}
	}

	return spec, nil
}

func (r *Router) typeToSchema(t reflect.Type) map[string]any {
	switch t.Kind() {
	case reflect.Struct:
		properties := make(map[string]any)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = strings.ToLower(field.Name)
			}
			properties[jsonTag] = r.typeToSchema(field.Type)
		}
		return map[string]any{
			"type":       "object",
			"properties": properties,
		}
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Int, reflect.Int64:
		return map[string]any{"type": "integer"}
	default:
		return map[string]any{"type": "string"}
	}
}
