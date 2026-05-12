// internal/delivery/grpc/handler.go
package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	paymentpb "github.com/GorkyiChocolate/smart-parking-proto/gen/go/payment"

	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/usecase"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentUseCase *usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	input := usecase.CreatePaymentInput{
		BookingID:     req.BookingId,
		UserID:        req.UserId,
		Amount:        req.Amount,
		PaymentMethod: domain.PaymentMethod(req.PaymentMethod),
		UserEmail:     req.UserEmail,
	}

	payment, err := h.paymentUseCase.CreatePayment(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return &paymentpb.CreatePaymentResponse{
		Id:     payment.ID,
		Status: string(payment.Status),
	}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	payment, err := h.paymentUseCase.GetPayment(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &paymentpb.GetPaymentResponse{
		Payment: &paymentpb.Payment{
			Id:            payment.ID,
			BookingId:     payment.BookingID,
			UserId:        payment.UserID,
			Amount:        payment.Amount,
			Status:        string(payment.Status),
			PaymentMethod: string(payment.PaymentMethod),
			CreatedAt:     timestamppb.New(payment.CreatedAt),
			UpdatedAt:     timestamppb.New(payment.UpdatedAt),
		},
	}, nil
}
