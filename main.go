package main

import (
	"os"
	"time"

	tui "github.com/gizak/termui/v3"

	"github.com/sunghun7511/gotop/widgets"
)

var (
	dummyWidget widgets.Widget
)

func initWidgets() {
	dummyWidget = widgets.NewDummyWidget()
}

func render() {
	grid := tui.NewGrid()
	termWidth, termHeight := tui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		tui.NewRow(1.0/2, dummyWidget.GetUI()),
		tui.NewRow(1.0/2,
			tui.NewCol(1.0/2, dummyWidget.GetUI()),
			tui.NewCol(1.0/2, dummyWidget.GetUI()),
		),
	)

	tui.Render(grid)
}

func handleSignal(e tui.Event) {
	if e.ID == "q" || e.ID == "<C-c>" {
		os.Exit(0)
	}

	dummyWidget.HandleSignal(e)
}

func updateWidgets() {
	dummyWidget.Update()
}

func handleEvents() {
	uiEvents := tui.PollEvents()
	updateStatTicker := time.NewTicker(1 * time.Second)

	for {
		// 항상 키보드 입력을 우선시 합니다.
		select {
		case e := <-uiEvents:
			handleSignal(e)
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
