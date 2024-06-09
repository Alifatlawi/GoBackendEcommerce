package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func load() {
	// Load configuration from file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
