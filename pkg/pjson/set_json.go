package pjson

import (
	"github.com/tidwall/sjson"
)

func (jd *JsonDocument) SetPathValue(path string, value interface{}) (string, error) {
	s, err := sjson.Set(jd.doc, path, value)
	if err != nil {
		return s, err
	}
	jd.doc = s
	return s, nil
}
