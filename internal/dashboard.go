package admin

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/internal/template"
)

func dashboardHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name

	content := template.Page("Test", template.Dashboard(name))
	return content.Render(c.Request().Context(), c.Response().Writer)
}
