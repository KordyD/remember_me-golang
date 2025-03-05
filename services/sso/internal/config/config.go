package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	GRPCConfig `yaml:"grpc"`
	DBConfig   `yaml:"db" env-required:"true"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	GRPCPort int           `yaml:"port"`
	Timeout  time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string `yaml:"db_host" env-default:"localhost"`
	Port     int    `yaml:"db_port" env-default:"5432"`
	User     string `yaml:"db_user" env-default:"postgres"`
	Password string `yaml:"db_password" env-default:"postgres"`
	Name     string `yaml:"db_name" env-default:"postgres"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
