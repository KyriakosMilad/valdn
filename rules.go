package valdn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func GetErrMsg(ruleName string, ruleVal string, name string, val interface{}) string {
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

// requiredRule checks if val exists, and it's not empty.
// It returns error if val is not exist or empty.
func requiredRule(name string, val interface{}, ruleVal string) error {
	if IsEmpty(val) {
		return errors.New(GetErrMsg("required", ruleVal, name, val))
	}
	return nil
}

// kindRule checks if val's kind equals ruleVal.
// It returns error if val's kind does not equal ruleVal.
func kindRule(name string, val interface{}, ruleVal string) error {
	if !IsKind(val, ruleVal) {
		return errors.New(GetErrMsg("kind", ruleVal, name, val))
	}
	return nil
}

// kindInRule checks if val's kind is one of ruleVal[].
// It returns error if val's kind is not one of ruleVal[].
func kindInRule(name string, val interface{}, ruleVal string) error {
	if !IsKindIn(val, strings.Split(ruleVal, ",")) {
		return errors.New(GetErrMsg("kindIn", ruleVal, name, val))
	}
	return nil
}

// kindNotInRule checks if val's kind is not one of ruleVal[].
// It returns error if val's kind is one of ruleVal[].
func kindNotInRule(name string, val interface{}, ruleVal string) error {
	if IsKindIn(val, strings.Split(ruleVal, ",")) {
		return errors.New(GetErrMsg("kindNotIn", ruleVal, name, val))
	}
	return nil
}

// typeRule checks if val's type equals ruleVal.
// It returns error if val's type does not equal ruleVal.
func typeRule(name string, val interface{}, ruleVal string) error {
	if !IsType(val, ruleVal) {
		return errors.New(GetErrMsg("type", ruleVal, name, val))
	}
	return nil
}

// typeInRule checks if val's type is one of ruleVal[].
// It returns error if val's type is not one of ruleVal[].
func typeInRule(name string, val interface{}, ruleVal string) error {
	if !IsTypeIn(val, strings.Split(ruleVal, ",")) {
		return errors.New(GetErrMsg("typeIn", ruleVal, name, val))
	}
	return nil
}

// typeNotInRule checks if val's type is not one of ruleVal[].
// It returns error if val's type is one of ruleVal[].
func typeNotInRule(name string, val interface{}, ruleVal string) error {
	if IsTypeIn(val, strings.Split(ruleVal, ",")) {
		return errors.New(GetErrMsg("typeNotIn", ruleVal, name, val))
	}
	return nil
}

// equalRule checks if val equals ruleVal.
// It returns error if val does not equal ruleVal.
func equalRule(name string, val interface{}, ruleVal string) error {
	if toString(val) != ruleVal {
		return errors.New(GetErrMsg("equal", ruleVal, name, val))
	}
	return nil
}

// intRule checks if val is integer.
// It returns error if val is not an integer.
func intRule(name string, val interface{}, ruleVal string) error {
	if !IsInteger(val) {
		return errors.New(GetErrMsg("int", ruleVal, name, val))
	}
	return nil
}

// uintRule checks if val is unsigned integer.
// It returns error if val is not an unsigned integer.
func uintRule(name string, val interface{}, ruleVal string) error {
	if !IsUnsignedInteger(val) {
		return errors.New(GetErrMsg("uint", ruleVal, name, val))
	}
	return nil
}

// complexRule checks if val is complex number.
// It returns error if val is not a complex number.
func complexRule(name string, val interface{}, ruleVal string) error {
	if !IsComplex(val) {
		return errors.New(GetErrMsg("complex", ruleVal, name, val))
	}
	return nil
}

// floatRule checks if val is float.
// It returns error if val is not a float.
func floatRule(name string, val interface{}, ruleVal string) error {
	if !IsFloat(val) {
		return errors.New(GetErrMsg("float", ruleVal, name, val))
	}
	return nil
}

// ufloatRule checks if val is unsigned float.
// It returns error if val is not an unsigned float.
func ufloatRule(name string, val interface{}, ruleVal string) error {
	if !IsUnsignedFloat(val) {
		return errors.New(GetErrMsg("ufloat", ruleVal, name, val))
	}
	return nil
}

// numericRule checks if val is numeric.
// It returns error if val is not a numeric.
func numericRule(name string, val interface{}, ruleVal string) error {
	if !IsNumeric(val) {
		return errors.New(GetErrMsg("numeric", ruleVal, name, val))
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
	vFloat, err := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by betweenRule")
	}

	ruleValSpliced := strings.Split(ruleVal, ",")
	if len(ruleValSpliced) != 2 {
		panic(fmt.Errorf("betweenRule expects two numeric values as min and max, got: %v", len(ruleValSpliced)))
	}
	min, err := stringToFloat(ruleValSpliced[0])
	if err != nil {
		panic(fmt.Errorf("betweenRule: min must be an integer or a float, got: %v", ruleValSpliced[0]))
	}
	max, err := stringToFloat(ruleValSpliced[1])
	if err != nil {
		panic(fmt.Errorf("betweenRule: max must be an integer or a float, got: %v", ruleValSpliced[1]))
	}

	if vFloat < min || vFloat > max {
		return errors.New(GetErrMsg("between", ruleVal, name, val))
	}
	return nil
}

// minRule checks if val is lower than ruleVal.
// It panics if val is not an integer or a float.
// It panics if ruleVal is empty.
// It panics if ruleVal is not an integer or a float.
// It returns error if val is lower than ruleVal.
func minRule(name string, val interface{}, ruleVal string) error {
	vFloat, err := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by minRule")
	}
	min, err := stringToFloat(ruleVal)
	if err != nil {
		panic(fmt.Errorf("minRule: min must be an integer or a float, got: %v", ruleVal))
	}

	if vFloat < min {
		return errors.New(GetErrMsg("min", ruleVal, name, val))
	}
	return nil
}

// maxRule checks if val is greater than ruleVal.
// It panics if val is not an integer or a float.
// It panics if ruleVal is empty.
// It panics if ruleVal is not an integer or a float.
// It returns error if val is greater than ruleVal.
func maxRule(name string, val interface{}, ruleVal string) error {
	vFloat, err := interfaceToFloat(val)
	if err != nil {
		panic(name + " must be an integer or a float to be validated by minRule")
	}
	max, err := stringToFloat(ruleVal)
	if err != nil {
		panic(fmt.Errorf("maxRule: max must be an integer or a float, got: %v", ruleVal))
	}

	if vFloat > max {
		return errors.New(GetErrMsg("max", ruleVal, name, val))
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
		return errors.New(GetErrMsg("in", ruleVal, name, val))
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
		return errors.New(GetErrMsg("notIn", ruleVal, name, val))
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
	vLen, err := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen != int(l) {
		return errors.New(GetErrMsg("len", ruleVal, name, val))
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
	vLen, err := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen < int(l) {
		return errors.New(GetErrMsg("minLen", ruleVal, name, val))
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
	vLen, err := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	if vLen > int(l) {
		return errors.New(GetErrMsg("maxLen", ruleVal, name, val))
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
	l, err := getLen(val)
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
		return errors.New(GetErrMsg("lenBetween", ruleVal, name, val))
	}
	return nil
}

// lenInRule checks if val's length equals one of ruleVal[] items.
// It panics if val is not array, slice, map, string, integer or float.
// It panics if one of ruleVal items is not an integer.
// It returns error if val's length doesn't equal any item in ruleVal[].
func lenInRule(name string, val interface{}, ruleVal string) error {
	vLen, err := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		l, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("length must be an integer")
		}
		if vLen == int(l) {
			in = true
			break
		}
	}
	if !in {
		return errors.New(GetErrMsg("lenIn", ruleVal, name, val))
	}
	return nil
}

// lenNotInRule checks if val's length doesn't equal any item in ruleVal[].
// It panics if val is not array, slice, map, string, integer or float.
// It panics if one of ruleVal items is not an integer.
// It returns error if val's length equals any item in ruleVal[].
func lenNotInRule(name string, val interface{}, ruleVal string) error {
	vLen, err := getLen(val)
	if err != nil {
		panic(err.Error())
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		l, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("length must be an integer")
		}
		if vLen == int(l) {
			in = true
			break
		}
	}
	if in {
		return errors.New(GetErrMsg("lenNotIn", ruleVal, name, val))
	}
	return nil
}

// regexRule checks if val matches ruleVal regular expression.
// It panics if val is not a string.
// It panics if ruleVal is not a valid regular expression.
// It returns error if val doesn't match ruleVal regular expression.
func regexRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with regexRule", name))
	}
	r, err := regexp.Compile(ruleVal)
	if err != nil {
		panic(fmt.Errorf("%v is not a valid regex", ruleVal))
	}
	match := r.MatchString(toString(val))
	if !match {
		return errors.New(GetErrMsg("regex", ruleVal, name, val))
	}
	return nil
}

// notRegexRule checks if val doesn't match ruleVal regular expression.
// It panics if val is not a string.
// It panics if ruleVal is not a valid regular expression.
// It returns error if val matches ruleVal regular expression.
func notRegexRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with notRegexRule", name))
	}
	r, err := regexp.Compile(ruleVal)
	if err != nil {
		panic(fmt.Errorf("%v is not a valid regex", ruleVal))
	}
	match := r.MatchString(toString(val))
	if match {
		return errors.New(GetErrMsg("notRegex", ruleVal, name, val))
	}
	return nil
}

// emailRule checks if val is a valid email address.
// It panics if val is not a string.
// It returns error if val is not a valid email address.
func emailRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with emailRule", name))
	}
	ok := IsEmail(toString(val))
	if !ok {
		return errors.New(GetErrMsg("email", ruleVal, name, val))
	}
	return nil
}

// jsonRule checks if val is a valid json.
// It panics if val is not a string.
// It returns error if val is not a valid json.
func jsonRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with jsonRule", name))
	}
	ok := IsJSON(toString(val))
	if !ok {
		return errors.New(GetErrMsg("json", ruleVal, name, val))
	}
	return nil
}

// ipv4Rule checks if val is a valid IPv4.
// It panics if val is not a string.
// It returns error if val is not a valid IPv4.
func ipv4Rule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with ipv4Rule", name))
	}
	ok := IsIPv4(toString(val))
	if !ok {
		return errors.New(GetErrMsg("ipv4", ruleVal, name, val))
	}
	return nil
}

// ipv6Rule checks if val is a valid IPv6.
// It panics if val is not a string.
// It returns error if val is not a valid IPv6.
func ipv6Rule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with ipv6Rule", name))
	}
	ok := IsIPv6(toString(val))
	if !ok {
		return errors.New(GetErrMsg("ipv6", ruleVal, name, val))
	}
	return nil
}

// ipRule checks if val is a valid IP address.
// It panics if val is not a string.
// It returns error if val is not a valid IP address.
func ipRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with ipRule", name))
	}
	ok := IsIP(toString(val))
	if !ok {
		return errors.New(GetErrMsg("ip", ruleVal, name, val))
	}
	return nil
}

// macRule checks if val is a valid mac address.
// It panics if val is not a string.
// It returns error if val is not a valid mac address.
func macRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with macRule", name))
	}
	ok := IsMAC(toString(val))
	if !ok {
		return errors.New(GetErrMsg("mac", ruleVal, name, val))
	}
	return nil
}

// urlRule checks if val is a valid URL.
// It panics if val is not a string.
// It returns error if val is not a valid URL.
func urlRule(name string, val interface{}, ruleVal string) error {
	if !IsString(val) {
		panic(fmt.Errorf("%v must be a string to be validated with macRule", name))
	}
	ok := IsURL(toString(val))
	if !ok {
		return errors.New(GetErrMsg("url", ruleVal, name, val))
	}
	return nil
}

// timeRule checks if val is type of time.Time.
// It returns error if val is not type of time.Time.
func timeRule(name string, val interface{}, ruleVal string) error {
	if _, ok := val.(time.Time); !ok {
		return errors.New(GetErrMsg("time", ruleVal, name, val))
	}
	return nil
}

// timeFormatRule checks if val's format matches ruleVal.
// It returns error if val's format doesn't match ruleVal.
func timeFormatRule(name string, val interface{}, ruleVal string) error {
	_, err := time.Parse(ruleVal, toString(val))
	if err != nil {
		return errors.New(GetErrMsg("timeFormat", ruleVal, name, val))
	}
	return nil
}

// timeFormatInRule checks if val's format matches any of ruleVal[].
// Use [] to split between two formats.
// It returns error if val's format doesn't match any of ruleVal[].
func timeFormatInRule(name string, val interface{}, ruleVal string) error {
	stringVal := toString(val)
	ruleValSpliced := strings.Split(ruleVal, "[]")
	in := false
	for _, v := range ruleValSpliced {
		_, err := time.Parse(v, stringVal)
		if err == nil {
			in = true
			break
		}
	}
	if !in {
		return errors.New(GetErrMsg("timeFormatIn", ruleVal, name, val))
	}
	return nil
}

// timeFormatNotInRule checks if val's format doesn't match any of ruleVal[].
// Use [] to split between two formats.
// It returns error if val's format matches any of ruleVal[].
func timeFormatNotInRule(name string, val interface{}, ruleVal string) error {
	stringVal := toString(val)
	ruleValSpliced := strings.Split(ruleVal, "[]")
	in := false
	for _, v := range ruleValSpliced {
		_, err := time.Parse(v, stringVal)
		if err == nil {
			in = true
			break
		}
	}
	if in {
		return errors.New(GetErrMsg("timeFormatNotIn", ruleVal, name, val))
	}
	return nil
}

// fileRule checks if val is a valid file.
// It returns error if val is not a valid file.
func fileRule(name string, val interface{}, ruleVal string) error {
	if !IsFile(val) {
		return errors.New(GetErrMsg("file", ruleVal, name, val))
	}
	return nil
}

// sizeRule checks if val's size equals ruleVal.
// it panics if val is not a valid file.
// it panics if ruleVal is not an integer.
// It returns error if val's size doesn't equal ruleVal.
func sizeRule(name string, val interface{}, ruleVal string) error {
	size, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("size must be an integer")
	}
	fileSize, err := getFileSize(val)
	if err != nil {
		panic(err)
	}
	if size != fileSize {
		return errors.New(GetErrMsg("size", ruleVal, name, val))
	}
	return nil
}

// sizeMinRule checks if val's size greater than or equal ruleVal or not.
// it panics if val is not a valid file.
// it panics if ruleVal is not an integer.
// It returns error if val's size is lower than ruleVal.
func sizeMinRule(name string, val interface{}, ruleVal string) error {
	size, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("size must be an integer")
	}
	fileSize, err := getFileSize(val)
	if err != nil {
		panic(err)
	}
	if fileSize < size {
		return errors.New(GetErrMsg("sizeMin", ruleVal, name, val))
	}
	return nil
}

// sizeMaxRule checks if val's size lower than or equal ruleVal or not.
// it panics if val is not a valid file.
// it panics if ruleVal is not an integer.
// It returns error if val's size is greater than ruleVal.
func sizeMaxRule(name string, val interface{}, ruleVal string) error {
	size, err := strconv.ParseInt(ruleVal, 10, 64)
	if err != nil {
		panic("size must be an integer")
	}
	fileSize, err := getFileSize(val)
	if err != nil {
		panic(err)
	}
	if fileSize > size {
		return errors.New(GetErrMsg("sizeMax", ruleVal, name, val))
	}
	return nil
}

// sizeBetweenRule checks if val's size is between ruleVal[0] and ruleVal[1].
// it panics if val is not a valid file.
// It panics if min or max is not set.
// It panics if min is not an integer.
// It panics if max is not an integer.
// It returns error if val's size is not between ruleVal[0] and ruleVal[1].
func sizeBetweenRule(name string, val interface{}, ruleVal string) error {
	fileSize, err := getFileSize(val)
	if err != nil {
		panic(err)
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	if len(ruleValSpliced) != 2 {
		panic(fmt.Errorf("sizeBetweenRule expects two integer values as min and max, got: %v", len(ruleValSpliced)))
	}
	min, err := strconv.ParseInt(ruleValSpliced[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("sizeBetweenRule: min must be an integer, got: %v", ruleValSpliced[0]))
	}
	max, err := strconv.ParseInt(ruleValSpliced[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("sizeBetweenRule: max must be an integer, got: %v", ruleValSpliced[1]))
	}
	if fileSize < min || fileSize > max {
		return errors.New(GetErrMsg("sizeBetween", ruleVal, name, val))
	}
	return nil
}

// extRule checks if val's extension equals ruleVal.
// it panics if val is not a valid file.
// It returns error if val's extension doesn't equal ruleVal.
func extRule(name string, val interface{}, ruleVal string) error {
	ext, err := getFileExt(val)
	if err != nil {
		panic(err)
	}
	if ruleVal[0] != '.' {
		ruleVal = "." + ruleVal
	}
	if ruleVal != ext {
		return errors.New(GetErrMsg("ext", ruleVal, name, val))
	}
	return nil
}

// notExtRule checks if val's extension does not equal ruleVal.
// it panics if val is not a valid file.
// It returns error if val's extension equals ruleVal.
func notExtRule(name string, val interface{}, ruleVal string) error {
	ext, err := getFileExt(val)
	if err != nil {
		panic(err)
	}
	if ruleVal[0] != '.' {
		ruleVal = "." + ruleVal
	}
	if ruleVal == ext {
		return errors.New(GetErrMsg("notExt", ruleVal, name, val))
	}
	return nil
}

// extInRule checks if val's extension equals one of ruleVal[] items.
// It panics if val is not a valid file.
// It returns error if val's extension doesn't equal any item in ruleVal[].
func extInRule(name string, val interface{}, ruleVal string) error {
	ext, err := getFileExt(val)
	if err != nil {
		panic(err)
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		if v[0] != '.' {
			v = "." + v
		}
		if v == ext {
			in = true
			break
		}
	}
	if !in {
		return errors.New(GetErrMsg("extIn", ruleVal, name, val))
	}
	return nil
}

// extNotInRule checks if val's extension doesn't equal one of ruleVal[] items.
// It panics if val is not a valid file.
// It returns error if val's extension equals any item in ruleVal[].
func extNotInRule(name string, val interface{}, ruleVal string) error {
	ext, err := getFileExt(val)
	if err != nil {
		panic(err)
	}
	ruleValSpliced := strings.Split(ruleVal, ",")
	var in bool
	for _, v := range ruleValSpliced {
		if v[0] != '.' {
			v = "." + v
		}
		if v == ext {
			in = true
			break
		}
	}
	if in {
		return errors.New(GetErrMsg("extNotIn", ruleVal, name, val))
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
	AddRule("ufloat", ufloatRule, "[name] must be an unsigned float")
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
	AddRule("lenIn", lenInRule, "[name]'s length must be in these values: [ruleVal]")
	AddRule("lenNotIn", lenNotInRule, "[name]'s length must not be in these values: [ruleVal]")
	AddRule("regex", regexRule, "[name]'s format is not valid")
	AddRule("notRegex", notRegexRule, "[name]'s format is not valid")
	AddRule("email", emailRule, "[name] must be a valid email address")
	AddRule("json", jsonRule, "[name] must be a valid json")
	AddRule("ipv4", ipv4Rule, "[name] must be a valid ipv4")
	AddRule("ipv6", ipv6Rule, "[name] must be a valid ipv6")
	AddRule("ip", ipRule, "[name] must be a valid ip address")
	AddRule("mac", macRule, "[name] must be a valid mac address")
	AddRule("url", urlRule, "[name] must be a valid url")
	AddRule("time", timeRule, "[name] must be type of time.Time")
	AddRule("timeFormat", timeFormatRule, "[name]'s format must match [ruleVal]")
	AddRule("timeFormatIn", timeFormatInRule, "[name]'s format must match at least one of [ruleVal]")
	AddRule("timeFormatNotIn", timeFormatNotInRule, "[name]'s format must not match any of [ruleVal]")
	AddRule("file", fileRule, "[name] must be a valid file")
	AddRule("size", sizeRule, "[name]'s size doesn't equal [ruleVal]")
	AddRule("sizeMin", sizeMinRule, "[name]'s size must be greater than or equal [ruleVal]")
	AddRule("sizeMax", sizeMaxRule, "[name]'s size must be lower than or equal [ruleVal]")
	AddRule("sizeBetween", sizeBetweenRule, "[name]'s size must be between [ruleVal]")
	AddRule("ext", extRule, "[name]'s extension must be [ruleVal]")
	AddRule("notExt", notExtRule, "[name]'s extension must not be [ruleVal]")
	AddRule("extIn", extInRule, "[name]'s extension must be one of [ruleVal]")
	AddRule("extNotIn", extNotInRule, "[name]'s extension must not be one of [ruleVal]")
}
