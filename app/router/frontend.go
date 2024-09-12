//go:build frontend

package router

import "github.com/alioth-center/akasha-whisper/frontend"

func init() {
	_ = frontend.ManagementModule
}
