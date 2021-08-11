package widgets

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	tui "github.com/gizak/termui/v3"
	tWidgets "github.com/gizak/termui/v3/widgets"

	"github.com/sunghun7511/gotop/core"
	"github.com/sunghun7511/gotop/model"
)

func getFormattedString(pid, user, cpu, mem, cmd string, cursor int) string {
	if len(user) > 8 {
		user = user[0:6] + ".."
	}

	str := fmt.Sprintf("%7s  %-8s  %6s  %6s  %s", pid, user, cpu, mem, cmd)
	if cursor >= len(str) {
		return ""
	}
	return str[cursor:]
}

func getString(process *model.Process, cursor int) string {
	return getFormattedString(
		process.Pid,
		process.User,
		fmt.Sprintf("%2.1f%%", process.CPUUsage),
		fmt.Sprintf("%2.1f%%", process.MemUsage),
		process.Cmd,
		cursor,
	)
}

// ProcessWidget process widget
type ProcessWidget struct {
	listWidget *tWidgets.List

	cpuStats         model.CpuStats
	totalMem         uint64
	pageSizeKB       uint64
	processList      []*model.Process
	cursor           int
	horizontalCursor int
}

// NewProcessWidget get new process widget
func NewProcessWidget() Widget {
	listWidget := tWidgets.NewList()
	listWidget.Title = "Process List"
	listWidget.TextStyle = tui.NewStyle(tui.ColorYellow)
	listWidget.Rows = make([]string, 0)

	cpuStats, err := core.GetCPUStats()
	if err != nil {
		log.Fatal(err)
	}

	totalMem := uint64(core.ReadMemoryInformation().Total)

	// KB로 단위를 맞추기 위해 1024를 나눠줍니다.
	pageSizeKB := os.Getpagesize() / 1024
	return &ProcessWidget{
		listWidget:       listWidget,
		cpuStats:         cpuStats,
		totalMem:         totalMem,
		pageSizeKB:       uint64(pageSizeKB),
		processList:      make([]*model.Process, 0),
		cursor:           1,
		horizontalCursor: 0,
	}
}

// Update update process data
func (widget *ProcessWidget) Update() {
	// update cpu data
	curCPUStat, err := core.GetCPUStats()
	if err != nil {
		return
	}

	var totalTimeDiff uint64
	for i := 0; i < widget.cpuStats.Cores; i++ {
		totalTimeDiff += curCPUStat.Stats[i].TotalTime - widget.cpuStats.Stats[i].TotalTime
		widget.cpuStats.Stats[i].TotalTime = curCPUStat.Stats[i].TotalTime
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
		} else if widget.cursor -1 == 0{
			// For showing header
			// termui list widget does not support handle topRow in other package
			newListWidget := tWidgets.NewList()
			newListWidget.Title = widget.listWidget.Title
			newListWidget.TextStyle = widget.listWidget.TextStyle
			newListWidget.Rows = widget.listWidget.Rows

			widget.listWidget = newListWidget
		}
	case "<Down>":
		if widget.cursor+1 < len(widget.processList) {
			widget.cursor++
		}
	case "<Left>":
		if widget.horizontalCursor > 0 {
			widget.horizontalCursor--;
		}
	case "<Right>":
		widget.horizontalCursor++;
	case "K":
		_ = widget.killProcess()
	}
}

func (widget *ProcessWidget) GetUI() tui.Drawable {
	drawWidget := widget.listWidget

	drawWidget.SelectedRow = widget.cursor
	drawWidget.Rows = widget.getRows()
	return drawWidget
}

func (widget *ProcessWidget) parseProcessList(files []fs.FileInfo, totalTime uint64) []*model.Process {
	processList := make([]*model.Process, 0)
	for _, file := range files {
		pid := file.Name()
		_, err := strconv.Atoi(pid)
		if err != nil {
			continue
		}

		cmd, err := core.GetCommand(pid)
		if err != nil {
			continue
		}

		curCPUUsage, err := core.GetCPUUsage(pid)
		if err != nil {
			continue
		}

		var prevCPUUsage uint64
		prevProcess, err := widget.findProcess(pid)
		if err == nil {
			prevCPUUsage = prevProcess.TotalCPUUsage
		}
		cpuUsage := float64((curCPUUsage-prevCPUUsage)*uint64(widget.cpuStats.Cores)*100) / float64(totalTime)

		resident, err := core.GetMemUsage(pid)
		if err != nil {
			continue
		}
		memUsage := (float64)(resident*widget.pageSizeKB) / (float64)(widget.totalMem) * 100.0

		user, err := core.GetUser(pid)
		if err != nil {
			continue
		}

		process := &model.Process{
			Pid:           file.Name(),
			User:          user,
			Cmd:           cmd,
			TotalCPUUsage: curCPUUsage,
			CPUUsage:      cpuUsage,
			MemUsage:      memUsage,
		}
		processList = append(processList, process)
	}
	return processList
}

func (widget *ProcessWidget) findProcess(pid string) (*model.Process, error) {
	for _, process := range widget.processList {
		if process.Pid == pid {
			return process, nil
		}
	}
	return nil, errors.New("Not Found")
}

func (widget *ProcessWidget) getRows() []string {
	rows := make([]string, len(widget.processList)+1)
	rows[0] = getFormattedString("PID", "USER", "CPU", "MEM", "COMMAND", widget.horizontalCursor)

	sort.Slice(widget.processList, func(i int, j int) bool {
		return widget.processList[i].CPUUsage > widget.processList[j].CPUUsage
	})
	for i, process := range widget.processList {
		rows[i+1] = getString(process, widget.horizontalCursor)
	}
	return rows
}

func (widget *ProcessWidget) killProcess() error {
	if widget.cursor > len(widget.processList) {
		return errors.New("Cursor index out of range")
	}

	// ignore error, handle number string in parseProcessList
	pid, _ := strconv.Atoi(widget.processList[widget.cursor-1].Pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}
