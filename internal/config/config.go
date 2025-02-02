package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"environment" env-required:"true"`
	TemplatesPath string `yaml:"templates_path" env-required:"true"`
	Dsn           string `yaml:"dsn" env-required:"true"`
	JWTSecret     string `yaml:"jwt_secret" env-required:"true"`
	HTTPServer    `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

func MustConfigParsing() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
