package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ContentDealProposalParameters struct {
	ID                 int64     `gorm:"primaryKey"`
	Content            int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Label              string    `json:"label,omitempty"`
	Duration           int64     `json:"duration,omitempty"`
	StartEpoch         int64     `json:"start_epoch,omitempty"`
	EndEpoch           int64     `json:"end_epoch,omitempty"`
	TransferParams     string    `json:"transfer_params,omitempty"`
	RemoveUnsealedCopy bool      `json:"remove_unsealed_copy,omitempty"`
	SkipIPNIAnnounce   bool      `json:"skip_ipni_announce,omitempty"`
	CreatedAt          time.Time `json:"created_at" json:"created-at"`
	UpdatedAt          time.Time `json:"updated_at" json:"updated-at"`
}

func (u *ContentDealProposalParameters) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDealProposalParameters Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDealProposalParameters %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDealProposalParameters) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDealProposalParameters Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDealProposalParameters %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDealProposalParameters) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentDealProposalParameters Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentDealProposalParameters %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
