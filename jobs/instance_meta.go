package jobs

import (
	"delta/core"
	model "delta/models"
	"runtime"
	"syscall"
)

// InstanceMetaProcessor It's a struct that contains a pointer to a DeltaNode.
// @property LightNode - This is the node that is being processed.
type InstanceMetaProcessor struct {
	LightNode *core.DeltaNode
}

// NewInstanceMetaProcessor > This function creates a new instance of the `InstanceMetaProcessor` struct and returns it as an `IProcessor` interface
func NewInstanceMetaProcessor(ln *core.DeltaNode, contentDeal model.ContentDeal) IProcessor {
	return &InstanceMetaProcessor{
		LightNode: ln,
	}
}

// Run It's checking if the CPU or Mem is above the meta set.
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
