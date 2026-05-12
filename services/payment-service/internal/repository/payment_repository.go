// internal/repository/payment_repository.go
package repository

import (
	"context"

	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error
}
