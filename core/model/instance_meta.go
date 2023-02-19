package model

import (
	"gorm.io/gorm"
)

type InstanceMeta struct {
	// gorm id
	ID                               uint64 `gorm:"primary_key" json:"id"`
	MemoryLimit                      uint64 `json:"memory_limit"`
	CpuLimit                         uint64 `json:"cpu_limit"`
	StorageLimit                     uint64 `json:"storage_limit"`
	DisableRequest                   bool   `json:"disable_requests"`
	DisableCommitmentPieceGeneration bool   `json:"disable_commitment_piece_generation"`
	DisableStorageDeal               bool   `json:"disable_storage_deal"`
	DisableOnlineDeals               bool   `json:"disable_online_deals"`
	DisableOfflineDeals              bool   `json:"disable_offline_deals"`
}

func (u *InstanceMeta) BeforeSave(tx *gorm.DB) (err error) {
	return
}

func (u *InstanceMeta) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (u *InstanceMeta) AfterSave(tx *gorm.DB) (err error) {
	return
}
