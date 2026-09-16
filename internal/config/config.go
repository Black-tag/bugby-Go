package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	AppEnv string
	DBURL  string
	Port   string
	Secret string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")

	}
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			slog.Warn("no .env.development file found, using environment variables")

		}
	}
	cfg := Config{
		AppEnv: os.Getenv("APP_ENV"),
		DBURL:  os.Getenv("DB_URL"),
		Port:   os.Getenv("PORT"),
		Secret: os.Getenv("SECRET"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil

}
