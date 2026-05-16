package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment"
	grpchandler "github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/delivery/grpc"
	"github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/infrastructer/rabbitmq"
	"github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/infrastructer/repository"
	"github.com/Gorkyichocolate/smart-parking/services/payment-service/internal/usecase"
)

func main() {
	databaseURL := os.Getenv("PAYMENT_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://payment_service_user:123456@localhost:5432/payment_service_db?sslmode=disable"
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	paymentRepo := repository.NewPaymentRepositoryPG(db)
	invoiceRepo := repository.NewInvoiceRepositoryPG(db)
	rabbitPublisher, err := rabbitmq.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitPublisher.Close()

	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, invoiceRepo, rabbitPublisher)

	grpcServer := grpc.NewServer()
	paymentHandler := grpchandler.NewPaymentHandler(paymentUseCase)
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Payment Service gRPC server listening on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	log.Println("Payment Service stopped")
}
