package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxConns string

	JWTSecret           string
	JWTIssuer           string
	JWTAccessTTLMinutes string
	JWTRefreshTTLDays   string
	AllowedOrigins      string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	return Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBMaxConns: os.Getenv("DB_MAX_CONNS"),

		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTIssuer:           os.Getenv("JWT_ISSUER"),
		JWTAccessTTLMinutes: os.Getenv("JWT_ACCESS_TTL_MINUTES"),
		JWTRefreshTTLDays:   os.Getenv("JWT_REFRESH_TTL_DAYS"),
		AllowedOrigins:      os.Getenv("ALLOWED_ORIGINS"),
	}
}
