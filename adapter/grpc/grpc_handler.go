package grpc

import (
	"context"
	"inventoryService/core/entity"
	"inventoryService/core/ports"
	"inventoryService/proto/pb"
	inventorypb "inventoryService/proto/pb"
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

func (h *InventoryGRPCHandler) CheckAvailability(ctx context.Context, req *inventorypb.CheckAvailabilityRequest) (*inventorypb.CheckAvailabilityResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}
	available, err := h.service.CheckAvailability(req.RoomTypeId, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check availability: %v", err)
	}
	return &inventorypb.CheckAvailabilityResponse{Available: int32(available)}, nil
}

func (h *InventoryGRPCHandler) Reserve(ctx context.Context, req *inventorypb.ReserveRequest) (*inventorypb.ReserveResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}
	err = h.service.Reserve(req.RoomTypeId, date, int(req.Count))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reserve rooms: %v", err)
	}
	return &inventorypb.ReserveResponse{}, nil
}

func (h *InventoryGRPCHandler) Release(ctx context.Context, req *inventorypb.ReleaseRequest) (*inventorypb.ReleaseResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}
	err = h.service.Release(req.RoomTypeId, date, int(req.Count))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release rooms: %v", err)
	}
	return &inventorypb.ReleaseResponse{}, nil
}

func (h *InventoryGRPCHandler) AddInventory(ctx context.Context, req *inventorypb.AddInventoryRequest) (*inventorypb.AddInventoryResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %v", err)
	}
	inv := entity.RoomInventory{
		RoomTypeId: req.RoomTypeId,
		Date:       date,
		Available:  int(req.Available),
		Total:      int(req.Total),
	}
	err = h.service.AddInventory(inv)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add inventory: %v", err)
	}
	return &inventorypb.AddInventoryResponse{}, nil
}