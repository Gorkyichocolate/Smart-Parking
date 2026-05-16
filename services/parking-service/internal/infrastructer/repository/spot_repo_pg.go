package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type SpotRepositoryPG struct {
	db *sql.DB
}

func NewSpotRepositoryPG(db *sql.DB) *SpotRepositoryPG {
	return &SpotRepositoryPG{db: db}
}

func (r *SpotRepositoryPG) GetByID(ctx context.Context, id string) (*domain.ParkingSpot, error) {
	query := `SELECT id, zone_id, spot_number, status, price_per_hour, created_at FROM parking_spots WHERE id = $1`
	var s domain.ParkingSpot
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.ZoneID, &s.SpotNumber, &s.Status, &s.PricePerHour, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrSpotNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get spot: %w", err)
	}
	return &s, nil
}

func (r *SpotRepositoryPG) GetByZoneID(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error) {
	query := `SELECT id, zone_id, spot_number, status, price_per_hour, created_at FROM parking_spots WHERE zone_id = $1`
	rows, err := r.db.QueryContext(ctx, query, zoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spots: %w", err)
	}
	defer rows.Close()
	var spots []*domain.ParkingSpot
	for rows.Next() {
		var s domain.ParkingSpot
		err := rows.Scan(&s.ID, &s.ZoneID, &s.SpotNumber, &s.Status, &s.PricePerHour, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spot: %w", err)
		}
		spots = append(spots, &s)
	}
	return spots, nil
}

func (r *SpotRepositoryPG) GetFreeSpots(ctx context.Context, zoneID string) ([]*domain.ParkingSpot, error) {
	query := `SELECT id, zone_id, spot_number, status, price_per_hour, created_at FROM parking_spots WHERE zone_id = $1 AND status = 'free'`
	rows, err := r.db.QueryContext(ctx, query, zoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get free spots: %w", err)
	}
	defer rows.Close()
	var spots []*domain.ParkingSpot
	for rows.Next() {
		var s domain.ParkingSpot
		err := rows.Scan(&s.ID, &s.ZoneID, &s.SpotNumber, &s.Status, &s.PricePerHour, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spot: %w", err)
		}
		spots = append(spots, &s)
	}
	return spots, nil
}

func (r *SpotRepositoryPG) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE parking_spots SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
