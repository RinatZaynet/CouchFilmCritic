package config

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string `yaml:"environment" env-required:"true"`
	TemplatesPath   string `yaml:"templates_path" env-required:"true"`
	Dsn             string `yaml:"dsn" env-required:"true"`
	JWTSecret       string `yaml:"jwt_secret" env-required:"true"`
	HashPassOptions `yaml:"hash_pass_options"`
	HTTPServer      `yaml:"http_server"`
}

type HashPassOptions struct {
	Time    uint32 `yaml:"time" env-default:"5"`
	Memory  uint32 `yaml:"memory" env-default:"65536"`
	Threads uint8  `yaml:"threads" env-default:"2"`
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

	if cfg.Memory < 0 || cfg.Time <= 0 || cfg.Threads <= 0 ||
		cfg.Memory > math.MaxUint32 ||
		cfg.Time > math.MaxUint32 ||
		cfg.Threads > math.MaxUint8 {
		log.Fatalf("password hash options is invalid. Options: Memory: %d, Time: %d, Threads: %d.", cfg.Memory, cfg.Time, cfg.Threads)
	}
	return &cfg
}
