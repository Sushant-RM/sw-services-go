package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func Load() *Config {

	_ = godotenv.Load()

	return &Config{
		ServerPort: getOptionalEnv("SERVER_PORT", "3468"),

		DBHost:     getRequiredEnv("DB_HOST"),
		DBPort:     getRequiredEnv("DB_PORT"),
		DBName:     getRequiredEnv("DB_NAME"),
		DBUser:     getRequiredEnv("DB_USER"),
		DBPassword: getRequiredEnv("DB_PASSWORD"),
	}
}

func getOptionalEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getRequiredEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s is required", key)
	}

	return value
}
