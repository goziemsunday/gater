package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env               string
	Port              string
	DatabaseURL       string
	RedisURL          string
	CORSAllowedOrigin string
}

func (c *Config) validate() error {
	required := map[string]string{
		"DATABASE_URL":        c.DatabaseURL,
		"REDIS_URL":           c.RedisURL,
		"CORS_ALLOWED_ORIGIN": c.CORSAllowedOrigin,
	}

	for k, v := range required {
		if v == "" {
			return fmt.Errorf("missing required environment variable: %s", k)
		}
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	return nil
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Env:               os.Getenv("ENV"),
		Port:              os.Getenv("PORT"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		RedisURL:          os.Getenv("REDIS_URL"),
		CORSAllowedOrigin: os.Getenv("CORS_ALLOWED_ORIGIN"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
