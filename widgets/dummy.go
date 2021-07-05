package widgets

import (
	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

// TODO: 사용되지 않을 widget이라 삭제되어야 합니다.
type DummyWidget struct {
	value int
	gauge *tWidgets.Gauge
}

func NewDummyWidget() Widget {
	gauge := tWidgets.NewGauge()
	gauge.Title = "Slim Gauge"
	gauge.SetRect(0, 0, 30, 30)
	gauge.Percent = 0
	gauge.BarColor = tui.ColorRed
	gauge.BorderStyle.Fg = tui.ColorWhite
	gauge.TitleStyle.Fg = tui.ColorCyan

	return &DummyWidget{
		value: 0,
		gauge: gauge,
	}
}

func (widget *DummyWidget) Update() {
	widget.value += 1
}

func (widget *DummyWidget) HandleSignal(event tui.Event) {
	if event.ID == "<Space>" {
		widget.value = 0
	}
}

func (widget *DummyWidget) GetUI() tui.Drawable {
	widget.gauge.Percent = widget.value
	return widget.gauge
}
