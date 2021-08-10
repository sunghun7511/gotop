package widgets

import (
	"log"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"

	"github.com/sunghun7511/gotop/core"
	"github.com/sunghun7511/gotop/model"
	"github.com/sunghun7511/gotop/util"
)

type CpuWidget struct {
	cpuStats     model.CpuStats
	totalData    [][]float64
	data         [][]float64
	showEachCore bool
	plot         *tWidgets.Plot
}

var HorizontalScale = 3

func calculateDataLength() int {
	termWidth, _ := tui.TerminalDimensions()
	return termWidth/HorizontalScale + 1
}

func NewCpuWidget() Widget {
	plot := tWidgets.NewPlot()
	plot.Title = " CPU Usage "
	plot.AxesColor = tui.ColorWhite
	plot.HorizontalScale = HorizontalScale
	plot.ShowAxes = false
	plot.MaxVal = 100

	cpuStats, err := core.GetCPUStats()
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]float64, cpuStats.Cores)

	dataLength := calculateDataLength()
	for i := 0; i < cpuStats.Cores; i++ {
		data[i] = make([]float64, dataLength)
	}

	totalData := make([][]float64, 1)
	totalData[0] = make([]float64, dataLength)

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
	currentCpuStats, err := core.GetCPUStats()
	if err != nil {
		log.Print(err)
		return
	}

	previousCpuStats := widget.cpuStats

	totalCpuUsage := calculateCoreUsage(previousCpuStats.TotalStats, currentCpuStats.TotalStats)
	widget.totalData[0] = util.PushUsageData(widget.totalData[0], totalCpuUsage)

	cores := currentCpuStats.Cores
	for core := 0; core < cores; core++ {
		previousCoreStats := previousCpuStats.Stats[core]
		currentCoreStats := currentCpuStats.Stats[core]

		coreUsage := calculateCoreUsage(previousCoreStats, currentCoreStats)
		widget.data[core] = util.PushUsageData(widget.data[core], coreUsage)
	}

	widget.cpuStats = currentCpuStats

	widget.updateDataLength()
}

func (widget *CpuWidget) updateDataLength() {
	newLength := calculateDataLength()
	currentLength := len(widget.totalData[0])
	cores := widget.cpuStats.Cores

	if newLength > currentLength { // terminal expanded
		lengthDifference := newLength - currentLength

		widget.totalData[0] = append(make([]float64, lengthDifference), widget.totalData[0]...)
		for core := 0; core < cores; core++ {
			widget.data[core] = append(make([]float64, lengthDifference), widget.data[core]...)
		}
	} else if newLength < currentLength { // terminal shrinked
		lengthDifference := currentLength - newLength

		widget.totalData[0] = widget.totalData[0][lengthDifference:]
		for core := 0; core < cores; core++ {
			widget.data[core] = widget.data[core][lengthDifference:]
		}
	}
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

func calculateCoreUsage(previousCoreStats, currentCoreStats model.CpuCoreStats) float64 {
	deltaUserProcessTime := currentCoreStats.UserProcessTime - previousCoreStats.UserProcessTime
	deltaTotalTime := currentCoreStats.TotalTime - previousCoreStats.TotalTime

	return (float64(deltaUserProcessTime) / float64(deltaTotalTime)) * 100.0
}
