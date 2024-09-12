//go:build frontend

package frontend

import "embed"

//go:embed management/dist/*
var ManagementModule embed.FS
