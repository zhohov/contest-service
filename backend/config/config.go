package config

import (
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		DB   DB
	}

	App struct {
		Name    string `env:"BACKEND_APP_NAME,required"`
		Version string `env:"BACKEND_APP_VERSION,required"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	DB struct {
		Url string `env:"DB_URL,required"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
