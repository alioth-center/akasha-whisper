package global

import "github.com/alioth-center/infrastructure/logger"

var Logger logger.Logger

func initLogger() {
	Logger = logger.NewLoggerWithConfig(Config.Logger)
}
