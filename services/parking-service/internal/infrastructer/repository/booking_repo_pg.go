package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type BookingRepositoryPG struct {
	db *sql.DB
}

func NewBookingRepositoryPG(db *sql.DB) *BookingRepositoryPG {
	return &BookingRepositoryPG{db: db}
}

func (r *BookingRepositoryPG) Create(ctx context.Context, booking *domain.Booking) error {
	query := `INSERT INTO bookings (id, user_id, spot_id, start_time, end_time, total_price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, booking.ID, booking.UserID, booking.SpotID, booking.StartTime, booking.EndTime, booking.TotalPrice, booking.Status, now, now)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	booking.CreatedAt = now
	booking.UpdatedAt = now
	return nil
}

func (r *BookingRepositoryPG) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	query := `SELECT id, user_id, spot_id, start_time, end_time, total_price, status, created_at, updated_at FROM bookings WHERE id = $1`
	var b domain.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.UserID, &b.SpotID, &b.StartTime, &b.EndTime, &b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrBookingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return &b, nil
}

func (r *BookingRepositoryPG) GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error) {
	query := `SELECT id, user_id, spot_id, start_time, end_time, total_price, status, created_at, updated_at FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	defer rows.Close()
	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(&b.ID, &b.UserID, &b.SpotID, &b.StartTime, &b.EndTime, &b.TotalPrice, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *BookingRepositoryPG) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}
