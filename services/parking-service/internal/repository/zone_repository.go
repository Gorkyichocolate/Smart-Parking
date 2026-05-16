package repository

import (
	"context"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type ZoneRepository interface {
	GetAll(ctx context.Context) ([]*domain.ParkingZone, error)
	GetByID(ctx context.Context, id string) (*domain.ParkingZone, error)
}
