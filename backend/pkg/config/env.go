package config

import (
	"os"
	"strings"
)

type Config struct {
	DB_URL          string
	API_PORT        string
	CORS_CONFIG     []string
	JWT_SECRET      []byte
	JWT_AUTH_METHOD string
	LLM_URL         string
	Dynamic         map[string]string
}

func Init() *Config {
	corsConfig := strings.Split(os.Getenv("CORS_ORIGINS"), ",")

	// secret := sha512.New().Sum([]byte(time.Now().GoString()))
	secret := []byte("a")

	config := &Config{
		DB_URL:          os.Getenv("DB_URL"),
		API_PORT:        os.Getenv("API_PORT"),
		LLM_URL:         os.Getenv("LLM_URL"),
		CORS_CONFIG:     corsConfig,
		JWT_AUTH_METHOD: os.Getenv("JWT_AUTH_METHOD"),
		JWT_SECRET:      secret,
		Dynamic:         map[string]string{},
	}
	return config
}
