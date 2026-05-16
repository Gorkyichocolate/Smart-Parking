package repository

import (
	"context"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type SpotRepository interface {
	GetByID(ctx context.Context, id string) (*domain.ParkingSpot, error)
	GetByZoneID(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error)
	GetFreeSpots(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
