package pmap

import (
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// Map2Struct convert map into a struct with 'tagName'
//
//  eg:
//  self := MapStr{"testName":"testvalue"}
//  targetStruct := &struct{
//      Name string `field:"testName"`
//  }
//  After call the function Map2Struct(self, targetStruct, "field", false)
//  the targetStruct.Name value will be 'testvalue'
func Map2Struct(sourceMap map[string]interface{}, targetStruct interface{}, tagName string, weaklyTyped bool, v ...interface{}) error {
	config := &mapstructure.DecoderConfig{
		// pointer to the struct
		Result:   			targetStruct,
		// defaults to "mapstructure"
		TagName:  			tagName,
		// do weak conversion
		WeaklyTypedInput: 	weaklyTyped,
	}
	if len(v) > 0 {
		// custom data decode
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(v[0])
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	if err := decoder.Decode(sourceMap); err != nil {
		return err
	}
	return nil
}

func Map2StructDefault(sourceMap map[string]interface{}, targetStruct interface{}) error {
	return Map2Struct(sourceMap, targetStruct, "json", false)
}

// Struct2Map convert struct into a map
//
//  eg:
//  sourceStruct := &struct{
//      Name string `field:"testName"`
//  }
//  After call the function Struct2Map(sourceStruct, "field")
//  will return map info
func Struct2Map(sourceStruct interface{}, tagName string) (map[string]interface{}, error) {
	s := structs.New(sourceStruct)
	// defaults to "structs"
	s.TagName = tagName
	return s.Map(), nil
}

func Struct2MapDefault(sourceStruct interface{}) (map[string]interface{}, error) {
	return Struct2Map(sourceStruct, "json")
}
