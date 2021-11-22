package validation

import (
	"errors"
	"reflect"
	"strings"
)

type RuleFunc func(fieldName string, fieldValue interface{}, ruleValue string) error

type rule struct {
	fn     RuleFunc
	errMsg string
}

var registeredRules = make(map[string]*rule)

// AddRule registers a new rule.
// It panics if the rule is already registered.
func AddRule(name string, fn RuleFunc, errMsg string) {
	_, ruleExist := registeredRules[name]
	if ruleExist {
		panic("rule already registered")
	}
	r := &rule{
		fn:     fn,
		errMsg: errMsg,
	}
	registeredRules[name] = r
}

// OverwriteRule registers a new rule.
// If there is a rule already registered with that name it will be overwritten by the new rule.
func OverwriteRule(name string, fn RuleFunc, errMsg string) {
	r := &rule{
		fn:     fn,
		errMsg: errMsg,
	}
	registeredRules[name] = r
}

// SetErrMsg sets errMsg to ruleName.
// It panics if rule does not exist.
func SetErrMsg(ruleName string, errMsg string) {
	r, ok := registeredRules[ruleName]
	if !ok {
		panic("cannot set error message to rule does not exist: " + ruleName)
	}
	r.errMsg = errMsg
}

func getErrMsg(ruleName string, ruleVal string, name string, val interface{}) string {
	errMsg := registeredRules[ruleName].errMsg
	errMsg = strings.ReplaceAll(errMsg, "[name]", name)
	errMsg = strings.ReplaceAll(errMsg, "[val]", toString(val))
	errMsg = strings.ReplaceAll(errMsg, "[ruleVal]", ruleVal)
	return errMsg
}

func getRuleInfo(r string) (string, string, RuleFunc, bool) {
	rName, rValue := splitRuleNameAndRuleValue(r)
	val, rExist := registeredRules[rName]
	var rFunc RuleFunc
	if rExist {
		rFunc = val.fn
	}
	return rName, rValue, rFunc, rExist
}

// requiredRule checks if val is empty.
// It returns error if val IsEmpty().
func requiredRule(name string, val interface{}, ruleVal string) error {
	if IsEmpty(val) {
		return errors.New(getErrMsg("required", ruleVal, name, val))
	}
	return nil
}

// kindRule checks if val's kind equals ruleVal.
// It returns error if val's kind does not equal ruleVal.
func kindRule(name string, val interface{}, ruleVal string) error {
	if k := reflect.TypeOf(val).Kind(); toString(k) != ruleVal {
		return errors.New(getErrMsg("kind", ruleVal, name, val))
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
		return errors.New(getErrMsg("type", ruleVal, name, val))
	}
	return nil
}

// equalRule checks if val equals ruleVal.
// It returns error if val does not equal ruleVal.
func equalRule(name string, val interface{}, ruleVal string) error {
	if toString(val) != ruleVal {
		return errors.New(getErrMsg("equal", ruleVal, name, val))
	}
	return nil
}

// numericRule checks if val is numeric.
// It returns error if val is not numeric.
func numericRule(name string, val interface{}, ruleVal string) error {
	if !IsNumeric(val) {
		return errors.New(getErrMsg("numeric", ruleVal, name, val))
	}
	return nil
}

func init() {
	AddRule("required", requiredRule, "[name] is required")
	AddRule("type", typeRule, "[name] must be type of [ruleVal]")
	AddRule("kind", kindRule, "[name] must be kind of [ruleVal]")
	AddRule("equal", equalRule, "[name] does not equal [ruleVal]")
	AddRule("numeric", numericRule, "[name] must be numeric")
}
