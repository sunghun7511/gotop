package widgets

import (
	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"

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

	termWidth, _ := tui.TerminalDimensions()
	return &MemoryWidget{
		history: make([]float64, termWidth/2-2),
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

func (widget *MemoryWidget) GetUI() tui.Drawable {
	widget.widget.Data = widget.history
	return widget.group
}
