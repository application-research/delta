package db_models

import (
	"time"
)

// BatchImport create an entry first
type BatchImport struct {
	ID        int64     `gorm:"primaryKey"`
	Uuid      string    `json:"uuid" gorm:"index:,option:CONCURRENTLY"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//func (u *BatchImport) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	// get the latest instance uuid
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	tx.Model(&InstanceMeta{}).Where("id > 0").First(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := BatchImportLog{
//		Uuid:                u.Uuid,
//		Status:              u.Status,
//		NodeInfo:            GetHostname(),
//		RequesterInfo:       ip,
//		DeltaNodeUuid:       instanceFromDb.InstanceUuid,
//		SystemBatchImportId: u.ID,
//		CreatedAt:           time.Now(),
//		UpdatedAt:           time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "BatchImportLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	if err != nil {
//		return err
//	}
//	producer.Publish(messageBytes)
//	return
//}

// BatchContent associate the content to a batch
type BatchImportContent struct {
	ID            int64     `gorm:"primaryKey"`
	BatchImportID int64     `json:"batch_import_id" gorm:"index:,option:CONCURRENTLY"`
	ContentID     int64     `json:"content_id" gorm:"index:,option:CONCURRENTLY"` // check status of the content
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

//func (u *BatchImportContent) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	// get the latest instance uuid
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	tx.Model(&InstanceMeta{}).Where("id > 0").First(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := BatchImportContentLog{
//		BatchImportID:        u.BatchImportID,
//		ContentID:            u.ContentID,
//		NodeInfo:             GetHostname(),
//		RequesterInfo:        ip,
//		DeltaNodeUuid:        instanceFromDb.InstanceUuid,
//		SystemBatchContentId: u.ID,
//		CreatedAt:            time.Now(),
//		UpdatedAt:            time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "BatchImportContentLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	if err != nil {
//		return err
//	}
//	producer.Publish(messageBytes)
//	return
//}
