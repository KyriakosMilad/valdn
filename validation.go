package validation

import (
	"encoding/json"
	"errors"
	"reflect"
)

type (
	Rules  map[string][]string
	Errors map[string]string
)

func ValidateField(fieldName string, fieldValue interface{}, fieldRules []string) (error, string) {
	for _, rule := range fieldRules {
		if rule == "" {
			continue
		}

		ruleName, ruleValue := splitRuleNameAndRuleValue(rule)
		ruleFunc, ruleExists := rules[ruleName]
		if !ruleExists {
			err := errors.New("unknown validation rule: " + ruleName)
			return err, ""
		}
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

func ValidateStruct(structData interface{}, validationRules Rules, parentName string) (error, Errors) {
	validationErrors := make(Errors)

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

func ValidateJson(jsonData string, validationRules Rules) (error, Errors) {
	var decodedJson map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &decodedJson)
	if err != nil {
		return err, nil
	}

	return ValidateMap(decodedJson, validationRules, "")
}

func ValidateMap(mapData map[string]interface{}, validationRules Rules, parentName string) (error, Errors) {
	validationErrors := make(Errors)

	err := validateMapFields(mapData, parentName, validationRules, validationErrors)
	if err != nil {
		return err, nil
	}

	fieldsExists := make(map[string]bool)
	registerMapFields(mapData, parentName, fieldsExists)
	validateNonExistRequiredFields(validationRules, fieldsExists, validationErrors)

	return nil, validationErrors
}
