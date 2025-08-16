package root

import (
	"github.com/panzerit/runway/router"
)

func (a *RootApp) Routes() map[string]router.HandlerFunc {
	return map[string]router.HandlerFunc{
		"/": a.Get,
	}
}
