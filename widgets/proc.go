package widgets

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"
)

func getFormattedString(pid, cmd, cpu string) string {
	if len(cmd) > 20 {
		cmd = cmd[:17] + "..."
	}
	return fmt.Sprintf("%7s  %20s  %4s%%", pid, cmd, cpu)
}

// Process process info
type Process struct {
	pid           string
	cmd           string
	totalCPUUsage uint64
	cpuUsage      float64
}

func (process *Process) getString() string {
	return getFormattedString(
		process.pid,
		process.cmd,
		fmt.Sprintf("%2.1f", process.cpuUsage),
	)
}

// TODO: cpuStats refactoring is needed
// ProcessWidget process widget
type ProcessWidget struct {
	listWidget *tWidgets.List

	cpuStats    CpuStats
	processList []*Process
	cursor      int
}

// NewProcessWidget get new process widget
func NewProcessWidget() Widget {
	listWidget := tWidgets.NewList()
	listWidget.Title = "Process List"
	listWidget.TextStyle = tui.NewStyle(tui.ColorYellow)
	listWidget.Rows = make([]string, 0)

	cpuStats, err := getCpuStats()
	if err != nil {
		log.Fatal(err)
	}

	return &ProcessWidget{
		listWidget:  listWidget,
		cpuStats:    cpuStats,
		processList: make([]*Process, 0),
		cursor:      1,
	}
}

// Update update process data
func (widget *ProcessWidget) Update() {
	// update cpu data
	curCPUStat, err := getCpuStats()
	if err != nil {
		return
	}

	var totalTimeDiff uint64
	for i := 0; i < widget.cpuStats.cores; i++ {
		totalTimeDiff += curCPUStat.stats[i].totalTime - widget.cpuStats.stats[i].totalTime
		widget.cpuStats.stats[i].totalTime = curCPUStat.stats[i].totalTime
	}

	// update process data
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return
	}
	processList := widget.parseProcessList(files, totalTimeDiff)
	widget.processList = processList
}

func (widget *ProcessWidget) HandleSignal(event tui.Event) {
	switch event.ID {
	case "<Up>":
		if widget.cursor-1 > 0 {
			widget.cursor--
		}
	case "<Down>":
		if widget.cursor+1 < len(widget.processList) {
			widget.cursor++
		}
	}
}

func (widget *ProcessWidget) GetUI() tui.Drawable {
	drawWidget := widget.listWidget

	drawWidget.SelectedRow = widget.cursor
	drawWidget.Rows = widget.getRows()
	return drawWidget
}

func (widget *ProcessWidget) parseProcessList(files []fs.FileInfo, totalTime uint64) []*Process {
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

		cmd := strings.TrimSpace(string(cmdBytes))
		statBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/stat", pid))
		if err != nil {
			continue
		}

		statStrings := strings.Split(strings.TrimSpace(string(statBytes)), " ")
		curCPUUsage, err := strconv.ParseUint(statStrings[13], 10, 64)
		if err != nil {
			continue
		}

		var prevCPUUsage uint64
		prevProcess, err := widget.findProcess(pid)
		if err == nil {
			prevCPUUsage = prevProcess.totalCPUUsage
		}

		cpuUsage := float64((curCPUUsage-prevCPUUsage)*uint64(widget.cpuStats.cores)*100) / float64(totalTime)
		process := &Process{
			pid:           file.Name(),
			cmd:           cmd,
			totalCPUUsage: curCPUUsage,
			cpuUsage:      cpuUsage,
		}
		processList = append(processList, process)
	}
	return processList
}

func (widget *ProcessWidget) findProcess(pid string) (*Process, error) {
	for _, process := range widget.processList {
		if process.pid == pid {
			return process, nil
		}
	}
	return nil, errors.New("Not Found")
}

func (widget *ProcessWidget) getRows() []string {
	rows := make([]string, len(widget.processList)+1)
	rows[0] = getFormattedString("PID", "COMMAND", "CPU")

	sort.Slice(widget.processList, func(i int, j int) bool {
		return widget.processList[i].cpuUsage > widget.processList[j].cpuUsage
	})
	for i, process := range widget.processList {
		rows[i+1] = process.getString()
	}
	return rows
}
