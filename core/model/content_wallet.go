package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ContentWallet struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Wallet    string    `json:"wallet_meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *ContentWallet) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "Content Wallet Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentWallet %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentWallet) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMiner Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentWallet %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentWallet) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentWallet Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentWallet %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
