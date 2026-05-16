package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	authpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/auth"
	parkingpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/parking"
	paymentpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/payment"
)

type Gateway struct {
	authClient    authpb.AuthServiceClient
	parkingClient parkingpb.ParkingServiceClient
	paymentClient paymentpb.PaymentServiceClient
}

func NewGateway(authAddr, parkingAddr, paymentAddr string) *Gateway {
	authConn, _ := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	parkingConn, _ := grpc.NewClient(parkingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	paymentConn, _ := grpc.NewClient(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &Gateway{
		authClient:    authpb.NewAuthServiceClient(authConn),
		parkingClient: parkingpb.NewParkingServiceClient(parkingConn),
		paymentClient: paymentpb.NewPaymentServiceClient(paymentConn),
	}
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (g *Gateway) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.authClient.Register(ctx, &authpb.RegisterRequest{
		Email: req.Email, Password: req.Password, FullName: req.FullName, Role: req.Role,
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.authClient.Login(ctx, &authpb.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		jsonError(w, err.Error(), 401)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetZones(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetZones(ctx, &parkingpb.GetZonesRequest{})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetZone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetZone(ctx, &parkingpb.GetZoneRequest{Id: id})
	if err != nil {
		jsonError(w, "not found", 404)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetSpots(w http.ResponseWriter, r *http.Request) {
	zoneID := r.URL.Query().Get("zone_id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetSpots(ctx, &parkingpb.GetSpotsRequest{ZoneId: zoneID})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetFreeSpots(w http.ResponseWriter, r *http.Request) {
	zoneID := r.URL.Query().Get("zone_id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetFreeSpots(ctx, &parkingpb.GetFreeSpotsRequest{ZoneId: zoneID})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string `json:"user_id"`
		SpotID    string `json:"spot_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.CreateBooking(ctx, &parkingpb.CreateBookingRequest{
		UserId: req.UserID, SpotId: req.SpotID,
		StartTime: timestamppb.New(startTime), EndTime: timestamppb.New(endTime),
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetBooking(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetBooking(ctx, &parkingpb.GetBookingRequest{Id: id})
	if err != nil {
		jsonError(w, "not found", 404)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.GetUserBookings(ctx, &parkingpb.GetUserBookingsRequest{UserId: userID})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) CancelBooking(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.parkingClient.CancelBooking(ctx, &parkingpb.CancelBookingRequest{Id: id})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookingID     string  `json:"booking_id"`
		UserID        string  `json:"user_id"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		UserEmail     string  `json:"user_email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.paymentClient.CreatePayment(ctx, &paymentpb.CreatePaymentRequest{
		BookingId: req.BookingID, UserId: req.UserID, Amount: req.Amount,
		PaymentMethod: req.PaymentMethod, UserEmail: req.UserEmail,
	})
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	resp, err := g.paymentClient.GetPayment(ctx, &paymentpb.GetPaymentRequest{Id: id})
	if err != nil {
		jsonError(w, "not found", 404)
		return
	}
	jsonOK(w, resp)
}

func (g *Gateway) Health(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"status": "ok"})
	fmt.Println("health check")
}
