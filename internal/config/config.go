package config

import "github.com/caarlos0/env/v11"

type Config struct {
	HTTPAddr string `env:"HTTP_ADDR" envDefault:"0.0.0.0:8080"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
