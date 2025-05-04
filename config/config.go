package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type (
	Config struct {
		HttpServer HttpServer `yaml:"http_server"`
		Log        Log        `yaml:"log"`
	}

	HttpServer struct {
		Address     string        `yaml:"address" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}

	Log struct {
		Level                string `yaml:"level" env-default:"debug"`
		FilePath             string `yaml:"file_path" env-required:"true"`
		CallerSkipFrameCount int    `yaml:"caller_skip_frame_count" env-default:"3"`
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
		log.Fatal().Err(err).Msg("could not get current directory")
	}

	var env Env
	dotenvPath := fmt.Sprintf("%s/.env", currentDir)
	if err := cleanenv.ReadConfig(dotenvPath, &env); err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	configPath := fmt.Sprintf("%s/config/%s.yaml", currentDir, env.AppEnv)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	return &env, &cfg
}
