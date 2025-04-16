package helper

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

// GenUniqueId .
func GenUniqueId() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

// Marshal .
func Marshal(data interface{}) interface{} {
	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	for i := 0; i < typeof.Elem().NumField(); i++ {
		if valueof.Elem().Field(i).IsZero() {
			def := typeof.Elem().Field(i).Tag.Get("default")
			if def != "" {
				switch typeof.Elem().Field(i).Type.String() {
				case "int":
					result, _ := strconv.Atoi(def)
					valueof.Elem().Field(i).SetInt(int64(result))
				case "uint":
					result, _ := strconv.ParseUint(def, 10, 64)
					valueof.Elem().Field(i).SetUint(result)
				case "string":
					valueof.Elem().Field(i).SetString(def)
				case "interface {}":
					valueof.Elem().Field(i).SetZero()
				}
			}
		}
	}
	return data
}

// MarshalJson .
func MarshalJson(data interface{}) ([]byte, error) {
	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	for i := 0; i < typeof.Elem().NumField(); i++ {
		if valueof.Elem().Field(i).IsZero() {
			def := typeof.Elem().Field(i).Tag.Get("default")
			if def != "" {
				switch typeof.Elem().Field(i).Type.String() {
				case "int":
					result, _ := strconv.Atoi(def)
					valueof.Elem().Field(i).SetInt(int64(result))
				case "uint":
					result, _ := strconv.ParseUint(def, 10, 64)
					valueof.Elem().Field(i).SetUint(result)
				case "string":
					valueof.Elem().Field(i).SetString(def)
				}
			}
		}
	}
	return json.Marshal(data)
}
