package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI        string
	MongoDBName     string
	Port            string
	JWTSecret       string
	JWTTTL          time.Duration
	RefreshTokenTTL time.Duration
	AllowedOrigins  []string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:     getEnv("MONGO_DB_NAME", "agriculture"),
		Port:            getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTTTL:          2 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		AllowedOrigins:  splitAndTrim(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
	}
}

// splitAndTrim turns a comma-separated env value like
// "https://app.vercel.app,http://localhost:3000" into a clean slice.
func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
