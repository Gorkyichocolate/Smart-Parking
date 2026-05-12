// services/payment-service/internal/delivery/grpc/handler.go
package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	paymentpb "github.com/GorkyiChocolate/smart-parking-proto/gen/go/payment"
	"github.com/GorkyiChocolate/smart-parking/pkg/metrics"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/usecase"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentUseCase *usecase.PaymentUseCase
	metrics        *metrics.Metrics
}

func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		metrics:        metrics.GetMetrics(),
	}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	start := time.Now()

	// Record metrics
	h.metrics.GrpcRequestsTotal.WithLabelValues("CreatePayment", "pending").Inc()

	input := usecase.CreatePaymentInput{
		BookingID:     req.BookingId,
		UserID:        req.UserId,
		Amount:        req.Amount,
		PaymentMethod: domain.PaymentMethod(req.PaymentMethod),
		UserEmail:     req.UserEmail,
	}

	payment, err := h.paymentUseCase.CreatePayment(ctx, input)
	if err != nil {
		h.metrics.GrpcRequestErrors.WithLabelValues("CreatePayment").Inc()
		h.metrics.GrpcRequestsTotal.WithLabelValues("CreatePayment", "error").Inc()
		h.metrics.PaymentFailedTotal.Inc()
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	// Record duration
	duration := time.Since(start).Seconds()
	h.metrics.GrpcRequestDuration.WithLabelValues("CreatePayment").Observe(duration)
	h.metrics.GrpcRequestsTotal.WithLabelValues("CreatePayment", "success").Inc()
	h.metrics.PaymentCreatedTotal.Inc()

	return &paymentpb.CreatePaymentResponse{
		Id:     payment.ID,
		Status: string(payment.Status),
	}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	start := time.Now()

	h.metrics.GrpcRequestsTotal.WithLabelValues("GetPayment", "pending").Inc()

	payment, err := h.paymentUseCase.GetPayment(ctx, req.Id)
	if err != nil {
		h.metrics.GrpcRequestErrors.WithLabelValues("GetPayment").Inc()
		h.metrics.GrpcRequestsTotal.WithLabelValues("GetPayment", "error").Inc()
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	duration := time.Since(start).Seconds()
	h.metrics.GrpcRequestDuration.WithLabelValues("GetPayment").Observe(duration)
	h.metrics.GrpcRequestsTotal.WithLabelValues("GetPayment", "success").Inc()

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
