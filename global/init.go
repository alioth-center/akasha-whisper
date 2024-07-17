package global

func Init() {
	initConfig()
	initLogger()
	initDatabase()
	initEngine()
}
