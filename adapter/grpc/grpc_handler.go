package grpc

import (
	"context"
	"inventoryService/core/ports"
	"inventoryService/proto/pb"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryGRPCHandler struct {
	service ports.InventoryService
	pb.UnimplementedInventoryServiceServer
}

func NewInventoryGRPCHandler(service ports.InventoryService) *InventoryGRPCHandler {
	return &InventoryGRPCHandler{service: service}
}

// CheckAvailability handles the CheckAvailability RPC call
func (h *InventoryGRPCHandler) CheckAvailability(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	if req.RoomTypeId == "" || req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "room_type_id and date are required")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date format, expected YYYY-MM-DD")
	}

	available, err := h.service.CheckAvailability(req.RoomTypeId, date)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check availability: "+err.Error())
	}

	return &pb.CheckResponse{
		AvailableRooms: int32(available),
	}, nil
}

// ReserveRoom handles the ReserveRoom RPC call
func (h *InventoryGRPCHandler) ReserveRoom(ctx context.Context, req *pb.ReserveRequest) (*pb.ReserveResponse, error) {
	if req.BookingId == "" || req.RoomTypeId == "" || req.CheckInDate == "" || req.CheckOutDate == "" || req.RoomCount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "booking_id, room_type_id, check_in_date, check_out_date, and valid room_count are required")
	}

	checkInDate, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid check_in_date format, expected YYYY-MM-DD")
	}

	checkOutDate, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid check_out_date format, expected YYYY-MM-DD")
	}

	if checkInDate.After(checkOutDate) || checkInDate.Equal(checkOutDate) {
		return nil, status.Error(codes.InvalidArgument, "check_in_date must be before check_out_date")
	}

	// Reserve rooms for each date in the range
	currentDate := checkInDate
	for currentDate.Before(checkOutDate) {
		err := h.service.Reserve(req.RoomTypeId, currentDate, int(req.RoomCount))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to reserve rooms: "+err.Error())
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &pb.ReserveResponse{
		Success: true,
		Message: "Successfully reserved " + strconv.Itoa(int(req.RoomCount)) + " room(s) for booking " + req.BookingId,
	}, nil
}

// ReleaseRoom handles the ReleaseRoom RPC call
func (h *InventoryGRPCHandler) ReleaseRoom(ctx context.Context, req *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	if req.BookingId == "" || req.RoomTypeId == "" || req.RoomCount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "booking_id, room_type_id, and valid room_count are required")
	}

	// For now, we'll release rooms for today's date
	// In a real application, you'd need to track booking dates
	currentDate := time.Now()
	
	err := h.service.Release(req.RoomTypeId, currentDate, int(req.RoomCount))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to release rooms: "+err.Error())
	}

	return &pb.ReleaseResponse{
		Success: true,
	}, nil
}


