package main

import (
	"github.com/alioth-center/akasha-whisper/api"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/service"
	"github.com/alioth-center/infrastructure/exit"
)

func main() {
	global.Init()
	service.InitService()

	api.BindChatCompletion()
	api.BindListModel()
	api.BindCreateUser()
	global.Engine.ServeAsync(global.Config.ServeAt, make(chan struct{}))

	exit.BlockedUntilTerminate()
}
