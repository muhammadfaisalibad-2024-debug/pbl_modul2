package config

import (
	"os"
	"strconv"

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

	MongoURI      string
	MongoDatabase string

	JWTSecret           string
	JWTIssuer           string
	JWTAccessTTLMinutes string
	JWTRefreshTTLDays   string
	AllowedOrigins      string
}

func Load() Config {
	// Try multiple relative paths for .env
	for _, path := range []string{".env", "../.env", "../../.env", "../../../.env"} {
		if err := godotenv.Overload(path); err == nil {
			break
		}
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "123456"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "praktikum_backend"
	}
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	mongoDB := os.Getenv("MONGO_DATABASE")
	if mongoDB == "" {
		mongoDB = "praktikum_backend"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "65b458dc012be6c16665db71ee6ec7f53dc9726237bbbd4e7288bacd789955cf"
	}
	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "siakad-mini"
	}

	return Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbSSLMode,
		DBMaxConns: os.Getenv("DB_MAX_CONNS"),

		MongoURI:      mongoURI,
		MongoDatabase: mongoDB,

		JWTSecret:           jwtSecret,
		JWTIssuer:           jwtIssuer,
		JWTAccessTTLMinutes: os.Getenv("JWT_ACCESS_TTL_MINUTES"),
		JWTRefreshTTLDays:   os.Getenv("JWT_REFRESH_TTL_DAYS"),
		AllowedOrigins:      os.Getenv("ALLOWED_ORIGINS"),
	}
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
