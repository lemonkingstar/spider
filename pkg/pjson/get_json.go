package pjson

import (
	"github.com/tidwall/gjson"
)

type JsonDocument struct {
	doc  	string
}

func NewJsonDocument(json string) *JsonDocument {
	return &JsonDocument{
		doc: json,
	}
}

func (jd *JsonDocument) GetJson() string {
	return jd.doc
}

func (jd *JsonDocument) SetJson(json string) {
	jd.doc = json
}

func (jd *JsonDocument) GetPathString(path string) string {
	return gjson.Get(jd.doc, path).String()
}

func (jd *JsonDocument) GetPathInt(path string) int64 {
	return gjson.Get(jd.doc, path).Int()
}

func (jd *JsonDocument) GetPathBool(path string) bool {
	return gjson.Get(jd.doc, path).Bool()
}
