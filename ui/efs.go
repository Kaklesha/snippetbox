package ui

import "embed"

//go:embed "static/css/*.css" "static/img" "static/js"
var Files embed.FS
