package validation

import (
	"encoding/json"
	"errors"
	"reflect"
)

func ValidateField(fieldName string, fieldValue interface{}, fieldRules []string) (error, string) {
	for _, rule := range fieldRules {
		if rule == "" {
			continue
		}

		ruleFunc, ruleExists := rules[rule]
		if !ruleExists {
			err := errors.New("unknown validation rule: " + rule)
			return err, ""
		}
		ruleValue := getRuleValue(rule)
		err, validationError := ruleFunc(fieldName, fieldValue, ruleValue)
		if err != nil {
			return err, ""
		}
		if validationError != "" {
			return nil, validationError
		}
	}
	return nil, ""
}

func ValidateStruct(structData interface{}, validationRules map[string][]string, parentName string) (error, map[string]string) {
	validationErrors := make(map[string]string)

	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)
	if t.Kind() != reflect.Struct {
		err := errors.New("can only proceed `struct` kind")
		return err, nil
	}

	addValidationTagRules(t, validationRules, parentName)

	err := validateStructFields(t, v, parentName, validationRules, validationErrors)
	if err != nil {
		return err, nil
	}

	fieldsExists := make(map[string]bool)
	registerStructFields(structData, parentName, fieldsExists)
	validateNonExistRequiredFields(validationRules, fieldsExists, validationErrors)

	return nil, validationErrors
}

func ValidateJson(jsonData string, validationRules map[string][]string) (error, map[string]string) {
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	return ValidateMap(decodedJson, validationRules, "")
}

func ValidateMap(mapData map[string]interface{}, validationRules map[string][]string, parentName string) (error, map[string]string) {
	validationErrors := make(map[string]string)

	err := validateMapFields(mapData, parentName, validationRules, validationErrors)
	if err != nil {
		return err, nil
	}

	fieldsExists := make(map[string]bool)
	registerMapFields(mapData, parentName, fieldsExists)
	validateNonExistRequiredFields(validationRules, fieldsExists, validationErrors)

	return nil, validationErrors
}
