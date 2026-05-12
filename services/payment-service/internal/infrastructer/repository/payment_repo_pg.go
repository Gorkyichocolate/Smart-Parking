// internal/infrastructer/repository/payment_repo_pg.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/GorkyiChocolate/smart-parking/services/payment-service/internal/domain"
)

type PaymentRepositoryPG struct {
	db *sql.DB
}

func NewPaymentRepositoryPG(db *sql.DB) *PaymentRepositoryPG {
	return &PaymentRepositoryPG{db: db}
}

func (r *PaymentRepositoryPG) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, booking_id, user_id, amount, status, payment_method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.BookingID,
		payment.UserID,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	payment.CreatedAt = now
	payment.UpdatedAt = now

	return nil
}

func (r *PaymentRepositoryPG) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `
		SELECT id, booking_id, user_id, amount, status, payment_method, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepositoryPG) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error {
	query := `UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}

func (r *PaymentRepositoryPG) GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error) {
	query := `
        SELECT id, booking_id, user_id, amount, status, payment_method, created_at, updated_at
        FROM payments
        WHERE booking_id = $1
    `

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, bookingID).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Не ошибка, просто нет платежа
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by booking: %w", err)
	}

	return &payment, nil
}
