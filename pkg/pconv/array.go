package pconv

import (
	"reflect"
	"strings"
)

// ArrayContains
// 检测 set数组中是否存在指定 item
func ArrayContains(set []string, item string) bool {
	for _, s := range set {
		if s == item { return true }
	}
	return false
}

func ArrayContainsInt(set []int, item int) bool {
	for _, s := range set {
		if s == item { return true }
	}
	return false
}

func ArrayContainsInt64(set []int64, item int64) bool {
	for _, s := range set {
		if s == item { return true }
	}
	return false
}

// InArray
// 检测 target集合中是否包含 obj
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// ArrayUnique
// 数组去重
func ArrayUnique(a interface{}) (ret []interface{}) {
	ret = make([]interface{}, 0)
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if !InArray(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

func strArrayUnique(a []string, delEmpty bool) (ret []string) {
	unique := make(map[string]struct{})
	for _, v := range a {
		if delEmpty && strings.TrimSpace(v) == "" { continue }
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

// StrArrayUnique
// 字符串数组去重
func StrArrayUnique(a []string) (ret []string) {
	ret = strArrayUnique(a, true)
	return
}

// IntArrayUnique get unique int array
func IntArrayUnique(a []int) (ret []int) {
	unique := make(map[int]struct{})
	for _, v := range a {
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

func Int64ArrayUnique(a []int64) (ret []int64) {
	unique := make(map[int64]struct{})
	for _, v := range a {
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

// BoolArrayUnique get unique int array
func BoolArrayUnique(a []bool) (ret []bool) {
	ret = make([]bool, 0)
	trueExist := false
	falseExist := false
	for _, item := range a {
		if item == true { trueExist = true }
		if item == false { falseExist = true }
	}
	if trueExist { ret = append(ret, true) }
	if falseExist { ret = append(ret, false) }
	return ret
}

// StrArrDiff 查询差集
func StrArrDiff(slice1 []string, slice2 []string) []string {
	diffStr := make([]string, 0)
	for _, i := range slice1 {
		isIn := false
		for _, j := range slice2 {
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

// IntArrIntersection 查询交集
func IntArrIntersection(slice1 []int64, slice2 []int64) []int64 {
	intersectInt := make([]int64, 0)
	intMap := make(map[int64]bool)
	for _, i := range slice1 {
		intMap[i] = true
	}
	for _, j := range slice2 {
		if _, ok := intMap[j]; ok == true {
			intersectInt = append(intersectInt, j)
		}
	}
	return intersectInt
}

func StrArrIntersection(slice1 []string, slice2 []string) []string {
	intersectStr := make([]string, 0)
	strMap := make(map[string]bool)
	for _, i := range slice1 {
		strMap[i] = true
	}
	for _, j := range slice2 {
		if _, ok := strMap[j]; ok == true {
			intersectStr = append(intersectStr, j)
		}
	}
	return intersectStr
}

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
