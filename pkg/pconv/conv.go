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

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StrToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FloatToStr(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}
