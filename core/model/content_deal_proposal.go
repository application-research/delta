package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ContentDealProposal struct {
	ID      int64  `gorm:"primaryKey"`
	Content int64  `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Meta    string `json:"meta"`
}

func (u *ContentDealProposal) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDealProposal Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDealProposalParameters %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDealProposal) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDealProposal Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDealProposalParameters %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDealProposal) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentDealProposalParameters Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentDealProposalParameters %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
