package widgets

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

func getFormattedString(pid, cmd string) string {
	if len(cmd) > 20 {
		cmd = cmd[:17] + "..."
	}
	return fmt.Sprintf("%7s  %20s", pid, cmd)
}

type Process struct {
	pid string
	cmd string
}

func (process *Process) getString() string {
	return getFormattedString(process.pid, process.cmd)
}

type ProcessWidget struct {
	processList []*Process
	cursor      int
	listWidget  *tWidgets.List
}

func NewProcessWidget() Widget {
	listWidget := tWidgets.NewList()
	listWidget.Title = "Process List"
	listWidget.TextStyle = tui.NewStyle(tui.ColorYellow)
	listWidget.Rows = make([]string, 0)

	return &ProcessWidget{
		processList: make([]*Process, 0),
		cursor:      1,
		listWidget:  listWidget,
	}
}

func (widget *ProcessWidget) Update() {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return
	}

	processList := make([]*Process, 0)
	for _, file := range files {
		pid := file.Name()
		_, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}
		cmdBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/comm", pid))
		if err != nil {
			continue
		}

		process := &Process{
			pid: file.Name(),
			cmd: strings.TrimSpace(string(cmdBytes)),
		}
		processList = append(processList, process)
	}
	widget.processList = processList
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
