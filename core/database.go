package core

import (
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
	db.AutoMigrate(&Content{}, &ContentDeal{}, &PieceCommitment{}, &MinerInfo{}, &MinerPrice{})
}

type ContentSplitRequest struct {
	ID        int64  `gorm:"primaryKey"`
	ContentId int64  `json:"content_id"`
	Cid       string `json:"cid"`
	Size      int64  `json:"size"`
	Name      string `json:"name"`
	Origins   string `json:"origins,omitempty"`
	ChunkSize int64  `json:"chunk_size"`
}

// main content record
type Content struct {
	ID                int64     `gorm:"primaryKey"`
	Name              string    `json:"name"`
	Size              int64     `json:"size"`
	Cid               string    `json:"cid"`
	RequestingApiKey  string    `json:"requesting_api_key,omitempty"`
	PieceCommitmentId int64     `json:"piece_commitment_id,omitempty"`
	Status            string    `json:"status"`
	Origins           string    `json:"origins,omitempty"`
	Created_at        time.Time `json:"created_at"`
	Updated_at        time.Time `json:"updated_at"`
}

type ContentStatus struct {
	ID            int64     `gorm:"primaryKey"`
	ContentId     int64     `json:"estuary_content_id"`
	CreatedAt     time.Time `json:"createdAtOnEstuary"`
	UpdatedAt     time.Time `json:"updatedAtOEstuary"`
	Cid           string    `json:"cid"`
	Name          string    `json:"name"`
	UserID        int       `json:"userId"`
	Description   string    `json:"description"`
	Size          int       `json:"size"`
	Type          int       `json:"type"`
	Active        bool      `json:"active"`
	Offloaded     bool      `json:"offloaded"`
	Replication   int       `json:"replication"`
	AggregatedIn  int       `json:"aggregatedIn"`
	Aggregate     bool      `json:"aggregate"`
	Pinning       bool      `json:"pinning"`
	PinMeta       string    `json:"pinMeta"`
	Replace       bool      `json:"replace"`
	Origins       string    `json:"origins"`
	Failed        bool      `json:"failed"`
	Location      string    `json:"location"`
	DagSplit      bool      `json:"dagSplit"`
	SplitFrom     int       `json:"splitFrom"`
	PinningStatus string    `json:"pinningStatus"`
	DealStatus    string    `json:"dealStatus"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type ContentDeal struct {
	gorm.Model
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
}

// buckets are aggregations of contents. It can either generate a car or just aggregate.
type Bucket struct {
	ID               int64     `gorm:"primaryKey"`
	Name             string    `json:"name"`
	UUID             string    `json:"uuid"`
	Status           string    `json:"status"`
	Cid              string    `json:"cid"`
	RequestingApiKey string    `json:"requesting_api_key,omitempty"`
	EstuaryContentId int64     `json:"estuary_content_id"`
	Created_at       time.Time `json:"created_at"`
	Updated_at       time.Time `json:"updated_at"`
}

type PieceCommitment struct {
	ID              int64  `gorm:"primaryKey"`
	Cid             string `json:"cid"`
	Piece           string `json:"piece"`
	Size            int64  `json:"size"`
	PaddedPieceSize uint64 `json:"padded_piece_size"`
	Status          string `json:"status"` // open, in-progress, completed (closed).
	Created_at      time.Time
	Updated_at      time.Time
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
