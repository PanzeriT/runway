package main

import (
	"context"
	"log"
	"net/http"
)

func HomeHandler(ctx context.Context, req CreateUserRequest) (Responder, error) {
	return HTMLResponse(http.StatusOK, "<h1>Welcome to My API</h1>"), nil
}

type CreateUserRequest struct {
	Name  string `form:"name" validate:"required,min=2"`
	Email string `form:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserHandler(ctx context.Context, req CreateUserRequest) (Responder, error) {
	return JSONResponse(http.StatusOK, UserResponse{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
	}), nil
}

func HealthHandler(ctx context.Context) (Responder, error) {
	return NoContentResponse(), nil
}

func main() {
	router := NewRouter()

	router.GET("/", HomeHandler)
	router.POST("/users", UserHandler)
	router.DELETE("/users", UserHandler)
	router.GET("/health", HealthHandler)
	router.GET("/openapi.json", OpenAPIHandler(router))

	router.PrintRoutes()

	if err := http.ListenAndServe(":1291", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func OpenAPIHandler(r *Router) Handler {
	return func(ctx context.Context, req CreateUserRequest) (Responder, error) {
		spec, err := r.GenerateOpenAPIJSON()
		if err != nil {
			return nil, err
		}

		return JSONResponse(http.StatusOK, spec), nil
	}
}
