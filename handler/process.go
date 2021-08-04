package handler

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// GetCommand pid 에 해당하는 프로세스 커맨드를 가져옴
func GetCommand(pid string) (string, error) {
	cmdBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/comm", pid))
	if err != nil {
		return "", err
	}
	cmd := strings.TrimSpace(string(cmdBytes))
	return cmd, nil
}

// GetCPUUsage pid 에 해당하는 프로세스의 cpu 사용량을 가져옴
func GetCPUUsage(pid string) (uint64, error) {
	statBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/stat", pid))
	if err != nil {
		return 0, err
	}
	statStrings := strings.Split(strings.TrimSpace(string(statBytes)), " ")
	curCPUUsage, err := strconv.ParseUint(statStrings[13], 10, 64)
	if err != nil {
		return 0, err
	}
	return curCPUUsage, nil
}

// GetMemUsage pid 에 해당하는 프로세스의 memory 사용량을 가져옴
func GetMemUsage(pid string) (uint64, error) {
	statmBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/statm", pid))
	if err != nil {
		return 0, nil
	}
	statmStrings := strings.Split(strings.TrimSpace(string(statmBytes)), " ")
	resident, err := strconv.ParseUint(statmStrings[1], 10, 64)
	if err != nil {
		return 0, nil
	}
	return resident, nil
}
