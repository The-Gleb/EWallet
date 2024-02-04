package config

import (
	"github.com/num30/config"
)

type Config struct {
	RunAddress string   `default:"0.0.0.0:8080"`
	LogLevel   string   `default:"debug" flag:"loglevel"`
	DB         Database `default:"{}"`
}

type Database struct {
	Host     string `default:"localhost" envvar:"DB_HOST" validate:"required"`
	Password string `default:"emoney" validate:"required" envvar:"DB_PASS"`
	DbName   string `default:"emoney"`
	Username string `default:"emoney"`
	Port     int    `default:"5434" envvar:"DB_PORT"`
}

func BuildConfig(cfgFile string) *Config {
	var conf Config

	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
