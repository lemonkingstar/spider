package pconv

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func Type2Str(obj interface{}) (string, error) {
	switch v := obj.(type) {
	case nil:
		return "", errors.New("nil object")
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case string:
		return v, nil
	case json.Number:
		return v.String(), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func Type2Int64(obj interface{}) (int64, error) {
	switch v := obj.(type) {
	case nil:
		return 0, errors.New("nil object")
	case int:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 0, 64)
	case []byte:
		return strconv.ParseInt(string(v), 0, 64)
	default:
		return 0, errors.New("invalid num")
	}
}

func Type2Float(obj interface{}) (float64, error) {
	switch v := obj.(type) {
	case nil:
		return 0, errors.New("nil object")
	case int:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, errors.New("invalid num")
	}
}

func Type2StrArr(obj interface{}) ([]string, error) {
	switch v := obj.(type) {
	case nil:
		return nil, errors.New("nil object")
	case []string:
		return v, nil
	case []interface{}:
		arr := []string{}
		for _, val := range v {
			if str, err := Type2Str(val); err != nil {
				return nil, err
			} else {
				arr = append(arr, str)
			}
		}
		return arr, nil
	default:
		return nil, errors.New("invalid arr")
	}
}

func Type2Int64Arr(obj interface{}) ([]int64, error) {
	switch v := obj.(type) {
	case nil:
		return nil, errors.New("nil object")
	case []int64:
		return v, nil
	case []interface{}:
		arr := []int64{}
		for _, val := range v {
			if i, err := Type2Int64(val); err != nil {
				return nil, err
			} else {
				arr = append(arr, i)
			}
		}
		return arr, nil
	default:
		return nil, errors.New("invalid arr")
	}
}
