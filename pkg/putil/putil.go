package putil

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mohae/deepcopy"
)

// Clone 深拷贝
// Usage: obj := Clone(input)
// obj.(StructObj)
func Clone(input interface{}) interface{} {
	return deepcopy.Copy(input)
}

// GetObjString 字符串转换
func GetObjString(obj interface{}) (string, error) {
	switch t := obj.(type) {
	case nil:
		return "", nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case map[string]interface{}, []interface{}:
		rest, err := json.Marshal(t)
		if nil != err {
			return "", err
		}
		return string(rest), nil
	case json.Number:
		return t.String(), nil
	case string:
		return t, nil
	default:
		return fmt.Sprintf("%v", t), nil
	}
}
