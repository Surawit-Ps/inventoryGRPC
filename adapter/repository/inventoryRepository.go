package repository

import (
	"inventoryService/core/entity"
	"strconv"
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

func (r *InventoryRepository) GetAvailability(roomTypeId string, date time.Time) (int, error) {
	var inventoryModel InventoryModel
	roomid,err :=  strconv.Atoi(roomTypeId)
	if err != nil {
		return 0, err
	}
	result := r.db.Where("id = ? OR date = ?", roomid, date.Format("2006-01-02")).First(&inventoryModel)
	if result.Error != nil {
		return 0, result.Error
	}
	parsedDate, err := time.Parse("2006-01-02", inventoryModel.Date)
	if err != nil {
		return 0, err
	}
	inventory := &entity.RoomInventory{
		ID:         inventoryModel.ID,
		RoomTypeId: inventoryModel.RoomTypeId,
		Date:       parsedDate,
		Available:  inventoryModel.Available,
		Total:      inventoryModel.Total,
	}
	return inventory.Available, nil
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

func (r *InventoryRepository) CreateInventory(inventory entity.RoomInventory) error {
	inventoryModel := InventoryModel{
		RoomTypeId: inventory.RoomTypeId,
		Date:       inventory.Date.Format("2006-01-02"),
		Available:  inventory.Available,
		Total:      inventory.Total,
	}
	result := r.db.Create(&inventoryModel)
	return result.Error
}

func (r *InventoryRepository) ReserveRooms(roomTypeId string, date time.Time, count int) error {
	var inventoryModel InventoryModel
	result := r.db.Where("room_type_id = ? AND date = ?", roomTypeId, date.Format("2006-01-02")).First(&inventoryModel)
	if result.Error != nil {
		return result.Error
	}
	if inventoryModel.Available < count {
		return nil // Not enough rooms available
	}

	inventoryModel.Available -= count
	result = r.db.Save(&inventoryModel)
	return result.Error
}

func (r *InventoryRepository) ReleaseRooms(roomTypeId string, date time.Time, count int) error {
	var inventoryModel InventoryModel
	result := r.db.Where("room_type_id = ? AND date = ?", roomTypeId, date.Format("2006-01-02")).First(&inventoryModel)
	if result.Error != nil {
		return result.Error
	}
	inventoryModel.Available += count
	if inventoryModel.Available > inventoryModel.Total {
		inventoryModel.Available = inventoryModel.Total // Ensure available does not exceed total
	}
	result = r.db.Save(&inventoryModel)
	return result.Error
}