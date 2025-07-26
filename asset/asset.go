package asset

import "embed"

//go:embed favicon.ico css/* image/* js/*.min.js
var FS embed.FS
