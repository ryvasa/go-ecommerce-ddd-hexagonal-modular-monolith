package casbin

import "embed"

//go:embed model.conf policy.csv
var FS embed.FS
