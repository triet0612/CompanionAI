package config

import (
	"crypto/sha512"
	"os"
	"strings"
	"time"
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

	config := &Config{
		DB_URL:          os.Getenv("DB_URL"),
		API_PORT:        os.Getenv("API_PORT"),
		LLM_URL:         os.Getenv("LLM_URL"),
		CORS_CONFIG:     corsConfig,
		JWT_AUTH_METHOD: os.Getenv("JWT_AUTH_METHOD"),
		JWT_SECRET:      sha512.New().Sum([]byte(time.Now().GoString())),
		Dynamic:         map[string]string{},
	}
	return config
}
