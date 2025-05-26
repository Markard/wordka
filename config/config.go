package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type (
	Setup struct {
		Config *Config
		Env    *Env
	}

	Config struct {
		HttpServer HttpServer `yaml:"http_server"`
		Log        Log        `yaml:"log"`
	}

	HttpServer struct {
		Address     string        `yaml:"address" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	}

	Log struct {
		Level string `yaml:"level" env-default:"debug"`
	}

	Env struct {
		AppEnv          string
		PgDSN           string
		ES256PrivateKey string
		ES256PublicKey  string
	}
)

func MustLoad() *Setup {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get current directory")
	}

	var env Env
	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatal().Err(errLoadEnv).Msg("Error loading .env file")
	}
	env.AppEnv = os.Getenv("APP_ENV")
	env.PgDSN = os.Getenv("PG_DSN")
	env.ES256PrivateKey = os.Getenv("ES256_PRIVATE_KEY")
	env.ES256PublicKey = os.Getenv("ES256_PUBLIC_KEY")

	configPath := fmt.Sprintf("%s/config/%s.yaml", currentDir, env.AppEnv)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	return &Setup{&cfg, &env}
}
