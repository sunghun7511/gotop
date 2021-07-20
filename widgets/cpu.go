package widgets

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

type CpuCoreStats struct {
	userProcessTime uint64
	totalTime       uint64
}

type CpuStats struct {
	cores int
	stats []CpuCoreStats
}

type CpuWidget struct {
	cpuStats CpuStats
	data     [][]float64
	plot     *tWidgets.Plot
}

var HORIZONTAL_SCALE = 3

func NewCpuWidget() Widget {
	plot := tWidgets.NewPlot()
	plot.Title = " CPU Usage "
	plot.AxesColor = tui.ColorWhite
	plot.HorizontalScale = HORIZONTAL_SCALE
	plot.ShowAxes = false
	plot.MaxVal = 100

	cpuStats := getCpuStats()

	data := make([][]float64, cpuStats.cores)

	termWidth, _ := tui.TerminalDimensions()
	for i := 0; i < cpuStats.cores; i++ {
		data[i] = make([]float64, termWidth/HORIZONTAL_SCALE+1)
	}

	plot.Data = data

	return &CpuWidget{
		cpuStats: cpuStats,
		data:     data,
		plot:     plot,
	}
}

func (widget *CpuWidget) Update() {
	currentCpuStats := getCpuStats()
	previousCpuStats := widget.cpuStats

	cores := currentCpuStats.cores
	for core := 0; core < cores; core++ {
		previousCoreStats := previousCpuStats.stats[core]
		currentCoreStats := currentCpuStats.stats[core]

		coreUsage := calculateCoreUsage(previousCoreStats, currentCoreStats)
		widget.pushCoreUsageData(core, coreUsage)
	}

	widget.cpuStats = currentCpuStats
}

func (widget *CpuWidget) HandleSignal(event tui.Event) {
	// is there anything to handle?
}

func (widget *CpuWidget) GetUI() tui.Drawable {
	widget.plot.Data = widget.data
	return widget.plot
}

// read and parse cpu usage data from /proc/stat
func getCpuStats() CpuStats {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // ignore first line

	cores := runtime.NumCPU()
	stats := make([]CpuCoreStats, cores)
	for core := 0; core < cores; core++ {
		scanner.Scan()
		cpuCoreStats := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		stats[core] = parseCpuCoreStats(cpuCoreStats)
	}

	return CpuStats{
		cores: cores,
		stats: stats,
	}
}

// parse usage data of each cpu core
func parseCpuCoreStats(cpuCoreStats string) CpuCoreStats {
	cpuTimes := strings.Fields(cpuCoreStats)[1:]

	userProcessTime, _ := strconv.ParseUint(cpuTimes[0], 10, 64)

	totalTime := uint64(0)
	for _, data := range cpuTimes {
		time, _ := strconv.ParseUint(data, 10, 64)
		totalTime += time
	}

	return CpuCoreStats{
		userProcessTime: userProcessTime,
		totalTime:       totalTime,
	}
}

func calculateCoreUsage(previousCoreStats, currentCoreStats CpuCoreStats) float64 {
	deltaUserProcessTime := currentCoreStats.userProcessTime - previousCoreStats.userProcessTime
	deltaTotalTime := currentCoreStats.totalTime - previousCoreStats.totalTime

	return (float64(deltaUserProcessTime) / float64(deltaTotalTime)) * 100.0
}

func (widget *CpuWidget) pushCoreUsageData(core int, coreUsage float64) {
	widget.data[core] = append(widget.data[core], coreUsage)
	widget.data[core] = widget.data[core][1:]
}
