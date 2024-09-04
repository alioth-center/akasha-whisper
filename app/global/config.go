package global

import "github.com/alioth-center/infrastructure/database/postgres"

var Config WhisperConfig

type WhisperConfig struct {
	ServeURL  string          `yaml:"serve_url"`
	ServeAddr string          `yaml:"serve_addr"`
	Database  postgres.Config `yaml:"database"`
	LogDir    string          `yaml:"log_dir"`
}
