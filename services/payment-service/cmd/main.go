// services/payment-service/cmd/main.go
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

	paymentpb "github.com/GorkyiChocolate/smart-parking-proto/gen/go/payment"

	grpchandler "github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/delivery/grpc"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/infrastructer/rabbitmq"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/infrastructer/repository"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/usecase"

	"github.com/GorkyiChocolate/smart-parking/pkg/config"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Убираем логирование конфиденциальных данных
	// Просто лог о подключении к БД без показа данных
	log.Println("📊 Initializing database connection...")

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.PaymentDatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepositoryPG(db)
	invoiceRepo := repository.NewInvoiceRepositoryPG(db)

	// Initialize RabbitMQ publisher
	rabbitPublisher, err := rabbitmq.NewPublisher(cfg.RABBITMQURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitPublisher.Close()
	log.Println("✅ Connected to RabbitMQ")

	// Initialize usecase
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, invoiceRepo, rabbitPublisher)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	paymentHandler := grpchandler.NewPaymentHandler(paymentUseCase)
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	go func() {
		log.Println("🚀 Payment Service gRPC server listening on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gracefully...")
	grpcServer.GracefulStop()
	log.Println("✅ Payment Service stopped")
}
