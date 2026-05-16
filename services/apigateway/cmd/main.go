package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Gorkyichocolate/smart-parking/services/apigateway/internal/handler"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" { port = "8080" }
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" { authAddr = "localhost:50051" }
	parkingAddr := os.Getenv("PARKING_SERVICE_ADDR")
	if parkingAddr == "" { parkingAddr = "localhost:50053" }
	paymentAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentAddr == "" { paymentAddr = "localhost:50052" }

	gw := handler.NewGateway(authAddr, parkingAddr, paymentAddr)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", gw.Health)
	mux.HandleFunc("/api/auth/register", gw.Register)
	mux.HandleFunc("/api/auth/login", gw.Login)
	mux.HandleFunc("/api/auth/validate", gw.ValidateToken)
	mux.HandleFunc("/api/parking/zones", gw.GetZones)
	mux.HandleFunc("/api/parking/zone", gw.GetZone)
	mux.HandleFunc("/api/parking/spots", gw.GetSpots)
	mux.HandleFunc("/api/parking/free-spots", gw.GetFreeSpots)
	mux.HandleFunc("/api/parking/booking", gw.CreateBooking)
	mux.HandleFunc("/api/parking/booking/get", gw.GetBooking)
	mux.HandleFunc("/api/parking/booking/user", gw.GetUserBookings)
	mux.HandleFunc("/api/parking/booking/cancel", gw.CancelBooking)
	mux.HandleFunc("/api/payment/create", gw.CreatePayment)
	mux.HandleFunc("/api/payment/get", gw.GetPayment)

	log.Printf("API Gateway listening on :%s", port)
	http.ListenAndServe(":"+port, corsMiddleware(mux))
}
