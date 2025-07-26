package handler

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/template/page"
)

type Handler interface {
	Register(e *echo.Group)
}

func render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func render404(c echo.Context, appName string) {
	code := 404
	err := page.NewHttpError(
		code,
		page.WithMessage("Page Not Found"),
		page.WithDescription("The page you are looking for does not exist."),
	)
	render(c, code, page.Error(appName, nil, err))
}

type param interface {
	~int | ~string | ~bool
}

func getParamWithDefault[T param](c echo.Context, name string, value T) T {
	v := c.Param(name)
	if v == "" {
		return value
	}

	var result T
	switch any(result).(type) {
	case int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return value
		}
		return any(i).(T)
	case bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return value
		}
		return any(b).(T)
	case string:
		return any(v).(T)
	default:
		return value
	}
}
