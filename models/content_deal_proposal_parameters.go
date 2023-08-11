package db_models

import (
	"time"
)

type ContentDealProposalParameters struct {
	ID                     int64     `gorm:"primaryKey"`
	Content                int64     `json:"content" gorm:"index:,option:CONCURRENTLY"`
	Label                  string    `json:"label,omitempty"`
	Duration               int64     `json:"duration,omitempty"`
	StartEpoch             int64     `json:"start_epoch,omitempty"`
	EndEpoch               int64     `json:"end_epoch,omitempty"`
	TransferParams         string    `json:"transfer_parameters,omitempty"`
	RemoveUnsealedCopy     bool      `json:"remove_unsealed_copy"`
	SkipIPNIAnnounce       bool      `json:"skip_ipni_announce"`
	VerifiedDeal           bool      `json:"verified_deal"`
	UnverifiedDealMaxPrice string    `json:"unverified_deal_max_price"`
	CreatedAt              time.Time `json:"created_at" json:"created-at"`
	UpdatedAt              time.Time `json:"updated_at" json:"updated-at"`
}

//func (u *ContentDealProposalParameters) AfterSave(tx *gorm.DB) (err error) {
//
//	var instanceFromDb InstanceMeta
//	tx.Raw("SELECT * FROM instance_meta ORDER BY id DESC LIMIT 1").Scan(&instanceFromDb)
//
//	if instanceFromDb.ID == 0 {
//		return
//	}
//
//	var contentDealProposalParams ContentDealProposalParameters
//	tx.Model(&ContentDealProposalParameters{}).Where("id = ?", u.ID).First(&contentDealProposalParams)
//
//	if contentDealProposalParams.ID == 0 {
//		return
//	}
//
//	// get instance info
//	ip, err := GetPublicIP()
//	if err != nil {
//		return
//	}
//	log := ContentDealProposalParametersLog{
//		Content:                               u.Content,
//		Label:                                 u.Label,
//		Duration:                              u.Duration,
//		StartEpoch:                            u.StartEpoch,
//		EndEpoch:                              u.EndEpoch,
//		TransferParams:                        u.TransferParams,
//		RemoveUnsealedCopy:                    u.RemoveUnsealedCopy,
//		SkipIPNIAnnounce:                      u.SkipIPNIAnnounce,
//		VerifiedDeal:                          u.VerifiedDeal,
//		UnverifiedDealMaxPrice:                u.UnverifiedDealMaxPrice,
//		NodeInfo:                              GetHostname(),
//		RequesterInfo:                         ip,
//		DeltaNodeUuid:                         instanceFromDb.InstanceUuid,
//		SystemContentDealProposalParametersId: u.ID,
//		CreatedAt:                             time.Now(),
//		UpdatedAt:                             time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "ContentDealProposalParametersLog",
//		Object:     log,
//	}
//
//	messageBytes, err := json.Marshal(deltaMetricsBaseMessage)
//	producer.Publish(messageBytes)
//
//	return
//}
