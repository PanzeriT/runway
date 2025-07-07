package runway

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/asset"
	"github.com/panzerit/runway/template/page"
	"gorm.io/gorm"
)

type App struct {
	Name      string
	JWTSecret string
	DB        *gorm.DB
	Server    *echo.Echo
}

func New(name, jwtSecret string, db *gorm.DB) *App {
	server := echo.New()

	MustMeetSecretCriteria(jwtSecret)

	app := &App{
		Name:      name,
		JWTSecret: jwtSecret,
		Server:    server,
	}

	app.Server.HTTPErrorHandler = app.customHTTPErrorHandler
	app.Server.StaticFS("/", echo.MustSubFS(asset.FS, "./"))
	app.addPublicRoutes()
	app.addPrivateRoutes()

	return app
}

func (a *App) addPublicRoutes() {
	a.Server.GET("/", a.introHandler)

	a.Server.GET("/login", a.getLoginHandler)
	a.Server.POST("/login", a.postLoginHandler)
}

func (a *App) addPrivateRoutes() {
	r := a.Server.Group("/admin")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  []byte(a.JWTSecret),
		TokenLookup: "cookie:token",
	}

	r.Use(echojwt.WithConfig(config))
	r.Use(JWTExtractor)

	r.GET("", a.dashboardHandler)
	r.GET("/logout", a.logoutHandler)

	r.GET("/table", a.tableHandler)
}

func (a *App) Start() {
	addr := ":1291"

	s := http.Server{
		Addr:        addr,
		Handler:     a.Server,
		ReadTimeout: 30 * time.Second,
	}

	fmt.Printf("Runway serving '%s' is running on port %s.\n", a.Name, addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (a *App) customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		// redirect to login page if the cookie is not valid
		if he.Code == 404 {
			err := page.NewHttpError(
				he.Code,
				page.WithMessage("Page Not Found"),
				page.WithDescription("The page you are looking for does not exist."),
			)
			Render(c, he.Code, page.Error(a.Name, nil, err))
			return
		}
		if he.Internal.Error() == "missing value in cookies" {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			fmt.Println("bad cookie")
		}
		code = he.Code
		return
	}

	// log all unexpected errors
	c.Logger().Error(err)

	if err := Render(c, code, page.Error(a.Name, nil, page.NewHttpError(code))); err != nil {
		c.Logger().Error(err)
	}
}

func (a *App) introHandler(c echo.Context) error {
	return Render(c, 200, page.Intro(a.Name, nil, time.Now().Year()))
}

func MustMeetSecretCriteria(secret string) {
	if len(secret) < 16 {
		Terminate(ErrSecretToShort)
	}
}
