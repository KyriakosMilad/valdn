package validation

import (
	"reflect"
)

var rules = make(map[string]func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string))

func AddRule(ruleName string, ruleFunc func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string)) {
	_, ruleExists := rules[ruleName]
	if ruleExists {
		panic("rule already registered")
	}

	rules[ruleName] = ruleFunc

	return
}

func init() {
	AddRule("required", func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string) {
		if reflect.ValueOf(fieldValue).IsZero() {
			validationError = field + " is required"
			return
		}
		return
	})

	AddRule("string", func(field string, fieldValue interface{}, ruleValue string) (err error, validationError string) {
		if reflect.ValueOf(fieldValue).Kind() != reflect.String {
			validationError = field + " must be a string"
			return
		}
		return
	})
}
