package db_models

import (
	"time"
)

type RepairRequest struct {
	ID        int64     `gorm:"primaryKey"`
	ObjectId  int64     `json:"object_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
