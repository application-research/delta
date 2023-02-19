package model

import (
	"fmt"
	"github.com/libp2p/go-libp2p/core/protocol"
	"gorm.io/gorm"
	"time"
)

type ContentDeal struct {
	ID                  int64       `gorm:"primaryKey"`
	Content             int64       `json:"content" gorm:"index:,option:CONCURRENTLY"`
	PropCid             string      `json:"propCid"`
	DealUUID            string      `json:"dealUuid"`
	Miner               string      `json:"miner"`
	DealID              int64       `json:"dealId"`
	Failed              bool        `json:"failed"`
	Verified            bool        `json:"verified"`
	Slashed             bool        `json:"slashed"`
	FailedAt            time.Time   `json:"failedAt,omitempty"`
	DTChan              string      `json:"dtChan" gorm:"index"`
	TransferStarted     time.Time   `json:"transferStarted"`
	TransferFinished    time.Time   `json:"transferFinished"`
	OnChainAt           time.Time   `json:"onChainAt"`
	SealedAt            time.Time   `json:"sealedAt"`
	LastMessage         string      `json:"lastMessage"`
	DealProtocolVersion protocol.ID `json:"deal_protocol_version"`
	MinerVersion        string      `json:"miner_version,omitempty"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

func (u *ContentDeal) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDeal Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDeal %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDeal) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentDeal Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentDeal %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentDeal) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentDeal Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentDeal %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}
