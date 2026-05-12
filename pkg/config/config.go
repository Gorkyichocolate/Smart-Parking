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
	SMTP_PORT     int
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
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Получаем значения из переменных окружения (БЕЗ ДЕФОЛТНЫХ ЗНАЧЕНИЙ!)
	dbHost := os.Getenv("PAYMENT_DATABASE_URL")
	dbUser := os.Getenv("PAYMENT_DATABASE_USER")
	dbPassword := os.Getenv("PAYMENT_DATABASE_PASSWORD")
	dbName := os.Getenv("PAYMENT_DATABASE_NAME")

	// Проверяем, что все обязательные переменные заданы
	if dbHost == "" {
		return nil, fmt.Errorf("PAYMENT_DATABASE_URL is not set in .env file")
	}
	if dbUser == "" {
		return nil, fmt.Errorf("PAYMENT_DATABASE_USER is not set in .env file")
	}
	if dbPassword == "" {
		return nil, fmt.Errorf("PAYMENT_DATABASE_PASSWORD is not set in .env file")
	}
	if dbName == "" {
		return nil, fmt.Errorf("PAYMENT_DATABASE_NAME is not set in .env file")
	}

	// Формируем полный connection string для PostgreSQL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbName)

	// Получаем остальные переменные
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	maxRetriesStr := os.Getenv("MAX_RETRIES")
	environment := os.Getenv("ENVIRONMENT")

	// Конвертируем порт
	smtpPort := 587
	if smtpPortStr != "" {
		if port, err := strconv.Atoi(smtpPortStr); err == nil {
			smtpPort = port
		}
	}

	// Конвертируем max retries
	maxRetries := 3
	if maxRetriesStr != "" {
		if retries, err := strconv.Atoi(maxRetriesStr); err == nil {
			maxRetries = retries
		}
	}

	config := &Config{
		// SMTP settings
		SMTP_HOST:     smtpHost,
		SMTP_PORT:     smtpPort,
		SMTP_USERNAME: os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		SMTP_FROM:     os.Getenv("SMTP_FROM"),

		// RabbitMQ
		RABBITMQURL: rabbitmqURL,

		// Application settings
		MaxRetries:  maxRetries,
		Environment: environment,

		// Payment database
		PaymentDatabaseURL:      dbURL,
		PaymentDatabaseUser:     dbUser,
		PaymentDatabasePassword: dbPassword,
		PaymentDatabaseName:     dbName,
	}

	// Валидация обязательных полей
	if config.RABBITMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is not set in .env file")
	}

	return config, nil
}
