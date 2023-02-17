package core

import (
	"fmt"
	"time"

	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/libp2p/go-libp2p/core/protocol"
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
	db.AutoMigrate(&Content{}, &ContentDeal{}, &PieceCommitment{}, &MinerInfo{}, &MinerPrice{}, &LogEvent{}, &ContentMinerAssignment{}, &ProcessContentCounter{}, &ContentWalletAssignment{}, &ContentDealProposalParameters{}, &Wallet{})
}

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

type ProcessContentCounter struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Counter   int64     `json:"counter"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ContentMinerAssignment struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Miner     string    `json:"miner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *ContentMinerAssignment) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMinerAssignment Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentMinerAssignment %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentMinerAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMinerAssignment Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentMinerAssignment %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentMinerAssignment) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentMinerAssignment Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentMinerAssignment %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

type ContentWalletAssignment struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Wallet    string    `json:"wallet_meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *ContentWalletAssignment) BeforeSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentWalletAssignment Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentWalletAssignment %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentWalletAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "ContentMinerAssignment Create",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("ContentWalletAssignment %d create", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

func (u *ContentWalletAssignment) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&LogEvent{}).Save(&LogEvent{
		EventType:  "After ContentWalletAssignment Save",
		LogEventId: u.ID,
		LogEvent:   fmt.Sprintf("After ContentWalletAssignment %d saved", u.ID),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	})
	return
}

type ContentDealProposalParameters struct {
	ID             int64     `gorm:"primaryKey"`
	Content        int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Label          string    `json:"label,omitempty"`
	Duration       int64     `json:"duration,omitempty"`
	StartEpoch     int64     `json:"start-epoch,omitempty"`
	EndEpoch       int64     `json:"end-epoch,omitempty"`
	TransferParams string    `json:"transfer-params,omitempty"`
	CreatedAt      time.Time `json:"created_at" json:"created-at"`
	UpdatedAt      time.Time `json:"updated_at" json:"updated-at"`
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

type RepairRequest struct {
	ID        int64     `gorm:"primaryKey"`
	ObjectId  int64     `json:"object_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//type LogEventCommp
//
//log_event_commp
//log_event_deal
//log_event_miner
//log_event_api_key
//log_event_jobs

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
	DealUUID  string    `json:"deal_uuid"`
	Count     int64     `json:"count"`
	LastError string    `json:"last_error"`
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
