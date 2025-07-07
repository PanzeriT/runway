package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetProject(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", map[string]interface{}{
		"Title":           "Runway Admin Portal",
		"ContentTemplate": "project",
	})
}

func GetSubscription(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", map[string]interface{}{
		"Title":           "Runway Admin Portal",
		"ContentTemplate": "subscription",
	})
}
