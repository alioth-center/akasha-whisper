package router

import "github.com/alioth-center/infrastructure/network/http"

func init() {
	engine := http.NewEngine("")

	engine.AddEndPoints(OpenAiCompatibleRouterGroup...)
	engine.AddEndPoints(ManagementRouterGroup...)

	engine.ServeAsync("0.0.0.0:8881", make(chan struct{}, 1))
}
