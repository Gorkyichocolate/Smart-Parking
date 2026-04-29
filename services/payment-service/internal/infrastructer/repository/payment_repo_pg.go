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
	_, err := r.db.Exec("INSERT INTO payments (id, amount, name, email) VALUES ($1, $2, $3, $4)",
		p.ID, p.Amount, p.Name, p.Email)
	return err
}
