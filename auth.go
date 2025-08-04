package runway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/panzerit/runway/template/page"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (a *Runway) getLoginHandler(c echo.Context) error {
	// check if the setup is complete
	if a.setupComplete == false {
		c.Redirect(http.StatusTemporaryRedirect, "/setup")
	}

	// show the dashboard, if there is a valid JWT token in the cookie
	tokenCookie, err := c.Cookie("token")
	if err == nil && tokenCookie != nil && tokenCookie.Value != "" {
		jwtToken, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(a.jwtSecret), nil
		})
		if err == nil && jwtToken.Valid {
			return c.Redirect(http.StatusPermanentRedirect, "/admin")
		}
	}

	// in all other cases, show the login page
	return Render(c, http.StatusOK, page.Login(a.name, nil))
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

func (a *Runway) postLoginHandler(c echo.Context) error {
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

	t, err := token.SignedString([]byte(a.jwtSecret))
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

func (a *Runway) logoutHandler(c echo.Context) error {
	// clear the token cookie
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})

	return Render(c, http.StatusOK, page.Logout(a.name, nil))
}
