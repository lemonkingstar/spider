package pmongo

import (
	"context"
	"reflect"
)

func ConvertToInterfaceSlice(value interface{}) []interface{} {
	rflVal := reflect.ValueOf(value)
	for rflVal.CanAddr() {
		rflVal = rflVal.Elem()
	}
	if rflVal.Kind() != reflect.Slice {
		return []interface{}{value}
	}

	result := make([]interface{}, 0)
	for i := 0; i < rflVal.Len(); i++ {
		if rflVal.Index(i).CanInterface() {
			result = append(result, rflVal.Index(i).Interface())
		}
	}

	return result
}

const (
	HTTPReadReference = "Cc_Read_Preference"
)

// GetDBReadPreference
func GetDBReadPreference(ctx context.Context) ReadPreferenceMode {
	val := ctx.Value(HTTPReadReference)
	if val != nil {
		mode, ok := val.(string)
		if ok {
			return ReadPreferenceMode(mode)
		}
	}
	return NilMode
}
