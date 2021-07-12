package widgets

import (
	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

type CpuWidget struct {
	data [][]float64
	plot *tWidgets.Plot
}

func NewCpuWidget() Widget {
	plot := tWidgets.NewPlot()
	plot.Title = "CPU Usage"

	plot.AxesColor = tui.ColorWhite
	plot.LineColors[0] = tui.ColorBlue

	plot.ShowAxes = false
	plot.MaxVal = 100

	data := makeData()
	plot.Data = data

	return &CpuWidget{
		data: data,
		plot: plot,
	}
}

func makeData() [][]float64 {
	data := make([][]float64, 1)
	n := 100
	data[0] = make([]float64, n)

	for i := 0; i < n; i++ {
		data[0][i] = float64(i) + 5
	}
	return data
}

func (widget *CpuWidget) Update() {
	// read data from /proc/stat
}

func (widget *CpuWidget) HandleSignal(event tui.Event) {
	// is there anything to handle?
}

func (widget *CpuWidget) GetUI() tui.Drawable {
	widget.plot.Data = widget.data
	return widget.plot
}
