package runway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/panzerit/runway/template/page"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (a *App) getLoginHandler(c echo.Context) error {
	// if there is a token cookie, check if it's valid
	// and if so, redirect to dashboard
	tokenCookie, err := c.Cookie("token")
	if err == nil && tokenCookie != nil && tokenCookie.Value != "" {
		// Validate JWT token
		jwtToken, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(a.JWTSecret), nil
		})
		if err == nil && jwtToken.Valid {
			return c.Redirect(http.StatusPermanentRedirect, "/admin")
		}
	}

	return Render(c, http.StatusOK, page.Login(a.Name, nil))
}

type jwtCustomClaims struct {
	Name           string `json:"name"`
	Admjwtcustomin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type postLoginHandlerRequest struct {
	UserName string `form:"username"`
	Password string `form:"password"`
}

type postLoginHandlerResponse struct {
	Token string `json:"token"`
}

func (a *App) postLoginHandler(c echo.Context) error {
	var req postLoginHandlerRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	if req.UserName != "user" || req.Password != "pw" {
		// TODO: wrong pw causes a nil pointer dereference in the template
		// Replace with real user
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		"John Doe",
		true,
		jwt.RegisteredClaims{
			Subject:   "john@doe.com",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(a.JWTSecret))
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   false, // set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return c.Redirect(http.StatusMovedPermanently, "/admin")
}

func (a *App) logoutHandler(c echo.Context) error {
	// clear the token cookie
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})

	return Render(c, http.StatusOK, page.Logout(a.Name, nil))
}
