package db_models

import (
	"time"
)

// Wallet time series log events
type Wallet struct {
	ID         int64     `gorm:"primaryKey"`
	UuId       string    `json:"uuid"`
	Addr       string    `json:"addr"`
	Owner      string    `json:"owner"`
	KeyType    string    `json:"key_type"`
	PrivateKey string    `json:"private_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

//func (7u *Wallet) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var walletFromDb Wallet
//	tx.Model(&Wallet{}).Where("id = ?", u.ID).First(&walletFromDb)
//
//	if walletFromDb.ID == 0 {
//		return
//	}
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := WalletLog{
//		UuId:           u.UuId,
//		Addr:           u.Addr,
//		Owner:          u.Owner,
//		KeyType:        "REDACTED",
//		PrivateKey:     "REDACTED",
//		NodeInfo:       GetHostname(),
//		RequesterInfo:  ip,
//		SystemWalletId: u.ID,
//		DeltaNodeUuid:  instanceFromDb.InstanceUuid,
//		CreatedAt:      time.Now(),
//		UpdatedAt:      time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "WalletLog",
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
