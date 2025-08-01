// Package config provides loading config data from
// external sources like env variables.
package config

import (
	"fmt"
	"slices"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		Server
		DB
	}

	App struct {
		LogLevel        string `env:"LOG_LEVEL" env-default:"info" env-description:"info/warn/error"`
		LogFormat       string `env:"LOG_FORMAT" env-default:"text" env-description:"text/json"`
		CoingeckoAPIKey string `env-required:"true" env:"COINGECKO_API_KEY"`
	}

	Server struct {
		Port            string        `env:"SERVER_PORT" env-default:"8000"`
		ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"5s"`
	}

	DB struct {
		MigrationsURL string `env:"MIGRATIONS_URL" env-default:"file://migrations"`
		User          string `env-required:"true" env:"POSTGRES_USER"`
		Password      string `env-required:"true" env:"POSTGRES_PASSWORD"`
		Host          string `env-required:"true" env:"POSTGRES_HOST"`
		Port          string `env-required:"true" env:"POSTGRES_PORT"`
		Name          string `env-required:"true" env:"POSTGRES_DB"`
		ConnString    string
		ConnURL       string
	}
)

var (
	_acceptedLogFormats = []string{"text", "json"}
	_acceptedLogLevels  = []string{"info", "warn", "error"}
)

// New returns app config loaded from ENV-vars.
func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("load env variables: %w", err)
	}
	// if invalid log format
	if !slices.Contains(_acceptedLogFormats, cfg.App.LogFormat) {
		return nil, fmt.Errorf(
			"invalid log format %s. Accepted formats: %v",
			cfg.App.LogFormat, _acceptedLogFormats,
		)
	}
	// if invalid log level
	if !slices.Contains(_acceptedLogLevels, cfg.App.LogLevel) {
		return nil, fmt.Errorf(
			"invalid log level %s. Accepted levels: %v",
			cfg.App.LogLevel, _acceptedLogLevels,
		)
	}

	cfg.DB.ConnString = fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable connect_timeout=10",
		cfg.DB.User, cfg.DB.Password,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.Name,
	)
	cfg.DB.ConnURL = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		cfg.DB.User, cfg.DB.Password,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.Name,
	)
	return cfg, nil
}
