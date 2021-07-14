package widgets

import (
	"runtime"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

var POINTS = 30

type CpuWidget struct {
	cores            int
	prevUserProcTime []uint64
	prevTotalTime    []uint64
	data             [][]float64
	plot             *tWidgets.Plot
}

func NewCpuWidget() Widget {
	cores := runtime.NumCPU()

	prevUserProcTime := make([]uint64, cores)
	prevTotalTime := make([]uint64, cores)
	data := make([][]float64, cores)

	plot := tWidgets.NewPlot()
	plot.Title = " CPU Usage "

	termWidth, _ := tui.TerminalDimensions()
	for i := 0; i < cores; i++ {
		data[i] = make([]float64, termWidth/3+1)
	}

	plot.AxesColor = tui.ColorWhite

	plot.HorizontalScale = 3

	plot.ShowAxes = false
	plot.MaxVal = 100

	plot.Data = data

	return &CpuWidget{
		cores:            cores,
		prevUserProcTime: prevUserProcTime,
		prevTotalTime:    prevTotalTime,
		data:             data,
		plot:             plot,
	}
}

func (widget *CpuWidget) Update() {
	// read data from /proc/stat

	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan() // ignore first line

	for core := 0; core < widget.cores; core++ {
		scanner.Scan()
		cpuTimeReport := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		cpuTimes := strings.Fields(cpuTimeReport)[1:]

		userProcTime, _ := strconv.ParseUint(cpuTimes[0], 10, 64)

		totalTime := uint64(0)
		for _, token := range cpuTimes {
			time, _ := strconv.ParseUint(token, 10, 64)
			totalTime += time
		}

		cpuTimeChange := userProcTime - widget.prevUserProcTime[core]
		totalTimeChange := totalTime - widget.prevTotalTime[core]
		cpuUsage := (float64(cpuTimeChange) / float64(totalTimeChange)) * 100.0

		widget.data[core] = append(widget.data[core], cpuUsage)
		widget.data[core] = widget.data[core][1:]
		
		widget.prevUserProcTime[core] = userProcTime
		widget.prevTotalTime[core] = totalTime
	}
	file.Close()
}

func (widget *CpuWidget) HandleSignal(event tui.Event) {
	// is there anything to handle?
}

func (widget *CpuWidget) GetUI() tui.Drawable {
	widget.plot.Data = widget.data
	return widget.plot
}
