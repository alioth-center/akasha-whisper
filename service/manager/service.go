package manager

import (
	"github.com/alioth-center/akasha-whisper/dao"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
)

type AkashaManagerSrv interface {
	CreateUser(ctx http.Context[*model.CreateUserRequest, *model.BaseResponse])
}

func NewAkashaManagerSrv(db dao.DatabaseAccessor, log logger.Logger) AkashaManagerSrv {
	return &akashaManagerSrvImpl{db: db, log: log}
}

type akashaManagerSrvImpl struct {
	db  dao.DatabaseAccessor
	log logger.Logger
}
