package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	PORT                string
	DB_HOST             string
	DB_PORT             int
	DB_USER             string
	DB_PASS             string
	OIDC_ISSUER_URL     string
	OIDC_CLIENT_ID      string
	OIDC_ALLOWED_GROUPS []string
	CORS_ALLOWED_ORIGIN string
	ENABLE_AUTH         bool
)

func LoadConf() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso variabili d'ambiente di sistema", err)
	}

	PORT = getEnv("PORT", ":5001")
	DB_HOST = getEnv("DB_HOST", "localhost")
	DB_PORT = getEnvInt("DB_PORT", 3306)
	DB_USER = getEnv("DB_USER", "")
	DB_PASS = getEnv("DB_PASS", "")

	//OAuth config
	OIDC_ISSUER_URL = getEnv("OIDC_ISSUER_URL", "")
	OIDC_CLIENT_ID = getEnv("OIDC_CLIENT_ID", "")
	OIDC_ALLOWED_GROUPS = strings.Split(getEnv("OIDC_ALLOWED_GROUPS", ""), ",")
	CORS_ALLOWED_ORIGIN = getEnv("CORS_ALLOWED_ORIGIN", "")
	ENABLE_AUTH = getEnvBool("ENABLE_AUTH", "false")
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

func getEnvBool(key string, fallback string) bool {
	value := getEnv(key, fallback)
	if value == "true" {
		return true
	} else {
		return false
	}
}
