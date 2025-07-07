package runway

import (
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/template/page"
)

func (a *App) dashboardHandler(c echo.Context) error {
	uc := c.(*data.UserContext)

	return Render(c, 200, page.Dashboard(a.Name, &uc.User, 0))
}
