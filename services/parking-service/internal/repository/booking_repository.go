package repository

import (
	"context"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
