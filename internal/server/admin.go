package admin

import "embed"

//go:embed template/*.tmpl
var Templates embed.FS
