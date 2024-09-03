package static

import (
	"embed"
	_ "embed"
)

//go:embed css
//go:embed img
//go:embed js
var Assets embed.FS
