package runway

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// TODO: remove this and use the Render function in ./handler/handler.go

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
