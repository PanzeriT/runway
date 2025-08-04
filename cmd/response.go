package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter)
}

type Response struct {
	statusCode int
}

type htmlResponse struct {
	Response
	html string
}

func HTMLResponse(statusCode int, html string) *htmlResponse {
	return &htmlResponse{
		Response: Response{statusCode: statusCode},
		html:     html,
	}
}

func (r *htmlResponse) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(r.statusCode)

	// Write the HTML content
	if _, err := fmt.Fprint(w, r.html); err != nil {
		log.Println("Error writing HTML response:", err)
	}
}

type jsonResponse struct {
	Response
	data any
}

func JSONResponse(statusCode int, data any) *jsonResponse {
	return &jsonResponse{
		Response: Response{statusCode: statusCode},
		data:     data,
	}
}

func (r *jsonResponse) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(r.data)
	if err != nil {
		fmt.Fprintf(w, `{"error":"internal server error"}`)
		return
	}

	fmt.Fprintf(w, string(j))
}

type noContentResponse struct {
	Response
}

func NoContentResponse() *noContentResponse {
	return &noContentResponse{
		Response: Response{statusCode: http.StatusNoContent},
	}
}

func (r *noContentResponse) Respond(w http.ResponseWriter) {
	w.WriteHeader(r.statusCode)
}

func (r *noContentResponse) WithStatusCode(statusCode int) *noContentResponse {
	r.statusCode = statusCode
	return r
}
