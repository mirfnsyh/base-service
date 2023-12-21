package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type AppConfiguration struct {
	Host    string
	Env     string
	Port    int
	Code    string
	Version string
}

type DatabaseConfiguration struct {
	Driver               string
	Name                 string
	User                 string
	Password             string
	Host                 string
	Port                 int
	AdditionalParameters string
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
	Debug                bool
}

type Configuration struct {
	App      AppConfiguration
	Database DatabaseConfiguration
}

var configuration *Configuration
var once sync.Once

func GetConfig() *Configuration {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../../../.")
		viper.AddConfigPath("../../../../.")

		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&configuration); err != nil {
			logrus.Fatalf("Unable to decode into struct, %v", err)
		}
	})

	return configuration
}
