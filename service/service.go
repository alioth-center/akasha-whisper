package service

import (
	"github.com/alioth-center/akasha-whisper/dao"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/service/agent"
	"github.com/alioth-center/akasha-whisper/service/manager"
)

var (
	AgentService   agent.AkashaAgentSrv
	ManagerService manager.AkashaManagerSrv
)

func InitService() {
	accessor := dao.NewDatabaseAccessor(global.Database)
	AgentService = agent.NewAkashaAgentSrv(accessor, global.Logger)
	ManagerService = manager.NewAkashaManagerSrv(accessor, global.Logger)
}
