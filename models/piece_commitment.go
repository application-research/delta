package db_models

import (
	"time"
)

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

//func (u *PieceCommitment) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var pieceComm PieceCommitment
//	tx.Model(&PieceCommitment{}).Where("id = ?", u.ID).First(&pieceComm)
//
//	if pieceComm.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//	log := PieceCommitmentLog{
//		Cid:                            u.Cid,
//		Piece:                          u.Piece,
//		Size:                           u.Size,
//		PaddedPieceSize:                u.PaddedPieceSize,
//		UnPaddedPieceSize:              u.UnPaddedPieceSize,
//		Status:                         u.Status,
//		LastMessage:                    u.LastMessage,
//		NodeInfo:                       GetHostname(),
//		RequesterInfo:                  ip,
//		SystemContentPieceCommitmentId: u.ID,
//		DeltaNodeUuid:                  instanceFromDb.InstanceUuid,
//		CreatedAt:                      time.Now(),
//		UpdatedAt:                      time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "PieceCommitmentLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	producer.Publish(messageBytes)
//
//	return
//}
