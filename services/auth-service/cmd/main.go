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

	authpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth"
	grpchandler "github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/delivery/grpc"
	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/infrastructer/repository"
	"github.com/Gorkyichocolate/smart-parking/services/auth-service/internal/usecase"
)

func main() {
	databaseURL := os.Getenv("AUTH_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://auth_service_user:123456@localhost:5433/auth_db?sslmode=disable"
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

	userRepo := repository.NewUserRepositoryPG(db)
	authUseCase := usecase.NewAuthUseCase(userRepo)

	grpcServer := grpc.NewServer()
	authHandler := grpchandler.NewAuthHandler(authUseCase)
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Auth Service gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	log.Println("Auth Service stopped")
}
