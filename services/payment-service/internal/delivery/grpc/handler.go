package grpc

import (
	"context"

	"payment-service/internal/domain"
	"payment-service/internal/usecase"
	pb "smart-parking/proto"
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
		ID:        req.Id,
		BookingID: req.BookingId,
		UserID:    req.UserId,
		Amount:    req.Amount,
	}

	err := h.usecase.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		Success: true,
	}, nil
}
