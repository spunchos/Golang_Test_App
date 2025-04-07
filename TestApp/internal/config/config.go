package config

import (
	"TestApp/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	}
	Storage StorageConfig `yaml:"storage"`
	MongoDB string        `yaml:"mongodb"`
}

var instance *Config
var once sync.Once

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("reading app config")
		instance = &Config{}
		err := cleanenv.ReadConfig("config.yml", instance)
		if err != nil {
			info, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(info)
			logger.Fatal(err)
		}
	})
	return instance
}
