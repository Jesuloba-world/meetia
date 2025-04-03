package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"APP_ENV"`
	Port string `mapstructure:"APP_PORT"`

	// Database
	DBUrl string `mapstructure:"DATABASE_URL"`

	// Auth
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func Load() *Config {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Sprintf("Error loading .env file: %v", err))
		}
	}

	viper.AutomaticEnv()

	// set defaults
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/meetia")
	viper.SetDefault("JWT_SECRET", "default-secret-please-change")

	// create config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}

	if cfg.JWTSecret == "default-secret-please-change" && cfg.Env == "production" {
		panic("JWT_SECRET must be set in production environment")
	}

	return &cfg
}
