package validation

import (
	"reflect"
)

var rules = make(map[string]func(fieldName string, fieldValue interface{}, ruleValue string) (error, string))

func CustomRule(ruleName string, ruleFunc func(fieldName string, fieldValue interface{}, ruleValue string) (error, string)) {
	rules[ruleName] = ruleFunc
}

func requiredRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if IsZero(fieldValue) {
		validationError := fieldName + " is required"
		return nil, validationError
	}
	return nil, ""
}

func kindRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	if k := reflect.TypeOf(fieldValue).Kind(); toString(k) != ruleValue {
		validationError := fieldName + " must be kind of " + ruleValue
		return nil, validationError
	}
	return nil, ""
}

func typeRule(fieldName string, fieldValue interface{}, ruleValue string) (error, string) {
	var typeInString string
	if t := reflect.TypeOf(fieldValue); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	if typeInString != ruleValue {
		validationError := fieldName + " must be type of " + ruleValue
		return nil, validationError
	}
	return nil, ""
}

func init() {
	CustomRule("required", requiredRule)
	CustomRule("type", typeRule)
	CustomRule("kind", kindRule)
}
