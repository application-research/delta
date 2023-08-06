package db_models

import (
	"fmt"
	"github.com/application-research/delta-db/messaging"
	"gorm.io/gorm/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var producer *messaging.DeltaMetricsMessageProducer

type DeltaMetricsBaseMessage struct {
	ObjectType string      `json:"object_type"`
	Object     interface{} `json:"object"`
}

func init() {
	producer = messaging.NewDeltaMetricsMessageProducer()
}

func OpenDatabase(dbDsn string) (*gorm.DB, error) {
	// use postgres
	var DB *gorm.DB
	var err error

	if dbDsn[:8] == "postgres" {
		DB, err = gorm.Open(postgres.Open(dbDsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		DB, err = gorm.Open(sqlite.Open(dbDsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	// generate new models.
	ConfigureModels(DB) // create models.

	if err != nil {
		return nil, err
	}
	return DB, nil
}

func ConfigureModels(db *gorm.DB) {
	db.AutoMigrate(&Content{}, &ContentDeal{}, &PieceCommitment{}, &MinerInfo{}, &MinerPrice{}, &messaging.LogEvent{}, &ContentMiner{}, &ProcessContentCounter{}, &ContentWallet{}, &ContentDealProposalParameters{}, &Wallet{}, &ContentDealProposal{}, &InstanceMeta{}, &RetryDealCount{}, &BatchImport{}, &BatchImportContent{})
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

type AdminUser struct {
	ID        int64     `gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
