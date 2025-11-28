package pjson

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var _json = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	UseNumber:              true,
}.Froze()

// RegisterFuzzyDecoders
// 字符串转换/空数组兼容
func RegisterFuzzyDecoders() { extra.RegisterFuzzyDecoders() }

func Marshal(v interface{}) ([]byte, error) {
	return _json.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return _json.MarshalToString(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return _json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return _json.Unmarshal(data, v)
}

func UnmarshalFromString(str string, v interface{}) error {
	return _json.UnmarshalFromString(str, v)
}
