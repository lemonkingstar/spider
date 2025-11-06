package pconv

import (
	"strconv"
	"strings"
)

func SliceStrToInt64(sliceStr []string) ([]int64, error) {
	sliceInt := make([]int64, 0)
	for _, str := range sliceStr {
		if strings.TrimSpace(str) == "" {
			continue
		}

		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		sliceInt = append(sliceInt, id)
	}
	return sliceInt, nil
}

func SliceStrToInt(sliceStr []string) ([]int, error) {
	sliceInt := make([]int, 0)
	for _, str := range sliceStr {
		if strings.TrimSpace(str) == "" {
			continue
		}

		id, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}
		sliceInt = append(sliceInt, id)
	}
	return sliceInt, nil
}

func Str2Int(str string) (int, error) {
	return strconv.Atoi(str)
}

func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Str2Float(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func Int2Str(i int) string {
	return strconv.Itoa(i)
}

func Int642Str(i int64) string {
	return strconv.FormatInt(i,10)
}

func Float2Str(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}
