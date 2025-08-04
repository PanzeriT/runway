package runway

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (r *Runway) GET(path string, fn func(c Context) error) {
	r.server.GET(path, fn)
}

func (r *Runway) Group(path string) *echo.Group {
	return r.server.Group(path)
}

func (r *Runway) StaticFS(path string, fs fs.FS) {
	r.Group(path).Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       path,
		Filesystem: http.FS(fs),
	}))
}
