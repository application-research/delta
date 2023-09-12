package db_models

import (
	"time"
)

type ContentDealProposal struct {
	ID        int64     `gorm:"primaryKey"`
	Content   int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Unsigned  string    `json:"unsigned"`
	Signed    string    `json:"signed"`
	Meta      string    `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//func (u *ContentDealProposal) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var contentDealProposal ContentDealProposal
//	tx.Model(&ContentDealProposal{}).Where("id = ?", u.ID).First(&contentDealProposal)
//
//	if contentDealProposal.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//
//	log := ContentDealProposalLog{
//		Content:                     u.Content,
//		Unsigned:                    u.Unsigned,
//		Signed:                      u.Signed,
//		Meta:                        u.Meta,
//		NodeInfo:                    GetHostname(),
//		RequesterInfo:               ip,
//		DeltaNodeUuid:               instanceFromDb.InstanceUuid,
//		SystemContentDealProposalId: u.ID,
//		CreatedAt:                   time.Now(),
//		UpdatedAt:                   time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "ContentDealProposalLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	producer.Publish(messageBytes)
//
//	return
//}
