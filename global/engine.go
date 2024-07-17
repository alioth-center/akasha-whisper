package global

import "github.com/alioth-center/infrastructure/network/http"

var Engine *http.Engine

func initEngine() {
	Engine = http.NewEngine(Config.BaseUrl)
}
