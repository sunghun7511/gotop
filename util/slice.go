package util

// PushUsageData datas 배열에 data를 추가하고 맨 앞 데이터를 삭제한다.
func PushUsageData(datas []float64, data float64) []float64 {
	datas = append(datas, data)
	datas = datas[1:]
	return datas
}
