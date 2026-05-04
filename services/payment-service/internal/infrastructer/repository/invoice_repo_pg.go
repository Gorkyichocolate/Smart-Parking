package repository

import (
	"database/sql"
	"payment-service/internal/domain"
)

type invoiceRepo struct {
	db *sql.DB
}

// Create implements [repository.InvoiceRepository].
func (r *invoiceRepo) Create(invoice domain.Invoice) error {
	panic("unimplemented")
}

func NewInvoiceRepo(db *sql.DB) *invoiceRepo {
	return &invoiceRepo{db: db}
}

func (r *invoiceRepo) CreateInvoice(inv domain.Invoice) error {
	_, err := r.db.Exec(
		`INSERT INTO invoices 
		(id, payment_id, user_id, amount, pdf_url) 
		VALUES ($1, $2, $3, $4, $5)`,
		inv.ID, inv.PaymentID, inv.UserID, inv.Amount, inv.PDFUrl,
	)
	return err
}
