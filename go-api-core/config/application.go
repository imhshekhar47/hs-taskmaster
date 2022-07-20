package config

import (
	"encoding/json"

	"github.com/imhshekhar47/hs-taskmaster/go-api-core/utils"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Name        string           `yaml:"name",mapstructure:"name"`
	Version     string           `yaml:"version",mapstructure:"version"`
	Description string           `yaml:"description",maptructure:"description"`
	Server      ServerConfig     `yaml:"server",mapstructure:"server"`
	Datasource  DatasourceConfig `yaml:"datasource",mapstructure:"datasource"`
}

func GetApplicationConfig() *ApplicationConfig {

	return &ApplicationConfig{
		Name:        utils.GetEnvOrElse("APP_NAME", ""),
		Version:     utils.GetEnvOrElse("APP_VERSION", ""),
		Description: "",
		Server:      GetServerConfig(),
		Datasource:  GetDatasourceConfig(),
	}
}

func (s *ApplicationConfig) Json() string {
	jsonStr, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}

	return string(jsonStr)
}

func (s *ApplicationConfig) Yaml() string {
	yamlStr, err := yaml.Marshal(s)
	if err != nil {
		return err.Error()
	}

	return string(yamlStr)
}
