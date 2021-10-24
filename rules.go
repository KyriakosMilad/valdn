package validation

import (
	"reflect"
)

var rules = make(map[string]func(field string, fieldValue interface{}, ruleValue string) (error, string))

func AddRule(ruleName string, ruleFunc func(field string, fieldValue interface{}, ruleValue string) (error, string)) {
	_, ruleExists := rules[ruleName]
	if ruleExists {
		panic("rule already registered")
	}

	rules[ruleName] = ruleFunc

	return
}

func init() {
	AddRule("required", func(field string, fieldValue interface{}, ruleValue string) (error, string) {
		if reflect.ValueOf(fieldValue).IsZero() {
			validationError := field + " is required"
			return nil, validationError
		}
		return nil, ""
	})

	AddRule("string", func(field string, fieldValue interface{}, ruleValue string) (error, string) {
		if reflect.ValueOf(fieldValue).Kind() != reflect.String {
			validationError := field + " must be a string"
			return nil, validationError
		}
		return nil, ""
	})
}
