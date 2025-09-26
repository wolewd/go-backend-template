package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}
}

func GetEnv(key, fallback string) string {
	return getEnvWithFallback(key, fallback)
}

func GetEnvBytes(key, fallback string) []byte {
	return []byte(getEnvWithFallback(key, fallback))
}

func GetEnvInt(key string, fallback int) int {
	valStr := getEnvWithFallback(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid %s, defaulting to %d: %v", key, fallback, err)
		return fallback
	}
	return val
}

func getEnvWithFallback(key, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		log.Printf("Warning: %s using default value", key)
		return fallback
	}
	return val
}
