package widgets

import (
	"fmt"
	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
	"io/ioutil"
	"strings"
)

type MemoryWidget struct {
	information map[string]int64
	gauge       *tWidgets.Gauge
}

func NewMemoryWidget() Widget {
	gauge := tWidgets.NewGauge()
	gauge.Title = "Memory Usage"
	gauge.Percent = 0
	gauge.BarColor = tui.ColorRed
	gauge.BorderStyle.Fg = tui.ColorWhite
	gauge.TitleStyle.Fg = tui.ColorCyan

	return &MemoryWidget{
		information: make(map[string]int64),
		gauge:       gauge,
	}
}

func (widget *MemoryWidget) Update() {
	widget.information = readMemoryInformation()
}

func (widget *MemoryWidget) HandleSignal(event tui.Event) {
	if event.ID == "<Space>" {
		widget.GetUI()
	}
}

func (widget *MemoryWidget) GetUI() tui.Drawable {
	total := widget.information["MemTotal"]
	available := widget.information["MemAvailable"]

	widget.gauge.Percent = int(float64(total-available) / float64(total) * 100)
	return widget.gauge
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
