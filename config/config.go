package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	SETA_GTFS_FTP string
	FTP_USERNAME  string
	FTP_PASSWORD  string

	DB_HOST     string
	DB_PORT     int
	DB_USERNAME string
	DB_PASSWORD string

	OIDC_ISSUER_URL     string
	OIDC_CLIENT_ID      string
	OIDC_ALLOWED_GROUPS []string
	CORS_ALLOWED_ORIGIN []string
	ENABLE_AUTH         bool

	IS_PRIMARY bool

	PORT string
)

func LoadConf() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso variabili d'ambiente di sistema", err)
	}

	SETA_GTFS_FTP = getEnv("SETA_GTFS_FTP", ":5001")
	FTP_USERNAME = getEnv("FTP_USERNAME", ":5001")
	FTP_PASSWORD = getEnv("FTP_PASSWORD", ":5001")

	DB_HOST = getEnv("DB_HOST", "localhost")
	DB_PORT = getEnvInt("DB_PORT", 3306)
	DB_USERNAME = getEnv("DB_USERNAME", "")
	DB_PASSWORD = getEnv("DB_PASSWORD", "")

	//OAuth config
	OIDC_ISSUER_URL = getEnv("OIDC_ISSUER_URL", "")
	OIDC_CLIENT_ID = getEnv("OIDC_CLIENT_ID", "")
	OIDC_ALLOWED_GROUPS = strings.Split(getEnv("OIDC_ALLOWED_GROUPS", ""), ",")
	CORS_ALLOWED_ORIGIN = strings.Split(getEnv("CORS_ALLOWED_ORIGIN", ""), ",")
	ENABLE_AUTH = getEnvBool("ENABLE_AUTH", "false")

	IS_PRIMARY = getEnvBool("IS_PRIMARY", "false")

	PORT = getEnv("PORT", ":5001")
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
