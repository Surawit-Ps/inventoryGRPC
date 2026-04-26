package service

import (
	"inventoryService/core/entity"
	"inventoryService/core/ports"
	"github.com/go-redis/redis/v8"
	"time"
)

type InventoryService struct {
	repo ports.InventoryRepository
	rds  *redis.Client
}

func NewInventoryService(repo ports.InventoryRepository, rds *redis.Client) *InventoryService {
	return &InventoryService{repo: repo, rds: rds}
}

func (s *InventoryService) CheckAvailability(roomTypeID string, date time.Time) (int, error) {
	return s.repo.GetAvailability(roomTypeID, date)
}

func (s *InventoryService) Reserve(roomTypeID string, date time.Time, count int) error {
	lockkey := "lock:" + roomTypeID + ":" + date.Format("2006-01-02")
	err := s.rds.SetNX(s.rds.Context(), lockkey, "locked", 5*time.Second).Err()
	if err != nil {
		return err
	}
	defer s.rds.Del(s.rds.Context(), lockkey)

	var available int
	available, err = s.repo.GetAvailability(roomTypeID, date)
	if err != nil {
		return err
	}
	if available < count {
		//// rollback 
		return nil // Not enough rooms available
	}

	err = s.repo.ReserveRooms(roomTypeID, date, count)
	if err != nil {
		//// rollback
		return err
	}

	return nil
}

func (s *InventoryService) Release(roomTypeID string, date time.Time, count int) error {

	return s.repo.ReleaseRooms(roomTypeID, date, count)
}

func (s *InventoryService) AddInventory(inv entity.RoomInventory) error {
	return s.repo.CreateInventory(inv)
}
