package core

import (
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"gorm.io/gorm"
	"runtime"
	"syscall"
	"time"
)

// ScanHostComputeResources Setting the global node meta.
// > This function sets the global node metadata for the given node
func ScanHostComputeResources(ln *DeltaNode, repo string) *model.InstanceMeta {

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	totalMemory := memStats.Sys
	totalMemory80 := totalMemory * 90 / 100

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Total system memory: %v bytes\n", m.Sys)
	fmt.Printf("Total heap memory: %v bytes\n", m.HeapSys)
	fmt.Printf("Heap in use: %v bytes\n", m.HeapInuse)
	fmt.Printf("Stack in use: %v bytes\n", m.StackInuse)

	// get the 80% of the total disk usage
	var stat syscall.Statfs_t
	syscall.Statfs(repo, &stat) // blockstore size
	totalStorage := stat.Blocks * uint64(stat.Bsize)
	totalStorage90 := totalStorage * 90 / 100

	// set the number of CPUs
	numCPU := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n", numCPU)
	runtime.GOMAXPROCS(numCPU / (1200 / 1000))

	// delete all data from the instance meta table
	ln.DB.Model(&model.InstanceMeta{}).Delete(&model.InstanceMeta{}, "id > ?", 0)
	// re-create
	instanceMeta := &model.InstanceMeta{
		MemoryLimit:                      totalMemory80,  // 80%
		StorageLimit:                     totalStorage90, // 90%
		NumberOfCpus:                     uint64(numCPU),
		StorageInBytes:                   totalStorage,
		BytesPerCpu:                      11000000000,
		SystemMemory:                     totalMemory,
		HeapMemory:                       m.HeapSys,
		HeapInUse:                        m.HeapInuse,
		StackInUse:                       m.StackInuse,
		DisableRequest:                   false,
		DisableCommitmentPieceGeneration: false,
		DisableStorageDeal:               false,
		DisableOnlineDeals:               false,
		DisableOfflineDeals:              false,
		CreatedAt:                        time.Now(),
		UpdatedAt:                        time.Now(),
		InstanceStart:                    time.Now(),
	}
	ln.DB.Model(&model.InstanceMeta{}).Create(instanceMeta)
	ln.MetaInfo = instanceMeta
	return instanceMeta

}

// CleanUpContentAndPieceComm It updates the status of all the content and piece commitments that were in the process of being transferred or computed
// to failed
func CleanUpContentAndPieceComm(ln *DeltaNode) {

	// if the transfer was started upon restart, then we need to update the status to failed
	ln.DB.Transaction(func(tx *gorm.DB) error {

		tx.Model(&model.Content{}).Where("status in (?,?)", utils.DEAL_STATUS_TRANSFER_STARTED, utils.CONTENT_PIECE_COMPUTING).Updates(
			model.Content{
				Status:      utils.DEAL_STATUS_TRANSFER_FAILED,
				UpdatedAt:   time.Now(),
				LastMessage: "Transfer failed due to node restart",
			})
		return nil
	})
}
