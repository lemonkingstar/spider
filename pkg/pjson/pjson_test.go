package pjson

import (
	"fmt"
	"testing"
)

func TestJsonIter(t *testing.T) {
	iter := map[string]interface{}{
		"package": "json-iterator",
		"feature": "high-performance",
		"year": 2022,
	}

	b, _ := MarshalToString(iter)
	fmt.Println(b)
}

func TestJsonDocument(t *testing.T) {
	json := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	jsonDoc := NewJsonDocument(json)
	fmt.Println(jsonDoc.GetPathString("name.first"))

	addr := "Shanghai"
	json2, _ := jsonDoc.SetPathValue("address", addr)
	fmt.Println(json2)
	jsonDoc.SetJson(json2)
	json3, _ := jsonDoc.SetPathValue("name.second", "yy")
	fmt.Println(json3)
}
