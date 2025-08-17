package user

import (
	"github.com/panzerit/runway/router"
)

func (a *UserApp) Routes() map[string]router.HandlerFunc {
	return map[string]router.HandlerFunc{
		"/":      a.Get,
		"/login": a.Get,
	}
}
