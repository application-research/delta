package db_models

import (
	"time"
)

type InstanceMeta struct {
	// gorm id
	ID                               int64     `gorm:"primary_key" json:"id"`
	InstanceUuid                     string    `json:"instance_uuid"`
	InstanceHostName                 string    `json:"instance_host_name"`
	InstanceNodeName                 string    `json:"instance_node_name"`
	OSDetails                        string    `json:"os_details"`
	PublicIp                         string    `json:"public_ip"`
	MemoryLimit                      uint64    `json:"memory_limit"`
	CpuLimit                         uint64    `json:"cpu_limit"`
	StorageLimit                     uint64    `json:"storage_limit"`
	DisableRequest                   bool      `json:"disable_requests"`
	DisableCommitmentPieceGeneration bool      `json:"disable_commitment_piece_generation"`
	DisableStorageDeal               bool      `json:"disable_storage_deal"`
	DisableOnlineDeals               bool      `json:"disable_online_deals"`
	DisableOfflineDeals              bool      `json:"disable_offline_deals"`
	NumberOfCpus                     uint64    `json:"number_of_cpus"`
	StorageInBytes                   uint64    `json:"storage_in_bytes"`
	SystemMemory                     uint64    `json:"system_memory"`
	HeapMemory                       uint64    `json:"heap_memory"`
	HeapInUse                        uint64    `json:"heap_in_use"`
	StackInUse                       uint64    `json:"stack_in_use"`
	InstanceStart                    time.Time `json:"instance_start"`
	BytesPerCpu                      uint64    `json:"bytes_per_cpu"`
	CreatedAt                        time.Time `json:"created_at"`
	UpdatedAt                        time.Time `json:"updated_at"`
}

//func (u *InstanceMeta) AfterSave(tx *gorm.DB) (err error) {
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
//	log := InstanceMetaLog{
//		InstanceUuid:                     u.InstanceUuid,
//		InstanceHostName:                 u.InstanceHostName,
//		InstanceNodeName:                 u.InstanceNodeName,
//		OSDetails:                        u.OSDetails,
//		PublicIp:                         u.PublicIp,
//		MemoryLimit:                      u.MemoryLimit,
//		CpuLimit:                         u.CpuLimit,
//		StorageLimit:                     u.StorageLimit,
//		DisableRequest:                   u.DisableRequest,
//		DisableCommitmentPieceGeneration: u.DisableCommitmentPieceGeneration,
//		DisableStorageDeal:               u.DisableStorageDeal,
//		DisableOnlineDeals:               u.DisableOnlineDeals,
//		DisableOfflineDeals:              u.DisableOfflineDeals,
//		NumberOfCpus:                     u.NumberOfCpus,
//		StorageInBytes:                   u.StorageInBytes,
//		SystemMemory:                     u.SystemMemory,
//		HeapMemory:                       u.HeapMemory,
//		HeapInUse:                        u.HeapInUse,
//		StackInUse:                       u.StackInUse,
//		InstanceStart:                    u.InstanceStart,
//		BytesPerCpu:                      u.BytesPerCpu,
//		NodeInfo:                         GetHostname(),
//		RequesterInfo:                    ip,
//		DeltaNodeUuid:                    u.InstanceUuid,
//		SystemInstanceMetaId:             u.ID,
//		CreatedAt:                        time.Now(),
//		UpdatedAt:                        time.Now(),
//	}
//
//	deltaMetricsBaseMessage := DeltaMetricsBaseMessage{
//		ObjectType: "InstanceMetaLog",
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
