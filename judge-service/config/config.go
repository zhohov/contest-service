package config

import (
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App App
	}

	App struct {
		Name    string `env:"JUDGE_APP_NAME,required"`
		Version string `env:"JUDGE_APP_VERSION,required"`
	}
)

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
