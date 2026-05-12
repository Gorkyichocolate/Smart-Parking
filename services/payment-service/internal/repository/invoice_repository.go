// internal/repository/invoice_repository.go
package repository

import (
	"context"

	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByPaymentID(ctx context.Context, paymentID string) (*domain.Invoice, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error)
}
