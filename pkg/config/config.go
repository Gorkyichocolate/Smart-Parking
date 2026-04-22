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
	}

	return config, nil
}