package frontend

import "embed"

// ManagementModule is the embedded frontend module for management.
//
// nolint
//
//go:embed management/dist/*
var ManagementModule embed.FS
