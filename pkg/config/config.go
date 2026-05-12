// pkg/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// SMTP settings
	SMTP_HOST     string
	SMTP_PORT     int // <-- должен быть int
	SMTP_USERNAME string
	SMTP_PASSWORD string
	SMTP_FROM     string

	// RabbitMQ
	RABBITMQURL string

	// Application settings
	MaxRetries  int
	Environment string

	// Payment service database
	PaymentDatabaseURL      string
	PaymentDatabaseUser     string
	PaymentDatabasePassword string
	PaymentDatabaseName     string
}

func LoadConfig() (*Config, error) {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Конвертируем SMTP_PORT из string в int
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 587 // default если не задан
	}

	// Конвертируем MAX_RETRIES
	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		maxRetries = 3
	}

	config := &Config{
		// SMTP settings
		SMTP_HOST:     os.Getenv("SMTP_HOST"),
		SMTP_PORT:     smtpPort,
		SMTP_USERNAME: os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		SMTP_FROM:     os.Getenv("SMTP_FROM"),

		// RabbitMQ
		RABBITMQURL: os.Getenv("RABBITMQ_URL"),

		// Application settings
		MaxRetries:  maxRetries,
		Environment: os.Getenv("ENVIRONMENT"),

		// Payment database
		PaymentDatabaseURL:      os.Getenv("PAYMENT_DATABASE_URL"),
		PaymentDatabaseUser:     os.Getenv("PAYMENT_DATABASE_USER"),
		PaymentDatabasePassword: os.Getenv("PAYMENT_DATABASE_PASSWORD"),
		PaymentDatabaseName:     os.Getenv("PAYMENT_DATABASE_NAME"),
	}

	return config, nil
}
