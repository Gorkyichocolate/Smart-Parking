package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/repository"
)

type ParkingUseCase struct {
	zoneRepo    repository.ZoneRepository
	spotRepo    repository.SpotRepository
	bookingRepo repository.BookingRepository
}

func NewParkingUseCase(zoneRepo repository.ZoneRepository, spotRepo repository.SpotRepository, bookingRepo repository.BookingRepository) *ParkingUseCase {
	return &ParkingUseCase{zoneRepo: zoneRepo, spotRepo: spotRepo, bookingRepo: bookingRepo}
}

func (uc *ParkingUseCase) GetZones(ctx context.Context) ([]*domain.ParkingZone, error) {
	return uc.zoneRepo.GetAll(ctx)
}

func (uc *ParkingUseCase) GetZone(ctx context.Context, id string) (*domain.ParkingZone, error) {
	return uc.zoneRepo.GetByID(ctx, id)
}

func (uc *ParkingUseCase) GetSpots(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error) {
	return uc.spotRepo.GetByZoneID(ctx, zoneID)
}

func (uc *ParkingUseCase) GetFreeSpots(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error) {
	return uc.spotRepo.GetFreeSpots(ctx, zoneID)
}

type CreateBookingInput struct {
	UserID    string
	SpotID    string
	StartTime time.Time
	EndTime   time.Time
}

func (uc *ParkingUseCase) CreateBooking(ctx context.Context, input CreateBookingInput) (*domain.Booking, error) {
	spot, err := uc.spotRepo.GetByID(ctx, input.SpotID)
	if err != nil {
		return nil, err
	}
	if spot.Status != "free" {
		return nil, domain.ErrSpotNotFree
	}
	if input.EndTime.Before(input.StartTime) {
		return nil, domain.ErrInvalidData
	}

	hours := input.EndTime.Sub(input.StartTime).Hours()
	totalPrice := spot.PricePerHour * hours

	booking := &domain.Booking{
		ID:         uuid.New().String(),
		UserID:     input.UserID,
		SpotID:     input.SpotID,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		TotalPrice: totalPrice,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.bookingRepo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	uc.spotRepo.UpdateStatus(ctx, input.SpotID, "occupied")

	return booking, nil
}

func (uc *ParkingUseCase) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	return uc.bookingRepo.GetByID(ctx, id)
}

func (uc *ParkingUseCase) GetUserBookings(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return uc.bookingRepo.GetByUserID(ctx, userID)
}

func (uc *ParkingUseCase) CancelBooking(ctx context.Context, id string) error {
	booking, err := uc.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.bookingRepo.UpdateStatus(ctx, id, "cancelled"); err != nil {
		return err
	}
	uc.spotRepo.UpdateStatus(ctx, booking.SpotID, "free")
	return nil
}
