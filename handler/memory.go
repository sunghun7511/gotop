package handler

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sunghun7511/gotop/model"
)

// ReadMemoryInformation 메모리 정보를 파일에서 읽어옴
func ReadMemoryInformation() model.MemoryStat {
	m := make(map[string]int64)

	content, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		var key string
		var value int64

		if len(line) == 0 {
			continue
		}

		if _, err := fmt.Sscanf(line, "%s %d", &key, &value); err != nil {
			panic(err)
		}
		key = strings.TrimRight(key, ":")

		m[key] = value
	}

	return model.MemoryStat{
		Total:     m["MemTotal"],
		Available: m["MemAvailable"],
	}
}
