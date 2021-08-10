package util

// PushUsageData data 배열에 value를 추가하고 맨 앞 데이터를 삭제한다.
func PushUsageData(data []float64, value float64) []float64 {
	data = append(data, value)
	data = data[1:]
	return data
}

// 2D slice 의 각 row 의 마지막 원소를 원하는 개수만큼 가져옵니다.
func GetLastElementsOfEachRow(data [][]float64, length int) [][]float64 {
	rows := len(data)
	lastElements := make([][]float64, rows)

	for i, row := range data {
		lastElements[i] = row[len(row)-length:]
	}
	return lastElements
}

// slice 의 마지막 원소를 원하는 개수만큼 가져옵니다.
func GetLastElements(data []float64, length int) []float64 {
	return data[len(data)-length:]
}
