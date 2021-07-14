package widgets

import (
	"fmt"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

func getFormattedString(pid, cmd string) string {
	if len(cmd) > 20 {
		cmd = cmd[:17] + "..."
	}
	return fmt.Sprintf("%5s  %20s", pid, cmd)
}

type Process struct {
	pid int
	cmd string
}

func (process *Process) getString() string {
	return getFormattedString(fmt.Sprint(process.pid), process.cmd)
}

type ProcessWidget struct {
	processList []Process
	cursor      int
	listWidget  *tWidgets.List
}

func NewProcessWidget() Widget {
	listWidget := tWidgets.NewList()
	listWidget.Title = "Process List"
	listWidget.TextStyle = tui.NewStyle(tui.ColorYellow)
	listWidget.Rows = make([]string, 0)

	return &ProcessWidget{
		processList: make([]Process, 0),
		cursor:      1,
		listWidget:  listWidget,
	}
}

func (widget *ProcessWidget) Update() {
	// TODO: delete dummy
	widget.processList = []Process{
		{
			pid: 1,
			cmd: "12345",
		},
		{
			pid: 1342,
			cmd: "456",
		},
		{
			pid: 12346,
			cmd: "23",
		},
		{
			pid: 12,
			cmd: "1",
		},
	}
}

func (widget *ProcessWidget) HandleSignal(event tui.Event) {
	switch event.ID {
	case "<Up>":
		widget.cursor--
		widget.handleCursorOutBound()
	case "<Down>":
		widget.cursor++
		widget.handleCursorOutBound()
	}
}

func (widget *ProcessWidget) handleCursorOutBound() {
	if widget.cursor < 1 {
		widget.cursor = 1
	} else if widget.cursor > len(widget.processList) {
		widget.cursor = len(widget.processList)
	}
}

func (widget *ProcessWidget) GetUI() tui.Drawable {
	drawWidget := widget.listWidget

	drawWidget.SelectedRow = widget.cursor
	drawWidget.Rows = widget.getRows()
	return drawWidget
}

func (widget *ProcessWidget) getRows() []string {
	rows := make([]string, len(widget.processList)+1)

	rows[0] = getFormattedString("PID", "COMMAND")
	for i, process := range widget.processList {
		rows[i+1] = process.getString()
	}
	return rows
}
