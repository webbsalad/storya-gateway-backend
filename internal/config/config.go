package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AllowedOrigins []string

	OtpAddr      string
	PassportAddr string
	ContentAddr  string
	RecsAddr     string
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("get env: %v", err)
	}

	allowedOriginsEnv := mustGetEnv("ALLOWED_ORIGINS")
	allowedOrigins := splitEnvList(allowedOriginsEnv)

	return Config{
		AllowedOrigins: allowedOrigins,

		OtpAddr:      mustGetEnv("OTP_ADR"),
		PassportAddr: mustGetEnv("PASSPORT_ADR"),
		ContentAddr:  mustGetEnv("CONTENT_ADR"),
		RecsAddr:     mustGetEnv("RECS_ADR"),
	}
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
		panic("Missing required environment variable: " + key)
	}
	return value
}

func splitEnvList(envValue string) []string {
	if envValue == "" {
		return []string{}
	}
	items := strings.Split(envValue, ",")
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}
