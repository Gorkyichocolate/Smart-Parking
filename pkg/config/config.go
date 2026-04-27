package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTPAddr          string
	SMTPDatabaseURL    string
	SMTPDatabaseUser   string
	SMTPDatabasePassword string
	SMTPDatabaseName   string
	PaymentDatabaseURL    string
	PaymentDatabaseUser   string
	PaymentDatabasePassword string
	PaymentDatabaseName   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	config := &Config{
		FromEmail:         os.Getenv("FROM_EMAIL"),
		FromEmailPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
		FromEmailSMTP:     os.Getenv("FROM_EMAIL_SMTP"),
		SMTPAddr:          os.Getenv("SMTP_ADDR"),
		SMTPDatabaseURL:    os.Getenv("SMTP_DATABASE_URL"),
		SMTPDatabaseUser:   os.Getenv("SMTP_DATABASE_USER"),
		SMTPDatabasePassword: os.Getenv("SMTP_DATABASE_PASSWORD"),
		SMTPDatabaseName:   os.Getenv("SMTP_DATABASE_NAME"),
		PaymentDatabaseURL:    os.Getenv("PAYMENT_DATABASE_URL"),
		PaymentDatabaseUser:   os.Getenv("PAYMENT_DATABASE_USER"),
		PaymentDatabasePassword: os.Getenv("PAYMENT_DATABASE_PASSWORD"),
		PaymentDatabaseName:   os.Getenv("PAYMENT_DATABASE_NAME"),
	}

	return config, nil
}