package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	AppPort         int      `mapstructure:"APP_PORT"`
	EnableZapLogger bool     `mapstructure:"ENABLE_ZAP_LOGGER"`
	DbConfig        DbConfig `mapstructure:"DB_CONFIG"`
	TestDbConfig    DbConfig `mapstructure:"TEST_DB_CONFIG"`
}

type DbConfig struct {
	DbHost     string `mapstructure:"DB_HOST"`
	DbName     string `mapstructure:"DB_NAME"`
	DbPort     int    `mapstructure:"DB_PORT"`
	DbUsername string `mapstructure:"DB_USERNAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
}

var config *Config

func init() {
	config = &Config{}

	_, b, _, _ := runtime.Caller(0)
	path := filepath.Dir(b)
	viper.AddConfigPath(path + "/../../config")
	viper.SetConfigName("application")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("todoapp")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic(err)
	}
}

func GetAppPort() int {
	return config.AppPort
}

func GetEnableZapLogger() bool {
	return config.EnableZapLogger
}

func GetDbConfig() DbConfig {
	return config.DbConfig
}

func GetTestDbConfig() DbConfig {
	return config.TestDbConfig
}
