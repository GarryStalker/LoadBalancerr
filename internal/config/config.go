package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "dev"
	EnvProd  Environment = "prod"
)

type Config struct {
	Env     Environment `yaml:"env" env-default:"local"`
	Port    string      `yaml:"port" env-default:":443"`
	CDNHost string      `yaml:"cdn_host" env-default:"https://cdn.example.com"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty!")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/local.yaml", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
