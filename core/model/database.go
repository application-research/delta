package model

import (
	"fmt"
	"time"

	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDatabase(dbDsn string) (*gorm.DB, error) {
	// use postgres
	var DB *gorm.DB
	var err error

	if dbDsn[:8] == "postgres" {
		DB, err = gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open(dbDsn), &gorm.Config{})
	}

	// generate new models.
	ConfigureModels(DB) // create models.

	if err != nil {
		return nil, err
	}
	return DB, nil
}

func ConfigureModels(db *gorm.DB) {
	db.AutoMigrate(&Content{}, &ContentDeal{}, &PieceCommitment{}, &MinerInfo{}, &MinerPrice{}, &LogEvent{}, &ContentMiner{}, &ProcessContentCounter{}, &ContentWallet{}, &ContentDealProposalParameters{}, &Wallet{}, &ContentDealProposal{}, &InstanceMeta{}, &RetryDealCount{})
}

type ProcessContentCounter struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Counter   int64     `json:"counter"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MinerInfo struct {
	ID              int64     `gorm:"primaryKey"`
	Addr            string    `json:"addr"` // same as Miner from MinerPrice
	Name            string    `json:"name"`
	Suspended       bool      `json:"suspended"`
	Version         string    `json:"version"`
	ChainInfo       string    `json:"chain_info"`
	SuspendedReason string    `json:"suspendedReason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MinerPrice struct {
	ID            int64     `gorm:"primaryKey"`
	Miner         string    `json:"miner"`
	Price         string    `json:"price"`
	VerifiedPrice string    `json:"verifiedPrice"`
	MinPieceSize  int64     `json:"minPieceSize"`
	MaxPieceSize  int64     `json:"maxPieceSize"`
	MinerVersion  string    `json:"miner_version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Wallet struct {
	ID         int64     `gorm:"primaryKey"`
	Addr       string    `json:"addr"`
	Owner      string    `json:"owner"`
	KeyType    string    `json:"key_type"`
	PrivateKey string    `json:"private_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AdminUser struct {
	ID        int64     `gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogEvent struct {
	ID         int64     `gorm:"primaryKey"`   // auto increment
	EventType  string    `json:"log_event"`    // content, deal, piece_commitment, upload, miner, info
	LogEventId int64     `json:"log_event_id"` // object id
	LogEvent   string    `json:"log_event"`    // description
	CreatedAt  time.Time `json:"created_at"`   // auto set
	UpdatedAt  time.Time `json:"updated_at"`
}

type RetryDealCount struct {
	ID        int64     `gorm:"primaryKey"`
	Type      string    `json:"type"`
	OldId     int64     `json:"old_id"`
	NewId     int64     `json:"new_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNoChannelID = fmt.Errorf("no data transfer channel id in deal")

func (cd ContentDeal) ChannelID() (datatransfer.ChannelID, error) {
	if cd.DTChan == "" {
		return datatransfer.ChannelID{}, ErrNoChannelID
	}

	chid, err := filclient.ChannelIDFromString(cd.DTChan)
	if err != nil {
		err = fmt.Errorf("incorrectly formatted data transfer channel ID in contentDeal record: %w", err)
		return datatransfer.ChannelID{}, err
	}
	return *chid, nil
}
