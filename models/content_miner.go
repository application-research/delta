package db_models

import (
	"time"
)

type ContentMiner struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Miner     string    `json:"miner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//func (u *ContentMiner) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var contentMiner ContentMiner
//	tx.Model(&ContentMiner{}).Where("id = ?", u.ID).First(&contentMiner)
//
//	if contentMiner.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := ContentMinerLog{
//		Content:              u.Content,
//		Miner:                u.Miner,
//		NodeInfo:             GetHostname(),
//		RequesterInfo:        ip,
//		SystemContentMinerId: u.ID,
//		DeltaNodeUuid:        instanceFromDb.InstanceUuid,
//		CreatedAt:            time.Now(),
//		UpdatedAt:            time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "ContentMinerLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	producer.Publish(messageBytes)
//
//	return
//}
