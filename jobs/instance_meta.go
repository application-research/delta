package jobs

//// Get the total amount of system memory in bytes
//memStats := &runtime.MemStats{}
//runtime.ReadMemStats(memStats)
//totalMemory := memStats.Sys
//
//// Get the number of available CPUs
//numCPUs := runtime.NumCPU()
//
//// Get the total amount of storage space in bytes
//var stat syscall.Statfs_t
//syscall.Statfs("/", &stat)
//totalStorage := stat.Blocks * uint64(stat.Bsize)
//
//fmt.Printf("Total Memory: %v bytes\n", totalMemory)
//fmt.Printf("Number of CPUs: %v\n", numCPUs)
//fmt.Printf("Total Storage: %v bytes\n", totalStorage)

// check if CPU or Mem is above the meta set
// if CPU or Mem is above, set the disable upload to false.
//
