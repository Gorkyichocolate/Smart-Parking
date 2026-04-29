package grpc

import (
	"context"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

type Handler struct {
	usecase *usecase.PaymentUseCase
}

func NewHandler(u *usecase.PaymentUseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) CreatePayment(ctx context.Context) error {
	payment := domain.Payment{
		ID:     "123",
		Amount: 100.0,
		Name:   "John Doe",
		Email:  "john.doe@example.com",
	}

	return h.usecase.CreatePayment(payment)
}