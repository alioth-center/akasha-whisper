package global

import (
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
)

var (
	Logger logger.Logger
	Client http.Client
)
