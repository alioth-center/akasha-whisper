package global

import "github.com/alioth-center/infrastructure/database/postgres"

var Config WhisperConfig

type WhisperConfig struct {
	HttpEngine  HttpEngineConfig  `yaml:"http_engine"`
	Logger      LogConfig         `yaml:"logger"`
	BloomFilter BloomFilterConfig `yaml:"bloom_filter"`
	Database    postgres.Config   `yaml:"database"`
	App         AppConfig         `yaml:"app"`
}

type HttpEngineConfig struct {
	ServeURL  string `yaml:"serve_url"`
	ServeAddr string `yaml:"serve_addr"`
}

type BloomFilterConfig struct {
	Enable     bool    `yaml:"enable"`
	FilterSize int     `yaml:"filter_size"`
	FalseRate  float64 `yaml:"false_rate"`
}

type LogConfig struct {
	LogToFile    bool   `yaml:"log_to_file"`
	LogSplit     bool   `yaml:"log_split"`
	LogDirectory string `yaml:"log_directory"`
	LogLevel     string `yaml:"log_level"`
}

type AppConfig struct {
	MaxToken        int    `yaml:"max_token"`
	ManagementToken string `yaml:"management_token"`
}
