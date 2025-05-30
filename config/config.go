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
	}

	HttpServer struct {
		Address     string        `yaml:"address" env-required:"true"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	}

	Env struct {
		AppEnv          string
		ES256PrivateKey string
		ES256PublicKey  string
		PgDb            string
		PgUser          string
		PgPass          string
		PgPort          string
		PgHost          string
		PgDSN           string
	}
)

func newEnv(appEnv, eS256PrivateKey, eS256PublicKey, pgDb, pgUser, pgPass, pgPort, pgHost string) *Env {
	env := &Env{
		AppEnv:          appEnv,
		ES256PrivateKey: eS256PrivateKey,
		ES256PublicKey:  eS256PublicKey,
		PgDb:            pgDb,
		PgUser:          pgUser,
		PgPass:          pgPass,
		PgPort:          pgPort,
		PgHost:          pgHost,
	}
	env.PgDSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.PgUser, env.PgPass, env.PgHost, env.PgPort, env.PgDb)

	return env
}

func MustLoad() *Setup {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get current directory")
	}

	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatal().Err(errLoadEnv).Msg("Error loading .env file")
	}
	env := newEnv(
		os.Getenv("APP_ENV"),
		os.Getenv("ES256_PRIVATE_KEY"),
		os.Getenv("ES256_PUBLIC_KEY"),
		os.Getenv("PG_DB"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_HOST"),
	)

	configPath := fmt.Sprintf("%s/config/%s.yaml", currentDir, env.AppEnv)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	return &Setup{&cfg, env}
}
