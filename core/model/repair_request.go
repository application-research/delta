package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type RepairRequest struct {
	ID        int64     `gorm:"primaryKey"`
	ObjectId  int64     `json:"object_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *RepairRequest) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "Content Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("Content %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *RepairRequest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "Content Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("Content %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *RepairRequest) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After Content Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After content %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
