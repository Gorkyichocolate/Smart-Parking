package main

import (
	"fmt"
	"log"
	"time"

	"payment-service/internal/infrastructer/rabbitmq"
	"smart-parking/pkg/config"
	"smart-parking/pkg/postgres"
	"smart-parking/pkg/rabbitmq/connection"
	"smart-parking/pkg/rabbitmq/publisher"
)

func main() {
	// config
	config, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	if err != nil {
		panic("Error loading .env file")
	}

	// database
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.PaymentDatabaseUser,
		config.PaymentDatabasePassword,
		config.PaymentDatabaseURL,
		config.PaymentDatabaseName,
	)
	

	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	_ = db

	// rabbitmq
	
	conn := connection.New("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()

	pub := publisher.New(ch, "events")
	paymentPub := rabbitmq.NewPaymentPublisher(pub)

	for {
		paymentPub.PublishPaymentCreated(rabbitmq.PaymentCreated{
			ID:     "123",
			Amount: 50,
		})

		time.Sleep(5 * time.Second)
	}
}
