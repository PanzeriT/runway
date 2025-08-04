package runway

import (
	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/template/page"
)

func (a *Runway) dashboardHandler(c echo.Context) error {
	uc := c.(*data.UserContext)

	for _, m := range a.service.GetModelNames() {
		logger.Info("model", "name", m)
	}

	return Render(c, 200, page.Dashboard(a.name, &uc.User, 0, a.service.GetModelNames()))
}
