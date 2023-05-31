package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AppPort    int    `yaml:"APP_PORT"`
	DbHost     string `yaml:"DB_HOST"`
	DbName     string `yaml:"DB_NAME"`
	DbPort     int    `yaml:"DB_PORT"`
	DbUsername string `yaml:"DB_USERNAME"`
	DbPassword string `yaml:"DB_PASSWORD"`
}

func GetConfig() (*Config, error) {
	config := &Config{}

	path, err := os.Getwd()
	if err != nil {
		return config, err
	}

	file, err := os.Open(path + "/config/application.yml")
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
