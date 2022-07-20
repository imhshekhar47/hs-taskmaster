package config

import "github.com/imhshekhar47/hs-taskmaster/go-api-core/utils"

type DatasourceConfig struct {
	Host     string `yaml:"host",mapstructure:"host"`
	Port     string `yaml:"port",mapstructure:"port"`
	Database string `yaml:"database",mapstructure:"database"`
	Username string `yaml:"username",mapstructure:"username"`
	Password string `yaml:"password",mapstructure:"password"`
}

func GetDatasourceConfig() DatasourceConfig {
	return DatasourceConfig{
		Host:     utils.GetEnvOrElse("DB_HOST", "localhost"),
		Port:     utils.GetEnvOrElse("DB_PORT", "0"),
		Database: utils.GetEnvOrElse("DB_NAME", ""),
		Username: utils.GetEnvOrElse("DB_USERNAME", ""),
		Password: utils.GetEnvOrElse("DB_PASSWORD", ""),
	}
}
