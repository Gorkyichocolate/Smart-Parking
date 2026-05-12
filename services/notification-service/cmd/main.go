// services/notification-service/cmd/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GorkyiChocolate/smart-parking/pkg/config"

	"github.com/GorkyiChocolate/smart-parking/pkg/metrics"

	"github.com/GorkyiChocolate/smart-parking/services/notification-service/internal/email"
	"github.com/GorkyiChocolate/smart-parking/services/notification-service/internal/queue"
)

func main() {
	// Initialize metrics
	metrics.InitMetrics("notification_service")
	metrics.StartMetricsServer("9090")

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Initialize email sender
	emailSender := email.NewSMTPSender(
		cfg.SMTP_HOST,
		cfg.SMTP_PORT,
		cfg.SMTP_USERNAME,
		cfg.SMTP_PASSWORD,
		cfg.SMTP_FROM,
	)

	// Create email handler with metrics
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

	// Create consumer
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
	log.Printf("   Metrics endpoint: http://localhost:9090/metrics")
	log.Println("   Press Ctrl+C to stop")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutting down gracefully...")
}
