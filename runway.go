package runway

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/asset"
	"github.com/panzerit/runway/handler"
	"github.com/panzerit/runway/model"
	"github.com/panzerit/runway/service"
	"github.com/panzerit/runway/template/page"
	"gorm.io/gorm"
)

type AppOption func(*App) *App

type App struct {
	name      string
	jwtSecret string
	port      string
	service   service.Service
	server    *echo.Echo
}

func init() {
	logger = NewRunwayLogger()

	modles := []any{
		model.User{},
		model.Role{},
	}

	for _, m := range modles {
		err := model.Register(m)
		if err != nil {
			Terminate(NewAppError(ErrCannotRegisterModel, err))
		}
	}
}

func New(name, jwtSecret string, db *gorm.DB) *App {
	MustMeetSecretCriteria(jwtSecret)

	db.Config.Logger = gormLogger{}

	server := echo.New()

	svc := service.New(db, model.GetRegisteredModels)

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

func (a *App) SetPort(port int) *App {
	a.port = fmt.Sprintf(":%d", port)
	return a
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

	handler.NewTableHandler(a.service, logger.Logger, a.name).Register(r)
}

func (a *App) Start() {
	s := http.Server{
		Addr:        a.port,
		Handler:     a.server,
		ReadTimeout: 30 * time.Second,
	}

	// run graceful shutdown in a separate goroutine
	done := make(chan bool, 1)

	// start server
	go gracefulShutdown(a.server.Server, done)

	fmt.Printf("Runway serving '%s' is running on port %s.\n", a.name, a.port)
	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		Exit(ServerError, err)
	}

	// wait for the graceful shutdown to complete
	<-done

	Exit(0, nil)
}

func (a *App) customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		//	andle the case where the page cannot be found
		if he.Code == 404 {
			err := page.NewHttpError(
				he.Code,
				page.WithMessage("Page Not Found"),
				page.WithDescription("The page you are looking for does not exist."),
			)
			Render(c, he.Code, page.Error(a.name, nil, err))
			return
		}

		// handle the case where the user is not authenticated
		if he.Internal.Error() == "missing value in cookies" {
			// send 401 if it was an HTMX request
			if c.Request().Header.Get("HX-Request") == "true" {
				c.NoContent(http.StatusUnauthorized)
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			return
		}
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

func gracefulShutdown(s *http.Server, done chan bool) {
	// create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // allow Ctrl+C to force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	// notify the main goroutine that the shutdown is complete
	done <- true
}
