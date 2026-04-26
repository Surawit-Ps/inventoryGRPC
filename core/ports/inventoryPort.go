package ports

import (
	"inventoryService/core/entity"
	"time"
)

type InventoryRepository interface {
	GetAvailability(roomTypeID string, date time.Time) (int, error)
	ReserveRooms(roomTypeID string, date time.Time, count int) error
	ReleaseRooms(roomTypeID string, date time.Time, count int) error
	CreateInventory(inv entity.RoomInventory) error
}

type InventoryService interface {
	CheckAvailability(roomTypeID string, date time.Time) (int, error)
	Reserve(roomTypeID string, date time.Time, count int) error
	Release(roomTypeID string, date time.Time, count int) error
	AddInventory(inv entity.RoomInventory) error
}