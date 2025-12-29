package casbin

import "embed"

//go:embed model.conf policy.csv
var Files embed.FS
