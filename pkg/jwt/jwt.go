package jwt

import (
	"os"
	"time"
)

var (
	SecretKey  = []byte(getEnv("JWT_SECRET", "smart-parking-secret-key-2024"))
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 168 * time.Hour
	IssuerName = "smart-parking"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
