package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"192.168.0.158"`
		Port string `env:"SERVER_PORT" envDefault:"8080"`
		Env  string `env:"ENV" envDefault:""`
	}

	DB struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		Name     string `env:"DB_NAME" envDefault:"postgres"`
		SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	JWT struct {
		Secret     string `env:"JWT_SECRET" envDefault:"your-secret-key"`
		AccessTTL  int    `env:"ACCESS_JWT_TTL" envDefault:"24"` // hours
		AdminToken string `env:"ADMIN_TOKEN" envDefault:""`
		Issuer     string `env:"JWT_ISSUER" envDefault:"doc-storage"`
	}

	Doc struct {
		DocTTL int `env:"DOC_TTL" envDefault:"24"` // hours
	}
)

type Config struct {
	Server
	DB
	Redis
	JWT
	Doc
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}
	_ = env.Parse(cfg)

	return cfg
}
