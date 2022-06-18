package utils

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func InterfaceToMap(in interface{}) map[string]interface{} {

	newMap := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			val := v.MapIndex(k)
			newMap[k.String()] = val.Interface()
		}
	}
	return newMap
}

func ExistingKeyInMap(in interface{}, key string) bool {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			if key == k.String() {
				return true
			}
		}
	}
	return false
}

func Transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}
