package config

import (
	"os"
	"strconv"
)

type Config struct {
	SMTP_HOST     string
	SMTP_PORT     int
	SMTP_USERNAME string
	SMTP_PASSWORD string
	SMTP_FROM     string
	RABBITMQURL   string
	MaxRetries    int
	Environment   string
	PaymentDatabaseURL      string
	PaymentDatabaseUser     string
	PaymentDatabasePassword string
	PaymentDatabaseName     string
}

func LoadConfig() (*Config, error) {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	maxRetries, _ := strconv.Atoi(getEnv("MAX_RETRIES", "3"))

	return &Config{
		SMTP_HOST:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTP_PORT:     smtpPort,
		SMTP_USERNAME: getEnv("SMTP_USERNAME", ""),
		SMTP_PASSWORD: getEnv("SMTP_PASSWORD", ""),
		SMTP_FROM:     getEnv("SMTP_FROM", "notifications@smartparking.com"),
		RABBITMQURL:   getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		MaxRetries:    maxRetries,
		Environment:   getEnv("ENVIRONMENT", "development"),
		PaymentDatabaseURL: getEnv("PAYMENT_DATABASE_URL", "postgres://payment_service_user:123456@localhost:5432/payment_service_db?sslmode=disable"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
