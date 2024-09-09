package locales

import (
	"embed"
)

//go:embed en
//go:embed de
var Content embed.FS
