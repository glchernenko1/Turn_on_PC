package config

import (
	"Turn_on_PC/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"true"`
	Server  struct {
		BindIP string `yaml:"bind_ip" env:"BindIP" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env:"PortService" env-default:"1234"`
	} `yaml:"Server"`
	Postgres struct {
		Database string `yaml:"database" env:"Database"`
		Host     string `yaml:"host" env:"Host"`
		Port     string `yaml:"port" env:"Port"`
		Username string `yaml:"username" env:"Username"`
		Password string `yaml:"password" env:"Password"`
	} `yaml:"Postgres"`
	TokenPassword string `env:"token_password"`
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
