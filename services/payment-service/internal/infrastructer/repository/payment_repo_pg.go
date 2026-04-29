package repository

import (
	"database/sql"
	"payment-service/internal/domain"
)

type paymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *paymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) CreatePayment(p domain.Payment) error {
	_, err := r.db.Exec(
		`INSERT INTO payments 
		(id, booking_id, user_id, amount, status, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID,
		p.BookingID,
		p.UserID,
		p.Amount,
		p.Status,
		p.PaymentMethod,
	)
	return err
}
func (r *paymentRepo) GetPaymentByID(id string) (domain.Payment, error) {
	var p domain.Payment

	err := r.db.QueryRow(
		`SELECT id, booking_id, user_id, amount, status, payment_method, created_at, updated_at
		 FROM payments WHERE id=$1`,
		id,
	).Scan(
		&p.ID,
		&p.BookingID,
		&p.UserID,
		&p.Amount,
		&p.Status,
		&p.PaymentMethod,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}
