package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Validate(jsonData string, rules map[string]string) (err error, validationErrors map[string]string) {
	validationErrors = make(map[string]string)
	var decodedJson map[string]interface{}

	err = json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	for ruleKey, ruleValue := range rules {
		val, ok := decodedJson[ruleKey]

		for _, v := range strings.Split(fmt.Sprintf("%v", ruleValue), "|") {
			switch v {
			case "required":
				if !ok {
					validationErrors[ruleKey] = ruleKey + " is required"
					continue
				}
			case "string":
				if reflect.TypeOf(val) != reflect.TypeOf("") {
					validationErrors[ruleKey] = ruleKey + " must be a string"
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
