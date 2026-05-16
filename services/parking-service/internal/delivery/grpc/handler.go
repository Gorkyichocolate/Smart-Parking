package grpc

import (
	"context"
	

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	parkingpb "github.com/Gorkyichocolate/smart-parking-proto/gen/go/parking"
	"github.com/Gorkyichocolate/smart-parking/services/parking-service/internal/usecase"
)

type ParkingHandler struct {
	parkingpb.UnimplementedParkingServiceServer
	parkingUseCase *usecase.ParkingUseCase
}

func NewParkingHandler(parkingUseCase *usecase.ParkingUseCase) *ParkingHandler {
	return &ParkingHandler{parkingUseCase: parkingUseCase}
}

func (h *ParkingHandler) GetZones(ctx context.Context, req *parkingpb.GetZonesRequest) (*parkingpb.GetZonesResponse, error) {
	zones, err := h.parkingUseCase.GetZones(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get zones: %v", err)
	}

	var pbZones []*parkingpb.ParkingZone
	for _, z := range zones {
		pbZones = append(pbZones, &parkingpb.ParkingZone{
			Id:         z.ID,
			Name:       z.Name,
			Address:    z.Address,
			TotalSpots: int32(z.TotalSpots),
			CreatedAt:  timestamppb.New(z.CreatedAt),
		})
	}
	return &parkingpb.GetZonesResponse{Zones: pbZones}, nil
}

func (h *ParkingHandler) GetZone(ctx context.Context, req *parkingpb.GetZoneRequest) (*parkingpb.GetZoneResponse, error) {
	z, err := h.parkingUseCase.GetZone(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "zone not found: %v", err)
	}
	return &parkingpb.GetZoneResponse{Zone: &parkingpb.ParkingZone{
		Id: z.ID, Name: z.Name, Address: z.Address, TotalSpots: int32(z.TotalSpots), CreatedAt: timestamppb.New(z.CreatedAt),
	}}, nil
}

func (h *ParkingHandler) GetSpots(ctx context.Context, req *parkingpb.GetSpotsRequest) (*parkingpb.GetSpotsResponse, error) {
	spots, err := h.parkingUseCase.GetSpots(ctx, req.ZoneId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get spots: %v", err)
	}
	var pbSpots []*parkingpb.ParkingSpot
	for _, s := range spots {
		pbSpots = append(pbSpots, &parkingpb.ParkingSpot{
			Id: s.ID, ZoneId: s.ZoneID, SpotNumber: s.SpotNumber, Status: s.Status, PricePerHour: s.PricePerHour, CreatedAt: timestamppb.New(s.CreatedAt),
		})
	}
	return &parkingpb.GetSpotsResponse{Spots: pbSpots}, nil
}

func (h *ParkingHandler) GetFreeSpots(ctx context.Context, req *parkingpb.GetFreeSpotsRequest) (*parkingpb.GetFreeSpotsResponse, error) {
	spots, err := h.parkingUseCase.GetFreeSpots(ctx, req.ZoneId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get free spots: %v", err)
	}
	var pbSpots []*parkingpb.ParkingSpot
	for _, s := range spots {
		pbSpots = append(pbSpots, &parkingpb.ParkingSpot{
			Id: s.ID, ZoneId: s.ZoneID, SpotNumber: s.SpotNumber, Status: s.Status, PricePerHour: s.PricePerHour, CreatedAt: timestamppb.New(s.CreatedAt),
		})
	}
	return &parkingpb.GetFreeSpotsResponse{Spots: pbSpots}, nil
}

func (h *ParkingHandler) CreateBooking(ctx context.Context, req *parkingpb.CreateBookingRequest) (*parkingpb.CreateBookingResponse, error) {
	input := usecase.CreateBookingInput{
		UserID:    req.UserId,
		SpotID:    req.SpotId,
		StartTime: req.StartTime.AsTime(),
		EndTime:   req.EndTime.AsTime(),
	}

	booking, err := h.parkingUseCase.CreateBooking(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
	}

	return &parkingpb.CreateBookingResponse{Booking: &parkingpb.Booking{
		Id: booking.ID, UserId: booking.UserID, SpotId: booking.SpotID,
		StartTime: timestamppb.New(booking.StartTime), EndTime: timestamppb.New(booking.EndTime),
		TotalPrice: booking.TotalPrice, Status: booking.Status,
		CreatedAt: timestamppb.New(booking.CreatedAt), UpdatedAt: timestamppb.New(booking.UpdatedAt),
	}}, nil
}

func (h *ParkingHandler) GetBooking(ctx context.Context, req *parkingpb.GetBookingRequest) (*parkingpb.GetBookingResponse, error) {
	b, err := h.parkingUseCase.GetBooking(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "booking not found: %v", err)
	}
	return &parkingpb.GetBookingResponse{Booking: &parkingpb.Booking{
		Id: b.ID, UserId: b.UserID, SpotId: b.SpotID,
		StartTime: timestamppb.New(b.StartTime), EndTime: timestamppb.New(b.EndTime),
		TotalPrice: b.TotalPrice, Status: b.Status,
		CreatedAt: timestamppb.New(b.CreatedAt), UpdatedAt: timestamppb.New(b.UpdatedAt),
	}}, nil
}

func (h *ParkingHandler) GetUserBookings(ctx context.Context, req *parkingpb.GetUserBookingsRequest) (*parkingpb.GetUserBookingsResponse, error) {
	bookings, err := h.parkingUseCase.GetUserBookings(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get bookings: %v", err)
	}
	var pbBookings []*parkingpb.Booking
	for _, b := range bookings {
		pbBookings = append(pbBookings, &parkingpb.Booking{
			Id: b.ID, UserId: b.UserID, SpotId: b.SpotID,
			StartTime: timestamppb.New(b.StartTime), EndTime: timestamppb.New(b.EndTime),
			TotalPrice: b.TotalPrice, Status: b.Status,
			CreatedAt: timestamppb.New(b.CreatedAt), UpdatedAt: timestamppb.New(b.UpdatedAt),
		})
	}
	return &parkingpb.GetUserBookingsResponse{Bookings: pbBookings}, nil
}

func (h *ParkingHandler) CancelBooking(ctx context.Context, req *parkingpb.CancelBookingRequest) (*parkingpb.CancelBookingResponse, error) {
	if err := h.parkingUseCase.CancelBooking(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel booking: %v", err)
	}
	return &parkingpb.CancelBookingResponse{Success: true}, nil
}
