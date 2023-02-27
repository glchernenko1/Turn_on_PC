package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"Turn_on_PC/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"true"`
	Server  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"1234"`
	} `yaml:"Server"`
	Postgres struct {
		Database string `yaml:"database"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"Postgres"`
	Token_password string `env-default:"token: "thisIsTheJwtSecretPassword"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infoln("read app config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Infoln(help)
			logger.Fatalln(err)
		}
	})
	return instance
}
