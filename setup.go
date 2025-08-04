package runway

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/template/page"
)

func (r *Runway) getSetup(c echo.Context) error {
	if r.setupComplete {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	return Render(c, http.StatusOK, page.Setup(r.name, nil))
}
