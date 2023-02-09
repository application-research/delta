package core

import (
	"fmt"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func OpenDatabase() (*gorm.DB, error) {

	dbName, okHost := viper.Get("DB_NAME").(string)
	if !okHost {
		panic("DB_NAME not set")
	}
	DB, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	// generate new models.
	ConfigureModels(DB) // create models.

	if err != nil {
		return nil, err
	}
	return DB, nil
}

func ConfigureModels(db *gorm.DB) {
	db.AutoMigrate(&Content{}, &ContentDeal{}, &PieceCommitment{}, &MinerInfo{}, &MinerPrice{}, &LogEvent{})
}

type Content struct {
	ID                int64     `gorm:"primaryKey"`
	Name              string    `json:"name"`
	Size              int64     `json:"size"`
	Cid               string    `json:"cid"`
	RequestingApiKey  string    `json:"requesting_api_key,omitempty"`
	PieceCommitmentId int64     `json:"piece_commitment_id,omitempty"`
	Status            string    `json:"status"`
	Origins           string    `json:"origins,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

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
	DealProtocolVersion protocol.ID `json:"deal_protocol_version"`
	MinerVersion        string      `json:"miner_version"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

type PieceCommitment struct {
	ID              int64  `gorm:"primaryKey"`
	Cid             string `json:"cid"`
	Piece           string `json:"piece"`
	Size            int64  `json:"size"`
	PaddedPieceSize int64  `json:"padded_piece_size"`
	Status          string `json:"status"` // open, in-progress, completed (closed).
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MinerInfo struct {
	ID              int64  `gorm:"primaryKey"`
	Addr            string `json:"addr"` // same as Miner from MinerPrice
	Name            string `json:"name"`
	Suspended       bool   `json:"suspended"`
	Version         string `json:"version"`
	ChainInfo       string `json:"chain_info"`
	SuspendedReason string `json:"suspendedReason,omitempty"`
}

type MinerPrice struct {
	ID            int64  `gorm:"primaryKey"`
	Miner         string `json:"miner"`
	Price         string `json:"price"`
	VerifiedPrice string `json:"verifiedPrice"`
	MinPieceSize  int64  `json:"minPieceSize"`
	MaxPieceSize  int64  `json:"maxPieceSize"`
	MinerVersion  string `json:"miner_version"`
}

type Wallet struct {
	ID        int64     `gorm:"primaryKey"`
	Addr      string    `json:"addr"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogEvent struct {
	ID         int64     `gorm:"primaryKey"`   // auto increment
	EventType  string    `json:"log_event"`    // content, deal, piece_commitment, upload, miner, info
	LogEventId int64     `json:"log_event_id"` // object id
	LogEvent   string    `json:"log_event"`    // description
	CreatedAt  time.Time `json:"created_at"`   // auto set
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
