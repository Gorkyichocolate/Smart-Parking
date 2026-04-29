package main

import (
	"fmt"
	"log"

	"payment-service/internal/delivery/grpc"
	"payment-service/internal/infrastructer/rabbitmq"
	repoImpl "payment-service/internal/infrastructer/repository"
	"payment-service/internal/usecase"
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
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	pub := publisher.New(ch, "events")
	paymentPub := rabbitmq.NewPaymentPublisher(pub)

	// repository
	paymentRepo := repoImpl.NewPaymentRepo(db)
	invoiceRepo := repoImpl.NewInvoiceRepo(db)

	usecase := usecase.NewPaymentUsecase(
		paymentRepo,
		invoiceRepo,
		paymentPub,
	)

	// usecase
	usecase := usecase.NewPaymentUseCase(repo, paymentPub)

	// handler
	handler := grpc.NewHandler(usecase)

	// Simulate creating a payment
	err = handler.CreatePayment(nil)
	if err != nil {
		log.Fatalf("Failed to create payment: %v", err)
	}
}
