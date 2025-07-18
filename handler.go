package runway

func (a *App) GET(path string, fn func(c Context) error) {
	a.server.GET(path, fn)
}
