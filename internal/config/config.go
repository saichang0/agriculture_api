package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI    string
	MongoDBName string
	Port        string
	JWTSecret   string
	JWTTTL      time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "agriculture"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTTTL:      24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
