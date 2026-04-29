package repository

import "payment-service/internal/domain"

type InvoiceRepository interface {
	Create(invoice domain.Invoice) error
}
