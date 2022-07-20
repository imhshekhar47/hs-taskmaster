package config

import "github.com/imhshekhar47/hs-taskmaster/go-api-core/utils"

var (
	APP_MODE_DEVELOPMENT = "development"
	APP_MODE_PRODUCTION  = "production"
)

type ServerConfig struct {
	Mode string `yaml:"mode",mapstructure:"mode"`
	Port string `yaml:"port",mapstructure:"port"`
	Path string `yaml:"path",mapstructure:"path"`
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		Mode: utils.GetEnvOrElse("SERVER_MODE", APP_MODE_DEVELOPMENT),
		Port: utils.GetEnvOrElse("SERVER_PORT", "8080"),
		Path: utils.GetEnvOrElse("SERVER_BASE_PATH", "/"),
	}
}
