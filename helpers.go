package validation

import (
	"fmt"
	"reflect"
	"strings"
)

func splitRuleNameAndRuleValue(rule string) (string, string) {
	if strings.ContainsRune(rule, ':') {
		ruleSpliced := strings.Split(rule, ":")
		return ruleSpliced[0], ruleSpliced[1]
	}
	return rule, ""
}

func getRuleInfo(rule string) (string, string, RuleFunc, bool) {
	ruleName, ruleValue := splitRuleNameAndRuleValue(rule)
	ruleFunc, ruleExists := rules[ruleName]
	return ruleName, ruleValue, ruleFunc, ruleExists
}

func makeParentNameJoinable(parentName string) string {
	if parentName != "" && parentName[len(parentName)-1] != '.' {
		return parentName + "."
	}
	return parentName
}

func getParentName(name string) string {
	nameSpliced := strings.Split(name, ".")
	if len(nameSpliced) > 1 {
		return strings.Join(nameSpliced[:len(nameSpliced)-1], ".")
	}
	return ""
}

func isRuleExists(rules []string, rule string) bool {
	for _, v := range rules {
		if v == rule {
			return true
		}
	}
	return false
}

func isParentRequired(fieldName string, validationRules Rules) bool {
	parent := getParentName(fieldName)
	if parent == "" {
		return true
	}
	parentRules := validationRules[parent]
	return isRuleExists(parentRules, "required")
}

func addValidationTagRules(t reflect.Type, validationRules Rules, parentName string) {
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldType := field.Type
		fieldName := parentName + field.Name
		fieldValidationTag := field.Tag.Get("validation")

		_, fieldHasRules := validationRules[fieldName]
		if !fieldHasRules && fieldValidationTag != "" {
			var fieldRules []string
			for _, v := range strings.Split(fieldValidationTag, "|") {
				if (v == "required" && isParentRequired(fieldName, validationRules)) || v != "required" {
					fieldRules = append(fieldRules, v)
				}
			}
			validationRules[fieldName] = fieldRules
		}

		if fieldType.Kind() == reflect.Struct {
			addValidationTagRules(fieldType, validationRules, fieldName)
		}
	}
}

func getStructFieldInfo(fieldNumber int, parentType reflect.Type, parentValue reflect.Value, parentName string) (string, reflect.Type, reflect.Value) {
	field := parentType.Field(fieldNumber)
	fieldName := parentName + field.Name
	fieldType := field.Type
	fieldValue := parentValue.Field(fieldNumber)

	return fieldName, fieldType, fieldValue
}

func addValidationErrors(validationErrors Errors, newValidationErrors Errors) {
	for k, v := range newValidationErrors {
		validationErrors[k] = v
	}
}

func getNestedRules(validationRules Rules, structName string) Rules {
	structRules := make(Rules)
	for k, v := range validationRules {
		if strings.Contains(k, makeParentNameJoinable(structName)) {
			structRules[k] = v
		}
	}
	return structRules
}

func validateNestedStruct(fieldName string, fieldValue interface{}, validationRules Rules) (error, Errors) {
	fieldValidationErrors := make(Errors)
	fieldRules := validationRules[fieldName]

	err, structFieldError := ValidateField(fieldName, fieldValue, fieldRules)
	if err != nil {
		return err, nil
	}
	if structFieldError != "" {
		fieldValidationErrors[fieldName] = structFieldError
		return nil, fieldValidationErrors
	}

	structRules := getNestedRules(validationRules, fieldName)

	return ValidateStruct(fieldValue, structRules, fieldName)
}

func convertInterfaceToMap(value interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	if val, ok := value.(map[string]interface{}); ok {
		for k, v := range val {
			newMap[k] = v
		}
	}
	return newMap
}

func validateNestedMap(fieldName string, fieldValue interface{}, validationRules Rules) (error, Errors) {
	fieldValidationErrors := make(Errors)
	fieldRules := validationRules[fieldName]

	err, mapFieldError := ValidateField(fieldName, fieldValue, fieldRules)
	if err != nil {
		return err, nil
	}
	if mapFieldError != "" {
		fieldValidationErrors[fieldName] = mapFieldError
		return nil, fieldValidationErrors
	}

	mapRules := getNestedRules(validationRules, fieldName)
	mapData := convertInterfaceToMap(fieldValue)

	return ValidateMap(mapData, mapRules, fieldName)
}

func validateByType(fieldName string, fieldType reflect.Type, fieldValue interface{}, validationRules Rules) (error, Errors) {
	var err error
	fieldValidationErrors := make(Errors)

	fieldRules, fieldHasRules := validationRules[fieldName]
	if !fieldHasRules && fieldType.Kind() != reflect.Struct && fieldType != reflect.TypeOf(map[string]interface{}{}) {
		return nil, nil
	}

	switch {
	case fieldType.Kind() == reflect.Struct:
		err, fieldValidationErrors = validateNestedStruct(fieldName, fieldValue, validationRules)
	case fieldType == reflect.TypeOf(map[string]interface{}{}):
		err, fieldValidationErrors = validateNestedMap(fieldName, fieldValue, validationRules)
	default:
		if fieldType.Kind() == reflect.Map {
			return fmt.Errorf("error validating %v: type %v is not supported", fieldName, fieldType), nil
		}
		var fieldValidationError string
		err, fieldValidationError = ValidateField(fieldName, fieldValue, fieldRules)
		if fieldValidationError != "" {
			fieldValidationErrors[fieldName] = fieldValidationError
		}
	}

	return err, fieldValidationErrors
}

func validateStructFields(t reflect.Type, v reflect.Value, parentName string, validationRules Rules, validationErrors Errors) error {
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < t.NumField(); i++ {
		fieldName, fieldType, fieldValue := getStructFieldInfo(i, t, v, parentName)
		err, fieldValidationErrors := validateByType(fieldName, fieldType, fieldValue.Interface(), validationRules)
		if err != nil {
			return err
		}
		addValidationErrors(validationErrors, fieldValidationErrors)
	}
	return nil
}

func validateMapFields(mapData map[string]interface{}, parentName string, validationRules Rules, validationErrors Errors) error {
	parentName = makeParentNameJoinable(parentName)
	for fieldName, fieldValue := range mapData {
		fieldName = parentName + fieldName
		fieldType := reflect.TypeOf(fieldValue)
		err, fieldValidationErrors := validateByType(fieldName, fieldType, fieldValue, validationRules)
		if err != nil {
			return err
		}
		addValidationErrors(validationErrors, fieldValidationErrors)
	}
	return nil
}

func registerStructFields(structData interface{}, parentName string, fieldsExists map[string]bool) {
	t := reflect.TypeOf(structData)
	v := reflect.ValueOf(structData)
	if t.Kind() != reflect.Struct {
		return
	}
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < t.NumField(); i++ {
		fieldName, fieldType, fieldValue := getStructFieldInfo(i, t, v, parentName)
		fieldsExists[fieldName] = true
		registerNestedFieldsByType(fieldType, fieldValue.Interface(), fieldName, fieldsExists)
	}
}

func registerMapFields(mapData interface{}, parentName string, fieldsExists map[string]bool) {
	if reflect.TypeOf(map[string]interface{}{}) != reflect.TypeOf(mapData) {
		return
	}
	parentName = makeParentNameJoinable(parentName)
	mapFields := convertInterfaceToMap(mapData)
	for k, v := range mapFields {
		fieldName := parentName + k
		fieldType := reflect.TypeOf(v)
		fieldsExists[fieldName] = true
		registerNestedFieldsByType(fieldType, v, fieldName, fieldsExists)
	}
}

func registerNestedFieldsByType(fieldType reflect.Type, fieldValue interface{}, fieldName string, fieldsExists map[string]bool) {
	switch {
	case fieldType.Kind() == reflect.Struct:
		registerStructFields(fieldValue, fieldName, fieldsExists)
	case fieldType == reflect.TypeOf(map[string]interface{}{}):
		registerMapFields(fieldValue, fieldName, fieldsExists)
	}
}

func validateNonExistRequiredFields(validationRules Rules, fieldsExists map[string]bool, validationErrors Errors) {
	for k, v := range validationRules {
		for _, val := range v {
			if val == "required" {
				_, ok := fieldsExists[k]
				if !ok {
					validationErrors[k] = k + " is required"
				}
			}
		}
	}
}
