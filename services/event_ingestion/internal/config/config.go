package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env  string `yaml:"env" env-required:"true"`
	Port int    `yaml:"port" env-default:"8080"`
}

func MustSetup() Config {
	config, err := setup()
	if err != nil {
		panic(err)
	}
	return config
}

func setup() (Config, error) {
	var config Config
	configPath := fetchConfigPath()
	if configPath == "" {
		return Config{}, errors.New("config path is empty")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, errors.New("config file not found: " + configPath)
	}
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		return Config{}, errors.New("failed to read config: " + err.Error())
	}
	return config, nil
}

func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	return configPath
}
