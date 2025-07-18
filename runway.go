package runway

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/asset"
	"github.com/panzerit/runway/model"
	"github.com/panzerit/runway/service"
	"github.com/panzerit/runway/template/page"
	"gorm.io/gorm"
)

type AppOption func(*App) *App

type App struct {
	name      string
	jwtSecret string
	service   service.Service
	server    *echo.Echo
}

func init() {
	logger = NewRunwayLogger()

	registeredModels = make(map[string]reflect.Type)

	MustRegisterModel(model.User{})
	MustRegisterModel(model.Role{})
}

func New(name, jwtSecret string, db *gorm.DB) *App {
	MustMeetSecretCriteria(jwtSecret)

	db.Config.Logger = gormLogger{}

	server := echo.New()

	svc := service.New(db, GetRegisteredModels)

	app := &App{
		name:      name,
		jwtSecret: jwtSecret,
		service:   svc,
		server:    server,
	}

	app.server.HTTPErrorHandler = app.customHTTPErrorHandler
	app.server.StaticFS("/", echo.MustSubFS(asset.FS, "./"))
	app.addPublicRoutes()
	app.addPrivateRoutes()

	return app
}

func (a *App) addPublicRoutes() {
	a.server.GET("/", a.introHandler)

	a.server.GET("/login", a.getLoginHandler)
	a.server.POST("/login", a.postLoginHandler)
}

func (a *App) addPrivateRoutes() {
	r := a.server.Group("/admin")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  []byte(a.jwtSecret),
		TokenLookup: "cookie:token",
	}

	r.Use(echojwt.WithConfig(config))
	r.Use(JWTExtractor)

	r.GET("", a.dashboardHandler)
	r.GET("/logout", a.logoutHandler)

	r.GET("/model/:model", a.tableHandler)
	r.DELETE("/model/:model/:id", a.deleteRowHandler)
}

func (a *App) Start() {
	addr := ":1291"

	s := http.Server{
		Addr:        addr,
		Handler:     a.server,
		ReadTimeout: 30 * time.Second,
	}

	fmt.Printf("Runway serving '%s' is running on port %s.\n", a.name, addr)
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
			Render(c, he.Code, page.Error(a.name, nil, err))
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

	if err := Render(c, code, page.Error(a.name, nil, page.NewHttpError(code))); err != nil {
		c.Logger().Error(err)
	}
}

func (a *App) introHandler(c echo.Context) error {
	return Render(c, 200, page.Intro(a.name, nil, time.Now().Year()))
}

func MustMeetSecretCriteria(secret string) {
	if len(secret) < 16 {
		Terminate(ErrSecretToShort)
	}
}
