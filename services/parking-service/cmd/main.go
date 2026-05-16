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

	parkingpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/parking"
	grpchandler "github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/delivery/grpc"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/infrastructer/repository"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/usecase"
)

func main() {
	databaseURL := os.Getenv("PARKING_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://parking_service_user:123456@localhost:5434/parking_db?sslmode=disable"
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

	zoneRepo := repository.NewZoneRepositoryPG(db)
	spotRepo := repository.NewSpotRepositoryPG(db)
	bookingRepo := repository.NewBookingRepositoryPG(db)

	parkingUseCase := usecase.NewParkingUseCase(zoneRepo, spotRepo, bookingRepo)

	grpcServer := grpc.NewServer()
	parkingHandler := grpchandler.NewParkingHandler(parkingUseCase)
	parkingpb.RegisterParkingServiceServer(grpcServer, parkingHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Parking Service gRPC server listening on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	log.Println("Parking Service stopped")
}
