package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	AppName            string
	AppPort            string
	AppVersion         string
	CORSAllowedOrigins string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
}

func Load() Config {
	_ = godotenv.Load(".env")
	_ = godotenv.Overload("apps/api/.env")

	return Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppName:            getEnv("APP_NAME", "ikant-setop-us-api"),
		AppPort:            getEnv("APP_PORT", "8081"),
		AppVersion:         getEnv("APP_VERSION", "v1"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5438"),
		DBUser:             getEnv("DB_USER", "ikant_user"),
		DBPassword:         getEnv("DB_PASSWORD", "ikant_pass"),
		DBName:             getEnv("DB_NAME", "ikant_setop_us_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
