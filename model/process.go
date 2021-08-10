package model

// Process process info
type Process struct {
	Pid           string
	User          string
	Cmd           string
	TotalCPUUsage uint64
	CPUUsage      float64
	MemUsage      float64
}
