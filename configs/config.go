package configs

import (
	"os"
	"time"
)

type Config struct {
	JWTSecret   []byte
	JWTDuration time.Duration
}

func InitConfig() *Config {
	jwtDuration, err := time.ParseDuration(GetEnv("JWT_DURATION", "1h"))
	if err != nil {
		panic(err)
	}

	return &Config{
		JWTSecret:   []byte(GetEnv("JWT_SECRET", "secret")),
		JWTDuration: jwtDuration,
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
