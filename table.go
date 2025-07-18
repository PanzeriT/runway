package runway

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/template/page"
)

func (a *App) tableHandler(c echo.Context) error {
	uc := c.(*data.UserContext)

	return Render(c, http.StatusOK, page.Table(a.name, &uc.User))
}
