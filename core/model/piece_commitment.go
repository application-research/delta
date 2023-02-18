package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PieceCommitment struct {
	ID                int64     `gorm:"primaryKey"`
	Cid               string    `json:"cid"`
	Piece             string    `json:"piece"`
	Size              int64     `json:"size"`
	PaddedPieceSize   uint64    `json:"padded_piece_size"`
	UnPaddedPieceSize uint64    `json:"unnpadded_piece_size"`
	Status            string    `json:"status"` // open, in-progress, completed (closed).
	LastMessage       string    `json:"last_message"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (u *PieceCommitment) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMiner Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentMiner %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *PieceCommitment) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMiner Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentMiner %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *PieceCommitment) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentMiner Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentMiner %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
