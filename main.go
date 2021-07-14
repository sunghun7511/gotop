package main

import (
	"time"

	tui "github.com/gizak/termui/v3"

	"github.com/sunghun7511/gotop/widgets"
)

var (
	cpuWidget     widgets.Widget
	memoryWidget  widgets.Widget
	processWidget widgets.Widget
)

func initWidgets() {
	cpuWidget = widgets.NewCpuWidget()
	memoryWidget = widgets.NewMemoryWidget()
	processWidget = widgets.NewProcessWidget()
}

func render() {
	grid := tui.NewGrid()
	termWidth, termHeight := tui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		tui.NewRow(1.0/2, cpuWidget.GetUI()),
		tui.NewRow(1.0/2,
			tui.NewCol(1.0/2, memoryWidget.GetUI()),
			tui.NewCol(1.0/2, processWidget.GetUI()),
		),
	)

	tui.Render(grid)
}

func handleSignal(e tui.Event) bool {
	if e.ID == "q" || e.ID == "<C-c>" {
		return true
	}

	memoryWidget.HandleSignal(e)
	processWidget.HandleSignal(e)
	return false
}

func updateWidgets() {
	cpuWidget.Update()
	memoryWidget.Update()
	processWidget.Update()
}

func handleEvents() {
	uiEvents := tui.PollEvents()
	updateStatTicker := time.NewTicker(1 * time.Second)

	for {
		// 항상 키보드 입력을 우선시 합니다.
		select {
		case e := <-uiEvents:
			if handleSignal(e) {
				return
			}
			render()
		default:
		}

		select {
		case <-updateStatTicker.C:
			updateWidgets()
			render()
		default:
		}
	}
}

func main() {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	initWidgets()
	handleEvents()
}
