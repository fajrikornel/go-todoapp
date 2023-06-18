package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	AppPort  int      `yaml:"APP_PORT"`
	DbConfig DbConfig `yaml:"DB_CONFIG"`
}

type DbConfig struct {
	DbHost     string `yaml:"DB_HOST"`
	DbName     string `yaml:"DB_NAME"`
	DbPort     int    `yaml:"DB_PORT"`
	DbUsername string `yaml:"DB_USERNAME"`
	DbPassword string `yaml:"DB_PASSWORD"`
}

var config Config

func init() {
	config = Config{}

	_, b, _, _ := runtime.Caller(0)
	path := filepath.Dir(b)
	file, err := os.Open(path + "/../../config/application.yml")
	if err != nil {
		log.Fatalf("NO APPLICATION YML FILE, %v", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		log.Fatalf("CANNOT DECODE CONFIG, %v", err)
	}
}

func GetAppPort() int {
	return config.AppPort
}

func GetDbConfig() DbConfig {
	return config.DbConfig
}
