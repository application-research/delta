package db_models

import (
	"time"
)

type ContentWallet struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	WalletId  int64     `json:"wallet_id" gorm:"index:,option:CONCURRENTLY"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//func (u *ContentWallet) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var contentWallet ContentWallet
//	tx.Model(&ContentWallet{}).Where("id = ?", u.ID).First(&contentWallet)
//
//	if contentWallet.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//	log := ContentWalletLog{
//		Content:               u.Content,
//		WalletId:              u.WalletId,
//		NodeInfo:              GetHostname(),
//		RequesterInfo:         ip,
//		SystemContentWalletId: u.ID,
//		DeltaNodeUuid:         instanceFromDb.InstanceUuid,
//		CreatedAt:             time.Now(),
//		UpdatedAt:             time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "ContentWalletLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	producer.Publish(messageBytes)
//
//	return
//}
