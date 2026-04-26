package repository

import (
	"inventoryService/core/entity"
	"time"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

type InventoryModel struct {
	ID         int    `gorm:"primaryKey"`
	RoomTypeId string `gorm:"index"`
	Date       string `gorm:"index"`
	Available  int
	Total      int
}

func (r *InventoryRepository) GetInventoryByRoomTypeAndDate(roomTypeId string, date string) (*entity.RoomInventory, error) {
	var inventoryModel InventoryModel
	result := r.db.Where("room_type_id = ? AND date = ?", roomTypeId, date).First(&inventoryModel)
	if result.Error != nil {
		return nil, result.Error
	}
	parsedDate, err := time.Parse("2006-01-02", inventoryModel.Date)
	if err != nil {
		return nil, err
	}
	inventory := &entity.RoomInventory{
		ID:         inventoryModel.ID,
		RoomTypeId: inventoryModel.RoomTypeId,
		Date:       parsedDate,
		Available:  inventoryModel.Available,
		Total:      inventoryModel.Total,
	}
	return inventory, nil
}

func (r *InventoryRepository) UpdateInventory(inventory *entity.RoomInventory) error {
	inventoryModel := InventoryModel{
		ID:         inventory.ID,
		RoomTypeId: inventory.RoomTypeId,
		Date:       inventory.Date.Format("2006-01-02"),
		Available:  inventory.Available,
		Total:      inventory.Total,
	}
	result := r.db.Save(&inventoryModel)
	return result.Error
}
