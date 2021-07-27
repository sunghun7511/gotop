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
	cores      int
	stats      []CpuCoreStats
	totalStats CpuCoreStats
}

type CpuWidget struct {
	cpuStats     CpuStats
	totalData    [][]float64
	data         [][]float64
	showEachCore bool
	plot         *tWidgets.Plot
}

var HORIZONTAL_SCALE = 3

func NewCpuWidget() Widget {
	plot := tWidgets.NewPlot()
	plot.Title = " CPU Usage "
	plot.AxesColor = tui.ColorWhite
	plot.HorizontalScale = HORIZONTAL_SCALE
	plot.ShowAxes = false
	plot.MaxVal = 100

	cpuStats, err := getCpuStats()
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]float64, cpuStats.cores)

	termWidth, _ := tui.TerminalDimensions()
	for i := 0; i < cpuStats.cores; i++ {
		data[i] = make([]float64, termWidth/HORIZONTAL_SCALE+1)
	}

	totalData := make([][]float64, 1)
	totalData[0] = make([]float64, termWidth/HORIZONTAL_SCALE+1)

	plot.Data = totalData

	return &CpuWidget{
		cpuStats:     cpuStats,
		data:         data,
		totalData:    totalData,
		showEachCore: false,
		plot:         plot,
	}
}

func (widget *CpuWidget) Update() {
	currentCpuStats, err := getCpuStats()
	if err != nil {
		log.Print(err)
		return
	}

	previousCpuStats := widget.cpuStats

	totalCpuUsage := calculateCoreUsage(previousCpuStats.totalStats, currentCpuStats.totalStats)
	widget.totalData[0] = append(widget.totalData[0], totalCpuUsage)
	widget.totalData[0] = widget.totalData[0][1:]

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
	switch event.ID {
	case "1":
		widget.showEachCore = !widget.showEachCore
	}
}

func (widget *CpuWidget) GetUI() tui.Drawable {
	if widget.showEachCore {
		widget.plot.Data = widget.data
	} else {
		widget.plot.Data = widget.totalData
	}
	return widget.plot
}

// read and parse cpu usage data from /proc/stat
func getCpuStats() (CpuStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return CpuStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // total cpu usage
	cpuTotalStats := scanner.Text()
	if err := scanner.Err(); err != nil {
		return CpuStats{}, err
	}

	totalStats, err := parseCpuStats(cpuTotalStats)
	if err != nil {
		return CpuStats{}, err
	}

	cores := runtime.NumCPU()
	stats := make([]CpuCoreStats, cores)
	for core := 0; core < cores; core++ {
		scanner.Scan()
		cpuCoreStats := scanner.Text()
		if err := scanner.Err(); err != nil {
			return CpuStats{}, err
		}

		stats[core], err = parseCpuStats(cpuCoreStats)
		if err != nil {
			return CpuStats{}, err
		}
	}

	return CpuStats{
		cores:      cores,
		stats:      stats,
		totalStats: totalStats,
	}, nil
}

// parse usage data from /proc/stat
func parseCpuStats(cpuCoreStats string) (CpuCoreStats, error) {
	cpuTimes := strings.Fields(cpuCoreStats)[1:]

	userProcessTime, err := strconv.ParseUint(cpuTimes[0], 10, 64)
	if err != nil {
		return CpuCoreStats{}, err
	}

	totalTime := uint64(0)
	for _, data := range cpuTimes {
		time, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return CpuCoreStats{}, err
		}
		totalTime += time
	}

	return CpuCoreStats{
		userProcessTime: userProcessTime,
		totalTime:       totalTime,
	}, nil
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
