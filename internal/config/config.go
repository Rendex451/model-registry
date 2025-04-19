package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env"`
	ModelDir     string `yaml:"model_dir"`
	DatabasePath string `yaml:"database_path"`
	Port         int    `yaml:"port"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		configPath := os.Getenv("CONFIG_PATH")

		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file %s does not exist", configPath)
		}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			log.Fatalf("error reading config file %s: %s", configPath, err)
		}
	})

	return cfg
}
