package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Content struct {
	ID                int64     `gorm:"primaryKey"`
	Name              string    `json:"name"`
	Size              int64     `json:"size"`
	Cid               string    `json:"cid"`
	RequestingApiKey  string    `json:"requesting_api_key,omitempty"`
	PieceCommitmentId int64     `json:"piece_commitment_id,omitempty"`
	Status            string    `json:"status"`
	ConnectionMode    string    `json:"connection_mode"` // offline or online
	LastMessage       string    `json:"last_message"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (u *Content) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "Content Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("Content %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *Content) BeforeCreate(tx *gorm.DB) (err error) {
	// check if it has a miner set

	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "Content Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("Content %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *Content) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After Content Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After content %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
