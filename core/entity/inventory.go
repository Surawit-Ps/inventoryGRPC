package entity

import "time"

type RoomInventory struct {
	ID        int       `json:"id"`
	RoomTypeId string	`json:"room_type_id"`
	Date	  time.Time `json:"date"`
	Available int       `json:"available"`
	Total     int       `json:"total"`
}