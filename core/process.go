package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var userMap map[int]string

func init() {
	err := setUserMap()
	if err != nil {
		panic(err)
	}
}

func setUserMap() error {
	statBytes, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}

	userMap = make(map[int]string)
	for _, line := range strings.Split(string(statBytes), "\n") {
		if len(line) == 0 {
			continue
		}

		splited := strings.Split(line, ":")

		userName := splited[0]
		userId, err := strconv.ParseUint(splited[2], 10, 64)
		if err != nil {
			return err
		}

		userMap[int(userId)] = userName
	}
	return nil
}

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

// GetMemUsage pid 에 해당하는 프로세스의 user를 가져옴
func GetUser(pid string) (string, error) {
	statBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/status", pid))
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(statBytes), "\n") {
		if len(line) == 0 {
			continue
		}

		splited := strings.Split(line, "\t")
		if splited[0] == "Uid:" {
			uid, err := strconv.ParseUint(splited[1], 10, 64)
			if err != nil {
				return "", err
			}

			userName, ok := userMap[int(uid)]
			if !ok {
				return "", errors.New("Undefined user")
			}
			return userName, nil
		}
	}
	return "", errors.New("Cannot get user")
}
