package runway

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/data"
)

func JWTExtractor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		name := claims.Name

		cc := &data.UserContext{
			Context: c,
			User: data.User{
				Name:  name,
				Email: "test@test.com",
			},
		}

		return next(cc)
	}
}
