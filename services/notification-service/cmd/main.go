package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gorkyichocolate/smart-parking/services/notification-service/internal/email"
	"github.com/Gorkyichocolate/smart-parking/services/notification-service/internal/queue"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	smtpPort := 587
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "notifications@smartparking.com"
	}
	maxRetries := 3

	emailSender := email.NewSMTPSender(smtpHost, smtpPort, smtpUsername, smtpPassword, smtpFrom)
	emailHandler := queue.NewEmailHandler("booking.notifications", emailSender, maxRetries)

	consumerConfig := queue.ConsumerConfig{
		QueueName:            emailHandler.GetQueueName(),
		PrefetchCount:        1,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}

	consumer, err := queue.NewConsumer(rabbitURL, consumerConfig, emailHandler)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Stop()

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Notification Service started")
	log.Println("Waiting for messages from RabbitMQ...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
