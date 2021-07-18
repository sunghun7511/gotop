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

	endChan       chan struct{}
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

func handleSignal(e tui.Event) {
	if e.ID == "q" || e.ID == "<C-c>" {
		endChan <- struct{}{}
		return
	}

	memoryWidget.HandleSignal(e)
	processWidget.HandleSignal(e)
}

func updateWidgets() {
	cpuWidget.Update()
	memoryWidget.Update()
	processWidget.Update()
}

func handleEvents() {
	go func() {
		uiEvents := tui.PollEvents()
		for {
			e := <-uiEvents
			handleSignal(e)
			render()
		}
	}()
	go func() {
		updateStatTicker := time.NewTicker(1 * time.Second)
		for {
			<-updateStatTicker.C
			updateWidgets()
			render()
		}
	}()
}

func main() {
	if err := tui.Init(); err != nil {
		panic(err)
	}
	defer tui.Close()

	initWidgets()
	handleEvents()

	endChan = make(chan struct{})
	<-endChan
}
