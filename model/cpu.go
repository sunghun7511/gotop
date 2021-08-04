package model

type CpuCoreStats struct {
	UserProcessTime uint64
	TotalTime       uint64
}

type CpuStats struct {
	Cores      int
	Stats      []CpuCoreStats
	TotalStats CpuCoreStats
}
