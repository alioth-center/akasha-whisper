package api

import (
	"github.com/alioth-center/akasha-whisper/app/service"
)

func init() {
	CompatibleApi = compatibleApiImpl{service: service.NewCompatibleService()}
	ManagementApi = managementApiImpl{service: service.NewManagementService()}
}
