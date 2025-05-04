package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		HttpServer HttpServer `yaml:"http_server"`
	}

	HttpServer struct {
		Address     string        `yaml:"address" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}

	Env struct {
		AppName    string `env:"APP_NAME" env-required:"true"`
		AppVersion string `env:"APP_VERSION" env-required:"true"`
		AppEnv     string `env:"APP_ENV" env-required:"true"`
	}
)

func MustLoad() (*Env, *Config) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	var env Env
	dotenvPath := fmt.Sprintf("%s/.env", currentDir)
	if err := cleanenv.ReadConfig(dotenvPath, &env); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	configPath := fmt.Sprintf("%s/config/%s.yaml", currentDir, env.AppEnv)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &env, &cfg
}
