package validation

import (
	"encoding/json"
	"errors"
	"reflect"
)

func ValidateJson(jsonData string, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	var decodedJson map[string]interface{}

	err = json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	return ValidateMap(decodedJson, validationRules)
}

func ValidateMap(mapData map[string]interface{}, validationRules map[string][]string) (err error, validationErrors map[string]string) {
	validationErrors = make(map[string]string)

	for key, rules := range validationRules {
		val, ok := mapData[key]

		for _, v := range rules {
			switch v {
			case "required":
				if !ok {
					validationErrors[key] = key + " is required"
					continue
				}
			case "string":
				if reflect.ValueOf(val).Kind() != reflect.String {
					validationErrors[key] = key + " must be a string"
					continue
				}
			default:
				err = errors.New("unknown validation rule: " + v)
				return err, validationErrors
			}
		}
	}

	return
}
