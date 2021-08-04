package widgets

import (
	"fmt"
	"io/ioutil"
	"strings"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
	"github.com/sunghun7511/gotop/util"
)

type MemoryWidget struct {
	history []float64
	widget  *tWidgets.Sparkline
	group   *tWidgets.SparklineGroup
}

func NewMemoryWidget() Widget {
	widget := tWidgets.NewSparkline()
	widget.LineColor = tui.ColorGreen
	widget.MaxVal = 100

	group := tWidgets.NewSparklineGroup(widget)
	group.Title = "Memory Usage"

	termWidth, _ := tui.TerminalDimensions()
	return &MemoryWidget{
		history: make([]float64, termWidth/2-2),
		widget:  widget,
		group:   group,
	}
}

func (widget *MemoryWidget) Update() {
	information := readMemoryInformation()
	total := information["MemTotal"]
	available := information["MemAvailable"]

	value := float64(total-available) / float64(total) * 100
	util.PushUsageData(widget.history, value)
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

func readMemoryInformation() map[string]int64 {
	m := make(map[string]int64)

	content, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		var key string
		var value int64

		if len(line) == 0 {
			continue
		}

		if _, err := fmt.Sscanf(line, "%s %d", &key, &value); err != nil {
			panic(err)
		}
		key = strings.TrimRight(key, ":")

		m[key] = value
	}

	return m
}
