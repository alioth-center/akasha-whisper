package global

import (
	"github.com/alioth-center/infrastructure/config"
	"github.com/alioth-center/infrastructure/database/postgres"
)

var Config WhisperConfig

type WhisperConfig struct {
	BaseUrl        string          `yaml:"base_url"`
	ServeAt        string          `yaml:"serve_at"`
	AdminToken     string          `yaml:"admin_token"`
	MaxToken       int             `yaml:"max_token"`
	DefaultBalance float64         `yaml:"default_balance"`
	Database       postgres.Config `yaml:"database"`
}

func initConfig() {
	loadErr := config.LoadConfig(&Config, "./config/app.yaml")
	if loadErr != nil {
		panic(loadErr)
	}
}
