package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
		panic("rule is already registered")
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

// kindInRule checks if val's kind is one of ruleVal[].
// It returns error if val's kind is not one of ruleVal[].
func kindInRule(name string, val interface{}, ruleVal string) error {
	k := toString(reflect.TypeOf(val).Kind())
	in := false
	for _, v := range strings.Split(ruleVal, ",") {
		if v == k {
			in = true
			break
		}
	}
	if !in {
		return errors.New(getErrMsg("kindIn", ruleVal, name, val))
	}
	return nil
}

// kindNotInRule checks if val's kind is not one of ruleVal[].
// It returns error if val's kind is one of ruleVal[].
func kindNotInRule(name string, val interface{}, ruleVal string) error {
	k := toString(reflect.TypeOf(val).Kind())
	fmt.Println(k)
	in := false
	for _, v := range strings.Split(ruleVal, ",") {
		if v == k {
			in = true
			break
		}
	}
	if in {
		return errors.New(getErrMsg("kindNotIn", ruleVal, name, val))
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

// typeInRule checks if val's type is one of ruleVal[].
// It returns error if val's type is not one of ruleVal[].
func typeInRule(name string, val interface{}, ruleVal string) error {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	in := false
	for _, v := range strings.Split(ruleVal, ",") {
		if v == typeInString {
			in = true
			break
		}
	}
	if !in {
		return errors.New(getErrMsg("typeIn", ruleVal, name, val))
	}
	return nil
}

// typeNotInRule checks if val's type is not one of ruleVal[].
// It returns error if val's type is one of ruleVal[].
func typeNotInRule(name string, val interface{}, ruleVal string) error {
	var typeInString string
	if t := reflect.TypeOf(val); t.Kind() == reflect.Struct {
		typeInString = t.Name()
	} else {
		typeInString = toString(t)
	}
	in := false
	for _, v := range strings.Split(ruleVal, ",") {
		if v == typeInString {
			in = true
			break
		}
	}
	if in {
		return errors.New(getErrMsg("typeNotIn", ruleVal, name, val))
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

// intRule checks if val is integer.
// It returns error if val is not an integer.
func intRule(name string, val interface{}, ruleVal string) error {
	if !IsInteger(val) {
		return errors.New(getErrMsg("int", ruleVal, name, val))
	}
	return nil
}

// uintRule checks if val is unsigned integer.
// It returns error if val is not an unsigned integer.
func uintRule(name string, val interface{}, ruleVal string) error {
	if !IsUnsignedInteger(val) {
		return errors.New(getErrMsg("uint", ruleVal, name, val))
	}
	return nil
}

// complexRule checks if val is complex number.
// It returns error if val is not a complex number.
func complexRule(name string, val interface{}, ruleVal string) error {
	if !IsComplex(val) {
		return errors.New(getErrMsg("complex", ruleVal, name, val))
	}
	return nil
}

// floatRule checks if val is float.
// It returns error if val is not a float.
func floatRule(name string, val interface{}, ruleVal string) error {
	if !IsFloat(val) {
		return errors.New(getErrMsg("float", ruleVal, name, val))
	}
	return nil
}

// numericRule checks if val is numeric.
// It returns error if val is not a numeric.
func numericRule(name string, val interface{}, ruleVal string) error {
	if !IsNumeric(val) {
		return errors.New(getErrMsg("numeric", ruleVal, name, val))
	}
	return nil
}

// betweenRule checks if val is between min (ruleVal[0]) and max (ruleVal[1]).
// It panics if val is not an integer or a float.
// It panics if min or max is not set.
// It panics if min is not an integer or a float.
// It panics if max is not an integer or a float.
// It returns error if val is not between min and max.
func betweenRule(name string, val interface{}, ruleVal string) error {
	err, vFloat := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by betweenRule")
	}

	ruleValSpliced := strings.Split(ruleVal, ",")
	if len(ruleValSpliced) != 2 {
		panic(fmt.Errorf("betweenRule expects two numeric values as min and max, got: %v", len(ruleValSpliced)))
	}
	err, min := stringToFloat(ruleValSpliced[0])
	if err != nil {
		panic(fmt.Errorf("betweenRule: min must be an integer or a float, got: %v", ruleValSpliced[0]))
	}
	err, max := stringToFloat(ruleValSpliced[1])
	if err != nil {
		panic(fmt.Errorf("betweenRule: max must be an integer or a float, got: %v", ruleValSpliced[1]))
	}

	if vFloat < min || vFloat > max {
		return errors.New(getErrMsg("between", ruleVal, name, val))
	}
	return nil
}

// minRule checks if val is lower than ruleVal.
// It panics if val is not an integer or a float.
// It panics if ruleVal is empty.
// It panics if ruleVal is not an integer or a float.
// It returns error if val is lower than ruleVal.
func minRule(name string, val interface{}, ruleVal string) error {
	err, vFloat := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by minRule")
	}
	err, min := stringToFloat(ruleVal)
	if err != nil {
		panic(fmt.Errorf("minRule: min must be an integer or a float, got: %v", ruleVal))
	}

	if vFloat < min {
		return errors.New(getErrMsg("min", ruleVal, name, val))
	}
	return nil
}

// maxRule checks if val is greater than ruleVal.
// It panics if val is not an integer or a float.
// It panics if ruleVal is empty.
// It panics if ruleVal is not an integer or a float.
// It returns error if val is greater than ruleVal.
func maxRule(name string, val interface{}, ruleVal string) error {
	err, vFloat := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by minRule")
	}
	err, max := stringToFloat(ruleVal)
	if err != nil {
		panic(fmt.Errorf("maxRule: max must be an integer or a float, got: %v", ruleVal))
	}

	if vFloat > max {
		return errors.New(getErrMsg("max", ruleVal, name, val))
	}
	return nil
}

// inRule checks if val equals one of ruleVal[] items.
// It returns error if val doesn't equal any item in ruleVal[].
func inRule(name string, val interface{}, ruleVal string) error {
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		if v == toString(val) {
			in = true
			break
		}
	}
	if !in {
		return errors.New(getErrMsg("in", ruleVal, name, val))
	}
	return nil
}

// notInRule checks if val doesn't equal any item in ruleVal[].
// It returns error if val equals one of ruleVal[] items.
func notInRule(name string, val interface{}, ruleVal string) error {
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		if v == toString(val) {
			in = true
			break
		}
	}
	if in {
		return errors.New(getErrMsg("notIn", ruleVal, name, val))
	}
	return nil
}

// lenRule checks if val's length equals ruleVal.
// It panics if val is not array, slice, map, string, integer or float.
// It returns error if val's length doesn't equal ruleVal.
func lenRule(name string, val interface{}, ruleVal string) error {
	l, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("length must be an integer")
	}
	err, vLen := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen != int(l) {
		return errors.New(getErrMsg("len", ruleVal, name, val))
	}
	return nil
}

// minLenRule checks if val's length is greater than or equal ruleVal or not.
// It panics if val is not array, slice, map, string, integer or float.
// It returns error if val's length is lower than ruleVal.
func minLenRule(name string, val interface{}, ruleVal string) error {
	l, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("length must be an integer")
	}
	err, vLen := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen < int(l) {
		return errors.New(getErrMsg("minLen", ruleVal, name, val))
	}
	return nil
}

// maxLenRule checks if val's length is lower than or equal ruleVal or not.
// It panics if val is not array, slice, map, string, integer or float.
// It returns error if val's length is greater than ruleVal.
func maxLenRule(name string, val interface{}, ruleVal string) error {
	l, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("length must be an integer")
	}
	err, vLen := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen > int(l) {
		return errors.New(getErrMsg("maxLen", ruleVal, name, val))
	}
	return nil
}

// lenBetweenRule checks if val's length is between ruleVal[0] and ruleVal[1] or not.
// It panics if val is not array, slice, map, string, integer or float.
// It panics if min or max is not set.
// It panics if min is not an integer.
// It panics if max is not an integer.
// It returns error if val's length is not between ruleVal[0] and ruleVal[1].
func lenBetweenRule(name string, val interface{}, ruleVal string) error {
	err, l := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	if len(ruleValSpliced) != 2 {
		panic(fmt.Errorf("lenBetweenRule expects two integer values as min and max, got: %v", len(ruleValSpliced)))
	}
	min, err := strconv.ParseInt(ruleValSpliced[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("lenBetweenRule: min must be an integer, got: %v", ruleValSpliced[0]))
	}
	max, err := strconv.ParseInt(ruleValSpliced[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("lenBetweenRule: max must be an integer, got: %v", ruleValSpliced[1]))
	}
	if l < int(min) || l > int(max) {
		return errors.New(getErrMsg("lenBetween", ruleVal, name, val))
	}
	return nil
}

func init() {
	AddRule("required", requiredRule, "[name] is required")
	AddRule("type", typeRule, "[name] must be type of [ruleVal]")
	AddRule("typeIn", typeInRule, "[name]'s type must be one of [ruleVal]")
	AddRule("typeNotIn", typeNotInRule, "[name]'s type must not be one of [ruleVal]")
	AddRule("kind", kindRule, "[name] must be kind of [ruleVal]")
	AddRule("kindIn", kindInRule, "[name]'s kind must be one of [ruleVal]")
	AddRule("kindNotIn", kindNotInRule, "[name] must not be kind of [ruleVal]")
	AddRule("equal", equalRule, "[name] does not equal [ruleVal]")
	AddRule("int", intRule, "[name] must be an integer")
	AddRule("uint", uintRule, "[name] must be an unsigned integer")
	AddRule("complex", complexRule, "[name] must be a complex number")
	AddRule("float", floatRule, "[name] must be a float")
	AddRule("numeric", numericRule, "[name] must be a numeric")
	AddRule("between", betweenRule, "[name] must be between [ruleVal]")
	AddRule("min", minRule, "[name] must be greater than or equal [ruleVal]")
	AddRule("max", maxRule, "[name] must be lower than or equal [ruleVal]")
	AddRule("in", inRule, "[name] must be in these values: [ruleVal]")
	AddRule("notIn", notInRule, "[name] must not be in these values: [ruleVal]")
	AddRule("len", lenRule, "[name]'s length must equal: [ruleVal]")
	AddRule("minLen", minLenRule, "[name]'s length must be greater than or equal: [ruleVal]")
	AddRule("maxLen", maxLenRule, "[name]'s length must be lower than or equal: [ruleVal]")
	AddRule("lenBetween", lenBetweenRule, "[name]'s length must be between: [ruleVal]")
}
