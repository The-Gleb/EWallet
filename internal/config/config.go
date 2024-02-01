package config

import (
	"github.com/num30/config"
)

type Config struct {
	RunAddress string   `default:":8080"`
	LogLevel   string   `default:"debug" flag:"loglevel"`
	DB         Database `default:"{}"`
}

type Database struct {
	Host     string `default:"localhost" validate:"required"`
	Password string `default:"emoney" validate:"required" envvar:"DB_PASS"`
	DbName   string `default:"emoney"`
	Username string `default:"emoney"`
	Port     int    `default:"5434"`
}

func BuildConfig(cfgFile string) *Config {
	var conf Config
	err := config.NewConfReader(cfgFile).Read(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
