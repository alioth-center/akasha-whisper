package router

import (
	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/infrastructure/network/http"
)

func serveBackend() {
	engine := http.NewEngine(global.Config.HttpEngine.ServeURL)

	engine.AddEndPoints(OpenAiCompatibleRouterGroup...)
	engine.AddEndPoints(FrameworkRouterGroup...)

	if global.Config.HttpEngine.EnableManagementApis {
		engine.AddEndPoints(ManagementRouterGroup...)
	}

	engine.ServeAsync(global.Config.HttpEngine.ServeAddr, make(chan struct{}, 1))
}

func init() {
	go serveBackend()
}
