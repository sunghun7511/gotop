package core

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
	return getUsage(fmt.Sprintf("/proc/%s/stat", pid), 13)
}

// GetMemUsage pid 에 해당하는 프로세스의 memory 사용량을 가져옴
func GetMemUsage(pid string) (uint64, error) {
	return getUsage(fmt.Sprintf("/proc/%s/statm", pid), 1)
}

// must use 1 line file
func getUsage(fileName string, idx int) (uint64, error) {
	statBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	statStrings := strings.Split(strings.TrimSpace(string(statBytes)), " ")
	usage, err := strconv.ParseUint(statStrings[idx], 10, 64)
	if err != nil {
		return 0, err
	}
	return usage, nil
}
