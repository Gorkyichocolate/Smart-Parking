// internal/infrastructer/repository/invoice_repo_pg.go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
)

type InvoiceRepositoryPG struct {
	db *sql.DB
}

func NewInvoiceRepositoryPG(db *sql.DB) *InvoiceRepositoryPG {
	return &InvoiceRepositoryPG{db: db}
}

func (r *InvoiceRepositoryPG) Create(ctx context.Context, invoice *domain.Invoice) error {
	query := `
        INSERT INTO invoices (id, payment_id, user_id, amount, pdf_url, issued_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.db.ExecContext(ctx, query,
		invoice.ID,
		invoice.PaymentID,
		invoice.UserID,
		invoice.Amount,
		invoice.PDFURL,
		invoice.IssuedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	return nil
}

func (r *InvoiceRepositoryPG) GetByPaymentID(ctx context.Context, paymentID string) (*domain.Invoice, error) {
	query := `
        SELECT id, payment_id, user_id, amount, pdf_url, issued_at
        FROM invoices
        WHERE payment_id = $1
    `

	var invoice domain.Invoice
	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(
		&invoice.ID,
		&invoice.PaymentID,
		&invoice.UserID,
		&invoice.Amount,
		&invoice.PDFURL,
		&invoice.IssuedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

func (r *InvoiceRepositoryPG) GetByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error) {
	query := `
        SELECT id, payment_id, user_id, amount, pdf_url, issued_at
        FROM invoices
        WHERE user_id = $1
        ORDER BY issued_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.PaymentID,
			&invoice.UserID,
			&invoice.Amount,
			&invoice.PDFURL,
			&invoice.IssuedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}
