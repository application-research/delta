package jobs

import (
	"delta/core"
	"delta/core/model"
	"runtime"
	"syscall"
)

type InstanceMetaProcessor struct {
	LightNode *core.DeltaNode
}

func NewInstanceMetaProcessor(ln *core.DeltaNode, contentDeal model.ContentDeal) IProcessor {
	return &InstanceMetaProcessor{
		LightNode: ln,
	}
}

func (d InstanceMetaProcessor) Run() error {
	// check if CPU or Mem is above the meta set
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	totalMemory := memStats.Sys

	// Get the number of available CPUs
	//numCPUs := runtime.NumCPU()

	// Get the total amount of storage space in bytes
	var stat syscall.Statfs_t
	syscall.Statfs("/", &stat)
	totalStorage := stat.Blocks * uint64(stat.Bsize)

	// if it's above, set the api to read only
	if d.LightNode.MetaInfo.StorageLimit < totalStorage || d.LightNode.MetaInfo.MemoryLimit < totalMemory {
		d.LightNode.MetaInfo.DisableRequest = true
		d.LightNode.MetaInfo.DisableCommitmentPieceGeneration = true
		d.LightNode.MetaInfo.DisableStorageDeal = true
		d.LightNode.MetaInfo.DisableOnlineDeals = true
		d.LightNode.MetaInfo.DisableOfflineDeals = true
	}

	return nil
}
