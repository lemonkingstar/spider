package pconv

import (
	"reflect"
	"strings"
)

// ArrayContainsStr checks target is in the source array.
func ArrayContainsStr(sourceArr []string, target string) bool {
	for _, s := range sourceArr {
		if s == target {
			return true
		}
	}
	return false
}

// ArrayContainsInt checks target is in the source array.
func ArrayContainsInt(sourceArr []int, target int) bool {
	for _, s := range sourceArr {
		if s == target {
			return true
		}
	}
	return false
}

// ArrayContainsInt64 checks target is in the source array.
func ArrayContainsInt64(sourceArr []int64, target int64) bool {
	for _, s := range sourceArr {
		if s == target {
			return true
		}
	}
	return false
}

// InArray checks target is in the source array.
func InArray(sourceArr interface{}, target interface{}) bool {
	sourceArrValue := reflect.ValueOf(sourceArr)
	switch reflect.TypeOf(sourceArrValue).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < sourceArrValue.Len(); i++ {
			if sourceArrValue.Index(i).Interface() == target {
				return true
			}
		}
	case reflect.Map:
		if sourceArrValue.MapIndex(reflect.ValueOf(target)).IsValid() {
			return true
		}
	default:
	}
	return false
}

// ArrayUnique removes duplicates from array.
func ArrayUnique(arr interface{}) (ret []interface{}) {
	ret = make([]interface{}, 0)
	va := reflect.ValueOf(arr)
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if !InArray(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

// StrArrayUnique removes duplicates from a string array.
func StrArrayUnique(arr []string, delEmpty bool) (ret []string) {
	unique := make(map[string]struct{})
	for _, v := range arr {
		if delEmpty && strings.TrimSpace(v) == "" {
			continue
		}
		unique[v] = struct{}{}
	}
	ret = make([]string, len(unique))
	idx := 0
	for k := range unique {
		ret[idx] = k
		idx += 1
	}
	return ret
}

// StrArrayUniqueDelEmpty removes duplicates from a string array.
func StrArrayUniqueDelEmpty(arr []string) (ret []string) {
	ret = StrArrayUnique(arr, true)
	return
}

// IntArrayUnique removes duplicates from a integer array.
func IntArrayUnique(arr []int) (ret []int) {
	unique := make(map[int]struct{})
	for _, v := range arr {
		unique[v] = struct{}{}
	}
	ret = make([]int, len(unique))
	idx := 0
	for k := range unique {
		ret[idx] = k
		idx += 1
	}
	return ret
}

// Int64ArrayUnique removes duplicates from a integer array.
func Int64ArrayUnique(arr []int64) (ret []int64) {
	unique := make(map[int64]struct{})
	for _, v := range arr {
		unique[v] = struct{}{}
	}
	ret = make([]int64, len(unique))
	idx := 0
	for k := range unique {
		ret[idx] = k
		idx += 1
	}
	return ret
}

// BoolArrayUnique removes duplicates from a boolean array.
func BoolArrayUnique(arr []bool) (ret []bool) {
	ret = make([]bool, 0)
	trueExist := false
	falseExist := false
	for _, item := range arr {
		if item == true {
			trueExist = true
		}
		if item == false {
			falseExist = true
		}
	}
	if trueExist {
		ret = append(ret, true)
	}
	if falseExist {
		ret = append(ret, false)
	}
	return ret
}

// StrArrayDiff queries difference set between two string array.
func StrArrayDiff(firstArr []string, secondArr []string) []string {
	diffStr := make([]string, 0)
	for _, i := range firstArr {
		isIn := false
		for _, j := range secondArr {
			if i == j {
				isIn = true
				break
			}
		}
		if !isIn {
			diffStr = append(diffStr, i)
		}
	}
	return diffStr
}

// IntArrayIntersection queries intersection between two integer array.
func IntArrayIntersection(firstArr []int64, secondArr []int64) []int64 {
	intersectInt := make([]int64, 0)
	intMap := make(map[int64]bool)
	for _, i := range firstArr {
		intMap[i] = true
	}
	for _, j := range secondArr {
		if _, ok := intMap[j]; ok == true {
			intersectInt = append(intersectInt, j)
		}
	}
	return intersectInt
}

// StrArrayIntersection queries intersection between two string array.
func StrArrayIntersection(firstArr []string, secondArr []string) []string {
	intersectStr := make([]string, 0)
	strMap := make(map[string]bool)
	for _, i := range firstArr {
		strMap[i] = true
	}
	for _, j := range secondArr {
		if _, ok := strMap[j]; ok == true {
			intersectStr = append(intersectStr, j)
		}
	}
	return intersectStr
}

// ExistStrIntersection checks intersection between two string array.
func ExistStrIntersection(source []string, target []string) bool {
	for _, v := range target {
		for _, k := range source {
			if k == v {
				return true
			}
		}
	}
	return false
}

// ExistIntIntersection checks intersection between two integer array.
func ExistIntIntersection(source []int64, target []int64) bool {
	for _, v := range target {
		for _, k := range source {
			if k == v {
				return true
			}
		}
	}
	return false
}
