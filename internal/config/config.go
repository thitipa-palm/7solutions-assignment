package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	MongoDatabase string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() (*Config, error) {
	// ไม่คืน error เพราะ production อาจส่ง environment variables
	// เข้ามาโดยตรงโดยไม่มีไฟล์ .env
	_ = godotenv.Load()

	jwtExpiration, err := time.ParseDuration(
		getEnv("JWT_EXPIRATION", "24h"),
	)
	if err != nil {
		return nil, fmt.Errorf("parse JWT_EXPIRATION: %w", err)
	}

	config := &Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DATABASE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: jwtExpiration,
	}

	if config.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}

	if config.MongoDatabase == "" {
		return nil, fmt.Errorf("MONGO_DATABASE is required")
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
