package db_models

import (
	"time"
)

type Content struct {
	ID                int64     `gorm:"primaryKey"`
	Name              string    `json:"name"`
	Size              int64     `json:"size"`
	Cid               string    `json:"cid"`
	RequestingApiKey  string    `json:"requesting_api_key,omitempty"`
	PieceCommitmentId int64     `json:"piece_commitment_id,omitempty"`
	Status            string    `json:"status"`
	RequestType       string    `json:"request_type"`    // default signed, or unsigned
	ConnectionMode    string    `json:"connection_mode"` // offline or online
	AutoRetry         bool      `json:"auto_retry"`
	LastMessage       string    `json:"last_message"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

//func (u *Content) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	// get the latest instance info based on created_at
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var contentFromDb Content
//	tx.Model(&Content{}).Where("id = ?", u.ID).First(&contentFromDb)
//
//	if contentFromDb.ID == 0 {
//		return
//	}
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := ContentLog{
//		Name:              u.Name,
//		Size:              u.Size,
//		Cid:               u.Cid,
//		RequestingApiKey:  u.RequestingApiKey,
//		PieceCommitmentId: u.PieceCommitmentId,
//		Status:            u.Status,
//		ConnectionMode:    u.ConnectionMode,
//		LastMessage:       u.LastMessage,
//		AutoRetry:         u.AutoRetry,
//		NodeInfo:          GetHostname(),
//		RequesterInfo:     ip,
//		DeltaNodeUuid:     instanceFromDb.InstanceUuid,
//		SystemContentId:   u.ID,
//		CreatedAt:         time.Now(),
//		UpdatedAt:         time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "ContentLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	if err != nil {
//		return err
//	}
//	producer.Publish(messageBytes)
//
//	return
//}
