package validation

import (
	"reflect"
)

type RuleFunc func(fieldName string, fieldValue interface{}, ruleValue string) (error, string)

var registeredRules = make(map[string]RuleFunc)

func AddRule(name string, f RuleFunc) {
	_, ruleExist := registeredRules[name]
	if ruleExist {
		panic("rule already registered")
	}
	registeredRules[name] = f
}

func OverwriteRule(name string, f RuleFunc) {
	registeredRules[name] = f
}

func getRuleInfo(r string) (string, string, RuleFunc, bool) {
	rName, rValue := splitRuleNameAndRuleValue(r)
	rFunc, rExist := registeredRules[rName]
	return rName, rValue, rFunc, rExist
}

func requiredRule(name string, val interface{}, ruleVal string) (error, string) {
	if IsEmpty(val) {
		validationError := name + " is required"
		return nil, validationError
	}
	return nil, ""
}

func kindRule(name string, val interface{}, ruleVal string) (error, string) {
	if k := reflect.TypeOf(val).Kind(); toString(k) != ruleVal {
		validationError := name + " must be kind of " + ruleVal
		return nil, validationError
	}
	return nil, ""
}

func typeRule(name string, val interface{}, ruleVal string) (error, string) {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	if typeInString != ruleVal {
		validationError := name + " must be type of " + ruleVal
		return nil, validationError
	}
	return nil, ""
}

func init() {
	AddRule("required", requiredRule)
	AddRule("type", typeRule)
	AddRule("kind", kindRule)
}
