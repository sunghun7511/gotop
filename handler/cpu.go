package handler

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sunghun7511/gotop/model"
)

// GetCPUStats read and parse cpu usage data from /proc/stat
func GetCPUStats() (model.CpuStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return model.CpuStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// parse total cpu usage
	totalStats, err := parseCpuStats(scanner)
	if err != nil {
		return model.CpuStats{}, err
	}

	// parse usage data of each core
	cores := runtime.NumCPU()
	stats := make([]model.CpuCoreStats, cores)
	for core := 0; core < cores; core++ {
		stats[core], err = parseCpuStats(scanner)
		if err != nil {
			return model.CpuStats{}, err
		}
	}

	return model.CpuStats{
		Cores:      cores,
		Stats:      stats,
		TotalStats: totalStats,
	}, nil
}

// read each line from /proc/stat and parse it to CpuCoreStats
func parseCpuStats(scanner *bufio.Scanner) (model.CpuCoreStats, error) {
	scanner.Scan()
	cpuCoreStats := scanner.Text()
	if err := scanner.Err(); err != nil {
		return model.CpuCoreStats{}, err
	}

	cpuTimes := strings.Fields(cpuCoreStats)[1:]

	userProcessTime, err := strconv.ParseUint(cpuTimes[0], 10, 64)
	if err != nil {
		return model.CpuCoreStats{}, err
	}

	totalTime := uint64(0)
	for _, data := range cpuTimes {
		time, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return model.CpuCoreStats{}, err
		}
		totalTime += time
	}

	return model.CpuCoreStats{
		UserProcessTime: userProcessTime,
		TotalTime:       totalTime,
	}, nil
}
