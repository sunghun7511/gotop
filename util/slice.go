package util

// PushUsageData data 배열에 value를 추가하고 맨 앞 데이터를 삭제한다.
func PushUsageData(data []float64, value float64) []float64 {
	data = append(data, value)
	data = data[1:]
	return data
}
