package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/domain"
)

type ZoneRepositoryPG struct {
	db *sql.DB
}

func NewZoneRepositoryPG(db *sql.DB) *ZoneRepositoryPG {
	return &ZoneRepositoryPG{db: db}
}

func (r *ZoneRepositoryPG) GetAll(ctx context.Context) ([]*domain.ParkingZone, error) {
	query := `SELECT id, name, address, total_spots, created_at FROM parking_zones`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones: %w", err)
	}
	defer rows.Close()
	var zones []*domain.ParkingZone
	for rows.Next() {
		var z domain.ParkingZone
		err := rows.Scan(&z.ID, &z.Name, &z.Address, &z.TotalSpots, &z.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan zone: %w", err)
		}
		zones = append(zones, &z)
	}
	return zones, nil
}

func (r *ZoneRepositoryPG) GetByID(ctx context.Context, id string) (*domain.ParkingZone, error) {
	query := `SELECT id, name, address, total_spots, created_at FROM parking_zones WHERE id = $1`
	var z domain.ParkingZone
	err := r.db.QueryRowContext(ctx, query, id).Scan(&z.ID, &z.Name, &z.Address, &z.TotalSpots, &z.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrZoneNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}
	return &z, nil
}
