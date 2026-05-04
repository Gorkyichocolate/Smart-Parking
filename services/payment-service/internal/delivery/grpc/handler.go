package grpc

import (
	"context"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"

	pb "github.com/Gorkyichocolate/Smart-Parking-Proto/proto/payment"
)

type Handler struct {
	pb.UnimplementedPaymentServiceServer
	usecase *usecase.PaymentUsecase
}

func NewHandler(u *usecase.PaymentUsecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {

	payment := domain.Payment{
		BookingID: req.BookingId,
		UserID:    req.UserId,
		Amount:    req.Amount,
	}

	err := h.usecase.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		Id:     payment.ID,
		Status: "created",
	}, nil
}
