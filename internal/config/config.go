package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env                 string
	Host                string
	Port                int
	DatabaseURL         string
	CORSAllowedOrigins  []string
}

func Load() (Config, error) {
	port, err := strconv.Atoi(getenv("APP_PORT", "8080"))
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("invalid APP_PORT")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return Config{
		Env:                getenv("APP_ENV", "development"),
		Host:               getenv("APP_HOST", "0.0.0.0"),
		Port:               port,
		DatabaseURL:        databaseURL,
		CORSAllowedOrigins: splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost,http://127.0.0.1,http://localhost:3000,http://localhost:5173")),
	}, nil
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	var result []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			result = append(result, item)
		}
	}
	return result
}
