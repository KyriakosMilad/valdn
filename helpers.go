package validation

import (
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

func addValidationTagRules(t reflect.Type, parentName string) {
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
				fieldRules = append(fieldRules, v)
			}
			validationRules[fieldName] = fieldRules
		}

		if fieldType.Kind() == reflect.Struct {
			addValidationTagRules(fieldType, fieldName)
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

func convertInterfaceToMap(value interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	if val, ok := value.(map[string]interface{}); ok {
		for k, v := range val {
			newMap[k] = v
		}
	}
	return newMap
}

func validateByType(fieldName string, fieldType reflect.Type, fieldValue interface{}) error {
	var err error
	fieldRules := validationRules[fieldName]

	switch fieldType.Kind() {
	case reflect.Struct:
		err = validateStruct(fieldValue, fieldName)
	case reflect.Map:
		err = validateMap(fieldValue, fieldName)
	default:
		var fieldValidationError string
		err, fieldValidationError = ValidateField(fieldName, fieldValue, fieldRules)
		if fieldValidationError != "" {
			validationErrors[fieldName] = fieldValidationError
		}
	}
	return err
}

func validateStructFields(t reflect.Type, v reflect.Value, parentName string) error {
	parentName = makeParentNameJoinable(parentName)
	for i := 0; i < t.NumField(); i++ {
		fieldName, fieldType, fieldValue := getStructFieldInfo(i, t, v, parentName)
		fieldsExists[fieldName] = true
		err := validateByType(fieldName, fieldType, fieldValue.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func validateMapFields(mapData map[string]interface{}, parentName string) error {
	parentName = makeParentNameJoinable(parentName)
	for fieldName, fieldValue := range mapData {
		fieldName = parentName + fieldName
		fieldsExists[fieldName] = true
		fieldType := reflect.TypeOf(fieldValue)
		err := validateByType(fieldName, fieldType, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateNonExistRequiredFields() {
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

func copyRules(rules Rules) Rules {
	newMap := make(Rules)
	for k, v := range rules {
		newMap[k] = v
	}
	return newMap
}
