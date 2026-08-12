package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT    string
	DB_HOST string
	DB_PORT int
	DB_USER string
	DB_PASS string
)

func LoadConf() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso variabili d'ambiente di sistema")
	}

	PORT = getEnv("PORT", ":5001")
	DB_HOST = getEnv("DB_HOST", "localhost")
	DB_PORT = getEnvInt("DB_PORT", 3306)
	DB_USER = getEnv("DB_USER", "")
	DB_PASS = getEnv("DB_PASS", "")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}
