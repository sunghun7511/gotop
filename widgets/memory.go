package widgets

import (
	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"

	constants "github.com/sunghun7511/gotop/constants"
	"github.com/sunghun7511/gotop/core"
	"github.com/sunghun7511/gotop/util"
)

type MemoryWidget struct {
	history []float64
	widget  *tWidgets.Sparkline
	group   *util.FixedSparklineGroup
}

func NewMemoryWidget() Widget {
	widget := tWidgets.NewSparkline()
	widget.LineColor = tui.ColorGreen
	widget.MaxVal = 100

	group := util.FixedSparklineGroup{*tWidgets.NewSparklineGroup(widget)}
	group.Title = "Memory Usage"

	return &MemoryWidget{
		history: make([]float64, constants.MaxDataLength),
		widget:  widget,
		group:   &group,
	}
}

func (widget *MemoryWidget) Update() {
	information := core.ReadMemoryInformation()
	total := information.Total
	available := information.Available

	value := float64(total-available) / float64(total) * 100
	widget.history = util.PushUsageData(widget.history, value)
}

func (widget *MemoryWidget) HandleSignal(event tui.Event) {
	if event.ID == "<Space>" {
		widget.GetUI()
	}
}

func calculateMemoryWidgetDataLength() int {
	termWidth, _ := tui.TerminalDimensions()
	return termWidth/2 - 2
}

func (widget *MemoryWidget) GetUI() tui.Drawable {
	dataLength := calculateMemoryWidgetDataLength()
	widget.widget.Data = util.GetLastElements(widget.history, dataLength)
	return widget.group
}
