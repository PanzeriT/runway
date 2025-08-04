package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
)

type handlerWrapper struct {
	handlerValue reflect.Value
	inputType    reflect.Type
	outputType   reflect.Type
}

func (hw *handlerWrapper) call(ctx context.Context, input any) (Responder, error) {
	// Call the actual handler method with proper types
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(input),
	}

	results := hw.handlerValue.Call(args)

	// Extract output and error
	output := results[0].Interface()
	errValue := results[1]

	var err error
	if !errValue.IsNil() {
		err = errValue.Interface().(error)
	}

	return output.(Responder), err
}

// Route definition
type Route struct {
	Method  string
	Path    string
	Handler *handlerWrapper
}

// Router
type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{
		routes: make([]Route, 0),
	}
}

// Register POST route
func (r *Router) POST(path string, handler any) {
	r.registerRoute("POST", path, handler)
}

// Register GET route
func (r *Router) GET(path string, handler any) {
	r.registerRoute("GET", path, handler)
}

func (r *Router) registerRoute(method, path string, handler any) {
	handlerValue := reflect.ValueOf(handler)
	kind := handlerValue.Kind()

	// Check, if the handler is a function
	if kind != reflect.Func {
		panic("handler must be a function, not " + kind.String() + handlerValue.String())
	}

	var hasInput bool
	functionType := reflect.TypeOf(handler)
	if functionType.NumIn() == 2 {
		hasInput = true
	}
	if (functionType.NumIn() != 2 && functionType.NumIn() != 1) || functionType.NumOut() != 2 {
		panic(`handle method has an invalid signature. Use:
 - func(ctx context.Context) (Responder, error)
 - func(ctx context.Context, input I) (Responder, error)`)
	}

	// Extract input and output types
	var inputType reflect.Type
	if hasInput {
		inputType = functionType.In(1) // Skip receiver and context
	}
	outputType := functionType.Out(0) // First return value Validate error return type
	errorType := functionType.Out(1)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !errorType.Implements(errorInterface) {
		panic("Handle method must return error as second value")
	}

	wrapper := &handlerWrapper{
		handlerValue: handlerValue,
		inputType:    inputType,
		outputType:   outputType,
	}

	route := Route{
		Method:  method,
		Path:    path,
		Handler: wrapper,
	}

	r.routes = append(r.routes, route)
}

func (r *Router) PrintRoutes() {
	getInputLabel := func(input reflect.Type) string {
		if input != nil {
			return input.Name()
		}

		return "[no input]"
	}

	getValues := func(route Route) []string {
		return []string{
			route.Method,
			route.Path,
			getInputLabel(route.Handler.inputType),
			// TODO: set the correct output type
			reflect.TypeOf(route.Handler.outputType).Name(),
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	//
	// Headers
	headers := []string{"Method", "Path", "Input", "Output"}

	// Calculate max width for each column, with special case for input
	maxWidths := make([]int, len(headers))
	maxWidths[2] = 10

	// Check header lengths
	for i, header := range headers {
		maxWidths[i] = len(header)
	}

	// Check data lengths
	for _, route := range r.routes {
		for i, value := range getValues(route) {
			if len(value) > maxWidths[i] {
				maxWidths[i] = len(value)
			}
		}
	}

	// Print headers
	for i, header := range headers {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, header)
	}
	fmt.Fprintln(w)

	// Print separators with correct length
	for i, width := range maxWidths {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, strings.Repeat("-", width))
	}
	fmt.Fprintln(w)

	// Print data
	for _, route := range r.routes {
		vs := getValues(route)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", vs[0], vs[1], vs[2], vs[3])
	}

	// Flush to ensure it gets written
	w.Flush()
}

// HTTP handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find matching route
	route := r.findRoute(req.Method, req.URL.Path)
	if route == nil {
		http.NotFound(w, req)
		return
	}

	// Create input struct instance
	inputValue := reflect.New(route.Handler.inputType).Elem()
	input := inputValue.Addr().Interface()

	// Parse request into input struct
	if err := r.parseRequest(req, input); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call handler
	ctx := req.Context()
	output, err := route.Handler.call(ctx, inputValue.Interface())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	output.Respond(w)
}

// Find route by method and path
func (r *Router) findRoute(method, path string) *Route {
	for i, route := range r.routes {
		if route.Method == method && route.Path == path {
			return &r.routes[i]
		}
	}
	return nil
}

// Parse HTTP request into struct
func (r *Router) parseRequest(req *http.Request, input any) error {
	inputValue := reflect.ValueOf(input).Elem()
	inputType := inputValue.Type()

	// Parse form data for POST requests
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			return err
		}
	}

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldValue := inputValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Get form tag
		formTag := field.Tag.Get("form")
		if formTag == "" {
			formTag = strings.ToLower(field.Name)
		}

		// Get value from request
		var value string
		if req.Method == "POST" {
			value = req.Form.Get(formTag)
		} else {
			value = req.URL.Query().Get(formTag)
		}

		if value == "" {
			continue
		}

		// Set field value based on type
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int, reflect.Int64:
			if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
				fieldValue.SetInt(intVal)
			}
		case reflect.Float64:
			if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				fieldValue.SetFloat(floatVal)
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(value); err == nil {
				fieldValue.SetBool(boolVal)
			}
		}
	}

	return nil
}
