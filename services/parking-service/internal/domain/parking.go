package domain

import "time"

type ParkingZone struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	TotalSpots int       `json:"total_spots"`
	CreatedAt  time.Time `json:"created_at"`
}

type ParkingSpot struct {
	ID           string    `json:"id"`
	ZoneID       string    `json:"zone_id"`
	SpotNumber   string    `json:"spot_number"`
	Status       string    `json:"status"`
	PricePerHour float64   `json:"price_per_hour"`
	CreatedAt    time.Time `json:"created_at"`
}

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	SpotID     string    `json:"spot_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrZoneNotFound    = &DomainError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
	ErrSpotNotFound    = &DomainError{Code: "SPOT_NOT_FOUND", Message: "spot not found"}
	ErrBookingNotFound = &DomainError{Code: "BOOKING_NOT_FOUND", Message: "booking not found"}
	ErrSpotNotFree     = &DomainError{Code: "SPOT_NOT_FREE", Message: "spot is not free"}
	ErrInvalidData     = &DomainError{Code: "INVALID_DATA", Message: "invalid data"}
)
