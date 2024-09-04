package main

import (
	"github.com/alioth-center/akasha-whisper/api"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/service"
	"github.com/alioth-center/infrastructure/exit"

	"github.com/alioth-center/akasha-whisper/app/dao"
)

func main() {
	global.Init()
	service.InitService()

	dao.QueryTest()

	api.BindChatCompletion()
	api.BindListModel()
	api.BindCreateUser()
	global.Engine.ServeAsync(global.Config.ServeAt, make(chan struct{}, 1))

	exit.BlockedUntilTerminate()
}
