package main

import (
	"os"
	"time"

	tui "github.com/gizak/termui/v3"

	"github.com/sunghun7511/gotop/widgets"
)

var (
	widgetList []widgets.Widget

	dummyWidget widgets.Widget
)

func initWidgets() {
	dummyWidget = widgets.NewDummyWidget()

	widgetList = make([]widgets.Widget, 0)
	widgetList = append(widgetList, dummyWidget)
}

func render() {
	tui.Render(dummyWidget.GetUI())
}

func handleSignal(e tui.Event) {
	if e.ID == "q" || e.ID == "<C-c>" {
		os.Exit(0)
	}

	for _, v := range widgetList {
		v.HandleSignal(e)
	}
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
			for _, v := range widgetList {
				v.Update()
			}
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
