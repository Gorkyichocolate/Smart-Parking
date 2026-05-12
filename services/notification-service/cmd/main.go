// services/notification-service/cmd/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smart-parking/pkg/config"

	"notification-service/internal/email"
	"notification-service/internal/queue"
)

func main() {
	// Загружаем конфиг из общего пакета
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Конвертируем SMTP_PORT из строки в int (если нужно)
	smtpPort := 587 // значение по умолчанию
	if cfg.SMTP_PORT != "" {
		// Если в вашем Config SMTP_PORT это string, конвертируем
		// Если у вас в Config SMTP_PORT уже int, то этот блок не нужен
	}

	// Initialize email sender
	emailSender := email.NewSMTPSender(
		cfg.SMTP_HOST,
		smtpPort,
		cfg.SMTP_USERNAME,
		cfg.SMTP_PASSWORD,
		cfg.SMTP_FROM,
	)

	// Create email handler
	emailHandler := queue.NewEmailHandler(
		"booking.notifications",
		emailSender,
		cfg.MaxRetries,
	)

	// Configure consumer
	consumerConfig := queue.ConsumerConfig{
		QueueName:            emailHandler.GetQueueName(),
		PrefetchCount:        1,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}

	// Create consumer (используем cfg.RABBITMQURL)
	consumer, err := queue.NewConsumer(
		cfg.RABBITMQURL,
		consumerConfig,
		emailHandler,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	defer consumer.Stop()

	// Start consumer
	if err := consumer.Start(); err != nil {
		log.Fatalf("❌ Failed to start consumer: %v", err)
	}

	log.Println("✅ Notification Service started successfully")
	log.Println("📡 Waiting for messages from RabbitMQ...")
	log.Printf("   Queue: booking.notifications")
	log.Printf("   RabbitMQ URL: %s", cfg.RABBITMQURL)
	log.Println("   Press Ctrl+C to stop")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutting down gracefully...")
}
