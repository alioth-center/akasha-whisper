package frontend

import "embed"

var (
	//go:embed management/dist/*
	ManagementModule embed.FS
)
