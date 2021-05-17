package validation

import (
	"encoding/json"
	"errors"
	"strings"
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

FIELDS:
	for field, fieldRules := range validationRules {
		fieldVal, fieldExists := mapData[field]

		for _, rule := range fieldRules {
			ruleFunc, ruleExists := rules[rule]
			if !ruleExists {
				err = errors.New("unknown validation rule: " + rule)
				return err, validationErrors
			}

			var ruleVal string
			if strings.ContainsRune(rule, ':') {
				ruleVal = strings.Split(rule, ":")[1]
			}

			err, validationError := ruleFunc(field, fieldVal, fieldExists, ruleVal)
			if err != nil {
				return err, validationErrors
			}

			if validationError != "" {
				validationErrors[field] = validationError
				continue FIELDS
			}

		}
	}

	return
}
