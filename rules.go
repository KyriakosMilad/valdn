package validation

import (
	"errors"
	"reflect"
)

type RuleFunc func(fieldName string, fieldValue interface{}, ruleValue string) error

var registeredRules = make(map[string]RuleFunc)

// AddRule registers a new rule.
// It panics if the rule is already registered.
func AddRule(name string, f RuleFunc) {
	_, ruleExist := registeredRules[name]
	if ruleExist {
		panic("rule already registered")
	}
	registeredRules[name] = f
}

// OverwriteRule registers a new rule.
// If there is a rule already registered with that name it will be overwritten by the new rule.
func OverwriteRule(name string, f RuleFunc) {
	registeredRules[name] = f
}

func getRuleInfo(r string) (string, string, RuleFunc, bool) {
	rName, rValue := splitRuleNameAndRuleValue(r)
	rFunc, rExist := registeredRules[rName]
	return rName, rValue, rFunc, rExist
}

// requiredRule checks if val is empty.
// It returns error if val IsEmpty().
func requiredRule(name string, val interface{}, ruleVal string) error {
	if IsEmpty(val) {
		return errors.New(name + " is required")
	}
	return nil
}

// kindRule checks if val's kind equals ruleVal
// It returns error if val's kind does not equal ruleVal.
func kindRule(name string, val interface{}, ruleVal string) error {
	if k := reflect.TypeOf(val).Kind(); toString(k) != ruleVal {
		return errors.New(name + " must be kind of " + ruleVal)
	}
	return nil
}

// typeRule checks if val's type equals ruleVal
// It returns error if val's type does not equal ruleVal.
func typeRule(name string, val interface{}, ruleVal string) error {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	if typeInString != ruleVal {
		return errors.New(name + " must be type of " + ruleVal)
	}
	return nil
}

func init() {
	AddRule("required", requiredRule)
	AddRule("type", typeRule)
	AddRule("kind", kindRule)
}
