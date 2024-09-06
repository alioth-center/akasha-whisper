package frontend

import "embed"

var (
	//go:embed management/*
	ManagementModule embed.FS
)
