package global

var Config WhisperConfig

type WhisperConfig struct {
	HttpEngine  HttpEngineConfig  `yaml:"http_engine"`
	Logger      LogConfig         `yaml:"logger"`
	BloomFilter BloomFilterConfig `yaml:"bloom_filter"`
	Database    DatabaseConfig    `yaml:"database"`
	App         AppConfig         `yaml:"app"`
}

type HttpEngineConfig struct {
	ServeURL             string `yaml:"serve_url"`
	ServeAddr            string `yaml:"serve_addr"`
	EnableManagementApis bool   `yaml:"enable_management_apis"`
}

type BloomFilterConfig struct {
	Enable     bool    `yaml:"enable"`
	FilterSize int     `yaml:"filter_size"`
	FalseRate  float64 `yaml:"false_rate"`
}

type LogConfig struct {
	LogToFile     bool   `yaml:"log_to_file"`
	LogSplit      bool   `yaml:"log_split"`
	LogDirectory  string `yaml:"log_directory"`
	LogLevel      string `yaml:"log_level"`
	LogFileSuffix string `yaml:"log_file_suffix"`
}

type AppConfig struct {
	MaxToken        int    `yaml:"max_token"`
	ManagementToken string `yaml:"management_token"`
	PriceTokenUnit  int64  `yaml:"price_token_unit"`
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	SSL        bool   `yaml:"ssl"`
	Location   string `yaml:"location"`
	SyncModels bool   `yaml:"sync_models"`
}
