package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr     string        `env:"HTTP_ADDR" envDefault:"0.0.0.0:8080"`
	DatabasePath string        `env:"DATABASE_PATH" envDefault:"data/comercio_exterior.db"`
	JWTSecret    string        `env:"JWT_SECRET,required"`
	JWTTTL       time.Duration `env:"JWT_TTL" envDefault:"15m"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
